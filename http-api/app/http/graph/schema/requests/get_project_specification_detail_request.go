/**
 * @Desc    获取项目规格列表验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/25
 * @Listen  MIT
 */
package requests

import (
	"context"
	"http-api/app/http/graph/model"
)

func  ValidateGetProjectSpecificationDetailRequest(ctx context.Context, input model.GetProjectSpecificationDetailInput) error {
	steps := ValidateGetProjectSpecificationDetailRequestSteps{}
	// 检验项目是否存在
	if err := steps.CheckProjectExists(ctx, input.ProjectID); err != nil {
		return err
	}
	// 检验当前用户能不能查看这个项目
	if err := steps.CheckProjectLeader(ctx, input.ProjectID); err != nil {
		return err
	}

	return nil
}