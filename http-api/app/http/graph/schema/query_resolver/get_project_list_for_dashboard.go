/**
 * @Desc    获取项目列表()
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/19
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/errors"
	"http-api/app/models/projects"
	"http-api/app/models/roles"
	"http-api/pkg/model"
)

func (*QueryResolver) GetProjectListForDashboard(ctx context.Context) ([]*projects.Projects, error) {
	me := auth.GetUser(ctx)
	role, _ := me.GetRole()
	modeIn := model.DB.Model(&projects.Projects{})
	if role.ID != roles.RoleAdminId {
		modeIn.Where("company_id = ?", me.CompanyId)
	}
	var res []*projects.Projects
	if err := modeIn.Find(&res).Error; err != nil {
		return res, errors.ServerErr(ctx, err)
	}

	return res, nil
}
