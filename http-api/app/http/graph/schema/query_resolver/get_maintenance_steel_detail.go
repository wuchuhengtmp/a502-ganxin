/**
 * @Desc    获取维修厂详细信息
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/7
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"fmt"
	"http-api/app/http/graph/errors"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/maintenance"
	"http-api/app/models/maintenance_record"
	"http-api/app/models/specificationinfo"
	"http-api/app/models/steels"
	"http-api/pkg/model"
)

func (*QueryResolver)GetMaintenanceSteelDetail(ctx context.Context, input graphModel.GetMaintenanceSteelDetailInput) (*maintenance.GetSteelForOutOfMaintenanceDetailRes, error) {
	if err := requests.ValidateGetMaintenanceSteelDetailRequest(ctx, input); err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	var res maintenance.GetSteelForOutOfMaintenanceDetailRes
	steelTable := steels.Steels{}.TableName()
	recordItem := maintenance_record.MaintenanceRecord{}
	specificationTable := specificationinfo.SpecificationInfo{}.TableName()
	modelInstance := model.DB.Model(&recordItem).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.steel_id", steelTable, steelTable, recordItem.TableName())).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.specification_id", specificationTable, specificationTable, steelTable)).
		Where(fmt.Sprintf("%s.maintenance_id = ?", recordItem.TableName()), input.MaintenanceID)

	if input.State != nil {
		modelInstance = modelInstance.Where(fmt.Sprintf("%s.state = ?", recordItem.TableName()), *input.State)
	}
	if input.SpecificationID != nil {
		modelInstance = modelInstance.Where(fmt.Sprintf("%s.specification_id = ?", steelTable), *input.SpecificationID)
	}

	if err := modelInstance.Select(fmt.Sprintf("%s.*", recordItem.TableName())).Scan(&res.List).Error; err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	var weightInfo struct{
		Weight float64
	}
	if err := modelInstance.Select(fmt.Sprintf("sum(weight) as Weight")).Scan(&weightInfo).Error; err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	res.Weight = weightInfo.Weight
	res.Total = int64(len(res.List))

	return &res, nil
}


