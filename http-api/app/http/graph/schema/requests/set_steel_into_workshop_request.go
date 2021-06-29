/**
 * @Desc    型钢入场解析验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/24
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
	"http-api/app/models/project_leader"
	"http-api/app/models/projects"
	"http-api/app/models/steels"
	"http-api/pkg/model"
)

func ValidateSetSteelIntoWorkshopRequest(ctx context.Context, input graphModel.SetSteelIntoWorkshopInput) error {
	me := auth.GetUser(ctx)
	steps := ValidateGetProject2WorkshopDetailRequestSteps{}
	// 检验识码列表不能为空
	if err := steps.CheckIdentificationListMustBeEmpty(ctx, input.IdentifierList); err != nil {
		return err
	}
	// 检验订单
	o := orders.Order{}
	if err := model.DB.Model(&o).Where("id = ? AND company_id = ?", input.OrderID, me.CompanyId).First(&o).Error; err != nil {
		return err
	}
	// 检验订单是否归这个项目管理员名下
	projectLeaderItem := project_leader.ProjectLeader{}
	projectLeaderTable := project_leader.ProjectLeader{}.TableName()
	projectTable := projects.Projects{}.TableName()
	orderTable := orders.Order{}.TableName()
	err := model.DB.Model(&projectLeaderItem).
		Select(fmt.Sprintf("%s.*", projectLeaderTable)).
		Joins(fmt.Sprintf("join %s ON %s.id= %s.project_id", projectTable, projectTable, projectLeaderTable)).
		Joins(fmt.Sprintf("join %s ON %s.project_id = %s.id", orderTable, orderTable, projectTable)).
		Where(fmt.Sprintf("%s.uid = %d", projectLeaderTable, me.Id)).
		Where(fmt.Sprintf("%s.id = %d", orderTable, input.OrderID)).
		First(&projectLeaderItem).Error
	if err != nil {
		return fmt.Errorf("订单id为：%d 的项目管理员不是你,你无权操作", input.OrderID)
	}
	// 检验订单状态要为出库了
	if o.State != orders.StateSend && o.State != orders.StatePartOfReceipted {
		return fmt.Errorf("订单状态为: %s, 不能入库", orders.StateMapDesc[o.State])
	}
	// 标识码是否冗余
	if err := steps.CheckRedundancyIdentification(input.IdentifierList); err != nil {
		return err
	}
	// 检验型钢
	for _, steelIdentifier := range input.IdentifierList {
		steelItem := steels.Steels{}
		err := model.DB.Model(&steelItem).Where("identifier = ? AND company_id = ?", steelIdentifier, me.CompanyId).
			First(&steelItem).
			Error
		// 有没有
		if err != nil {
			return fmt.Errorf("没有识别码为%s的型钢", steelIdentifier)
		}
		// 状态对不对
		if steelItem.State != steels.StateRepository2Project {
			return fmt.Errorf("标识码为%s的型钢的状态为:%s, 不能入库",steelIdentifier, steels.StateCodeMapDes[steelItem.State])
		}
		// 这根型钢在不在订单中
		orderSpecificationSteelItem := order_specification_steel.OrderSpecificationSteel{}
		orderSpecificationSteelTable := order_specification_steel.OrderSpecificationSteel{}.TableName()
		orderSpecificationTable := order_specification.OrderSpecification{}.TableName()
		steelsTable := steels.Steels{}.TableName()
		orderTable := orders.Order{}.TableName()
		err = model.DB.Model(&orderSpecificationSteelItem).
			Select(fmt.Sprintf("%s.*", orderSpecificationSteelTable)).
			Joins(fmt.Sprintf("join %s ON %s.id = %s.steel_id", steelsTable, steelsTable, orderSpecificationSteelTable)).
			Joins(fmt.Sprintf("join %s ON %s.id = %s.order_specification_id", orderSpecificationTable, orderSpecificationTable, orderSpecificationSteelTable)).
			Joins(fmt.Sprintf("join %s ON %s.id = %s.order_id", orderTable, orderTable, orderSpecificationTable)).
			Where(fmt.Sprintf("%s.identifier = '%s'", steelsTable , steelIdentifier) ).
			First(&orderSpecificationSteelItem).
			Error
		if err != nil {
			return fmt.Errorf("型钢订单为:%d,的规格列表中不包含标识码表为:%s的型钢", o.Id, steelIdentifier)
		}
	}

	return nil
}
