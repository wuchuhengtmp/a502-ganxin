/**
 * @Desc    获取维修厂型钢请求验证器
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

func ValidateGetMaintenanceSteelRequest(ctx context.Context, input graphModel.GetMaintenanceSteelInput) error {
	steps := StepsForMaintenance{}
	// 检验有没有这家维修厂
	if err := steps.CheckHashMaintenance(ctx, input.MaintenanceID); err != nil {
		return err
	}

	return nil
}

