/**
 * @Desc    型钢入库验证器
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

func ValidateSetMaintenanceRequest(ctx context.Context, input graphModel.SetMaintenanceInput) error {
	steps := StepsForMaintenance{}
	for _, identifier := range input.IdentifierList {
		// 检验有没有这个根型钢
		if err := steps.CheckHasSteel(ctx, identifier); err != nil {
			return err
		}
		// 检验型钢能不能入库
		if err := steps.CheckIsEnterMaintenanceAccess(ctx, identifier); err != nil {
			return err
		}
		// 检验型钢是不是归属我
		if err := steps.CheckIsSteelBelong2Me(ctx, identifier); err != nil {
			return err
		}
	}

	return nil
}