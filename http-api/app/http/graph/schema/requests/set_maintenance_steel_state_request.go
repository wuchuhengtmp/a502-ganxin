/**
 * @Desc    修改维修型钢状态验证器
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

func ValidateSetMaintenanceSteelStateRequest(ctx context.Context, input graphModel.SetMaintenanceSteelStateInput) error {
	steps := StepsForMaintenance{}
	for _, identifier := range input.IdentifierList {
		// 检验有没有这个根型钢
		if err := steps.CheckHasSteel(ctx, identifier); err != nil {
			return err
		}
		// 检验是否归属我
		if err := steps.CheckIsSteelBelong2Me(ctx, identifier); err != nil {
			return err
		}
		// 检验能不能修改状态
		if err := steps.CheckIsChangedMaintenanceSteelAccess(ctx, identifier); err != nil {
			return err
		}
	}
	// 检验状态是否合法
	if err := steps.CheckStateForChanged(input.State); err != nil {
		return err
	}

	return nil
}
