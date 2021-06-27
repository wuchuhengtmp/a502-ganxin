/**
 * @Desc    获取待修改武钢信息解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/27
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"fmt"
	"http-api/app/http/graph/errors"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/order_specification_steel"
	"http-api/app/models/steels"
	"http-api/pkg/model"
)

func (*QueryResolver)GetProjectSteel2BeChange(ctx context.Context, input graphModel.GetProjectSteel2BeChangeInput) (*order_specification_steel.OrderSpecificationSteel, error) {
	if err := requests.GetProjectSteel2BeChangeRequest(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	var res order_specification_steel.OrderSpecificationSteel
	orderSpecificationSteelTable := order_specification_steel.OrderSpecificationSteel{}.TableName()
	steelTable := steels.Steels{}.TableName()
	err := model.DB.Model(&res).
		Select(fmt.Sprintf("%s.*",orderSpecificationSteelTable)).
		Joins(fmt.Sprintf("join %s ON %s.order_specification_steel_id = %s.id", steelTable, steelTable, orderSpecificationSteelTable)).
		Where(fmt.Sprintf("%s.identifier = ?", steelTable), input.Identifier).
		First(&res).
		Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	return &res, nil
}
