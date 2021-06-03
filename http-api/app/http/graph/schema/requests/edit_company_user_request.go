/**
 * @Desc    编辑公司员工请求验证
 * @Author  wuchuheng<wuchuheng@163.com>
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
	"http-api/app/models/roles"
	"http-api/app/models/users"
)

type EditCompanyUseRequest struct{}

func (EditCompanyUseRequest) ValidateEditCompanyUserRequest(input *model.EditCompanyUserInput) error {
	rules := govalidator.MapData{
		"phone": []string{"phone", "not_user_phone_exists"},
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
	roleModel := roles.Role{}
	err := roleModel.GetSelfById(input.RoleID)
	if err != nil {
		return fmt.Errorf("没有这个角色id")
	}
	userModel := users.Users{}
	if userModel.GetSelfById(input.ID) != nil {
		return fmt.Errorf("没有这个用户")
	}

	return nil
}
