/**
 * @Desc    获取待归库详情解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/30
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
	"http-api/pkg/model"
)

func (*QueryResolver)GetToBeEnterRepositoryDetail(ctx context.Context, input graphModel.GetToBeEnterRepositoryDetailInput) (res []*order_specification_steel.OrderSpecificationSteel, err error) {
	if err := requests.ValidateGetToBeEnterRepositoryDetailRequest(ctx, input); err != nil {
		return nil, errors.ServerErr(ctx,  err)
	}
	orderSpecificationSteelItem := order_specification_steel.OrderSpecificationSteel{}
	orderSpecificationTable := order_specification.OrderSpecification{}.TableName()
	orderTable := orders.Order{}.TableName()
	projectTable := projects.Projects{}.TableName()
	modelInstance := model.DB.Model(&orderSpecificationSteelItem).
		Select(fmt.Sprintf("%s.*", orderSpecificationSteelItem.TableName())).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.order_specification_id", orderSpecificationTable, orderSpecificationTable, orderSpecificationSteelItem.TableName())).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.order_id", orderTable, orderTable, orderSpecificationTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.project_id", projectTable, projectTable, orderTable)).
		Where(fmt.Sprintf("%s.id = ?", projectTable), input.ProjectID)
	if input.State != nil {
		modelInstance = modelInstance.Where(fmt.Sprintf("%s.state = ?", orderSpecificationSteelItem.TableName()), *input.State)
	}
	if input.SpecificationID != nil  {
		modelInstance = modelInstance.Where(fmt.Sprintf("%s.specification_id = ?", orderSpecificationTable), *input.SpecificationID)

	}
	err = modelInstance.Find(&res).Error
	if err != nil {
		return res, errors.ServerErr(ctx, err)
	}


	return
}

