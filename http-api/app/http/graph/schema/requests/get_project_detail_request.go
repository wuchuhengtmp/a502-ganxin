/**
 * @Desc    The requests is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/8
 * @Listen  MIT
 */
package requests

import (
	"context"
	graphModel "http-api/app/http/graph/model"
)

func ValidateGetProjectDetailRequest(ctx context.Context, input graphModel.GetProjectDetailInput) error {
	steps := StepsForProject{}
	if err := steps.CheckPagination(input.IsShowAll, input.PageSize, input.Page); err != nil {
		return err
	}
	// 检验有没有这个仓库
	if input.RepositoryID != nil {
		if err := steps.CheckHasRepository(ctx, *input.RepositoryID); err != nil {
			return err
		}
	}
	// 检验有没有这个规格
	if input.SpecificationID != nil {
		if err := steps.CheckHasSpecification(ctx, *input.SpecificationID); err != nil {
			return err
		}
	}
	// 检验状态
	if input.State != nil {
		if err := steps.CheckSteelState(*input.State); err != nil {
			return err
		}
	}

	return nil
}
