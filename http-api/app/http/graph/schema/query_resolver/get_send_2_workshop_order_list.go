/**
 * @Desc    获取可入场的订单列表解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/23
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"fmt"
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/errors"
	"http-api/app/models/orders"
	"http-api/app/models/project_leader"
	"http-api/app/models/projects"
	"http-api/pkg/model"
)

func (*QueryResolver) GetSend2WorkshopOrderList(ctx context.Context) ([]*orders.Order, error) {
	var orderList []*orders.Order
	orderTable := orders.Order{}.TableName()
	projectTable := projects.Projects{}.TableName()
	projectLeaderTable := project_leader.ProjectLeader{}.TableName()
	me := auth.GetUser(ctx)
	orderStateList := []orders.StateCode{
		orders.StateSend,
		orders.StatePartOfReceipted,
	}
	err := model.DB.Debug().Model(&orders.Order{}).
		Select(fmt.Sprintf("%s.*", orderTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.project_id", projectTable, projectTable, orderTable)).
		Joins(fmt.Sprintf("join %s ON %s.project_id = %s.id", projectLeaderTable, projectLeaderTable, projectTable)).
		Where(fmt.Sprintf("%s.uid = %d", projectLeaderTable, me.Id)).
		Where(fmt.Sprintf("%s.state IN ?", orderTable), orderStateList).
		Find(&orderList).
		Error
	if err != nil {
		return orderList, errors.ServerErr(ctx, err)
	}

	return orderList, nil
}
