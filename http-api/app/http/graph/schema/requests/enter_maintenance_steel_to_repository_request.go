/**
 * @Desc    型钢入库请求验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/12
 * @Listen  MIT
 */
package requests

import (
	"context"
	graphModel "http-api/app/http/graph/model"
)

func ValidateEnterMaintenanceSteelToRepositoryRequest(ctx context.Context, input graphModel.EnterMaintenanceSteelToRepositoryInput) error {
	steps := StepsForRepository{}
	for _, identifier := range input.IdentifierList {
		// 检验有没有这根型钢
		if err := steps.CheckHasSteel(ctx, identifier); err != nil {
			return err
		}
		// 检验能不能归库
		if err := steps.CheckIsEnterRepositoryFromMaintenanceAccess(ctx, identifier); err != nil {
			return err
		}
	}

	return nil
}