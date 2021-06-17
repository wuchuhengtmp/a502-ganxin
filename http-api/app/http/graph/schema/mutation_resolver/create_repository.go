/**
 * @Desc  	创建新仓库解析器
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/4
 * @Listen  MIT
 */
package mutation_resolver

import (
	"context"
	"http-api/app/http/graph/errors"
	"http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/http/graph/schema/services"
	"http-api/app/models/repositories"
)

func (*MutationResolver) CreateRepository(ctx context.Context, input model.CreateRepositoryInput) (*repositories.Repositories, error) {
	validator := requests.CreateRepositoryRequest{}
	if err := validator.ValidateCreateRepository(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	r, err := services.CreateRepository(ctx, input)
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	return r, nil
}
