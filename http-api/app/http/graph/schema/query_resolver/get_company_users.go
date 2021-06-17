/**
 * @Desc    获取公司人员解析器
 * @Author  wuchuheng<root@wuchuheng.com>
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
	model2 "http-api/pkg/model"
)

func (q *QueryResolver) GetCompanyUser(ctx context.Context, input *model.GetCompanyUserInput) ([]*users.Users, error) {
	me := auth.GetUser(ctx)
	res, err := companies.GetCompanyItems(me.CompanyId, input)

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
// :xxx 这里的查询按理说是要放到模型中操作的，但循环依赖了
func (UserItemResolver)Company(ctx context.Context, obj *users.Users) (*companies.Companies, error) {
	c := companies.Companies{}
	err := model2.DB.Model(&companies.Companies{}).Where("id = ?", obj.CompanyId).First(&c).Error
	if err != nil {
		return nil, err
	}

	return &c, nil
}
