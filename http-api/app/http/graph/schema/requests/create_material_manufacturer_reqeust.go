/**
 * @Desc    添加材料商请求验证器
 * @Author  wuchuheng<wuchuheng@163.com>
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
	"http-api/app/http/graph/model"
)

func ValidateCreateMaterialManufacturerRequest(ctx context.Context, input model.CreateMaterialManufacturerInput) error {
	rules := govalidator.MapData{
		"name": []string{"min:1"},
	}
	message := govalidator.MapData{
		"name": []string{
			"min:材料商不能为空",
		},
	}
	opt := govalidator.Options{
		Data: &input,
		Rules: rules,
		Messages: message,
		TagIdentifier: "json",
	}

	hasErr := govalidator.New(opt).ValidateStruct()
	if len(hasErr) > 0 {
		for _, fieldErrors := range hasErr{
			for _, err := range fieldErrors {
				return fmt.Errorf("%s", err)
			}
		}
	}

	return nil
}