/**
 * @Desc    获取订单型钢详情验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/29
 * @Listen  MIT
 */
package requests

import (
	"context"
	graphModel "http-api/app/http/graph/model"
)

func ValidateGetOrderSteelDetailRequest(ctx context.Context, input graphModel.GetOrderSteelDertailInput) error {
	steps := StepsForProject{}
	// 检验有没有这根型钢
	if err := steps.CheckHasSteel(ctx, input.Identifier); err != nil {
		return err
	}
	// 检验型钢是否归属我管
	if err := steps.CheckIsBelongMeByIdentifier(ctx, input.Identifier); err != nil {
		return err
	}
	return nil
}