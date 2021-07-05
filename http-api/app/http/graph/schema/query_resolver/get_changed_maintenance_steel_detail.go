/**
 * @Desc    待修改型钢详细信息解析器
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
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/errors"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/maintenance"
	"http-api/app/models/maintenance_record"
	"http-api/app/models/specificationinfo"
	"http-api/app/models/steels"
	"http-api/pkg/model"
)

func (*QueryResolver) GetChangedMaintenanceSteelDetail(ctx context.Context, input graphModel.GetChangedMaintenanceSteelDetailInput) (*maintenance.GetChangedMaintenanceSteelDetailRes, error) {
	var res maintenance.GetChangedMaintenanceSteelDetailRes
	if err := requests.ValidateGetChangedMaintenanceSteelDetailRequest(ctx, input); err != nil {

		return nil, errors.ValidateErr(ctx, err)
	}
	steelTable := steels.Steels{}.TableName()
	recordItem := maintenance_record.MaintenanceRecord{}
	me := auth.GetUser(ctx)
	modelInstance := model.DB.Model(&recordItem).
		Joins(fmt.Sprintf("join %s ON %s.maintenance_record_steel_id = %s.id", steelTable, steelTable, recordItem.TableName())).
		Where(fmt.Sprintf("%s.identifier IN ?", steelTable), input.IdentifierList).
		Where(fmt.Sprintf("%s.company_id = ?", steelTable), me.CompanyId)
	if err := modelInstance.Select(fmt.Sprintf("%s.*", recordItem.TableName())).Scan(&res.List).Error; err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	res.Total = int64(len(res.List))
	var weightInfo struct {
		Weight float64
	}

	specificationTable := specificationinfo.SpecificationInfo{}.TableName()
	err := modelInstance.Joins(fmt.Sprintf("join %s ON %s.id = %s.specification_id", specificationTable, specificationTable, steelTable)).
		Select("sum(weight) as Weight").Scan(&weightInfo).Error
	if err != nil {
		return &res, errors.ServerErr(ctx, err)
	}
	res.Weight = weightInfo.Weight

	return &res, nil
}
