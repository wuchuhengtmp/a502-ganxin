/**
 * @Desc    获取维修详情
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/9
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
	"http-api/app/http/graph/util/helper"
	"http-api/app/models/maintenance_record"
	"http-api/app/models/projects"
	"http-api/app/models/specificationinfo"
	"http-api/app/models/steels"
	"http-api/pkg/model"
)

func (*QueryResolver) GetMaintenanceDetail(ctx context.Context, input graphModel.GetMaintenanceDetailInput) (*projects.GetMaintenanceDetailRes, error) {
	if err := requests.ValidateGetMaintenanceDetailRequest(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	recordItem := maintenance_record.MaintenanceRecord{}
	steelTable := steels.Steels{}.TableName()
	specificationTable := specificationinfo.SpecificationInfo{}.TableName()
	recordTable := recordItem.TableName()
	me := auth.GetUser(ctx)
	modeIns := model.DB.Debug().Model(&recordItem).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.steel_id", steelTable, steelTable, recordTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.specification_id", specificationTable, specificationTable, steelTable)).
		Where(fmt.Sprintf("%s.company_id = ?", steelTable), me.CompanyId)
	//  维修状态
	if input.State != nil {
		modeIns = modeIns.Where(fmt.Sprintf("%s.state = ?", recordTable), *input.State)
	}
	//  仓库id
	if input.RepositoryID != nil {
		modeIns = modeIns.Where(fmt.Sprintf("%s.repository_id = ?", steelTable), *input.RepositoryID)
	}
	//  型钢编码
	if input.SpecificationID != nil {
		modeIns = modeIns.Where(fmt.Sprintf("%s.specification_id = ?", steelTable), *input.SpecificationID)
	}
	// 编码
	if input.Code != nil {
		modeIns = modeIns.Where(fmt.Sprintf("%s.code like ?", steelTable), "%"+*input.Code+"%")
	}
	//  入厂时间
	if input.EnteredMaintenanceAt != nil {
		s, e := helper.GetSecondBetween(*input.EnteredMaintenanceAt)
		modeIns = modeIns.Where(fmt.Sprintf("%s.entered_at between ? AND ?", recordTable), s, e)
	}
	//  出厂时间
	if input.OutMaintenanceAt != nil {
		s, e := helper.GetSecondBetween(*input.OutMaintenanceAt)
		modeIns = modeIns.Where(fmt.Sprintf("%s.outed_at between ? AND ?", recordTable), s, e)

	}
	res := projects.GetMaintenanceDetailRes{}
	if err := modeIns.Count(&res.Total).Error; err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	var weightInfo struct {
		Weight float64
	}
	if err := modeIns.Select("sum(weight) as Weight").Scan(&weightInfo).Error; err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	res.Weight = weightInfo.Weight
	if !input.IsShowAll {
		o := (*input.Page - 1) * *input.PageSize
		modeIns = modeIns.Offset(int(o)).Limit(int(*input.PageSize))
	}
	err := modeIns.Select(fmt.Sprintf("%s.*", recordTable)).
		Order(fmt.Sprintf("%s.id desc", recordTable)).
		Scan(&res.List).
		Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	return &res, nil
}
