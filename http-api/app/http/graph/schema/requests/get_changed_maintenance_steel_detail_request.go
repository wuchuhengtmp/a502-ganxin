/**
 * @Desc    待修改型钢详细信息请求验证器
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

func ValidateGetChangedMaintenanceSteelDetailRequest(ctx context.Context, input graphModel.GetChangedMaintenanceSteelDetailInput) error  {
	steps  := StepsForMaintenance{}
	for _, identifier := range input.IdentifierList {
		// 检验有没有这根型钢
		if err := steps.CheckHasSteel(ctx, identifier); err != nil {
			return err
		}
		// 检验型钢是否归属我
		if err := steps.CheckIsSteelBelong2Me(ctx, identifier); err != nil {
			return err
		}
		// 检验型钢能否修改状态
		if err := steps.CheckIsChangedMaintenanceSteelAccess(ctx, identifier); err != nil {
			return err
		}

	}

	return nil
}
