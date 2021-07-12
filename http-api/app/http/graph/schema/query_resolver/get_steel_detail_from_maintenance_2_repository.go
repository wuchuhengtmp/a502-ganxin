/**
 * @Desc    The query_resolver is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/12
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
	"http-api/app/models/maintenance_record"
	"http-api/app/models/projects"
	"http-api/app/models/specificationinfo"
	"http-api/app/models/steels"
	"http-api/pkg/model"
)

func (*QueryResolver) GetSteelDetailFromMaintenance2Repository(ctx context.Context, input graphModel.GetSteelDetailFromMaintenance2RepositoryInput) (*projects.GetMaintenanceDetailRes, error) {
	if err := requests.GetSteelDetailFromMaintenance2RepositoryRequest(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	me := auth.GetUser(ctx)
	steelsItem := steels.Steels{}
	steelTable := steelsItem.TableName()
	specificationTable := specificationinfo.SpecificationInfo{}.TableName()
	recordItem := maintenance_record.MaintenanceRecord{}
	recordTable := recordItem.TableName()
	modeIn := model.DB.Model(&recordItem).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.steel_id",steelTable, steelTable, recordTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.specification_id", specificationTable, specificationTable, steelTable)).
		Where(fmt.Sprintf("%s.identifier IN ?", steelTable), input.IdentifierList).
		Where(fmt.Sprintf("%s.company_id = ?", steelTable), me.CompanyId)
	if input.SpecificationID != nil {
		modeIn = modeIn.Where(fmt.Sprintf("%s.id = ?", specificationTable), *input.SpecificationID)
	}
	res := projects.GetMaintenanceDetailRes{}
	if err := modeIn.Select(fmt.Sprintf("%s.*", recordTable)).Find(&res.List).Error; err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	var weightInfo struct{ Weight float64 }
	if err := modeIn.Select("sum(weight) as Weight").Scan(&weightInfo).Error; err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	res.Weight = weightInfo.Weight
	res.Total = int64(len(res.List))

	return &res, nil
}
