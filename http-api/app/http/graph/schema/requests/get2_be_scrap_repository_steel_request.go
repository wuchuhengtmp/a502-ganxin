/**
 * @Desc    型钢待报废查询请求验证器
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

func ValidateGet2BeScrapRepositorySteelRequest(ctx context.Context, input graphModel.Get2BeScrapRepositorySteelInput)  error {
	steps := StepsForRepository{}
	// 检验有没有这根型钢
	if err := steps.CheckHasSteel(ctx, input.Identifier); err != nil {
		return err
	}
	// 检验能不能报废
	if err := steps.CheckIsScrapAccess(ctx, input.Identifier); err != nil {
		return err
	}

	return nil
}
