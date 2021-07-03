/**
 * @Desc    型钢维修出库验证器
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

func ValidateSetBatchOfMaintenanceSteelRequest(ctx context.Context, input graphModel.SetBatchOfMaintenanceSteelInput) error {
	steps := StepsForRepository{}
	for _, identifier := range input.IdentifierList {
		// 检验有没有这根型钢
		if err := steps.CheckHasSteel(ctx, identifier); err != nil {
			return err
		}
		// 检验是否归属我
		if err := steps.CheckIsSteelBeLongMe(ctx, identifier); err != nil {
			return err
		}
		// 检验能不能维修
		if err := steps.CheckIs2BeMaintainAccess(ctx, identifier); err != nil {
			return err
		}
	}
	// 检验有没有这家维修厂
	if err := steps.CheckHasMaintenance(ctx, input.MaintenanceID); err != nil {
		return err
	}

	return nil
}
