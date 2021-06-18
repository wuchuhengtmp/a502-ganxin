/**
 * @Desc    确认订单解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/18
 * @Listen  MIT
 */
package mutation_resolver

import (
	"context"
	"http-api/app/http/graph/errors"
	"http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/http/graph/schema/services"
	"http-api/app/models/orders"
)

func (*MutationResolver) ConfirmOrRejectOrder(ctx context.Context, input model.ConfirmOrderInput) (*orders.Order, error) {
	if err := requests.ValidateConfirmOrderRequest(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	o, err := services.ConfirmOrRejectOrder(ctx, input)
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	return o, nil
}
