/**
 * @Desc    创建公司角色请求验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/3
 * @Listen  MIT
 */
package requests

import (
	"fmt"
	"github.com/thedevsaddam/govalidator"
	"http-api/app/http/graph/model"
)

type CreateCompanyUserRequest struct{}

func (CreateCompanyUserRequest) ValidateCreateCompanyUserRequest(input model.CreateCompanyUserInput) error {
	rules := govalidator.MapData{
		"phone":    []string{"phone", "not_user_phone_exists"},
		"avatarId": []string{"fileExist"},
		"password": []string{"min:6"},
	}
	opts := govalidator.Options{
		Data:          &input,
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

	return nil
}
