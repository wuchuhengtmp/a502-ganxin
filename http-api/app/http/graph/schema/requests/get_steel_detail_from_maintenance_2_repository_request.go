/**
 * @Desc    The requests is part of http-api
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

func GetSteelDetailFromMaintenance2RepositoryRequest(ctx context.Context, input graphModel.GetSteelDetailFromMaintenance2RepositoryInput) error  {
	steps := StepsForRepository{}
	// 检验规格id
	if input.SpecificationID != nil {
		if err := steps.CheckHasSpecification(ctx, *input.SpecificationID); err != nil {
			return err
		}
	}
	for _, identifier := range input.IdentifierList  {
		// 检验有没有这根型钢
		if err := steps.CheckHasSteel(ctx, identifier); err != nil {
			return err
		}
		// 检验型钢能不能归库
		if err := steps.CheckIsEnterRepositoryFromMaintenanceAccess(ctx, identifier); err != nil {
			return err
		}
	}

	return nil
}
