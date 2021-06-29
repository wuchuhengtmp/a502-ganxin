/**
 * @Desc    项目归库的型钢查询验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/29
 * @Listen  MIT
 */
package requests

import (
	"context"
	graphModel "http-api/app/http/graph/model"
)

func ValidateGetEnterRepositorySteelDetailRequest(ctx context.Context, input graphModel.GetEnterRepostiroySteelDettailInput)  error {
	steps := StepsForProject{}
	// 检验有没有这根型钢
	if err := steps.CheckHasSteel(ctx, input.Identifier); err != nil {
		return err
	}
	// 检验有没有这个项目
	if err := steps.CheckHasProject(ctx, input.ProjectID); err != nil {
		return err
	}
	// 检验是否属于这个项目
	if err := steps.CheckIsBelongProject(ctx,input.ProjectID, input.Identifier); err != nil {
		return err
	}
	// 检验是否是入库状态
	if err := steps.CheckSteelIsEnterRepositoryState(ctx, input.Identifier); err != nil {
		return err
	}

	return nil
}