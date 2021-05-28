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
	"http-api/app/http/graph/model"
	companiesModel "http-api/app/models/companies"
	"http-api/app/models/files"
)

type QueryResolver struct{}

/**
 * 获取全部公司解析器
 */
func (q *QueryResolver) GetAllCompany(ctx context.Context) ([]*model.CreateCompanyRes, error) {
	companies := companiesModel.GetAll()
	var res []*model.CreateCompanyRes
	for _, company := range companies {
		signEl := model.CreateCompanyRes{}
		signEl.ID = int(company.ID)
		signEl.Name = company.Name
		signEl.PinYin = company.PinYin
		signEl.Symbol = company.Symbol
		logFile := files.File{ }
		_ = logFile.GetFileById(company.LogoFileId)
		signEl.LogoFile = &model.SingleUploadRes{
			ID: int(logFile.ID),
			URL: logFile.GetUrl(),
		}
		backgroundFile := files.File{}
		_ = backgroundFile.GetFileById(company.BackgroundFileId)
		signEl.BackgroundFile = &model.SingleUploadRes{
			ID: int(backgroundFile.ID),
			URL: backgroundFile.GetUrl(),
		}
		signEl.IsAble = company.IsAble
		signEl.Phone = company.Phone
		signEl.Wechat = company.Wechat
		signEl.StartedAt = company.StartedAt
		signEl.EndedAt = company.EndedAt
		adminUser, _ := company.GetAdmin()
		signEl.AdminName = adminUser.Name
		signEl.CreatedAt = company.CreatedAt
		res = append(res, &signEl)
	}

	return res, nil
}
