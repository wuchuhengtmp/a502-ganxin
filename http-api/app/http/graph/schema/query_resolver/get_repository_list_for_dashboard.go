/**
 * @Desc    获取仓库列表
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/12
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/errors"
	"http-api/app/models/repositories"
	"http-api/app/models/roles"
	"http-api/pkg/model"
)

func (*QueryResolver)GetRepositoryListForDashboard(ctx context.Context) (res []*repositories.Repositories, err error)  {
	me := auth.GetUser(ctx)
	role, _ := me.GetRole()
	repositoryItem := repositories.Repositories{}
	modelIn := model.DB.Model(&repositoryItem)
	// 不是超管
	if role.ID != roles.RoleAdminId {
		modelIn = modelIn.Where("company_id = ?", me.CompanyId)
	}
	if err = modelIn.Find(&res).Error; err != nil {
		return res, errors.ServerErr(ctx, err)
	}

	return
}

