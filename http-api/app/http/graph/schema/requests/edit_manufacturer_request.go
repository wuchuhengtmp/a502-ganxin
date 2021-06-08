/**
 * @Desc    编辑制造商请求验证器
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
	auth2 "http-api/app/http/graph/auth"
	"http-api/app/http/graph/model"
	"http-api/app/models/codeinfo"
)

func ValidateEditManufacturerRequest(ctx context.Context, input model.EditManufacturerInput) error {
	rules := govalidator.MapData{
		"type": []string{"min:1"},
		"id":   []string{"isCodeInfoId"},
	}
	message := govalidator.MapData{
		"type": []string{
			"min:类型不能为空",
		},
	}
	opts := govalidator.Options{
		Data: &input,
		Rules: rules,
		TagIdentifier: "json",
		Messages: message,
	}
	res := govalidator.New(opts).ValidateStruct()
	if len(res) > 0 {
		for _, fieldErrors := range res {
			for _, err := range fieldErrors {
				return fmt.Errorf("%s", err)
			}
		}
	}
	me := auth2.GetUser(ctx)
	c := codeinfo.CodeInfo{
		ID: input.ID,
	}
	_ = c.GetSelf()
	if c.CompanyId != me.CompanyId {
		return fmt.Errorf("当前修改的制造商与您不是归属同一家公司，您无权修改")
	}

	return nil
}
