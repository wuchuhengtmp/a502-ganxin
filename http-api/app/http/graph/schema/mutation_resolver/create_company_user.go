/**
 * @Desc    添加公司成员解析器
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

func (m *MutationResolver) CreateCompanyUser(ctx context.Context, input model.CreateCompanyUserInput) (*users.Users, error) {
	validate := requests.CreateCompanyUserRequest{}
	err := validate.ValidateCreateCompanyUserRequest(input)
	if err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	company := companies.Companies{}
	newUser, err := company.CreateUser(ctx, input)
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	return newUser, nil
}
