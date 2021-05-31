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
func (q *QueryResolver) GetAllCompany(ctx context.Context) ([]*model.CompanyItemRes, error) {
	companies := companiesModel.GetAll()
	var res []*model.CompanyItemRes
	for _, company := range companies {
		signEl := model.CompanyItemRes{}
		signEl.ID = int(company.ID)
		signEl.Name = company.Name
		signEl.PinYin = company.PinYin
		signEl.Symbol = company.Symbol
		logFile := files.File{ }
		_ = logFile.GetSelfById(company.LogoFileId)
		signEl.LogoFile = &model.SingleUploadRes{
			ID: int(logFile.ID),
			URL: logFile.GetUrl(),
		}
		backgroundFile := files.File{}
		_ = backgroundFile.GetSelfById(company.BackgroundFileId)
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
		signEl.CreatedAt = company.CreatedAt
		signEl.AdminName = adminUser.Name
		signEl.AdminPhone = adminUser.Phone
		signEl.Wechat = adminUser.Wechat
		adminAvatar := files.File{}
		adminAvatar.GetSelfById(adminUser.ID)
		signEl.AdminAvatar = &model.SingleUploadRes{
			ID: int( adminAvatar.ID),
			URL: adminAvatar.GetUrl(),
		}
		res = append(res, &signEl)
	}

	return res, nil
}
