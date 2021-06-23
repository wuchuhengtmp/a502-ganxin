/**
 * @Desc    获取待入场的订单的详情接口
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
	"http-api/app/http/graph/errors"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/order_specification"
	"http-api/app/models/order_specification_steel"
	"http-api/app/models/orders"
	"http-api/app/models/projects"
	"http-api/app/models/specificationinfo"
	"http-api/app/models/steels"
	"http-api/pkg/model"
)

func (*QueryResolver) GetSend2WorkshopOrderListDetail(ctx context.Context, input graphModel.GetProjectOrder2WorkshopDetailInput) (*projects.GetSend2WorkshopOrderListDetailRes, error) {
	var res projects.GetSend2WorkshopOrderListDetailRes
	if err := requests.ValidateGetSend2WorkshopOrderListDetail(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	orderSpecificationSteelTable := order_specification_steel.OrderSpecificationSteel{}.TableName()
	orderSpecificationTable := order_specification.OrderSpecification{}.TableName()
	steelTable := steels.Steels{}.TableName()
	orderTable := orders.Order{}.TableName()
	specificationTable := specificationinfo.SpecificationInfo{}.TableName()
	whereMap := fmt.Sprintf("%s.id = %d", orderTable, input.OrderID)
	if input.SpecificationID != nil {
		whereMap += fmt.Sprintf(" AND %s.specification_id = %d", steelTable, input.SpecificationID)
	}

	// 订单规格型钢列表
	err := model.DB.Model(&order_specification_steel.OrderSpecificationSteel{}).
		Select(fmt.Sprintf("%s.*", orderSpecificationSteelTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.order_Specification_id", orderSpecificationTable, orderSpecificationTable, orderSpecificationSteelTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.order_id", orderTable, orderTable, orderSpecificationTable)).
		Where(whereMap).
		Find(&res.List).
		Error
	if err != nil {
		return &res, errors.ServerErr(ctx, err)
	}
	// 数量
	res.Total = int64(len(res.List))
	// 重量
	var weightInfo struct {
		Weight float64
	}
	err = model.DB.Debug().Model(&order_specification_steel.OrderSpecificationSteel{}).
		Select(fmt.Sprintf("SUM(weight) as Weight")).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.order_Specification_id", orderSpecificationTable, orderSpecificationTable, orderSpecificationSteelTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.order_id", orderTable, orderTable, orderSpecificationTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.steel_id", steelTable, steelTable, orderSpecificationSteelTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.specification_id", specificationTable, specificationTable, steelTable)).
		Where(whereMap).
		Scan(&weightInfo).
		Error
	res.TotalWeight = weightInfo.Weight

	return &res, nil
}
