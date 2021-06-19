/**
 * @Desc    订单相关服务
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/17
 * @Listen  MIT
 */
package services

import (
	"context"
	"gorm.io/gorm"
	"http-api/app/http/graph/auth"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/models/orders"
	"http-api/app/models/projects"
	"http-api/pkg/model"
)

type OrderService orders.Order
/**
 * 把订单数量赋值给订单服务
 */
func (o *OrderService) OrderMoveIntoMe(guest *orders.Order) {
	o.Id = guest.Id
	o.ProjectId = guest.ProjectId
	o.RepositoryId = guest.RepositoryId
	o.State = guest.State
	o.ExpectedReturnAt = guest.ExpectedReturnAt
	o.PartList = guest.PartList
	o.CreateUid = guest.CreateUid
	o.ConfirmedAt = guest.ConfirmedAt
	o.ReceiveAt = guest.ReceiveAt
	o.OrderNo = guest.OrderNo
	o.Remark = guest.Remark
	o.DeletedAt = guest.DeletedAt
	o.CreatedAt = guest.CreatedAt
	o.UpdatedAt = guest.UpdatedAt
}

/**
 * 获取订单上关联的项目
 */
func (o *OrderService) GetProject() (p projects.Projects, err error) {
	err = model.DB.Model(&projects.Projects{}).Where("id = ?", o.ProjectId).First(&p).Error

	return
}
/**
 * 确认订单
 */
func ConfirmOrRejectOrder(ctx context.Context, input graphModel.ConfirmOrderInput) (*orders.Order, error) {
	o := orders.Order{Id: input.ID}
	err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 修改订单状态
		_ = o.GetSelf()
		if input.IsAccess {
			o.State = orders.StateConfirmed
		} else {
			o.State = orders.StateRejected
		}
		me := auth.GetUser(ctx)
		o.ConfirmedUid = me.Id
		err := tx.Model(&orders.Order{}).
			Where("id = ?", input.ID).
			Update("state", o.State).
			Update("confirmed_uid", o.ConfirmedUid).Error
		if err != nil {
			return err
		}
		// 添加消息
		if err := CreateConfirmOrRejectOrderMsg(tx, &o); err != nil {
			return err
		}
		return nil
	})

	return &o, err
}