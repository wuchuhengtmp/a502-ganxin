/**
 * @Desc    获取仓库列表的解析器
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/4
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/errors"
	"http-api/app/models/repositories"
)
func (*QueryResolver)GetRepositoryList(ctx context.Context) ([]*repositories.Repositories, error) {
	var res []*repositories.Repositories
	repositoryModel := repositories.Repositories{}
	me := auth.GetUser(ctx)
	repositoryList, err := repositoryModel.GetAllRepositoryByCompanyId(me.CompanyId)
	if err != nil {
		return res, errors.ServerErr(ctx, err)
	}

	return repositoryList, nil
}

type RepositoryItemResolver struct {}

// adminName 字段解析器
func (RepositoryItemResolver)AdminName(ctx context.Context, obj *repositories.Repositories) (string, error) {
	 user, err := obj.GetAdminUser()
	 if err != nil {
		return "", errors.ServerErr(ctx, err)
	 }

	 return user.Name, nil
}
// adminPhone字段解析器
func (RepositoryItemResolver)AdminPhone(ctx context.Context, obj *repositories.Repositories) (string, error) {
	user, err := obj.GetAdminUser()
	if err != nil {
		return "", errors.ServerErr(ctx, err)
	}

	return user.Phone, nil
}
// adminWechat字段解析器
func (RepositoryItemResolver)AdminWechat(ctx context.Context, obj *repositories.Repositories) (string, error) {
	user, err := obj.GetAdminUser()
	if err != nil {
		return "", errors.ServerErr(ctx, err)
	}

	return user.Wechat, nil
}
