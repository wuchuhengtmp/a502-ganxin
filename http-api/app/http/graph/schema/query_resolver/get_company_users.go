/**
 * @Desc    获取公司人员解析器
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/3
 * @Listen  MIT
 */

package query_resolver

import (
	"context"
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/model"
	"http-api/app/models/companies"
	"http-api/app/models/files"
	"http-api/app/models/roles"
	"http-api/app/models/users"
)

func (q *QueryResolver) GetCompanyUser(ctx context.Context) ([]*users.Users, error) {
	me := auth.GetUser(ctx)
	res, err := companies.GetCompanyItemsResById(me.CompanyId)

	return res, err
}

type UserItemResolver struct{}

// role 字段解析
func (UserItemResolver) Role(ctx context.Context, obj *users.Users) (*roles.Role, error) {
	r := roles.Role{}
	_ = r.GetSelfById(obj.RoleId)

	return &roles.Role{
		ID:   r.ID,
		Tag:  r.Tag,
		Name: r.Name,
	}, nil
}

// 用户的avatar 字段解析
func (UserItemResolver) Avatar(ctx context.Context, obj *users.Users) (*model.FileItem, error) {
	f := files.File{}
	_ = f.GetSelfById(obj.AvatarFileId)
	res := model.FileItem{
		ID:  f.ID,
		URL: f.GetUrl(),
	}

	return &res, nil
}
