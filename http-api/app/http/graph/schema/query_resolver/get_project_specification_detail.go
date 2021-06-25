/**
 * @Desc    获取项目规格列表解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/24
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

type GetProjectSpecificationDetailSteps struct{}

func (*QueryResolver) GetProjectSpecificationDetail(ctx context.Context, input graphModel.GetProjectSpecificationDetailInput) (*projects.GetProjectSpecificationDetailRes, error) {
	var res projects.GetProjectSpecificationDetailRes
	if err := requests.ValidateGetProjectSpecificationDetailRequest(ctx, input); err != nil {
		return &res, errors.ValidateErr(ctx, err)
	}
	steps := GetProjectSpecificationDetailSteps{}
	// 获取列表
	list, err := steps.GetSpecificationList(ctx, input)
	if err != nil {
		return &res, errors.ServerErr(ctx, err)
	}
	res.List = list
	res.Total = int64(len(res.List))
	// 获取重量
	totalWeight, err := steps.GetSpecificationTotalWeight(input)
	res.Weight = totalWeight

	return &res, nil
}

/**
 * 获取项目规格列表
 */
func (*GetProjectSpecificationDetailSteps) GetSpecificationList(ctx context.Context, input graphModel.GetProjectSpecificationDetailInput) (res []*order_specification.OrderSpecification, err error) {
	orderSpecificationTable := order_specification.OrderSpecification{}.TableName()
	orderTable := orders.Order{}.TableName()
	orderSpecification := order_specification.OrderSpecification{}
	projectTable := projects.Projects{}.TableName()
	err = model.DB.Model(&orderSpecification).
		Select(fmt.Sprintf("%s.*", orderSpecificationTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.order_id", orderTable, orderTable, orderSpecificationTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.project_id", projectTable, projectTable, orderTable)).
		Where(fmt.Sprintf("%s.id = ?", projectTable), input.ProjectID).
		Find(&res).
		Error

	return
}

/**
 * 获取总重量(已经扫描过的)
 */
func (*GetProjectSpecificationDetailSteps) GetSpecificationTotalWeight(input graphModel.GetProjectSpecificationDetailInput) (totalWeight float64, err error) {
	orderSpecificationTable := order_specification.OrderSpecification{}.TableName()
	orderTable := orders.Order{}.TableName()
	orderSpecification := order_specification.OrderSpecification{}
	projectTable := projects.Projects{}.TableName()
	orderSpecificationSteelTable := order_specification_steel.OrderSpecificationSteel{}.TableName()
	specificationInfoTable := specificationinfo.SpecificationInfo{}.TableName()
	steelTable := steels.Steels{}.TableName()
	var weightInfo struct {
		TotalWeight float64
	}
	beScanStateList := []int64{
		steels.StateInStore,              //【仓库】-在库
		steels.StateProjectWillBeUsed,    //【项目】-待使用
		steels.StateProjectInUse,         //【项目】-使用中
		steels.StateProjectException,     //【项目】-异常
		steels.StateProjectIdle,          //【项目】-闲置
		steels.StateProjectWillBeStore,   //【项目】-准备归库
		steels.StateProjectOnTheStoreWay, //【项目】-归库途中
	}
	err = model.DB.Model(&orderSpecification).
		Select("SUM(weight) as TotalWeight").
		Joins(fmt.Sprintf("join %s ON %s.id = %s.order_id", orderTable, orderTable, orderSpecificationTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.project_id", projectTable, projectTable, orderTable)).
		Joins(fmt.Sprintf("join %s ON %s.order_specification_id = %s.id", orderSpecificationSteelTable, orderSpecificationSteelTable, orderSpecificationTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.steel_id", steelTable, steelTable, orderSpecificationSteelTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.specification_id", specificationInfoTable, specificationInfoTable, steelTable)).
		Where(fmt.Sprintf("%s.id = ?", projectTable), input.ProjectID).
		Where(fmt.Sprintf("%s.state in ?", orderSpecificationSteelTable), beScanStateList).
		Scan(&weightInfo).
		Error
	totalWeight = weightInfo.TotalWeight

	return
}
