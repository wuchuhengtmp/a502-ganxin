/**
 * @Desc    型钢归库验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/30
 * @Listen  MIT
 */
package requests

import (
	"context"
	graphModel "http-api/app/http/graph/model"
)

func ValidateSetProjectSteelEnterRepositoryRequest(ctx context.Context, input graphModel.SetProjectSteelEnterRepositoryInput) error {
	steps := StepsForProject{}
	// 检验有没有这个项目
	if err := steps.CheckHasProject(ctx, input.ProjectID); err != nil {
		return err
	}
	for _, identifier := range input.IdentifierList {
		// 检验有没有这根型钢
		if err := steps.CheckHasSteel(ctx, identifier); err != nil {
			return err
		}
		// 检验是否是归库状态
		if err := steps.CheckIsToBeEnterRepositoryState(ctx, identifier); err != nil {
			return err
		}
		// 检验型钢是否归于我管理的仓库的
		if err := steps.CheckIsSteelEnterMyRepository(ctx, identifier); err != nil {
			return err
		}
	}

	return nil
}
