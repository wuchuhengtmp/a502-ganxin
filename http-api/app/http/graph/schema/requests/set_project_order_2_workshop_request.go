/**
 * @Desc    型钢出库验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/22
 * @Listen  MIT
 */
package requests

import (
	"context"
	"fmt"
	"http-api/app/http/graph/auth"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/models/codeinfo"
	"http-api/pkg/model"
)

func ValidateSetProjectOrder2WorkshopRequest(ctx context.Context, input graphModel.ProjectOrder2WorkshopInput) error {
	osl := ValidateGetProject2WorkshopDetailRequestSteps{}
	// 识别码列表不能为空
	if err := osl.IdentificationListMustBeEmpty(ctx,input.IdentifierList); err != nil {
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
	// 检验物流参数
	steps := ValidateSetProjectOrder2WorkshopRequestSteps{}
	if err := steps.CheckExpress(ctx, input); err != nil {
		return err
	}

	return nil
}


//  验证型钢出库的验证步骤合集
type ValidateSetProjectOrder2WorkshopRequestSteps struct{}

/**
 * 检验物流
 */
func (ValidateSetProjectOrder2WorkshopRequestSteps) CheckExpress(ctx context.Context, input graphModel.ProjectOrder2WorkshopInput) error {
	me := auth.GetUser(ctx)
	express := codeinfo.CodeInfo{}
	err := model.DB.Model(&express).Where("company_id = ?", me.CompanyId).
		Where("type = ?", codeinfo.ExpressCompany).
		Where("id = ?", input.ExpressCompanyID).
		First(&express).
		Error
	if err != nil {
		return fmt.Errorf("没有物流公司id为:%d的物流公司", input.ExpressCompanyID)
	}
	if len(input.ExpressNo) == 0 {
		return fmt.Errorf("物流编号不能为空")
	}

	return nil
}
