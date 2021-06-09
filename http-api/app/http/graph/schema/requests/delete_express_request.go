/**
 * @Desc    删除物流公司验证器
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/8
 * @Listen  MIT
 */
package requests

import (
	"context"
	"fmt"
	"github.com/thedevsaddam/govalidator"
	"http-api/app/http/graph/auth"
	"http-api/app/models/codeinfo"
	"http-api/app/models/orders"
)

func ValidateDeleteExpressRequest(ctx context.Context, id int64) error {
	rules := govalidator.MapData{
		"id": []string{"isCodeInfoId"},
	}
	opts := govalidator.Options{
		Data: &struct {
			Id int64 `json:"id"`
		}{id},
		Rules:         rules,
		TagIdentifier: "json",
	}
	res := govalidator.New(opts).ValidateStruct()
	if len(res) > 0 {
		for _, fieldErrors := range res {
			for _, err := range fieldErrors {

				return fmt.Errorf("%s", err)
			}
		}
	}
	me := auth.GetUser(ctx)
	c := codeinfo.CodeInfo{ID: id}
	_ = c.GetSelf()
	if me.CompanyId != c.CompanyId {
		return fmt.Errorf("要删除的物流与您不是归属于同一家公司，您无权删除")
	}
	if c.Type != codeinfo.ExpressCompany {
		return fmt.Errorf("这不是物流公司的数据，无法删除")
	}

	o := orders.Order{}
	os, _ := o.GetOrdersByExpressId(id)
	if len(os) > 0 {
		return fmt.Errorf("该物流信息已经在订单中使用了无法删除，无法删除")
	}

	return nil
}
