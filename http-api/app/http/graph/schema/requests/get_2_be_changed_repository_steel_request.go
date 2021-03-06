/**
 * @Desc    获取用于修改的仓库型钢请求验证器
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

func Get2BeChangedRepositorySteelRequest(ctx context.Context, input graphModel.Get2BeChangedRepositorySteelInput) error {
	steps := StepsForRepository{}
	// 检验有没有这根型钢
	if err := steps.CheckHasSteel(ctx, input.Identifier); err != nil {
		return err
	}
	// 检验是不是归属我
	if err := steps.CheckIsSteelBeLongMe(ctx, input.Identifier); err != nil {
		return err
	}
	return nil
}