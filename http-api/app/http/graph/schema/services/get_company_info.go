/**
 * @Desc    The services is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/14
 * @Listen  MIT
 */
package services

import (
	"context"
	"http-api/app/http/graph/auth"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/models/companies"
	"http-api/app/models/configs"
	"http-api/app/models/files"
	"http-api/pkg/model"
)

func GetCompanyInfo(ctx context.Context) (*graphModel.GetCompnayInfoRes, error) {
	fid := configs.GetVal(configs.TUTOR_FILE_NAME, ctx)
	fileItem := files.File{}
	if err := model.DB.Model(&fileItem).Where("id = ?", fid).First(&fileItem).Error; err != nil {
		return nil, err
	}
	res := graphModel.GetCompnayInfoRes{}
	res.Tutor = &graphModel.FileItem{
		ID:  fileItem.ID,
		URL: fileItem.GetUrl(),
	}
	companyItem := companies.Companies{}
	me := auth.GetUser(ctx)
	if err := model.DB.Model(&companyItem).Where("id = ?", me.CompanyId).First(&companyItem).Error; err != nil {
		return nil, err
	}

	res.Phone = companyItem.Phone
	res.Wechat = companyItem.Wechat
	return &res, nil
}
