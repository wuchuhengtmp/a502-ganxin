/**
 * @Desc     创建需求单解析器
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/16
 * @Listen  MIT
 */
package mutation_resolver

import (
	"context"
	"http-api/app/http/graph/model"
	"http-api/app/models/orders"
)

func (*MutationResolver)CreateOrder(ctx context.Context, input model.CreateOrderInput) (*orders.Order, error) {
	var o orders.Order

	return &o, nil
}
