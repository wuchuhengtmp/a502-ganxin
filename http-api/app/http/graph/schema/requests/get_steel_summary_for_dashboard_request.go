/**
 * @Desc    型钢概览-饼图请求验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/13
 * @Listen  MIT
 */
package requests

import (
	"context"
	"http-api/app/http/graph/auth"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/models/roles"
)

func ValidateGetSteelSummaryForDashboardRequest(ctx context.Context, input graphModel.GetSteelSummaryForDashboardInput) error {
	if input.RepositoryID != nil {
		me := auth.GetUser(ctx)
		role, _ := me.GetRole()
		if role.ID == roles.RoleAdminId {
			return nil
		}
	}

	return nil
}
