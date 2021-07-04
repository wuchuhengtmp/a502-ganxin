/**
 * @Desc    获取待入场详细信息解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/4
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

func (*QueryResolver) GetEnterMaintenanceSteelDetail(ctx context.Context, input graphModel.GetEnterMaintenanceSteelDetailInput) (*maintenance.GetEnterMaintenanceSteelDetailRes, error) {
	res := maintenance.GetEnterMaintenanceSteelDetailRes{}
	if err := requests.ValidateGetEnterMaintenanceSteelDetailRequest(ctx, input); err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	specificationTable := specificationinfo.SpecificationInfo{}.TableName()
	steelTable := steels.Steels{}.TableName()
	maintenanceRecordTable := maintenance_record.MaintenanceRecord{}.TableName()
	modelInstance := model.DB.Model(&maintenance_record.MaintenanceRecord{}).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.steel_id", steelTable, steelTable, maintenanceRecordTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.specification_id", specificationTable, specificationTable, steelTable)).
		Where(fmt.Sprintf("%s.identifier = ?", steelTable), input.IdentifierList)
	if input.SpecificationID != nil {
		modelInstance = modelInstance.Where(fmt.Sprintf("%s.specification_id = ?", steelTable), *input.SpecificationID)
	}
	if err := modelInstance.Scan(&res.List).Error; err != nil {
		return &res, errors.ServerErr(ctx, err)
	}
	res.Total = int64(len(res.List))
	var weightInfo struct {
		Weight float64
	}
	if err := modelInstance.Select("sum(weight) as Weight").Scan(&weightInfo).Error; err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	res.Weight = weightInfo.Weight

	return &res, nil
}
