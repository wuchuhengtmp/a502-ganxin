/**
 * @Desc    获取订单列表
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/17
 * @Listen  MIT
 */
package services

import (
	"context"
	"fmt"
	"http-api/app/http/graph/auth"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/models/orders"
	"http-api/app/models/project_leader"
	"http-api/app/models/projects"
	"http-api/app/models/roles"
	"http-api/pkg/model"
)

func GetOrderList(ctx context.Context, input graphModel.GetOrderListInput) (orderList []*orders.Order, err error) {
	me := auth.GetUser(ctx)
	role, err := me.GetRole()
	isDevice := auth.IsDevice(ctx)
	if err != nil {
		return
	}
	// 手持设备查看
	if isDevice {
		projectLeaderTable := project_leader.ProjectLeader{}.TableName()
		projectTable := projects.Projects{}.TableName()
		orderTable := orders.Order{}.TableName()
		whereMap := ""
		// 项目管理员的手持设备 只看到他自己项目下的订单
		if role.Tag == roles.RoleProjectAdmin {
			//  确认订单条件
			if *input.QueryType == graphModel.GetOrderListInputTypeConfirmOrder {
				whereMap = fmt.Sprintf("%s.state >= %d", orderTable, orders.StateConfirmed)
			} else {
				// 未确认订单条件
				whereMap = fmt.Sprintf("%s.state < %d", orderTable, orders.StateConfirmed)
			}
			err = model.DB.Model(&orders.Order{}).
				Select(fmt.Sprintf("%s.*", orders.Order{}.TableName())).
				Joins(fmt.Sprintf("join %s ON %s.id = %s.project_id", projectTable, projectTable, orderTable)).
				Joins(fmt.Sprintf("join %s ON %s.project_id = %s.id", projectLeaderTable, projectLeaderTable, projectTable)).
				Where(fmt.Sprintf("%s.uid = %d", projectLeaderTable, me.Id)).
				Where(fmt.Sprintf("%s.company_id = %d", projectTable, me.CompanyId)).
				Where(whereMap).
				Find(&orderList).
				Error
		} else {
			//  确认订单条件
			if *input.QueryType == graphModel.GetOrderListInputTypeConfirmOrder {
				whereMap = fmt.Sprintf("state >= %d", orders.StateConfirmed)
			} else {
				// 未确认订单条件
				whereMap = fmt.Sprintf("state < %d", orders.StateConfirmed)
			}
			err = model.DB.Model(&orders.Order{}).
				Where("company_id = ?", me.CompanyId).
				Where(whereMap).
				Find(&orderList).Error
		}
	} else {
		// 后台查看
		err = model.DB.Model(&orders.Order{}).
			Where("company_id = ?", me.CompanyId).
			Find(&orderList).Error
	}

	return
}
