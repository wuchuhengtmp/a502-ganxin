/**
 * @Desc    获取项目详情解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/25
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
	"http-api/pkg/model"
)

func (*QueryResolver) GetProjectSteelDetail(ctx context.Context, input graphModel.GetProjectSteelDetailInput) (*projects.GetProjectSteelDetailRes, error) {
	var p projects.GetProjectSteelDetailRes
	if err := requests.ValidateGetProjectSteelDetailRequest(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	steps := GetProjectSteelDetailSteps{}
	list, weight, err := steps.GetList(ctx, input)
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	p.List = list
	p.Total = int64(len(list))
	p.Weight = weight

	return &p, nil
}

/**
 * 获取项目详情解决步骤
 */
type GetProjectSteelDetailSteps struct{}

/**
 * 获取型钢详情列表
 */
func (*GetProjectSteelDetailSteps) GetList(ctx context.Context, input graphModel.GetProjectSteelDetailInput) (list []*order_specification_steel.OrderSpecificationSteel, weight float64, err error) {
	orderSpecificationSteelItem := order_specification_steel.OrderSpecificationSteel{}
	specificationInfoTable := specificationinfo.SpecificationInfo{}.TableName()
	orderSpecificationTable := order_specification.OrderSpecification{}.TableName()
	orderTable := orders.Order{}.TableName()
	projectTable := projects.Projects{}.TableName()
	orderSpecificationSteelTable := order_specification_steel.OrderSpecificationSteel{}.TableName()
	whereMap := fmt.Sprintf("%s.id = %d ", projectTable, input.ProjectID)
	if input.SpecificationID != nil {
		whereMap += fmt.Sprintf(" AND %s.specification_id = %d ", orderSpecificationTable, *input.SpecificationID)
	}
	if input.State != nil {
		whereMap += fmt.Sprintf(" AND %s.state = %d", orderSpecificationSteelTable, *input.State)
	}
	queryIns := model.DB.Model(&orderSpecificationSteelItem).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.order_specification_id", orderSpecificationTable, orderSpecificationTable, orderSpecificationSteelItem.TableName())).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.order_id", orderTable, orderTable, orderSpecificationTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.project_id", projectTable, projectTable, orderTable)).
		Where(whereMap)

	// 获取数量
	err = queryIns.
		Select(fmt.Sprintf("%s.*", orderSpecificationSteelTable)).
		Find(&list).
		Error
	if err != nil {
		return
	}
	//  获取重量
	var weightInfo struct{
		Weight float64
	}
	err = queryIns.
		Joins(fmt.Sprintf("join %s ON %s.id = %s.specification_id", specificationInfoTable, specificationInfoTable,orderSpecificationTable)).
		Select("sum(weight) as Weight").
		Scan(&weightInfo).
		Error
	if err != nil {
		return
	}
	weight = weightInfo.Weight

	return
}

