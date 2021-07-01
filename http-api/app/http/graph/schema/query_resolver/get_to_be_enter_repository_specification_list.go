/**
 * @Desc    获取待归库的尺寸列表解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/30
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"fmt"
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/errors"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/order_specification"
	"http-api/app/models/orders"
	"http-api/app/models/specificationinfo"
	"http-api/pkg/model"
)

func (*QueryResolver) GetToBeEnterRepositorySpecificationList(ctx context.Context, input graphModel.GetToBeEnterRepositorySpecificationListInput) (res []*specificationinfo.SpecificationInfo, err error) {
	if err := requests.ValidateGetToBeEnterRepositorySpecificationList(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	specificationInfoItem := specificationinfo.SpecificationInfo{}
	orderTable := orders.Order{}.TableName()
	orderSpecificationTable := order_specification.OrderSpecification{}.TableName()
	//err = model.DB.Model(&specificationInfoItem).
	//	Select(fmt.Sprintf("%s.*", specificationInfoItem.TableName())).
	//	Joins(fmt.Sprintf("join %s ON %s.specification_id = %s.id", orderSpecificationTable, orderSpecificationTable, specificationInfoItem.TableName())).
	//	Joins(fmt.Sprintf("join %s ON %s.id = %s.order_id", orderTable, orderTable, orderSpecificationTable)).
	//	Where(fmt.Sprintf("%s.project_id = ?", orderTable), input.ProjectID).
	//	Scan(&res).
	//	Error
	me := auth.GetUser(ctx)
	var orderSpecificationList []*order_specification.OrderSpecification
	err = model.DB.Model(&order_specification.OrderSpecification{}).
		Select(fmt.Sprintf("%s.*", order_specification.OrderSpecification{}.TableName())).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.id", orderTable, orderTable, orderSpecificationTable)).
		Where(fmt.Sprintf("%s.project_id = ?", orderTable), input.ProjectID).
		Where(fmt.Sprintf("%s.company_id = ?", orderTable), me.CompanyId).
		Scan(&orderSpecificationList).Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	var specificationIds []int64
	idMapId := make(map[int64]int64)
	for _, i := range orderSpecificationList {
		if _, ok := idMapId[i.SpecificationId]; !ok {
			specificationIds = append(specificationIds, i.SpecificationId)
			idMapId[i.SpecificationId] = i.SpecificationId
		}
	}
	err = model.DB.Model(&specificationInfoItem).Where("id IN ?", specificationIds).Find(&res).Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	return
}
