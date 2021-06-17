/**
 * @Desc    编辑材料商请求验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/5
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

func ValidateEditMaterialManufacturerRequest(ctx context.Context, input model.EditMaterialManufacturerInput) error {
	rules := govalidator.MapData{
		"name": []string{"min:1"},
		"id":   []string{"isCodeInfoId"},
	}
	message := govalidator.MapData{
		"name": []string{
			"min:材料商名称不能为空",
		},
	}
	opt := govalidator.Options{
		Data:          &input,
		Rules:         rules,
		TagIdentifier: "json",
		Messages:      message,
	}
	res := govalidator.New(opt).ValidateStruct()
	if len(res) > 0 {
		for _, fieldErrors := range res {
			for _, err := range fieldErrors {
				return fmt.Errorf("%s", err)
			}
		}
	}
	me := auth.GetUser(ctx)
	c := codeinfo.CodeInfo{ID: input.ID}
	_ = c.GetSelf()
	if c.CompanyId != me.CompanyId {
		return fmt.Errorf("材料商家与您不是归属于同一家公司，您无权操作")
	}

	return nil
}
