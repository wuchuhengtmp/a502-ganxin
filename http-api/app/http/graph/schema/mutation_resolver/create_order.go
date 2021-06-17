/**
 * @Desc    创建需求单解析器
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/16
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

func (*MutationResolver)CreateOrder(ctx context.Context, input model.CreateOrderInput) (*orders.Order, error) {
	if err := requests.ValidateCreateOrderValidate(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	o, err := services.CreateOrder(ctx, input)
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	return o, nil
}
