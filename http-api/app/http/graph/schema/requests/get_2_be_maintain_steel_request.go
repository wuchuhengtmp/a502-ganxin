/**
 * @Desc    获取维修出库型钢请求验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/3
 * @Listen  MIT
 */
package requests

import (
	"context"
	graphModel "http-api/app/http/graph/model"
)

func ValidateGet2BeMaintainSteelRequest(ctx context.Context, input graphModel.Get2BeMaintainSteelInput) error {
	steps := StepsForRepository{}
	// 检验有没有这个型钢
	if err := steps.CheckHasSteel(ctx, input.Identifier); err != nil {
		return err
	}
	// 检验型钢是否归属我
	if err := steps.CheckIsSteelBeLongMe(ctx, input.Identifier); err != nil {
		return err
	}
	// 检验型钢是否能出库维修
	if err := steps.CheckIs2BeMaintainAccess(ctx, input.Identifier); err != nil {
		return err
	}

	return nil
}
