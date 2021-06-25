/**
 * @Desc    获取项目详情请求验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/25
 * @Listen  MIT
 */
package requests

import (
	"context"
	graphModel "http-api/app/http/graph/model"
)

func ValidateGetProjectSteelDetailRequest(ctx context.Context, input graphModel.GetProjectSteelDetailInput) error {
	steps := StepsForProject{}
	// 检验有没有这个项目
	if err := steps.CheckHasProject(ctx, input.ProjectID); err != nil {
		return err
	}
	// 检验项目管理员是不是我
	if err := steps.CheckIsBelongMe(ctx, input.ProjectID); err != nil {
		return err
	}
	// 检验型钢状态
	if input.State != nil {
		if err := steps.CheckSteelState(ctx, *input.State); err != nil {
			return err
		}
	}
	// 检验规格
	if input.SpecificationID != nil {
		if err := steps.CheckSpecification(ctx, *input.SpecificationID, input.ProjectID); err != nil {
			return err
		}
	}

	return nil
}
