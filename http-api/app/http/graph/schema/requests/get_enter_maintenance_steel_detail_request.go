/**
 * @Desc    获取待入场详细信息请求解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/4
 * @Listen  MIT
 */
package requests

import (
	"context"
	graphModel "http-api/app/http/graph/model"
)

func ValidateGetEnterMaintenanceSteelDetailRequest(ctx context.Context, input graphModel.GetEnterMaintenanceSteelDetailInput) error {
	steps := StepsForMaintenance{}
	for _, identifier := range input.IdentifierList {
		// 检验有没有这根型钢
		if err := steps.CheckHasSteel(ctx, identifier); err != nil {
			return err
		}
		//检验型钢是否归属我
		if err := steps.CheckIsSteelBelong2Me(ctx, identifier); err != nil {
			return err
		}
		// 检验型钢能否入厂
		if err := steps.CheckIsEnterMaintenanceAccess(ctx, identifier); err != nil {
			return err
		}
	}
	// 检验规格列表
	if input.SpecificationID != nil{
		if err := steps.CheckSpecification(ctx, *input.SpecificationID); err != nil {
			return err
		}
	}

	return nil
}
