/**
 * @Desc    The requests is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/1
 * @Listen  MIT
 */
package requests

import (
	"context"
	graphModel "http-api/app/http/graph/model"
)

func ValidateEditMaintenanceRequest(ctx context.Context, input graphModel.EditMaintenanceInput) error {
	steps := StepsForMaintenance{}
	// 检验有没这个厂
	if err := steps.CheckHashMaintenance(ctx, input.ID); err != nil {
		return err
	}
	for _, uid := range input.AdminIDList {
		// 检验有没有这个管理员
		if err := steps.CheckHasUser(uid); err != nil {
			return err
		}
		// 检验是不是维修厂管理员
		if err := steps.CheckIsMaintenanceRole(uid); err != nil {
			return err
		}

	}
	// 检验用户id有没有冗余
	if err := steps.CheckRedundancyUid(input.AdminIDList); err != nil {
		return err
	}

	return nil
}
