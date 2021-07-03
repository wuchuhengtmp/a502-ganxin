/**
 * @Desc    获取待维修型钢详情请求验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/3
 * @Listen  MIT
 */
package requests

import (
	"context"
	graphModel "http-api/app/http/graph/model"
)

func ValidateGet2BeMaintainSteelDetailRequest(ctx context.Context, input graphModel.Get2BeMaintainSteelDetailInput) error {
	steps := StepsForRepository{}
	for _, identifier := range input.IdentifierList {
		// 检验有没有这个型钢
		if err := steps.CheckHasSteel(ctx, identifier); err != nil {
			return err
		}
		// 检验型钢是否我管的
		if err := steps.CheckIsSteelBeLongMe(ctx, identifier); err != nil {
			return err
		}
		//检验能否维修
		if err := steps.CheckIs2BeMaintainAccess(ctx, identifier); err != nil {
			return err
		}
	}
	if input.SpecificationID != nil && steps.CheckHasSpecification(ctx, *input.SpecificationID) != nil {
		return steps.CheckHasSpecification(ctx, *input.SpecificationID)
	}

	return nil
}
