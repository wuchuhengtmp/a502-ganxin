/**
 * @Desc    创建公司解析器
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/28
 * @Listen  MIT
 */
package mutation_resolver

import (
	"context"
	"http-api/app/http/graph/errors"
	"http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/http/graph/util/helper"
	"http-api/app/models/companies"
	"http-api/app/models/files"
	"http-api/app/models/roles"
	"http-api/app/models/users"
	globalHelper "http-api/pkg/helper"
)

func (m *MutationResolver) CreateCompany(ctx context.Context, input model.CreateCompanyInput) (*model.CreateCompanyRes, error) {
	CreateCompanyRequest := requests.CreateCompanyRequest{}
	err := CreateCompanyRequest.ValidateCreateCompanyRequest(input)
	if err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	startedAt, _ := helper.Str2Time(input.StartedAt)
	endedAt, _ := helper.Str2Time(input.EndedAt)
	company := companies.Companies{
		Name:             input.Name,
		PinYin:           input.PinYin,
		Symbol:           input.Symbol,
		LogoFileId:       int64(input.LogoFileID),
		BackgroundFileId: int64(input.BackgroundFileID),
		IsAble:           input.IsAble,
		Phone:            input.Phone,
		Wechat:           input.Wechat,
		StartedAt:        startedAt,
		EndedAt:          endedAt,
	}
	err = company.Create()
	user := users.Users{
		Name:         input.AdminName,
		Password:     globalHelper.GetHashByStr(input.AdminPassword),
		Phone:        input.AdminPhone,
		RoteId:       roles.RoleCompanyAdminId,
		Wechat:       input.AdminWechat,
		CompanyId:    company.ID,
		IsAble:       true,
		AvatarFileId: int64(input.AdminAvatarFileID),
	}
	err = user.Create()
	logoFile := files.File{}
	logoFile.GetFileById(company.LogoFileId)
	backgroundFile := files.File{}
	backgroundFile.GetFileById(company.BackgroundFileId)
	bf := model.SingleUploadRes{
		ID: int(backgroundFile.ID),
		URL: backgroundFile.GetUrl(),
	}

	// 响应数据
	res := model.CreateCompanyRes{
		ID:     int(company.ID),
		Name:   company.Name,
		PinYin: company.PinYin,
		Symbol: company.Symbol,
		LogoFile: &model.SingleUploadRes{
			ID: int( logoFile.ID),
			URL: logoFile.GetUrl(),
		},
		BackgroundFile: &bf,
		IsAble:    company.IsAble,
		Phone:     company.Phone,
		Wechat:    company.Wechat,
		StartedAt: company.StartedAt,
		EndedAt:   company.EndedAt,
		AdminName: user.Name,
		CreatedAt: company.CreatedAt,
	}

	return &res, nil
}
