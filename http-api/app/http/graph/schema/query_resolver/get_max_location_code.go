/**
 * @Desc    获取项目最大的安装码解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/26
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

func (*QueryResolver)GetMaxLocationCode(ctx context.Context, input graphModel.GetMaxLocationCodeInput) (int64, error) {
	if err := requests.ValidateGetMaxLocationCodeRequest(ctx, input); err != nil {
		return 0, errors.ServerErr(ctx, err)
	}
	orderSpecificationSteelTable := order_specification_steel.OrderSpecificationSteel{}.TableName()
	orderSpecificationTable := order_specification.OrderSpecification{}.TableName()
	orderTable := orders.Order{}.TableName()
	projectTable := projects.Projects{}.TableName()
	orderSpecificationSteelItem  := order_specification_steel.OrderSpecificationSteel{}
	err := model.DB.Model(&orderSpecificationSteelItem).
		Select(fmt.Sprintf("%s.*", orderSpecificationSteelTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.order_specification_id", orderSpecificationTable, orderSpecificationTable, orderSpecificationSteelTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.order_id", orderTable, orderTable, orderSpecificationTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.project_id", projectTable, projectTable, orderTable)).
		Order(fmt.Sprintf("%s.location_code desc", orderSpecificationSteelTable)).
		First(&orderSpecificationSteelItem).
		Error
	if err != nil && err.Error() == "record not found" {
		return 1, nil
	}
	if err != nil {
		return 0, errors.ServerErr(ctx, err)
	}

	return orderSpecificationSteelItem.LocationCode + 1, nil
}
