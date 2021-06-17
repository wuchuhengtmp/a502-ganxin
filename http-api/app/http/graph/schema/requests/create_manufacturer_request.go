/**
 * @Desc    创建制造商请求验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/7
 * @Listen  MIT
 */
package requests

import (
	"context"
	"fmt"
	"github.com/thedevsaddam/govalidator"
	"http-api/app/http/graph/model"
)

func ValidateCreateManufacturerRequest(ctx context.Context, input model.CreateManufacturerInput) error {
	rules := govalidator.MapData{
		"name": []string{"min:1"},
	}
	message := govalidator.MapData{
		"name": []string{
			"min:制造商不能为空",
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