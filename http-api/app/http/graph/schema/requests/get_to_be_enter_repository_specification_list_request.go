/**
 * @Desc    获取待归库的尺寸列表验证器
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

func ValidateGetToBeEnterRepositorySpecificationList(ctx context.Context, input graphModel.GetToBeEnterRepositorySpecificationListInput) error {
	steps := StepsForProject{}
	// 检验有没有这个项目
	if err := steps.CheckHasProject(ctx, input.ProjectID); err != nil {
		return err
	}
	// 检验项目是否包含仓库的型钢
	if err := steps.CheckIsIncludeMyRepository(ctx, input.ProjectID); err != nil {
		return err
	}

	return nil
}
