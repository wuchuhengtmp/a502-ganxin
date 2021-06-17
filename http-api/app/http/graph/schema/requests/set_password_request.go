/**
 * @Desc    重置密码验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/15
 * @Listen  MIT
 */
package requests

import (
	"context"
	"fmt"
	"github.com/thedevsaddam/govalidator"
	"http-api/app/http/graph/model"
)

func ValidateSetPasswordRequest(ctx context.Context, input *model.SetPasswordInput) error {
	rules := govalidator.MapData{
		"password":     []string{"required", "min:8"},
	}
	messages := govalidator.MapData{
		"password": []string{"min:密码不能少于8位"},
	}
	opt := govalidator.Options{
		Data:            input,
		Rules:           rules,
		TagIdentifier:   "json",
		Messages: messages,
	}
	res := govalidator.New(opt).ValidateStruct()
	if len(res) > 0 {
		for _, fieldErrors := range res {
			for _, err := range fieldErrors {
				return fmt.Errorf("%s", err)
			}
		}
	}
	if len(res) > 0 {
		for _, fieldErrors := range res {
			for _, err := range fieldErrors {

				return fmt.Errorf("%s", err)
			}
		}
	}

	return nil
}