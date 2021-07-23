/**
 * @Desc    The requests is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/23
 * @Listen  MIT
 */
package requests

import (
	"context"
	graphModel "http-api/app/http/graph/model"
)

func ValidateEditOrderRequest(ctx context.Context, input graphModel.EditOrderInput) error  {
	steps := StepsForOrder{}
	// 有没有这个订单
	if err := steps.CheckHasOrder(ctx, input.ID); err != nil {
		return err
	}
	// 检验开始时间
	if err := steps.CheckExpectedAt(ctx, input.ExpectedReturnAt); err != nil {
		return err
	}
	// 型钢规格需求验证
	if err := steps.CheckSteelSpecification(ctx, input.SteelList ); err != nil {
		return err
	}

	return nil
}
