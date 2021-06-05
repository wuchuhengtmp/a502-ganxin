/**
 * @Desc    编辑码表解析器
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/5
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

func (*MutationResolver)EditSpecification(ctx context.Context, input model.EditSpecificationInput) (*specificationinfo.SpecificationInfo, error) {
	if err := requests.ValidateEditSpecificationRequest(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	var s specificationinfo.SpecificationInfo
	if err := s.Edit(ctx, input); err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	return &s, nil
}
