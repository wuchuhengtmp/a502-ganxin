/**
 * @Desc    编辑公司请求验证
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/31
 * @Listen  MIT
 */
package requests

import (
	"fmt"
	"github.com/thedevsaddam/govalidator"
	"http-api/app/http/graph/model"
	"http-api/app/http/graph/util/helper"
	"http-api/app/models/companies"
	"http-api/app/models/users"
)

type EditCompanyRequest struct { }

func (data *EditCompanyRequest) ValidateEditCompanyRequest(input model.EditCompanyInput) error {
	rules := govalidator.MapData{
		"id":                []string{"required"},
		"name":              []string{"min:1"},
		"pinYin":            []string{"min:1"},
		"symbol":            []string{"min:1"},
		"phone":             []string{"phone"},
		"adminPhone":        []string{"phone"},
		"endedAt":           []string{"time"},
		"startedAt":         []string{"time"},
		"adminAvatarFileID": []string{"fileExist"},
		"isAble":            []string{"in:true,false"},
		"wechat":            []string{"min:2"},
		"logoFileID":        []string{"fileExist"},
		"backgroundFileID":  []string{"fileExist"},
		"adminName":         []string{"min:1"},
		"adminPassword":     []string{"min:6"},
		"adminWechat":       []string{"min:6"},
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
	startedAt, _ := helper.Str2Time(input.StartedAt)
	endedAt, _ := helper.Str2Time(input.EndedAt)
	if startedAt.Unix() >= endedAt.Unix() {
		return fmt.Errorf("开始时间不能大于或等于结束时间")
	}
	// 公司是否存在
	companyModel := companies.Companies{}
	_, err := companyModel.HasCompanyId(int64(input.ID))
	if err != nil {
		return fmt.Errorf("没有这家公司")
	}
	// 密码验证
	if input.AdminPassword != nil && len(*input.AdminPassword) < 6 {
		return fmt.Errorf("密码不能小于6位")
	}
	user := users.Users{}
	if user.IsChangeCompanyAdminPhone(int64(input.ID), input.Phone) && !user.IsUniPhone(input.Phone) {
		return fmt.Errorf("手机号已被占用")
	}

	return nil
}
