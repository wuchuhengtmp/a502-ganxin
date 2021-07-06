/**
 * @Desc    待出厂的型钢详情请求验证器
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

func ValidateGetSteelForOutOfMaintenanceDetailRequest(ctx context.Context, input graphModel.GetSteelForOutOfMaintenanceDetailInput) error{
	steps := StepsForMaintenance{}
	for _, identifier := range input.IdentifierList {
		// 检验是否有这根型钢
		if err := steps.CheckHasSteel(ctx, identifier); err != nil {
			return err
		}
		// 检验是否归属我
		if err := steps.CheckIsSteelBelong2Me(ctx, identifier); err != nil {
			return err
		}
		// 检验型钢能否出厂
		if err := steps.CheckIsOutOfMaintenanceAccess(ctx, identifier); err != nil {
			return err
		}
	}
	// 检验型钢尺寸
	if input.SpecificationID != nil {
		if err := steps.CheckSpecification(ctx, *input.SpecificationID); err != nil {
			return err
		}
	}

	return nil
}
