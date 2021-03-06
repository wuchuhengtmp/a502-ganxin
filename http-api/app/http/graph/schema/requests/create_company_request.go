/**
 * @Desc    创建公司请求验证
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/28
 * @Listen  MIT
 */
package requests

import (
	"fmt"
	"github.com/thedevsaddam/govalidator"
	"http-api/app/http/graph/model"
	"http-api/app/http/graph/util/helper"
)

type CreateCompanyRequest struct {
	Phone             string `valid:"Phone"`
	LogoFileId        int64    `valid:"LogoFileId"`
	BackgroundFileId  int64    `valid:"BackgroundFileId"`
	EndedAt           string `valid:"EndedAt"`
	StartedAt         string `valid:"StartedAt"`
	AdminPhone        string `valid:"AdminPhone"`
	AdminPassword     string `valid:"AdminPassword"`
	AdminAvatarFileId int64    `valid:"AdminAvatarFileId"`
}

func (data *CreateCompanyRequest) ValidateCreateCompanyRequest(input model.CreateCompanyInput) error {
	data.Phone = input.Phone
	data.LogoFileId = input.LogoFileID
	data.BackgroundFileId = input.BackgroundFileID
	data.EndedAt = input.EndedAt
	data.StartedAt = input.StartedAt
	data.AdminPhone = input.AdminPhone
	data.AdminPassword = input.AdminPassword
	data.AdminAvatarFileId = input.AdminAvatarFileID
	rules := govalidator.MapData{
		"Phone":             []string{"phone"},
		"AdminPhone":        []string{"phone", "not_user_phone_exists"},
		"AdminPassword":     []string{"min:6"},
		"EndedAt":           []string{"time"},
		"StartedAt":         []string{"time"},
		"AdminAvatarFileId": []string{"fileExist"},
		"LogoFileId":        []string{"fileExist"},
		"BackgroundFileId":  []string{"fileExist"},
	}
	opts := govalidator.Options{
		Data:          data,
		Rules:         rules,
		TagIdentifier: "valid",
	}
	res := govalidator.New(opts).ValidateStruct()
	if len(res) > 0 {
		for _, fieldErrors := range res {
			for _, err := range fieldErrors {
				return fmt.Errorf("%s", err)
			}
		}
	}
	startedAt, _ := helper.Str2Time(data.StartedAt)
	endedAt, _ := helper.Str2Time(data.EndedAt)
	if startedAt.Unix() >= endedAt.Unix() {
		return fmt.Errorf("开始时间不能大于或等于结束时间")
	}

	return nil
}
