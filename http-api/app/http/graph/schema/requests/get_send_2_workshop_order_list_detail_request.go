/**
 * @Desc    获取待入场的订单的详情接口验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/23
 * @Listen  MIT
 */
package requests

import (
	"context"
	"http-api/app/http/graph/model"
)

func ValidateGetSend2WorkshopOrderListDetail(ctx context.Context, input model.GetProjectOrder2WorkshopDetailInput) error {
	steps := ValidateGetProject2WorkshopDetailRequestSteps{}
	if err := steps.checkHasOrder(ctx, input.OrderID); err != nil {
		return err
	}
	if err := steps.CheckSpecification(ctx, input.OrderID, input.SpecificationID); err != nil {
		return err
	}

	return nil
}
