/**
 * @Desc    创建维修厂验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/1
 * @Listen  MIT
 */
package requests

import (
	graphModel "http-api/app/http/graph/model"
)

func ValidateCreateMaintenanceRequest(input graphModel.CreateMaintenanceInput) error {
	steps := StepsForMaintenance{}
	for _, uid := range input.UID {
		// 检验有没有这个用户
		if err := steps.CheckHasUser(uid); err != nil {
			return err
		}
		//检验用户角色
		if err := steps.CheckIsMaintenanceRole(uid); err != nil {
			return err
		}
	}

	return nil
}
