/**
 * @Desc    待出厂的型钢详情
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/5
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

func (*QueryResolver) GetSteelForOutOfMaintenanceDetail(ctx context.Context, input graphModel.GetSteelForOutOfMaintenanceDetailInput) (*maintenance.GetSteelForOutOfMaintenanceDetailRes, error) {
	var res maintenance.GetSteelForOutOfMaintenanceDetailRes
	if err := requests.ValidateGetSteelForOutOfMaintenanceDetailRequest(ctx, input); err != nil {
		return &res, err
	}
	steelTable := steels.Steels{}.TableName()
	record := maintenance_record.MaintenanceRecord{}
	modelIn := model.DB.Model(&record).
		Joins(fmt.Sprintf("join %s ON %s.maintenance_record_steel_id = %s.id", steelTable, steelTable, record.TableName())).
		Where(fmt.Sprintf("%s.identifier IN ?", steelTable), input.IdentifierList)
	if input.SpecificationID != nil {
		modelIn = modelIn.Where(fmt.Sprintf("%s.specification_id = ?", steelTable), *input.SpecificationID)
	}
	if err := modelIn.Select(fmt.Sprintf("%s.*", record.TableName())).Scan(&res.List).Error; err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	specificationTable := specificationinfo.SpecificationInfo{}.TableName()
	var weightInfo struct{ Weight float64 }
	err := modelIn.Select("sum(weight) as Weight").
		Joins(fmt.Sprintf("join %s ON %s.id = %s.specification_id", specificationTable, specificationTable, steelTable)).
		Scan(&weightInfo).
		Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	res.Weight = weightInfo.Weight
	res.Total = int64(len(res.List))

	return &res, nil
}

