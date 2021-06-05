/**
 * @Desc    删除规格解析器
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
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/specificationinfo"
)

func (*MutationResolver) DeleteSpecification(ctx context.Context, id int64) (bool, error) {
	if err := requests.ValidateDeleteSpecificationRequest(ctx, id); err != nil {
		return false, errors.ValidateErr(ctx, err)
	}
	s := specificationinfo.SpecificationInfo{ID: id}
	if err := s.DeleteSelf(ctx); err != nil {
		return false, errors.ServerErr(ctx, err)
	}

	return true, nil
}