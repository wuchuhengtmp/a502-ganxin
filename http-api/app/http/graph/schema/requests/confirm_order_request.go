/**
 * @Desc    确认订单验证器
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
	"github.com/thedevsaddam/govalidator"
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/model"
	"http-api/app/models/orders"
)

func ValidateConfirmOrderRequest(ctx context.Context, input model.ConfirmOrderInput) error {
	me := auth.GetUser(ctx)
	rules := govalidator.MapData{
		"id": []string{"required", "isCompanyOrder:" + fmt.Sprintf("%d", me.Id)},
	}
	opts := govalidator.Options{
		Data: &input,
		Rules: rules,
		TagIdentifier: "json",
	}
	if err := Validate(opts); err != nil {
		return err
	}
	o := orders.Order{Id: input.ID}
	_ = o.GetSelf()
	if o.State != orders.StateToBeConfirmed {
		return fmt.Errorf("当前订单状态为:%s 无法再次确认", orders.StateMapDesc[o.State] )
	}

	return nil
}
