/**
 * @Desc    获取型钢单根型钢出场详情解析验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/28
 * @Listen  MIT
 */
package requests

import (
	"context"
	"fmt"
	"http-api/app/http/graph/auth"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/models/order_specification_steel"
	"http-api/app/models/steels"
	"http-api/pkg/model"
)

func ValidateGetOutOfWorkshopProjectSteelDetailRequest(ctx context.Context, input graphModel.GetOutOfWorkshopProjectSteelDetail) error {
	steps := StepsForProject{}
	// 检验有没有这根型钢
	if err := steps.CheckHasSteel(ctx, input.Identifier); err != nil {
		return err
	}
	// 检验项目中有没有这根型钢
	if err := steps.CheckHasProjectByIdentifier(input.Identifier); err != nil {
		return err
	}
	// 检验这根型钢归不归我管理
	if err := steps.CheckIsBelongMeByIdentifier(ctx, input.Identifier); err != nil {
		return err
	}
	// 检验有没有这个项目
	if err := steps.CheckHasProject(ctx, input.ProjectID); err != nil {
		return err
	}
	// 检验型钢状态
	s := steels.Steels{}
	me := auth.GetUser(ctx)
	orderSpecificationSteelItem := order_specification_steel.OrderSpecificationSteel{}
	err := model.DB.Model(&orderSpecificationSteelItem).
		Select(fmt.Sprintf("%s.*", orderSpecificationSteelItem.TableName())).
		Joins(fmt.Sprintf("join %s ON %s.order_specification_steel_id = %s.id", s.TableName(), s.TableName(), orderSpecificationSteelItem.TableName())).
		Where(fmt.Sprintf("%s.identifier = ?", s.TableName()), input.Identifier).
		Where( fmt.Sprintf("%s.company_id = ?", s.TableName()), me.CompanyId).
		First(&orderSpecificationSteelItem).Error
	if err != nil {
		return err
	}
	if err := steps.CheckSteelState(orderSpecificationSteelItem.State); err != nil {
		return err
	}

	return nil
}
