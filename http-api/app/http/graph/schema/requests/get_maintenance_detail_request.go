/**
 * @Desc    获取维修详情验证
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/9
 * @Listen  MIT
 */
package requests

import (
	"context"
	graphModel "http-api/app/http/graph/model"
)

func ValidateGetMaintenanceDetailRequest(ctx context.Context, input graphModel.GetMaintenanceDetailInput) error {
	steps := StepsForMaintenance{}
	// 检验分页
	if err := steps.CheckPagination(input.IsShowAll, input.PageSize, input.Page); err != nil {
		return err
	}
	// 检验有没有这个仓库
	if input.RepositoryID != nil {
		if err := steps.CheckHashRepository(ctx, *input.RepositoryID); err != nil {
			return err
		}
	}
	if input.SpecificationID != nil {
		if err := steps.CheckSpecification(ctx, *input.SpecificationID); err != nil {
			return err
		}
	}

	return nil
}
