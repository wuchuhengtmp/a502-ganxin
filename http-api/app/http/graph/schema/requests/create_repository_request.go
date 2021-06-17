/**
 * @Desc    创建仓库请求验证
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
	"http-api/app/http/graph/model"
	"http-api/app/models/users"
)

type CreateRepositoryRequest struct { }

func  (CreateRepositoryRequest) ValidateCreateRepository(ctx context.Context, input model.CreateRepositoryInput) error {
	rules := govalidator.MapData{
		"repositoryAdminId": []string{"userExist"},
	}
	opts := govalidator.Options{
		Data: &input,
		Rules: rules,
		TagIdentifier: "json",
	}
	hasErr := govalidator.New(opts).ValidateStruct()
	if len(hasErr) > 0 {
		for _, fieldErrors := range hasErr{
			for _, err := range fieldErrors {
				return fmt.Errorf("%s", err)
			}
		}
	}
	user := users.Users{}
	_ = user.GetSelfById(input.RepositoryAdminID)
	me := auth.GetUser(ctx)
	if me.CompanyId != user.CompanyId {
		return fmt.Errorf("用户id:%d, 跟您不是同属于一家公司下的，您无权操作", user.ID)
	}

	return nil
}
