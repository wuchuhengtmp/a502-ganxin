/**
 * @Desc    型钢归库查询请求验证器
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

func ValidateGetSteelFromMaintenance2RepositoryRequest(ctx context.Context, input graphModel.GetSteelFromMaintenance2RepositoryInput) error  {
	steps := StepsForRepository{}
	// 检验有没有这根型钢
	if err := steps.CheckHasSteel(ctx, input.Identifer); err != nil {
		return err
	}
	// 这根型钢能不能归库
	if err := steps.CheckIsEnterRepositoryFromMaintenanceAccess(ctx, input.Identifer); err != nil {
		return err
	}

 return nil
}
