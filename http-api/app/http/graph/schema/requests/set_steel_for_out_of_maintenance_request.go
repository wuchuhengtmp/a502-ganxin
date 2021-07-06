/**
 * @Desc    型钢出厂请求验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/6
 * @Listen  MIT
 */
package requests

import (
	"context"
	graphModel "http-api/app/http/graph/model"
)

func ValidateSetSteelForOutOfMaintenanceRequest(ctx context.Context, input graphModel.SetSteelForOutOfMaintenanceInput) error {
	steps := StepsForMaintenance{}
	for _, identifier := range input.IdentifierList {
		// 检验有没有这根型钢
		if err := steps.CheckHasSteel(ctx, identifier); err != nil {
			return err
		}
		// 检验型钢是否归属我
		if err := steps.CheckIsSteelBelong2Me(ctx, identifier); err != nil {
			return err
		}
		// 检验能不能出厂
		if err := steps.CheckIsOutOfMaintenanceAccess(ctx, identifier); err != nil {
			return err
		}
	}

	return nil
}