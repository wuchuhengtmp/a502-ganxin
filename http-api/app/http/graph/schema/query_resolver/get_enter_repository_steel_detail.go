/**
 * @Desc   项目归库的型钢查询解析器
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
	"http-api/app/http/graph/errors"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/order_specification"
	"http-api/app/models/order_specification_steel"
	"http-api/app/models/orders"
	"http-api/app/models/projects"
	"http-api/app/models/steels"
	"http-api/pkg/model"
)

func (*QueryResolver)GetEnterRepositorySteelDetail(ctx context.Context, input graphModel.GetEnterRepositorySteelDetailInput) (*projects.GetEnterRepositorySteelDetailRes, error) {
	if err := requests.ValidateGetEnterRepositorySteelDetailRequest(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	res := projects.GetEnterRepositorySteelDetailRes{}
	orderSteelItem := order_specification_steel.OrderSpecificationSteel{}
	steelTable := steels.Steels{}.TableName()
	model.DB.Model(&orderSteelItem).
		Select(fmt.Sprintf("%s.*", orderSteelItem.TableName())).
		Joins(fmt.Sprintf("join %s ON %s.order_specification_steel_id = %s.id", steelTable, steelTable, orderSteelItem.TableName())).
		Where(fmt.Sprintf("%s.identifier = ?", steelTable), input.Identifier).
		Where(fmt.Sprintf("%s.state = ?", orderSteelItem.TableName()), steels.StateProjectOnTheStoreWay).
		First(&orderSteelItem)
	res.OrderSteel = orderSteelItem
	// 已归库数量
	orderTable := orders.Order{}.TableName()
	orderSpecificationTable := order_specification.OrderSpecification{}.TableName()
	orderSpecificationSteelTable := order_specification_steel.OrderSpecificationSteel{}.TableName()
	orderSpecificationSteelItem := order_specification_steel.OrderSpecificationSteel{}
	err := model.DB.Model(&orderSpecificationSteelItem).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.order_specification_id", orderSpecificationTable, orderSpecificationTable, orderSpecificationSteelTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.order_id", orderTable, orderTable, orderSpecificationTable)).
		Where( fmt.Sprintf("%s.project_id = ?", orderTable), input.ProjectID).
		Where(fmt.Sprintf("%s.state = ?", orderSpecificationSteelTable), steels.StateInStore).
		Count(&res.StoredTotal).Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	// 待归库数量
	err = model.DB.Model(&orderSpecificationSteelItem).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.order_specification_id", orderSpecificationTable, orderSpecificationTable, orderSpecificationSteelTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.order_id", orderTable, orderTable, orderSpecificationTable)).
		Where( fmt.Sprintf("%s.project_id = ?", orderTable), input.ProjectID).
		Where(fmt.Sprintf("%s.state != ?", orderSpecificationSteelTable), steels.StateInStore).
		Count(&res.ToBeStoreTotal).Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	return &res, nil
}
