/**
 * @Desc    创建物流公司验证器
 * @Author  wuchuheng<root@wuchuheng.com>
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
	"http-api/app/http/graph/model"
)

func ValidateCreateExpressRequest(ctx context.Context, input model.CreateExpressInput) error {
	rules := govalidator.MapData{
		"name": []string{"min:1"},
	}
	message := govalidator.MapData{
		"name": []string{"min:物流名不能为空"},
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

	return nil
}