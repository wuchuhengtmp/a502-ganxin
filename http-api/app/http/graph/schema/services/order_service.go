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
	o.ReceiveUid = guest.ReceiveUid
	o.ReceiveAt = guest.ReceiveAt
	o.ExpressCompanyId = guest.ExpressCompanyId
	o.ExpressNo = guest.ExpressNo
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
