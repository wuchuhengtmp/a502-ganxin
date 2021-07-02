/**
 * @Desc    The requests is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/2
 * @Listen  MIT
 */
package requests

import (
	"context"
	graphModel "http-api/app/http/graph/model"
)

func ValidateGetRepositorySteelRequest(ctx context.Context, input graphModel.GetRepositorySteelInput) error {
	steps := StepsForRepository{}
	// 检验有没有这个仓库
	if err := steps.CheckHasRepository(ctx, input.ReposirotyID); err != nil {
		return err
	}
	// 检验仓库是否归属我
	if err := steps.CheckRepositoryBelongMe(ctx, input.ReposirotyID); err != nil {
		return err
	}
	// 检验有没有这个状态
	if input.State != nil && steps.CheckHasState(*input.State) != nil {
		return steps.CheckHasState(*input.State)
	}
	// 检验有没有这个尺寸
	if input.SpecificationID != nil && steps.CheckHasSpecification(ctx, *input.SpecificationID) != nil {
		return steps.CheckHasSpecification(ctx, *input.SpecificationID)
	}

	return nil
}
