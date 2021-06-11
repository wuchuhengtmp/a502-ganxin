/**
 * @Desc    The query_resolver is part of http-api
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/28
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"http-api/app/http/graph/errors"
	"http-api/app/http/graph/model"
	companiesModel "http-api/app/models/companies"
	"http-api/app/models/files"
)

type QueryResolver struct{}

/**
 * 获取全部公司解析器
 */
func (q *QueryResolver) GetAllCompany(ctx context.Context) ([]*companiesModel.Companies, error) {
	c := companiesModel.Companies{}
	cs, err := c.GetAll()
	if  err != nil {
		return cs, errors.ServerErr(ctx, err)
	}

	return cs, nil
}

type CompanyItemResolver struct { }

func (CompanyItemResolver) LogoFile(ctx context.Context, obj *companiesModel.Companies) (*model.FileItem, error) {
	f := files.File{}
	if err := f.GetSelfById(obj.LogoFileId); err != nil {
		return nil, err
	}
	mf := model.FileItem{
		ID: f.ID,
		URL: f.GetUrl(),
	}

	return &mf, nil
}
func (CompanyItemResolver)BackgroundFile(ctx context.Context, obj *companiesModel.Companies) (*model.FileItem, error) {
	f := files.File{}
	if err := f.GetSelfById(obj.BackgroundFileId); err != nil {
		return nil, err
	}
	mf := model.FileItem{
		ID: f.ID,
		URL: f.GetUrl(),
	}

	return &mf, nil
}

func (CompanyItemResolver)AdminName(ctx context.Context, obj *companiesModel.Companies) (string, error) {
	u, err := obj.GetAdmin()
	if  err != nil {
		return "", err
	}

	return u.Name, nil
}

func (CompanyItemResolver)AdminPhone(ctx context.Context, obj *companiesModel.Companies) (string, error) {
	u, err := obj.GetAdmin()
	if  err != nil {
		return "", err
	}

	return u.Phone, nil
}
func (CompanyItemResolver)AdminWechat(ctx context.Context, obj *companiesModel.Companies) (string, error) {
	u, err := obj.GetAdmin()
	if  err != nil {
		return "", err
	}

	return u.Wechat, nil
}
func (CompanyItemResolver)AdminAvatar(ctx context.Context, obj *companiesModel.Companies) (*model.FileItem, error) {
	u, err := obj.GetAdmin()
	if err != nil {
		return nil, err
	}
	f := files.File{}
	if err := f.GetSelfById(u.AvatarFileId); err != nil {
		return nil, err
	}
	mf := model.FileItem{
		ID: f.ID,
		URL: f.GetUrl(),
	}

	return &mf, nil
}
