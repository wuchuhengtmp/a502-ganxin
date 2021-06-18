/**
 * @Desc    获取订单详情请求验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/18
 * @Listen  MIT
 */
package requests

import (
	"context"
	"fmt"
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/model"
	"http-api/app/models/orders"
)

func ValidateGetOrderDetailRequest(ctx context.Context, input model.GetOrderDetailInput) error {
	o := orders.Order{Id: input.ID}
	if err := o.GetSelf(); err != nil {
		return fmt.Errorf("没有这个订单")
	}
	me := auth.GetUser(ctx)
	if o.CompanyId != me.CompanyId {
		return fmt.Errorf("没有这个订单")
	}

	return nil
}