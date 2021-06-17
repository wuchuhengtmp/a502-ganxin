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
 * 获取订单上关联的项目
 */
func (o *OrderService)GetProject() (p projects.Projects, err error) {
	err = model.DB.Model(&projects.Projects{}).Where("id = ?", o.ProjectId).First(&p).Error

	return
}


