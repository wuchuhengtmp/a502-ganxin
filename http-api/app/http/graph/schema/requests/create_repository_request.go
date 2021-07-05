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
	"http-api/app/http/graph/model"
	"http-api/app/models/roles"
	"http-api/app/models/users"
)

type CreateRepositoryRequest struct { }

func  (CreateRepositoryRequest) ValidateCreateRepository(ctx context.Context, input model.CreateRepositoryInput) error {
	steps := StepsForRepository{}

	user := users.Users{}
	for _, adminUid := range input.RepositoryAdminID {
		if err := steps.CheckHasUser(ctx, adminUid); err != nil {
			return err
		}
		_ = user.GetSelfById(adminUid)
		if  user.RoleId != roles.RoleRepositoryAdminId {
			role, _ := user.GetRole()
			return fmt.Errorf("用户id为：%d, 的用户是%s 不是管理管理员角色", adminUid, role.Name)
		}
	}

	return nil
}
