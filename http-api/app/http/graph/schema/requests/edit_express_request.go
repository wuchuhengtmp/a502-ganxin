/**
 * @Desc    编辑物流公司请求验证器
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
	"http-api/app/http/graph/model"
	"http-api/app/models/codeinfo"
)

func ValidateEditExpressRequest(ctx context.Context, input model.EditExpressInput) error {
	rules := govalidator.MapData{
		"id":   []string{"isCodeInfoId"},
		"name": []string{"min:1"},
	}
	message := govalidator.MapData{
		"name": []string{"min:物流名不能为空"},
	}
	opts := govalidator.Options{
		Data:          &input,
		Messages:      message,
		TagIdentifier: "json",
		Rules:         rules,
	}
	res := govalidator.New(opts).ValidateStruct()
	if len(res) > 0 {
		for _, fieldErrors := range res {
			for _, err := range fieldErrors {

				return fmt.Errorf("%s", err)
			}
		}
	}
	c := codeinfo.CodeInfo{ID: input.ID}
	_ = c.GetSelf()
	me := auth.GetUser(ctx)
	if me.CompanyId != c.CompanyId {
		return fmt.Errorf("被编辑的物流公司与您是不是归属于同一家公司名下， 您无权操作")
	}
	if c.Type != codeinfo.ExpressCompany {
		return fmt.Errorf("这不是物流公司的id，无法修改")
	}

	return nil
}
