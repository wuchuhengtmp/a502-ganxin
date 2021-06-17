/**
 * @Desc    删除公司员工请求验证
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/4
 * @Listen  MIT
 */
package requests

import (
	"context"
	"fmt"
	"github.com/thedevsaddam/govalidator"
	"http-api/app/http/graph/auth"
	"http-api/app/models/users"
)

type DeleteCompanyUserRequest struct { }

func (DeleteCompanyUserRequest) ValidateDeleteCompanyUserRequest(ctx context.Context, uid int64) error {
	rules := govalidator.MapData{
		"uid": []string{"userExist"},
	}
	opts := govalidator.Options{
		Data: &struct { Uid int64 `json:"uid"` }{Uid: uid},
		Rules: rules,
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

	me := auth.GetUser(ctx)
	user := users.Users{}
	_ = user.GetSelfById(uid)
	if me.CompanyId != user.CompanyId {
		return fmt.Errorf("用户id%d,与您不是归属于同一家公司，您无权删除", uid)
	}

	return nil
}