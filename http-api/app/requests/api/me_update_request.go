/**
 * @Desc    添加用户属性
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @DATE    2021/4/28
 * @Listen  MIT
 */
package api

import (
	"encoding/json"
	"github.com/thedevsaddam/govalidator"
	"io"
)

type MeUpdateRequest struct {
	AvatarUrl string `valid:"avatarUrl" json:"avatarUrl"`
	Gender string `valid:"gender" json:"gender"`
	Nickname string `valid:"nickname" json:"nickname"`
}

func (data *MeUpdateRequest) Decode (body io.Reader)  {
	decoder := json.NewDecoder(body)
	_ = decoder.Decode(data)
}

func (data *MeUpdateRequest) ValidateAuthorizationLoginRequest(body io.Reader) map[string][]string {
	data.Decode(body)
	message := govalidator.MapData{
		"avatarUrl": []string {
			"required:头像不能为空",
		},
		"gender": []string {
			"required:性别不能为空",
		},
		"nickname": []string {
			"required:性别不能为空",
		},
	}
	rules := govalidator.MapData{
		"avatarUrl": []string{"required"},
		"gender":    []string{"required"},
		"nickname":  []string{"required"},
	}
	opts := govalidator.Options{
		Data: data,
		Rules: rules,
		TagIdentifier: "valid",
		Messages: message,
	}
	return govalidator.New(opts).ValidateStruct()
}
