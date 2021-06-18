/**
 * @Desc    获取订单详情
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/18
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"http-api/app/http/graph/errors"
	"http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/orders"
)

func (*QueryResolver)GetOrderDetail(ctx context.Context, input model.GetOrderDetailInput) (*orders.Order, error) {
	if err := requests.ValidateGetOrderDetailRequest(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	o := orders.Order{Id: input.ID}
	_ = o.GetSelf()

	return &o, nil
}
