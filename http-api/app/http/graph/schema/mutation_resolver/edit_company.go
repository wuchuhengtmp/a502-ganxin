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
	"http-api/app/http/graph/errors"
	"http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/companies"
)

/**
 * 更新公司
 */
func (*MutationResolver)EditCompany(ctx context.Context, input model.EditCompanyInput) (*companies.Companies, error) {
	editCompanyRequest := requests.EditCompanyRequest{}
	err := editCompanyRequest.ValidateEditCompanyRequest(input)
	if err !=  nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	company := companies.Companies{}
	if err := company.Update(ctx, input); err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	return  &company, nil
}
