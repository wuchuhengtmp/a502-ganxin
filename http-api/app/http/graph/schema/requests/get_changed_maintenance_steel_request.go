/**
 * @Desc    查询用于维修型钢的信息请求验证器
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

func ValidateGetChangedMaintenanceSteelRequest(ctx context.Context, input graphModel.GetChangedMaintenanceSteelInput) error {
	steps := StepsForMaintenance{}
	// 检验有没有这根型钢
	if err := steps.CheckHasSteel(ctx, input.Identifier); err != nil {
		return err
	}
	// 检验能否修改
	if err := steps.CheckIsChangedMaintenanceSteelAccess(ctx, input.Identifier); err != nil {
		return err
	}

	return nil
}