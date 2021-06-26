/**
 * @Desc    获取项目最大的安装码解析器
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
func ValidateGetMaxLocationCodeRequest(ctx context.Context, input graphModel.GetMaxLocationCodeInput) error{
	steps := StepsForProject{}
	// 检验项目管理员是不是我
	if err := steps.CheckHasProjectByIdentifier(input.Identifier); err != nil {
		return err
	}
	// 检验项目管理员是不是我
	if err := steps.CheckIsBelongMeByIdentifier(ctx, input.Identifier); err != nil  {
		return err
	}

	return nil
}

