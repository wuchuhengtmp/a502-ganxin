/**
 * @Desc 安装码是否可用请求验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/26
 * @Listen  MIT
 */
package requests

import (
	"context"
	graphModel "http-api/app/http/graph/model"
)

func ValidateIsAccessLocationCodeRequest(ctx context.Context, input graphModel.IsAccessLocationCodeInput) error {
	steps := StepsForProject{}
	// 检验有没有这个项目
	if err := steps.CheckHasProject(ctx, input.ProjectID); err != nil {
		return err
	}
	// 检验这是不是我管的项目
	if err := steps.CheckIsBelongMe(ctx, input.ProjectID); err != nil {
		return err
	}

	return nil
}