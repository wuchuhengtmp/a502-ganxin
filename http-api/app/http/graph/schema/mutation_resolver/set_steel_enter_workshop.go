/**
 *l@Desc    型钢入场解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/24
 * @Listen  MIT
 */
package mutation_resolver

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/errors"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/http/graph/util/helper"
	"http-api/app/models/msg"
	"http-api/app/models/order_express"
	"http-api/app/models/order_specification"
	"http-api/app/models/order_specification_steel"
	"http-api/app/models/orders"
	"http-api/app/models/project_leader"
	"http-api/app/models/projects"
	"http-api/app/models/specificationinfo"
	"http-api/app/models/steel_logs"
	"http-api/app/models/steels"
	"http-api/pkg/model"
	"time"
)

func (*MutationResolver) SetSteelEnterWorkshop(ctx context.Context, input graphModel.SetSteelIntoWorkshopInput) ([]*order_specification_steel.OrderSpecificationSteel, error) {
	var orderSpecificationSteelList []*order_specification_steel.OrderSpecificationSteel
	if err := requests.ValidateSetSteelIntoWorkshopRequest(ctx, input); err != nil {
		return orderSpecificationSteelList, errors.ValidateErr(ctx, err)
	}
	steps := SetSteelIntoWorkshopSteps{}
	err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 标记订单状态--收货或部分收货
		if err := steps.FlagOrder(tx, ctx, input); err != nil {
			return err
		}
		// 修改物流状态 收货人
		if err := steps.ChangeExpressState(tx, ctx, input); err != nil {
			return err
		}
		// 标记每根型钢状态--入场状态(待使用)入场时间和入场人
		if err := steps.FlatSteel(tx, ctx, input); err != nil {
			return err
		}
		// 添加消息
		if err := steps.PushMsg(tx, ctx, input); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return orderSpecificationSteelList, errors.ServerErr(ctx, err)
	}
	// 获取响应的列表
	if res, err := steps.GetRes(ctx, input); err == nil {
		return res, nil
	} else {
		return orderSpecificationSteelList, errors.ServerErr(ctx, err)
	}
}

/**
 * 型钢入场操作步骤
 */
type SetSteelIntoWorkshopSteps struct{}

/**
 * 标记每根型钢状态--入场状态(待使用)入场时间和入场人
 */
func (*SetSteelIntoWorkshopSteps) FlatSteel(tx *gorm.DB, ctx context.Context, input graphModel.SetSteelIntoWorkshopInput) error {
	me := auth.GetUser(ctx)
	orderSpecificationSteelTable := order_specification_steel.OrderSpecificationSteel{}.TableName()
	steelTable := steels.Steels{}.TableName()
	for _, identifier := range input.IdentifierList {
		orderSpecificationSteelItem := order_specification_steel.OrderSpecificationSteel{}
		err := tx.Model(&order_specification_steel.OrderSpecificationSteel{}).
			Select(fmt.Sprintf("%s.*", orderSpecificationSteelTable)).
			Joins(fmt.Sprintf("join %s ON %s.id = %s.steel_id", steelTable, steelTable, orderSpecificationSteelTable)).
			Where(fmt.Sprintf("%s.identifier = ?", steelTable), identifier).
			First(&orderSpecificationSteelItem).
			Error
		if err != nil {
			return err
		}
		// 修改订单型钢的状态和相关参数
		orderSpecificationSteelItem.EnterWorkshopAt = time.Now()
		orderSpecificationSteelItem.EnterRepositoryUid = me.Id
		err = tx.Model(&orderSpecificationSteelItem).Where("id = ?", orderSpecificationSteelItem.Id).
			// 入场用户
			Update("enter_workshop_uid", orderSpecificationSteelItem.EnterRepositoryUid).
			// 标记项目订单规格中的型钢的状态--待使用
			Update("state", steels.StateProjectWillBeUsed).
			// 入场时间
			Update("enter_workshop_at", orderSpecificationSteelItem.EnterRepositoryAt).
			Error
		if err != nil {
			return err
		}

		// 记录型钢日志
		steelLosItem := steel_logs.SteelLog{
			Uid:     me.Id,
			Type:    steel_logs.EnterWorkshopType,
			SteelId: orderSpecificationSteelItem.SteelId,
		}
		if err := tx.Create(&steelLosItem).Error; err != nil {
			return err
		}

	}
	// 标记型钢当前的状态--待使用
	err := tx.Model(&steels.Steels{}).
		Where("identifier in ?", input.IdentifierList).
		Update("state", steels.StateProjectWillBeUsed).Error

	if err != nil {
		return err
	}

	return nil
}

/**
 * 标记订单状态--收货或部分收货
 */
func (*SetSteelIntoWorkshopSteps) FlagOrder(tx *gorm.DB, ctx context.Context, input graphModel.SetSteelIntoWorkshopInput) error {
	// 修改订单状态 是否部分收货或全部收货
	orderItem := orders.Order{}
	if err := tx.Model(&orderItem).Where("id = ?", input.OrderID).Error; err != nil {
		return err
	}
	orderItem.State = orders.StateReceipted
	var orderSpecificationList []*order_specification.OrderSpecification
	if err := tx.Model(&order_specification.OrderSpecification{}).
		Where("order_id = ?", input.OrderID).
		Find(&orderSpecificationList).Error; err != nil {
		return err
	}
	for _, orderSpecificationItem := range orderSpecificationList {
		var orderSpecificationSteelList []*order_specification_steel.OrderSpecificationSteel
		err := tx.Model(&order_specification_steel.OrderSpecificationSteel{}).
			Where("order_specification_id = ?", orderSpecificationItem.Id).
			Find(&orderSpecificationSteelList).
			Error
		if err != nil {
			return err
		}
		if int64(len(orderSpecificationSteelList)) < orderSpecificationItem.Total {
			orderItem.State = orders.StatePartOfReceipted
			break
		}
	}
	err := tx.Model(&orders.Order{}).
		Where("id = ?", input.OrderID).
		Update("state", orderItem.State).Error
	if err != nil {
		return err
	}

	return nil
}

/**
 *  修改物流状态
 */
func (*SetSteelIntoWorkshopSteps) ChangeExpressState(tx *gorm.DB, ctx context.Context, input graphModel.SetSteelIntoWorkshopInput) error {
	orderSpecificationSteelItem := order_specification_steel.OrderSpecificationSteel{}
	steelTable := steels.Steels{}.TableName()
	orderSpecificationSteelTable := order_specification_steel.OrderSpecificationSteel{}.TableName()
	err := tx.Model(&orderSpecificationSteelItem).
		Select(fmt.Sprintf("%s.*", orderSpecificationSteelTable)).
		Joins(fmt.Sprintf("join %s on %s.id = %s.steel_id", steelTable, steelTable, orderSpecificationSteelTable)).
		Where(fmt.Sprintf("%s.identifier = %s", steelTable, input.IdentifierList[0])).
		First(&orderSpecificationSteelItem).
		Error
	if err != nil {
		return err
	}
	expressItem := order_express.OrderExpress{}
	err = tx.Model(&expressItem).Where("id = ?", orderSpecificationSteelItem.ToWorkshopExpressId).
		First(&expressItem).
		Error
	if err != nil {
		return err
	}
	me := auth.GetUser(ctx)
	expressItem.ReceiveUid = me.Id
	expressItem.ReceiveAt = time.Now()
	if err := tx.Model(&expressItem).Where("id = ?", expressItem.Id).Update("receive_uid", me.Id).Error; err != nil {
		return err
	}
	if err := tx.Model(&expressItem).Where("id = ?", expressItem.Id).Update("receive_at", time.Now()).Error; err != nil {
		return err
	}

	return nil
}

/**
 * 消息通知 其它项目管理员
 */
func (*SetSteelIntoWorkshopSteps) PushMsg(tx *gorm.DB, ctx context.Context, input graphModel.SetSteelIntoWorkshopInput) error {
	projectItem := projects.Projects{}
	projectTable := projects.Projects{}.TableName()
	orderTable := orders.Order{}.TableName()
	projectLeaderTable := project_leader.ProjectLeader{}.TableName()
	err := tx.Model(&projectItem).Select(fmt.Sprintf("%s.*", projectTable)).
		Joins(fmt.Sprintf("join %s on %s.project_id = %s.id", orderTable, orderTable, projectTable)).
		Where(fmt.Sprintf("%s.id = %d", orderTable, input.OrderID)).
		First(&projectItem).Error
	if err != nil {
		return err
	}
	orderItem := orders.Order{}
	if err := tx.Model(&orderItem).Where("id = ?", input.OrderID).First(&orderItem).Error; err != nil {
		return err
	}
	var weightInfo struct{ TotalWeight float64 }

	specificationInfoTable := specificationinfo.SpecificationInfo{}.TableName()
	steelsTable := steels.Steels{}.TableName()

	err = tx.Model(&steels.Steels{}).
		Select("sum(weight) as TotalWeight").
		Joins(fmt.Sprintf("join %s ON %s.id = %s.specification_id", specificationInfoTable, specificationInfoTable, steelsTable)).
		Where("identifier in ?", input.IdentifierList).
		Scan(&weightInfo).
		Error
	if err != nil {
		return err
	}

	content := fmt.Sprintf(
		"%s 于 %s 入库一批型钢，订单编号为:%s, 总数为: %d根, %.2f吨, 请准备安装",
		projectItem.Name,
		helper.Time2Str(time.Now()),
		orderItem.OrderNo,
		len(input.IdentifierList),
		weightInfo.TotalWeight,
	)

	var leaderList []*project_leader.ProjectLeader
	err = tx.Model(&project_leader.ProjectLeader{}).Select(fmt.Sprintf("%s.*", projectLeaderTable)).
		Joins(fmt.Sprintf("join %s on %s.id = %s.project_id", projectTable, projectTable, projectLeaderTable)).
		Joins(fmt.Sprintf("join %s on %s.project_id = %s.id", orderTable, orderTable, projectTable)).
		Where(fmt.Sprintf("%s.id = ?", orderTable), input.OrderID).
		Find(&leaderList).
		Error
	if err != nil {
		return err
	}

	for _, leaderItem := range leaderList {
		msgItem := msg.Msg{
			IsRead:  false,
			Content: content,
			Uid:     leaderItem.Uid,
			Type:    msg.EnterProject2Workshop,
		}
		if err := msgItem.CreateSelf(tx); err != nil {
			return err
		}
	}

	return nil
}

/**
 * 获取要响应的数据
 */
func (*SetSteelIntoWorkshopSteps) GetRes(ctx context.Context, input graphModel.SetSteelIntoWorkshopInput) (res []*order_specification_steel.OrderSpecificationSteel, err error) {
	orderSpecificationSteelTable := order_specification_steel.OrderSpecificationSteel{}.TableName()
	steelTable := steels.Steels{}.TableName()
	err = model.DB.Model(order_specification_steel.OrderSpecificationSteel{}).Select(fmt.Sprintf("%s.*", orderSpecificationSteelTable)).
		Joins(fmt.Sprintf("join %s on %s.id = %s.steel_id", steelTable, steelTable, orderSpecificationSteelTable)).
		Where(fmt.Sprintf("%s.identifier in ?", steelTable), input.IdentifierList).
		Find(&res).Error

	return
}
