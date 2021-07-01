/**
 * @Desc    获取可归库的项目列表解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/29
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"fmt"
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/errors"
	"http-api/app/models/order_specification"
	"http-api/app/models/order_specification_steel"
	"http-api/app/models/orders"
	"http-api/app/models/projects"
	"http-api/app/models/repositories"
	"http-api/app/models/repository_leader"
	"http-api/app/models/steels"
	"http-api/pkg/model"
)

func (*QueryResolver) GetEnterRepositoryProjectList(ctx context.Context) (res []*projects.Projects, err error) {
	me := auth.GetUser(ctx)
	projectTable := projects.Projects{}.TableName()
	orderTable := orders.Order{}.TableName()
	repositoryTable := repositories.Repositories{}.TableName()
	repositoryLeaderTable := repository_leader.RepositoryLeader{}.TableName()
	orderSpecificationTable := order_specification.OrderSpecification{}.TableName()
	orderSpecificationSteelTable := order_specification_steel.OrderSpecificationSteel{}.TableName()
	var orderItems []*orders.Order
	err = model.DB.Model(&orders.Order{}).
		Select(fmt.Sprintf("%s.*", orders.Order{}.TableName())).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.repository_id", repositoryTable, repositoryTable, orderTable)).
		Joins(fmt.Sprintf("join %s ON %s.repository_id = %s.id", repositoryLeaderTable, repositoryLeaderTable, repositoryTable)).
		Joins(fmt.Sprintf("join %s ON %s.order_id = %s.id", orderSpecificationTable, orderSpecificationTable, orderTable)).
		Joins(fmt.Sprintf("join %s ON %s.order_specification_id = %s.id", orderSpecificationSteelTable, orderSpecificationSteelTable, orderSpecificationTable)).
		Where(fmt.Sprintf("%s.uid = ?", repositoryLeaderTable), me.Id).
		Where(fmt.Sprintf("%s.state = ?", orderSpecificationSteelTable), steels.StateProjectOnTheStoreWay).
		Where(fmt.Sprintf("%s.company_id = ?", orderTable), me.CompanyId).
		Find(&orderItems).
		Error
	if err != nil {
		return res, errors.ServerErr(ctx, err)
	}
	orderIdMapId := make(map[int64]int64)
	var orderIds  []int64
	for _, o := range orderItems {
		if _, ok := orderIdMapId[o.Id]; !ok {
			orderIds = append(orderIds, o.Id)
			orderIdMapId[o.Id] = o.Id
		}
	}
	err = model.DB.Model(&projects.Projects{}).
		Joins(fmt.Sprintf("join %s ON %s.project_id = %s.id", orderTable, orderTable, projectTable)).
		Where(fmt.Sprintf("%s.id in ?", orderTable ), orderIds).
		Find(&res).Error
	if err != nil {
		return res, errors.ServerErr(ctx, err)
	}

	return
}
