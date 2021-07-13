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
	"http-api/app/http/graph/errors"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/models/configs"
	"http-api/app/models/files"
	"http-api/pkg/model"
)

func (*QueryResolver)GetCompanyInfo(ctx context.Context) (*graphModel.GetCompnayInfoRes, error) {
	fid := configs.GetVal(configs.TUTOR_FILE_NAME, ctx)
	fileItem := files.File{}
	if err := model.DB.Model(&fileItem).Where("id = ?", fid).First(&fileItem).Error; err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	res := graphModel.GetCompnayInfoRes{}
	res.Tutor = &graphModel.FileItem{
		ID: fileItem.ID,
		URL: fileItem.GetUrl(),
	}
	res.Phone = configs.GetVal(configs.PHONE_NAME, ctx)
	res.Wechat = configs.GetVal(configs.WECHAT_NAME, ctx)

	return &res, nil
}
