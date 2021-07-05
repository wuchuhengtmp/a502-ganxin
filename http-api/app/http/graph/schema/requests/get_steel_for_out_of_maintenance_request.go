/**
 * @Desc    获取可出厂的型钢解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/5
 * @Listen  MIT
 */
package requests

import (
	"context"
	graphModel "http-api/app/http/graph/model"
)

func ValidateGetSteelForOutOfMaintenanceRequest(ctx context.Context, input graphModel.GetSteelForOutOfMaintenanceInput) error {
	steps := StepsForMaintenance{}
	// 检验型钢是否存在
	if err := steps.CheckHasSteel(ctx, input.Identifier); err != nil {
		return err
	}
	// 检验有没有这根型钢能不能出厂
	if err := steps.CheckIsOutOfMaintenanceAccess(ctx, input.Identifier); err != nil {
		return err
	}

	return nil
}
