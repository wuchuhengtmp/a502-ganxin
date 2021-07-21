/**
 * @_desc    获取订单详情验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/9
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"fmt"
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/errors"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/order_specification"
	"http-api/app/models/order_specification_steel"
	"http-api/app/models/orders"
	"http-api/app/models/specificationinfo"
	"http-api/pkg/model"
)

func (*QueryResolver) GetOrderDetailForBackEnd(ctx context.Context, input graphModel.GetOrderDetailForBackEndInput) (*orders.GetOrderDetailForBackEndRes, error) {
	if err := requests.ValidateGetOrderDetailForBackEndRequest(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	orderItem := orders.Order{}
	orderTable := orderItem.TableName()
	orderSpecificationItem := order_specification.OrderSpecification{}
	orderSpecificationTable := order_specification.OrderSpecification{}.TableName()
	me := auth.GetUser(ctx)
	// 数量和重量
	orderSpecificationSteelItem := order_specification_steel.OrderSpecificationSteel{}
	recordTable := orderSpecificationSteelItem.TableName()
	specificationTable := specificationinfo.SpecificationInfo{}.TableName()
	summaryIn := model.DB.Model(&orderSpecificationSteelItem).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.order_specification_id", orderSpecificationTable, orderSpecificationTable, recordTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.order_id", orderTable, orderTable, orderSpecificationTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.specification_id", specificationTable, specificationTable, orderSpecificationTable)).
		Where(fmt.Sprintf("%s.company_id = ?", orderTable), me.CompanyId)
	// 订单列表
	modeIn := model.DB.Model(&orderSpecificationItem).
		Select(fmt.Sprintf("%s.*", orderSpecificationTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.order_id",orderTable,orderTable, orderSpecificationTable)).
		Where(fmt.Sprintf("%s.company_id = ?", orderTable), me.CompanyId)
	//  仓库过滤
	if input.RepositoryID != nil {
		modeIn = modeIn.Where(fmt.Sprintf("%s.repository_id = ?", orderTable), *input.RepositoryID)
		summaryIn = summaryIn.Where(fmt.Sprintf("%s.repository_id = ?", orderTable), *input.RepositoryID)
	}
	// 规格过滤
	if input.SpecificationID != nil {
		modeIn = modeIn.Where(fmt.Sprintf("%s.specification_id = ?", orderSpecificationTable), *input.SpecificationID)
		summaryIn = summaryIn.Where(fmt.Sprintf("%s.specification_id = ?", orderSpecificationTable), *input.SpecificationID)
	}
	// 项目过滤
	if input.ProjectID != nil {
		modeIn = modeIn.Where(fmt.Sprintf("%s.project_id = ?", orderTable), *input.ProjectID)
		summaryIn = summaryIn.Where(fmt.Sprintf("%s.project_id = ?", orderTable), *input.ProjectID)
	}
	// 订单编号过滤
	if input.OrderNo != nil {
		modeIn = modeIn.Where(fmt.Sprintf("%s.order_no = ?", orderTable), *input.OrderNo)
		summaryIn = summaryIn.Where(fmt.Sprintf("%s.order_no = ?", orderTable), *input.OrderNo)
	}
	// 分页
	if input.IsShowAll == false {
		o := (*input.Page - 1) * *input.PageSize
		modeIn = modeIn.Limit(int(*input.PageSize)).Offset(int(o))
	}
	res := orders.GetOrderDetailForBackEndRes{}
	if err := modeIn.Scan(&res.List).Error; err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	if err := summaryIn.Count(&res.Total).Error; err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	var weightInfo struct{ Weight float64 }
	if err := summaryIn.Select("sum(weight) as Weight").Scan(&weightInfo).Error; err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	res.Weight = weightInfo.Weight

	return &res, nil
}
