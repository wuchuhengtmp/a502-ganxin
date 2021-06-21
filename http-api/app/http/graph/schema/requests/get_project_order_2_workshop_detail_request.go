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
	"fmt"
	"http-api/app/http/graph/auth"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/models/order_specification"
	"http-api/app/models/order_specification_steel"
	"http-api/app/models/orders"
	"http-api/app/models/steels"
	"http-api/pkg/model"
)

// 订单各个规格对应的需求量
type OrderSpecificationList map[string]int64

func ValidateGetProject2WorkshopDetailRequest(ctx context.Context, input graphModel.ProjectOrder2WorkshopDetailInput) error {
	// 声明订单中各个规格需求量,用于检验输入数量和需求量是否超过上限
	osl := OrderSpecificationList{}
	// 识别码列表不能为空
	if err := osl.IdentificationListMustBeEmpty(ctx, input); err != nil {
		return err
	}
	// 检验有没有这个订单
	if err := osl.checkHasOrder(ctx, input); err != nil {
		return err
	}
	// 检验订单状态 只能是确认或部分发货才行
	if err := osl.checkOrderState(ctx, input); err != nil {
		return err
	}
	// 验证是否冗余识别码
	if err := osl.isRedundancyIdentification(input.IdentifierList); err != nil {
		return err
	}
	// 验证每根型钢的状态和规格是否满足订单要求，且数量也没超过上限
	if err := osl.CheckSteelList(ctx, input); err != nil {
		return err
	}
	// 检验规格
	if err := osl.CheckSpecification(ctx, input); err != nil {
		return err
	}

	return nil
}

// 有没有这个订单
func (OrderSpecificationList) checkHasOrder(ctx context.Context, input graphModel.ProjectOrder2WorkshopDetailInput) error {
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
func (OrderSpecificationList) checkOrderState(ctx context.Context, input graphModel.ProjectOrder2WorkshopDetailInput) error {
	o := orders.Order{}
	if e := model.DB.Model(&o).Where("id = ?", input.OrderID).First(&o).Error; e != nil {
		return fmt.Errorf("没有这个订单id:%d", input.OrderID)
	}
	if o.State != orders.StateConfirmed && o.State != orders.StatePartOfReceipted {
		return fmt.Errorf("当前订单状态为:%s, 不能接着发货", orders.StateMapDesc[o.State])
	}

	return nil
}

// 获取订单各规格的需求量上限
func (OrderSpecificationList) GetOrderSpecificationGroupTotal(ctx context.Context, input graphModel.ProjectOrder2WorkshopDetailInput) (map[string]int64, error) {
	list := make(map[string]int64)
	var osList []*order_specification.OrderSpecification
	if err := model.DB.Model(&order_specification.OrderSpecification{}).Where("order_id = ?", input.OrderID).Find(&osList).Error; err != nil {
		return list, nil
	}
	for _, item := range osList {
		oss := order_specification_steel.OrderSpecificationSteel{}
		var existsTotal int64
		if err := model.DB.Model(&oss).Where("order_specification_id = ?", item.Id).Count(&existsTotal).Error; err != nil {
			return list, err
		}
		list[item.Specification] = item.Total - existsTotal
	}

	return list, nil
}

/**
 * 检验是否有冗余识别码
 */
func (OrderSpecificationList) isRedundancyIdentification(list []string) error {
	identificationMapTotal := make(map[string]int64)
	for _, item := range list {
		if _, ok := identificationMapTotal[item]; ok {
			return fmt.Errorf("识别码出现冗余，%s 不能输入多个同样的", item)
		} else {
			identificationMapTotal[item] = 1
		}
	}

	return nil
}

/*
 * 识别码不能为空
 */
func (OrderSpecificationList) IdentificationListMustBeEmpty(ctx context.Context, input graphModel.ProjectOrder2WorkshopDetailInput) error {
	if len(input.IdentifierList) == 0 {
		return fmt.Errorf("识别码列表不能为空")
	}

	return nil
}

func (OrderSpecificationList) CheckSteelList(ctx context.Context, input graphModel.ProjectOrder2WorkshopDetailInput) error {
	me := auth.GetUser(ctx)
	// 订单规格合集
	var orderSpecificationList []*order_specification.OrderSpecification
	orderSpecificationSpecificationMapTotal := make(map[string]int64) // 当前同一规格统计量 用于比较上限
	var orderSpecificationIdList []int64 // 订单要求的规格id集合，用于检验型钢的规格是否在这个合集中
	err := model.DB.Model(&order_specification.OrderSpecification{}).Where("order_id = ?", input.OrderID).
		Find(&orderSpecificationList).
		Error
	if err != nil {
		return err
	}
	for _, item := range orderSpecificationList {
		orderSpecificationIdList = append(orderSpecificationIdList, item.SpecificationId)
	}
	// 获取订单各规格需求上限
	osl, err := OrderSpecificationList{}.GetOrderSpecificationGroupTotal(ctx, input)
	if err != nil {
		return nil
	}
	// 检验每根型钢
	for _, identification := range input.IdentifierList {
		s := steels.Steels{}
		// 检验型钢状态能否满足订单要求
		err := model.DB.Model(&steels.Steels{}).
			Where("identifier = ?", identification).
			Where("state = ?", steels.StateInStore).
			Where("company_id = ?", me.CompanyId).
			First(&s).
			Error
		if err != nil {
			return fmt.Errorf("仓库中没有 %s 标识码的型钢在仓库中", identification)
		}
		// 检验型钢的规格能否满足订单的要求
		if err := func() error{
			for _, specificationId := range orderSpecificationIdList {
				if specificationId == s.SpecificationId  {
					return nil
				}
			}
			return fmt.Errorf("订单中,要求的规格id为:%v, 而标识码的为%s的型钢的规格id为%d, 并不能满足订单的要求", orderSpecificationIdList, identification, s.SpecificationId)
		}(); err != nil {
			return err
		}
		specificationInstance, err := s.GetSpecification()
		if err != nil {
			return fmt.Errorf("型钢规格不存在 id:%d ，请联系管理员", identification)
		}
		// 检验当前规格型钢的数量是否超过订单要求的上限
		key := specificationInstance.GetSelfSpecification()
		orderSpecificationSpecificationMapTotal[key] += 1
		// 上限比较
		if orderSpecificationSpecificationMapTotal[key] > osl[key] {
			return fmt.Errorf("当前规格%s， 已经超过订单要求的%d 数量了", key, osl[key])
		}
	}

	return nil
}
/**
 * 检验规格
 */
func (OrderSpecificationList) CheckSpecification(ctx context.Context, input graphModel.ProjectOrder2WorkshopDetailInput) error {
	if input.SpecificationID != nil {
		err := model.DB.
			Model(&order_specification.OrderSpecification{}).
			Where("order_id = ?", input.OrderID).
			Where("specification_id = ?", input.SpecificationID).
			First(&order_specification.OrderSpecification{}).
			Error
		if err != nil {
			return fmt.Errorf("订单中没有id为: %d 的规格", *input.SpecificationID)
		}
	}

	return nil
}
