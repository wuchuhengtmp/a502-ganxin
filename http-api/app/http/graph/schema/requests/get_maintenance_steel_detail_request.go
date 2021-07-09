/**
 * @Desc    获取维修厂详细信息请求验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/7
 * @Listen  MIT
 */
package requests

import (
	"context"
	graphModel "http-api/app/http/graph/model"
)

func ValidateGetMaintenanceSteelDetailRequest(ctx context.Context, input graphModel.GetMaintenanceSteelDetailInput) error {
	steps := StepsForMaintenance{}
	// 检验有没有这家维修厂
	if err := steps.CheckHashMaintenance(ctx, input.MaintenanceID); err != nil {
		return err
	}
	// 检验有没有这个规格
	if input.SpecificationID != nil && steps.CheckSpecification(ctx, *input.SpecificationID) != nil {
		return steps.CheckSpecification(ctx, *input.SpecificationID)
	}
	// 检验状态
	if input.State != nil && steps.CheckStateForDetail(*input.State) != nil {
		return steps.CheckStateForDetail(*input.State)
	}
	//for _, identifier := range input.IdentifierList{
	//	// 检验有没有这根型钢
	//	if err := steps.CheckHasSteel(ctx, identifier); err != nil {
	//		return err
	//	}
	//	// 检验型钢是否归属这个维修厂
	//	if err := steps.CheckSteelIsBelongMaintenance(ctx, input.MaintenanceID, identifier); err !=  nil {
	//		return err
	//	}
	//}

	return nil
}

