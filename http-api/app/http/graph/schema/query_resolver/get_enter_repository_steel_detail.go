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
	"http-api/app/models/order_specification_steel"
	"http-api/app/models/steels"
	"http-api/pkg/model"
)

func (*QueryResolver)GetEnterRepositorySteelDetail(ctx context.Context, input graphModel.GetEnterRepositorySteelDetailInput) (*order_specification_steel.OrderSpecificationSteel, error) {
	if err := requests.ValidateGetEnterRepositorySteelDetailRequest(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	res := order_specification_steel.OrderSpecificationSteel{}
	steelTable := steels.Steels{}.TableName()
	model.DB.Model(&res).
		Select(fmt.Sprintf("%s.*", res.TableName())).
		Joins(fmt.Sprintf("join %s ON %s.order_specification_steel_id = %s.id", steelTable, steelTable, res.TableName())).
		Where(fmt.Sprintf("%s.identifier = ?", steelTable), input.Identifier).
		Where(fmt.Sprintf("%s.state = ?", res.TableName()), steels.StateProjectOnTheStoreWay).
		First(&res)

	return &res, nil
}
