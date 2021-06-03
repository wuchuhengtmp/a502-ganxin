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
	"fmt"
	"http-api/app/http/graph/errors"
	"http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/companies"
)

func (m *MutationResolver) CreateCompany(ctx context.Context, input model.CreateCompanyInput) (*model.CompanyItemRes, error) {
	CreateCompanyRequest := requests.CreateCompanyRequest{}
	err := CreateCompanyRequest.ValidateCreateCompanyRequest(input)
	if err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	companyModel := companies.Companies{}
	//添加公司
	newCompany, err := companyModel.Create(ctx, input)
	if err != nil {
		return nil, errors.ServerErr(ctx, fmt.Errorf(errors.ServerErrorMsg))
	}
	// 获取解析器要的数据
	res, _ := companies.GetCompanyItemResById(newCompany.ID)

	return res, nil
}
