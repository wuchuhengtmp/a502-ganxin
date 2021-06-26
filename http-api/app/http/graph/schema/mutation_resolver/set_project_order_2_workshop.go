/**
 * @Desc    型钢出库解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/22
 * @Listen  MIT
 */
package mutation_resolver

import (
	"context"
	"encoding/json"
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
	"http-api/app/models/repositories"
	"http-api/app/models/specificationinfo"
	"http-api/app/models/steel_logs"
	"http-api/app/models/steels"
	"http-api/pkg/model"
	"time"
)

func (*MutationResolver) SetProjectOrder2Workshop(ctx context.Context, input graphModel.ProjectOrder2WorkshopInput) (*orders.Order, error) {
	if err := requests.ValidateSetProjectOrder2WorkshopRequest(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	o := orders.Order{}
	err := model.DB.Model(&o).Where("id = ?", input.OrderID).First(&o).Error
	if err != nil {
		return  nil, errors.ServerErr(ctx, err)
	}
	// 标记型钢为出库状态
	err = model.DB.Transaction(func(tx *gorm.DB) error {
		steps := SetProjectOrder2WorkshopSteps{}
		// 创建订单物流
		orderExpress, err := steps.CreateExpressOrder(tx, ctx, input)
		if err != nil {
			return err
		}
		var steelsList []*steels.Steels
		err = tx.Model(&steels.Steels{}).Where("identifier in ?", input.IdentifierList).Find(&steelsList).Error
		if err != nil {
			return err
		}
		for _, steelItem := range steelsList {
			// 在订单规格中创建新记录
			newOrderSpecificationSteel, err := steps.CreateOrderSpecificationSteel(tx, ctx, input, steelItem, orderExpress)
			if err != nil {
				return err
			}
			err = tx.Model(&steels.Steels{}).Where("id = ?", steelItem.ID).
				// 标记型钢为【仓库】-运送至项目途中 状态
				Update("state", steels.StateRepository2Project).
				// 标记型钢当前应用在哪个订单规格型钢下
				Update("order_specification_steel_id", newOrderSpecificationSteel.Id).
				Error
			if err != nil {
				return err
			}

			// 添加型钢日志
			if err := steps.CreateSteelLog(tx, ctx, steelItem); err != nil {
				return err
			}
		}
		// 标记订单状态
		if err := steps.SetOrderOutState(tx, ctx, input, &o); err != nil {
			return err
		}
		// 发出消息通知项目管理员
		err = steps.CreateOutOfRepositoryMsg(tx, ctx, input)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	return &o, nil
}

/**
 * 解析器的解决步骤合集
 */
type SetProjectOrder2WorkshopSteps struct{}

/**
 * 创建物流单号
 */
func (*SetProjectOrder2WorkshopSteps) CreateExpressOrder(tx *gorm.DB, ctx context.Context, input graphModel.ProjectOrder2WorkshopInput) (*order_express.OrderExpress, error) {
	me := auth.GetUser(ctx)
	orderExpress := order_express.OrderExpress{
		OrderId:          input.OrderID,
		ExpressCompanyId: input.ExpressCompanyID,
		ExpressNo:        input.ExpressNo,
		SenderUid:        me.Id,
		CompanyId:        me.CompanyId,
		Direction:        order_express.OrderExpressDirectionToWorkshop,
	}
	err := tx.Create(&orderExpress).Error

	return &orderExpress, err
}

/**
 * 在订单规格中创建新记录
 */
func (*SetProjectOrder2WorkshopSteps) CreateOrderSpecificationSteel(
	tx *gorm.DB,
	ctx context.Context,
	input graphModel.ProjectOrder2WorkshopInput,
	steelItem *steels.Steels,
	orderExpress *order_express.OrderExpress,
) (*order_specification_steel.OrderSpecificationSteel, error) {
	me := auth.GetUser(ctx)
	orderSpecificationModel := order_specification.OrderSpecification{}
	orderSpecificationRecord := order_specification.OrderSpecification{}
	err := tx.Model(&orderSpecificationModel).Where("order_id = ? AND specification_id = ?", input.OrderID, steelItem.SpecificationId).First(&orderSpecificationRecord).Error
	if err != nil {
		return nil, err
	}
	orderSpecificationSteel := order_specification_steel.OrderSpecificationSteel{
		SteelId:              steelItem.ID,
		OrderSpecificationId: orderSpecificationRecord.Id,
		ToWorkshopExpressId:  orderExpress.Id,
		State: steels.StateRepository2Project,
		OutRepositoryAt:      time.Now(),
		EnterRepositoryUid:   me.Id,
	}
	err = tx.Create(&orderSpecificationSteel).Error

	return &orderSpecificationSteel, err
}

/**
 * 创建型钢的操作日志
 */
func (*SetProjectOrder2WorkshopSteps) CreateSteelLog(
	tx *gorm.DB,
	ctx context.Context,
	steelItem *steels.Steels,
) error {
	me := auth.GetUser(ctx)
	steelLog := steel_logs.SteelLog{
		Type:    steel_logs.OutSteelType,
		SteelId: steelItem.ID,
		Uid:     me.Id,
	}

	return tx.Create(&steelLog).Error
}

/**
 * 标记订单状态
 */
func (*SetProjectOrder2WorkshopSteps) SetOrderOutState(tx *gorm.DB, ctx context.Context, input graphModel.ProjectOrder2WorkshopInput, o *orders.Order) error {
	// 如果订单是已确认的情况下，则标记为发货状态
	if o.State == orders.StateConfirmed {
		err := tx.Model(o).Where("id = ?", input.OrderID).Update("state", orders.StateSend).Error
		if err != nil {
			return err
		}
	}

	return nil
}

/**
 * 添加消息
 */
func (SetProjectOrder2WorkshopSteps) CreateOutOfRepositoryMsg(tx *gorm.DB, ctx context.Context, input graphModel.ProjectOrder2WorkshopInput) error {
	var leaderList []*project_leader.ProjectLeader
	leaderTable := project_leader.ProjectLeader{}.TableName()
	projectTable := projects.Projects{}.TableName()
	orderTable := orders.Order{}.TableName()
	err := tx.Model(&project_leader.ProjectLeader{}).
		Select(fmt.Sprintf("%s.*", leaderTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.project_id", projectTable, projectTable, leaderTable)).
		Joins(fmt.Sprintf("join %s ON %s.project_id = %s.id", orderTable, orderTable, projectTable)).
		Where(fmt.Sprintf("%s.id = %d", orderTable, input.OrderID)).
		Scan(&leaderList).
		Error
	if err != nil {
		return err
	}
	projectItem := projects.Projects{}
	err = tx.Model(&projectItem).
		Select(fmt.Sprintf("%s.*", projectTable)).
		Joins(fmt.Sprintf("join %s ON %s.project_id = %s.id", orderTable, orderTable, projectTable)).
		Where(fmt.Sprintf("%s.id = %d", orderTable, input.OrderID)).
		First(&projectItem).
		Error
	if err != nil {
		return err
	}
	repositoryTable := repositories.Repositories{}.TableName()
	repositoryItem := repositories.Repositories{}
	err = tx.Model(&repositoryItem).Select(fmt.Sprintf("%s.*", repositoryTable)).
		Joins(fmt.Sprintf("join %s ON %s.repository_id = %s.id", orderTable, orderTable, repositoryTable)).
		Where(fmt.Sprintf("%s.id = %d", orderTable, input.OrderID)).
		Scan(&repositoryItem).
		Error
	if err != nil {
		return err
	}
	var weightInfo struct {
		Weight float64
	}
	specificationTable := specificationinfo.SpecificationInfo{}.TableName()
	steelTable := steels.Steels{}.TableName()
	err = tx.Model(&steels.Steels{}).
		Select("SUM(weight) as Weight").
		Joins(fmt.Sprintf("join %s ON %s.id = %s.specification_id", specificationTable, specificationTable, steelTable)).
		Where("identifier in ?", input.IdentifierList).Scan(&weightInfo).Error
	if err != nil {
		return err
	}

	content := fmt.Sprintf(
		"%s 项目于%s, 在%s出库,总数: %d根, %.2f吨, 请准备接收",
		projectItem.Name,
		helper.Time2Str(time.Now()),
		repositoryItem.Name,
		len(input.IdentifierList),
		weightInfo.Weight,
	)
	projectItemJsonBytes, _ := json.Marshal(projectItem)

	for _, projectLeaer := range leaderList {
		msgInstance := msg.Msg{
			IsRead:  false,
			Content: content,
			Uid:     projectLeaer.Uid,
			Type:    msg.OutProject2Workshop,
			Extends: string(projectItemJsonBytes),
		}
		err = msgInstance.CreateSelf(tx)
		if err != nil {
			return err
		}
	}

	return nil
}
