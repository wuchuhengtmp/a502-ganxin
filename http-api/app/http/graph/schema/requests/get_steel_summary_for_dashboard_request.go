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
	"http-api/app/models/repositories"
	"http-api/app/models/roles"
	"http-api/pkg/model"
)

func ValidateGetSteelSummaryForDashboardRequest(ctx context.Context, input graphModel.GetSteelSummaryForDashboardInput) error {
	if input.RepositoryID != nil {
		steps := StepsForRepository{}
		me := auth.GetUser(ctx)
		r, _ := me.GetRole()
		// 超管
		if r.ID == roles.RoleAdminId {
			repositoryItem := repositories.Repositories{}
			if err := model.DB.Model(&repositoryItem).Where("id = ?", *input.RepositoryID).First(&repositoryItem).Error; err != nil {
				return err
			}
		} else if err := steps.CheckHasRepository(ctx, *input.RepositoryID); err != nil {
			return err
		}
	}

	return nil
}
