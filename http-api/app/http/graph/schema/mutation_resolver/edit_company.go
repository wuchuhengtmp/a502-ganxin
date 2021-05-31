/**
 * @Desc    The mutation_resolver is part of http-api
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/31
 * @Listen  MIT
 */
package mutation_resolver

import (
	"context"
	"fmt"
	"http-api/app/http/graph/errors"
	"http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/http/graph/util/helper"
	"http-api/app/models/companies"
	"http-api/app/models/files"
)

/**
 * 更新公司
 */
func (*MutationResolver)EditCompany(ctx context.Context, input model.EditCompanyInput) (*model.CreateCompanyRes, error) {
	editCompanyRequest := requests.EditCompanyRequest{}
	err := editCompanyRequest.ValidateEditCompanyRequest(input)
	if err !=  nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	companyModel := companies.Companies{}
	ok := companyModel.Update(input)
	if !ok {
		err = fmt.Errorf("操作失败, 请联系管理员")
		return nil, errors.ServerErr(ctx, err)
	}
	loginFile := files.File{ }
	_ = loginFile.GetSelfById(int64(input.LogoFileID))
	backgroundFile := files.File{}
	_ = backgroundFile.GetSelfById(int64(input.BackgroundFileID))
	endedAt, _ :=  helper.Str2Time(input.EndedAt)
	startAt, _ :=  helper.Str2Time(input.StartedAt)
	_ = companyModel.GetSelfById(int64(input.ID))
	py := model.CreateCompanyRes{
		Name: input.Name,
		PinYin: input.PinYin,
		Symbol: input.Symbol,
		LogoFile: &model.SingleUploadRes{
			ID:  int(loginFile.ID),
			URL: loginFile.GetUrl(),
		},
		BackgroundFile: &model.SingleUploadRes{
			ID:  int(backgroundFile.ID),
			URL: backgroundFile.GetUrl(),
		},
		IsAble: input.IsAble,
		Phone: input.Phone,
		Wechat: input.Wechat,
		StartedAt: startAt,
		EndedAt: endedAt,
		AdminName: input.AdminName,
		CreatedAt: companyModel.CreatedAt,
	}

	return  &py, nil
}
