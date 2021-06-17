/**
 * @Desc    The mutation_resolver is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/3
 * @Listen  MIT
 */
package mutation_resolver

import (
	"context"
	"http-api/app/http/graph/errors"
	"http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/companies"
	"http-api/app/models/users"
)

func (m *MutationResolver)EditCompanyUser(ctx context.Context, input *model.EditCompanyUserInput) (*users.Users, error) {
	requestValidation := requests.EditCompanyUseRequest{}
	err := requestValidation.ValidateEditCompanyUserRequest(ctx, input)
	if err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	res, err := companies.UpdateCompanyUser(ctx, input)
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	return res, nil
}
