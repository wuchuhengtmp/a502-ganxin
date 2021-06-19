/**
 * @Desc    获取出库的订单列表解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/19
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"fmt"
	"http-api/app/http/graph/auth"
	"http-api/app/models/orders"
	"http-api/app/models/project_leader"
	"http-api/app/models/projects"
	"http-api/pkg/model"
)

func (*QueryResolver) GetTobeSendWorkshopOrderList(ctx context.Context) (orderList []*orders.Order, err error) {
	orderModel := orders.Order{}
	me := auth.GetUser(ctx)
	orderTable := orders.Order{}.TableName()
	projectTable := projects.Projects{}.TableName()
	projectLeaderTable := project_leader.ProjectLeader{}.TableName()
	err = model.DB.Model(&orderModel).
		Select(fmt.Sprintf("%s.*", orderTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.project_id", projectTable, projectTable, orderTable)).
		Joins(fmt.Sprintf("join %s ON %s.project_id = %s.id", projectLeaderTable, projectLeaderTable, projectTable)).
		Where(fmt.Sprintf("%s.company_id = ? AND %s.state in ?", orderTable, orderTable), me.CompanyId, []int64{orders.StateConfirmed, orders.StatePartOfReceipted}).
		Where(fmt.Sprintf("%s.uid = %d", projectLeaderTable, me.Id)).
		Find(&orderList).
		Error

	return
}
