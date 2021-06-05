/**
 * @Desc    添加码表解析器
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
	"http-api/app/models/specificationinfo"
)

func (*MutationResolver) CreateSpecification(ctx context.Context, input model.CreateSpecificationInput) (*specificationinfo.SpecificationInfo, error) {
	if err := requests.ValidateCreateSpecificationRequest(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	s := specificationinfo.SpecificationInfo{
		Type:      input.Type,
		Length:    input.Length,
		Weight:    input.Weight,
		IsDefault: input.IsDefault,
	}
	if err := s.CreateSelf(ctx); err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	return &s, nil
}
