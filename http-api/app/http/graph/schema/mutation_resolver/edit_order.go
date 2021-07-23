/**
 * @Desc    The mutation_resolver is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/23
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
	"http-api/app/models/logs"
	"http-api/app/models/msg"
	"http-api/app/models/order_specification"
	"http-api/app/models/orders"
	"http-api/app/models/repository_leader"
	"http-api/app/models/roles"
	"http-api/app/models/specificationinfo"
	"http-api/pkg/model"
)

func (*MutationResolver) EditOrder(ctx context.Context, input graphModel.EditOrderInput) (*orders.Order, error) {
	if err := requests.ValidateEditOrderRequest(ctx, input); err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	me := auth.GetUser(ctx)
	myRole, _ := me.GetRole()
	orderItem := orders.Order{}
	err := model.DB.Transaction(func(tx *gorm.DB) error {
		steps := EditOrderSteps{}
		// 编辑订单
		if err := steps.EditOrder(ctx, tx, input); err != nil {
			return err
		}
		// 编辑详情
		oldTotal, oldWeight, newTotal, newWeight, err := steps.EditOrderDetail(tx, input)
		if err != nil {
			return err
		}
		if err := tx.Model(&orderItem).Where("id = ?", input.ID).First(&orderItem).Error; err != nil {
			return err
		}
		// 添加消息
		var leaders []*repository_leader.RepositoryLeader
		leaderItem := repository_leader.RepositoryLeader{}
		if err := tx.Model(&leaderItem).Where("repository_id = ?", orderItem.RepositoryId).Find(&leaders).Error; err != nil {
			return err
		}
		for _, i := range leaders {
			msgItem := msg.Msg{
				Content: fmt.Sprintf(
					"%s %s 编辑了订单:%s，总数量: %d根 -> %d根, %.3f吨 -> %.3f吨, 对方电话:%s",
					myRole.Name,
					me.Name,
					orderItem.OrderNo,
					oldTotal,
					newTotal,
					oldWeight,
					newWeight,
					me.Phone,
				),
				Uid:  i.Uid,
				Type: msg.EditOrder,
			}
			if err := tx.Create(&msgItem).Error; err != nil {
				return err
			}
		}
		// 添加操作记录
		if err := steps.CrateLog(ctx, tx, input); err != nil {
			return err
		}
		if err := tx.Model(&orders.Order{}).Where("id = ?", input.ID).First(&orderItem).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	return &orderItem, nil
}

type EditOrderSteps struct{}

/**
 * 编辑订单
 */
func (EditOrderSteps) EditOrder(ctx context.Context, tx *gorm.DB, input graphModel.EditOrderInput) error {
	orderItem := orders.Order{}
	err := tx.Model(&orderItem).Where("id = ?", input.ID).First(&orderItem).Error
	if err != nil {
		return err
	}
	err = tx.Model(&orderItem).Where("id = ?", input.ID).
		Update("expected_return_at", input.ExpectedReturnAt).Error
	if err != nil {
		return nil
	}
	if input.Remark != nil {
		err := tx.Model(&orders.Order{}).Where("id = ?", input.ID).
			Update("remark", *input.Remark).
			Error
		if err != nil {
			return err
		}
	}

	return nil
}

/**
 * 编辑订单详情
 */
func (EditOrderSteps) EditOrderDetail(tx *gorm.DB, input graphModel.EditOrderInput) (oldTotal int64, oldWeight float64, newTotal int64, newWeight float64, err error) {
	orderSpecificationItem := order_specification.OrderSpecification{}
	specificationTable := specificationinfo.SpecificationInfo{}.TableName()
	// 收集修改前的汇总
	var oldInfo struct {
		Total  int64
		Weight float64
	}
	err = tx.Model(&orderSpecificationItem).
		Select("sum(total) as Total, sum(weight) as Weight").
		Joins(fmt.Sprintf("join %s ON %s.id = %s.specification_id ", specificationTable, specificationTable, orderSpecificationItem.TableName())).
		Where(fmt.Sprintf("%s.order_id = ?", orderSpecificationItem.TableName()), input.ID).
		Scan(&oldInfo).
		Error
	if err != nil {
		return
	}
	oldTotal = oldInfo.Total
	oldWeight = oldInfo.Weight
	// 删除
	err = tx.Exec(fmt.Sprintf(
		"DELETE %s WHERE order_id = %d",
		orderSpecificationItem.TableName(),
		input.ID,
	)).Error
	if err != nil {
		return
	}
	// 添加新的
	for _, i := range input.SteelList {
		newTotal += i.Total
		specificationItem := specificationinfo.SpecificationInfo{}
		if err = tx.Model(&specificationItem).Where("id = ?", i.SpecificationID).First(&specificationItem).Error; err != nil {
			return
		}
		newWeight += specificationItem.Weight * float64(i.Total)
		err = tx.Create(&order_specification.OrderSpecification{
			SpecificationId: i.SpecificationID,
			Total:           i.Total,
			Specification:   specificationItem.GetSelfSpecification(),
			OrderId:         input.ID,
		}).Error
		if err != nil {
			return
		}
	}

	return
}
/**
 * 添加操作日志
 */
func (EditOrderSteps) CrateLog(ctx context.Context, tx *gorm.DB, input graphModel.EditOrderInput) error {
	me := auth.GetUser(ctx)
	roleItem := roles.Role{}
	if err := tx.Model(&roleItem).Where( "id = ?", me.RoleId).First(&roleItem).Error; err != nil {
		return err
	}
	orderItem := orders.Order{}
	if err := tx.Model(&orderItem).Where("id = ?", input.ID).First(&orderItem).Error; err != nil {
		return err
	}
	logItem := logs.Logos{
		Content: fmt.Sprintf("%s %s 编辑订单:%s",
			roleItem.Name,
			me.Name,
			orderItem.OrderNo,
		),
		Uid: me.Id,
	}
	if err := tx.Create(&logItem).Error; err != nil {
		return err
	}

	return nil
}
