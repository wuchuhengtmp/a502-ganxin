/**
 * @Desc    获取型钢入厂的型钢信息请求验证器
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

func ValidateGetEnterMaintenanceRequest(ctx context.Context, input graphModel.EnterMaintenanceInput) error {
	steps := StepsForMaintenance{}
	// 检验有没有这根型钢
	if err := steps.CheckHasSteel(ctx, input.Identifier); err != nil {
		return err
	}
	// 检验这根型钢是否属于我
	if err := steps.CheckIsSteelBelong2Me(ctx, input.Identifier); err != nil {
		return err
	}
	// 检验能否入厂
	if err := steps.CheckIsEnterMaintenanceAccess(ctx, input.Identifier); err != nil {
		return err
	}

	return nil
}
