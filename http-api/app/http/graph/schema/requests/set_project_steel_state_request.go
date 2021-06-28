/**
 * @Desc    The requests is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/28
 * @Listen  MIT
 */
package requests

import (
	"context"
	graphModel "http-api/app/http/graph/model"
)

func ValidateSetProjectSteelStateRequest(ctx context.Context, input graphModel.SetProjectSteelInput) error {
	steps := StepsForProject{}
	// 检验状态是否合法
	if err := steps.CheckSteelState(input.State); err != nil {
		return err
	}
	for _, identifier := range input.IdentifierList {
		// 检验有没这个型钢
		if err := steps.CheckHasSteel(ctx, identifier); err != nil {
			return err
		}
		// 检验是否归我管
		if err := steps.CheckIsBelongMeByIdentifier(ctx, identifier); err != nil {
			return err
		}
		// 检验型钢是否是项目的型钢
		if err := steps.CheckIsProjectSteel(ctx, identifier); err != nil {
			return err
		}
	}

	return nil
}