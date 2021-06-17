/**
 * @Desc    登录请求验证
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

type AuthorizationLoginRequest struct {
	Code string `valid:"code" json:"code"`
}

func (data *AuthorizationLoginRequest) Decode (body io.Reader)  {
	decoder := json.NewDecoder(body)
	_ = decoder.Decode(data)
}

func (data *AuthorizationLoginRequest) ValidateAuthorizationLoginRequest(body io.Reader) map[string][]string {
	data.Decode(body)
	message := govalidator.MapData{
		"code": []string {
			"required:code码不能为空",
		},
	}
	rules := govalidator.MapData{
		"code": []string{"required"},
	}
	opts := govalidator.Options{
		Data: data,
		Rules: rules,
		TagIdentifier: "valid",
		Messages: message,
	}
	return govalidator.New(opts).ValidateStruct()
}
