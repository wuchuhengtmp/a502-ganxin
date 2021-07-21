/**
 * @Desc    获取项目详情解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/8
 * @Listen  MIT
 */
package query_resolver

import "C"
import (
	"context"
	"fmt"
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/errors"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/http/graph/util/helper"
	"http-api/app/models/order_specification"
	"http-api/app/models/order_specification_steel"
	"http-api/app/models/orders"
	"http-api/app/models/projects"
	"http-api/app/models/specificationinfo"
	"http-api/app/models/steels"
	"http-api/pkg/model"
)

func (*QueryResolver)GetProjectDetail(ctx context.Context, input graphModel.GetProjectDetailInput) (*projects.GetProjectDetailRes, error) {
	if err := requests.ValidateGetProjectDetailRequest(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	res :=  projects.GetProjectDetailRes{}
	orderSteelItem := order_specification_steel.OrderSpecificationSteel{}
	projectTable := projects.Projects{}.TableName()
	orderTable := orders.Order{}.TableName()
	orderSpecificationTable := order_specification.OrderSpecification{}.TableName()
	me := auth.GetUser(ctx)
	steelTable := steels.Steels{}.TableName()
	modelIns := model.DB.Debug().Model(&orderSteelItem).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.order_specification_id", orderSpecificationTable, orderSpecificationTable,  orderSteelItem.TableName())).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.order_id", orderTable, orderTable, orderSpecificationTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.project_id", projectTable, projectTable, orderTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.steel_id", steelTable, steelTable, orderSteelItem.TableName())).
		Where(fmt.Sprintf("%s.company_id = ?", projectTable), me.CompanyId).
		Where(fmt.Sprintf("%s.id = ?", projectTable), input.ProjectID)

	// 订单id过滤
	if input.OrderID != nil {
		modelIns = modelIns.Where(fmt.Sprintf("%s.id = ?", orderTable), *input.OrderID)
	}
	//  出库仓库id过滤
	if input.RepositoryID != nil {
		modelIns = modelIns.Where(fmt.Sprintf("%s.repository_id = ?", steelTable), *input.RepositoryID)
	}
	//  规格尺寸id 过滤
	if input.SpecificationID != nil {
		modelIns = modelIns.Where(fmt.Sprintf("%s.specification_id = ?", steelTable), *input.SpecificationID)
	}
	// 状态
	if input.State != nil {
		modelIns = modelIns.Where(fmt.Sprintf("%s.state = ?", orderSteelItem.TableName()), *input.State)
	}
	// 出库时间
	if input.OutOfRepositoryAt != nil {
		s, e := helper.GetSecondBetween(*input.OutOfRepositoryAt)
		modelIns = modelIns.Where(fmt.Sprintf("%s.out_repository_at BETWEEN ? AND ?", orderSteelItem.TableName()), s, e)
	}
	// 入库时间
	if input.EnterRepositoryAt != nil {
		s, e := helper.GetSecondBetween(*input.EnterRepositoryAt)
		modelIns = modelIns.Where(fmt.Sprintf("%s.enter_repository_at BETWEEN ? AND ?", orderSteelItem.TableName()), s, e)
	}
	//  入场时间
	if input.EnteredWorkshopAt != nil {
		s, e := helper.GetSecondBetween(*input.EnteredWorkshopAt)
		modelIns = modelIns.Where(fmt.Sprintf("%s.enter_workshop_at BETWEEN ? AND ?", orderSteelItem.TableName()), s, e)
	}
	// 出场时间
	if input.OutOfWorkshopAt != nil {
		s, e := helper.GetSecondBetween(*input.OutOfWorkshopAt)
		modelIns = modelIns.Where(fmt.Sprintf("%s.out_workshop_at BETWEEN ? AND ?", orderSteelItem.TableName()), s, e)
	}
	// 安装编码
	if input.LocationCode != nil {
		modelIns = modelIns.Where(fmt.Sprintf("%s.location_code = ?", orderSteelItem.TableName()), *input.LocationCode)
	}
	if input.InstallationAt != nil  {
		modelIns = modelIns.Where(fmt.Sprintf("%s.installation_at = ?", orderSteelItem.TableName()), *input.InstallationAt)
	}
	// 数量
	if err := modelIns.Count(&res.Total).Error ; err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	// 重量
	var weightInfo struct{ Weight float64 }
	specificationTable := specificationinfo.SpecificationInfo{}.TableName()
	err := modelIns.Select(fmt.Sprintf("sum(weight) as Weight")).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.specification_id", specificationTable, specificationTable, steelTable)).
		Scan(&weightInfo).Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	res.Weight = weightInfo.Weight

	// 分页
	if !input.IsShowAll {
		modelIns = modelIns.Limit(int(*input.PageSize)).Offset(int((*input.Page - 1) * *input.PageSize))
	}
	err = modelIns.Select(fmt.Sprintf("%s.*", orderSteelItem.TableName())).Scan(&res.List).Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	return &res, nil
}