/**
 * @Desc    获取仓库列表的解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/4
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"fmt"
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/errors"
	"http-api/app/models/repositories"
	"http-api/app/models/repository_leader"
	"http-api/app/models/roles"
	"http-api/app/models/users"
	"http-api/pkg/model"
)
func (*QueryResolver)GetRepositoryList(ctx context.Context) ([]*repositories.Repositories, error) {
	var res []*repositories.Repositories
	repositoryModel := repositories.Repositories{}
	me := auth.GetUser(ctx)
	// 设备上的仓库管理员
	if auth.IsDevice(ctx) && me.RoleId == roles.RoleRepositoryAdminId  {
		repositoryLeaderTable := repository_leader.RepositoryLeader{}.TableName()
		repositoryTable := repositories.Repositories{}.TableName()
		var res []*repositories.Repositories
		err := model.DB.Model(&repositories.Repositories{}).
			Select(fmt.Sprintf("%s.*", repositoryTable)).
			Joins(fmt.Sprintf("join %s ON %s.repository_id = %s.id", repositoryLeaderTable, repositoryLeaderTable, repositoryTable)).
			Where(fmt.Sprintf("%s.uid = ?", repositoryLeaderTable), me.Id).
			Find(&res).
			Error
		if err != nil {
			return res, errors.ServerErr(ctx, err)
		}
		return res, nil
	} else {
		repositoryList, err := repositoryModel.GetAllRepositoryByCompanyId(me.CompanyId)
		if err != nil {
			return res, errors.ServerErr(ctx, err)
		}

		return repositoryList, nil
	}

}

type RepositoryItemResolver struct {}

func (RepositoryItemResolver) Leaders(ctx context.Context, obj *repositories.Repositories) ([]*users.Users, error) {
	return obj.GetLeaders()
}
