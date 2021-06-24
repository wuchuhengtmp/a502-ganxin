/**
 * @Desc    获取可待出库的订单型钢详情验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/19
 * @Listen  MIT
 */
package requests

import (
	"context"
	graphModel "http-api/app/http/graph/model"
)


func ValidateGetProject2WorkshopDetailRequest(ctx context.Context, input graphModel.ProjectOrder2WorkshopDetailInput) error {
	// 声明订单中各个规格需求量,用于检验输入数量和需求量是否超过上限
	osl := ValidateGetProject2WorkshopDetailRequestSteps{}
	// 识别码列表不能为空
	if err := osl.CheckIdentificationListMustBeEmpty(ctx, input.IdentifierList); err != nil {
		return err
	}
	// 检验有没有这个订单
	if err := osl.checkHasOrder(ctx, input.OrderID); err != nil {
		return err
	}
	// 检验订单状态 只能是确认或部分发货才行
	if err := osl.checkOrderState(ctx, input.OrderID); err != nil {
		return err
	}
	// 验证是否冗余识别码
	if err := osl.CheckRedundancyIdentification(input.IdentifierList); err != nil {
		return err
	}
	// 验证每根型钢的状态和规格是否满足订单要求，且数量也没超过上限
	if err := osl.CheckSteelList(ctx, input.OrderID, input.IdentifierList); err != nil {
		return err
	}
	// 检验规格
	if err := osl.CheckSpecification(ctx, input.OrderID, input.SpecificationID); err != nil {
		return err
	}

	return nil
}

