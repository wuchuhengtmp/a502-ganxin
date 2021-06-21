/**
 * @Desc    获取可待出库型钢列表验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/19
 * @Listen  MIT
 */
package requests

import (
	"context"
	"fmt"
	"http-api/app/http/graph/auth"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/models/order_specification"
	"http-api/app/models/order_specification_steel"
	"http-api/app/models/orders"
	"http-api/pkg/model"
)
// 订单各个规格对应的需求量
type OrderSpecificationList map[string]int64

// 有没有超出订单的要求数量
// 数据有没有重复
// 型钢要有仓库中， 且符合订单的规格中,且规格数量还没达到数量
func ValidateGetProject2WorkshopDetailRequest(ctx context.Context, input graphModel.ProjectOrder2WorkshopDetail) error {
	// 声明订单中各个规格需求量,用于检验输入数量和需求量是否超过上限
	osl := OrderSpecificationList{}
	// 检验有没有这个订单
	if err := osl.checkHasOrder(ctx, input); err != nil {
		return err
	}
	// 检验订单状态 只能是确认或部分发货才行
	if err := osl.checkOrderState(ctx, input); err != nil {
		return err
	}
	// 获取订单各规格需求上限
	//if err := osl.getOrderSpecificationGroupTotal(ctx, input); err != nil {
	//	return nil
	//}

	return nil
}

// 有没有这个订单
func (OrderSpecificationList) checkHasOrder(ctx context.Context, input graphModel.ProjectOrder2WorkshopDetail) error {
	me := auth.GetUser(ctx)
	o := orders.Order{}
	err := model.DB.Model(&o).Where("id = ?", input.OrderID).
		Where("company_id = ?", me.CompanyId).
		First(&o).
		Error
	if err != nil {
		return fmt.Errorf("没有这个订单id: %d", input.OrderID)
	}

	return nil
}

// 检验订单状态 只能是确认或部分发货才行
func (OrderSpecificationList) checkOrderState(ctx context.Context, input graphModel.ProjectOrder2WorkshopDetail) error {
	o := orders.Order{}
	if e := model.DB.Model(&o).Where("id = ?", o.Id).First(&o).Error; e != nil {
		return fmt.Errorf("没有这个订单id:%d", o.Id)
	}
	if o.State != orders.StateConfirmed && o.State != orders.StatePartOfReceipted {
		return fmt.Errorf("当前订单状态为:%s, 不能接着发货", orders.StateMapDesc[o.State])
	}

	return nil
}

// 获取订单各规格的需求量上限
func (OrderSpecificationList) int8(ctx context.Context, input graphModel.ProjectOrder2WorkshopDetail)  (OrderSpecificationList error){
	var osList  []*order_specification.OrderSpecification
	if err := model.DB.Model(&order_specification.OrderSpecification{}).Where("order_id = ?", input.OrderID).Find( &osList).Error; err != nil {
		return err
	}
	for _, item := range osList {
		oss := order_specification_steel.OrderSpecificationSteel{}
		var existsTotal int64
		if err := model.DB.Model(&oss).Where( "order_specification_id = ?", item.Id).Count(&existsTotal).Error; err != nil {
			return err
		}
			//osl[item.Specification] = item.Total - existsTotal
	}


	return nil
}
