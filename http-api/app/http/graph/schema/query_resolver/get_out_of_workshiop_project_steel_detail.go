/**
 * @Desc    获取型钢单根型钢出场详情解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/28
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"fmt"
	"http-api/app/http/graph/errors"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/order_specification"
	"http-api/app/models/order_specification_steel"
	"http-api/app/models/orders"
	"http-api/app/models/projects"
	"http-api/app/models/specificationinfo"
	"http-api/app/models/steels"
	"http-api/pkg/model"
)

func (*QueryResolver) GetOutOfWorkshopProjectSteelDetail(ctx context.Context, input graphModel.GetOutOfWorkshopProjectSteelDetail) (*projects.GetOutOfWorkshopProjectSteelDetailRes, error) {
	var res projects.GetOutOfWorkshopProjectSteelDetailRes
	if err := requests.ValidateGetOutOfWorkshopProjectSteelDetailRequest(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	item := order_specification_steel.OrderSpecificationSteel{}
	steelTable := steels.Steels{}.TableName()
	err := model.DB.Model(&item).
		Select(fmt.Sprintf("%s.*", item.TableName())).
		Joins(fmt.Sprintf("join %s ON %s.order_specification_steel_id = %s.id", steelTable, steelTable, item.TableName())).
		Where(fmt.Sprintf("%s.identifier = ?", steelTable), input.Identifier).
		First(&item).
		Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	res.OrderSteel = item

	specificationItem := specificationinfo.SpecificationInfo{}
	err = model.DB.Model(&specificationItem).
		Select(fmt.Sprintf("%s.*", specificationItem.TableName())).
		Joins(fmt.Sprintf("join %s ON %s.specification_id = %s.id", steelTable, steelTable, specificationItem.TableName())).
		First(&specificationItem).
		Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	// 已归库数量
	orderSpecificationSteelTable := order_specification_steel.OrderSpecificationSteel{}.TableName()
	orderSpecificationTable := order_specification.OrderSpecification{}.TableName()
	orderTable := orders.Order{}.TableName()
	err = model.DB.Model(&order_specification_steel.OrderSpecificationSteel{}).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.order_specification_id", orderSpecificationTable, orderSpecificationTable, orderSpecificationSteelTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.order_id", orderTable, orderTable, orderSpecificationTable)).
		Where(fmt.Sprintf("%s.specification_id = ?", orderSpecificationTable), specificationItem.ID).
		Where(fmt.Sprintf("%s.project_id = ?", orderTable), input.ProjectID).
		Where(fmt.Sprintf("%s.state = ?", orderSpecificationSteelTable), steels.StateInStore).
		Count(&res.StoreTotal).Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	// 待归库数量
	err = model.DB.Model(&order_specification_steel.OrderSpecificationSteel{}).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.order_specification_id", orderSpecificationTable, orderSpecificationTable, orderSpecificationSteelTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.order_id", orderTable, orderTable, orderSpecificationTable)).
		Where(fmt.Sprintf("%s.specification_id = ?", orderSpecificationTable), specificationItem.ID).
		Where(fmt.Sprintf("%s.project_id = ?", orderTable), input.ProjectID).
		Where(fmt.Sprintf("%s.state in ?", orderSpecificationSteelTable), []int64{
			steels.StateProjectWillBeUsed,     //【项目】-待使用
			steels.StateProjectInUse,          //【项目】-使用中
			steels.StateProjectException,      //【项目】-异常
			steels.StateProjectIdle,           //【项目】-闲置
			steels.StateProjectWillBeStore,    //【项目】-准备归库
			steels.StateProjectOnTheStoreWay,  //【项目】-归库途中
		}).
		Count(&res.ToBeStoreTotal).Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	return &res, nil
}
