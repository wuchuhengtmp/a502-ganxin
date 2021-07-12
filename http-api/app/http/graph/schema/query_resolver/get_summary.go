/**
 * @Desc    获取资产概况
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
	"http-api/app/http/graph/util/helper"
	"http-api/app/models/companies"
	"http-api/app/models/configs"
	"http-api/app/models/roles"
	"http-api/app/models/specificationinfo"
	"http-api/app/models/steels"
	"http-api/pkg/model"
	"strconv"
	"time"
)

type GetSummarySteps struct{}

func (*QueryResolver) GetSummary(ctx context.Context) (*graphModel.GetSummaryRes, error) {
	res := graphModel.GetSummaryRes{}
	me := auth.GetUser(ctx)
	cRole, _ := me.GetRole()
	steps := GetSummarySteps{}
	// 超管
	if cRole.ID == roles.RoleAdminId {
		var companyList []*companies.Companies
		err := model.DB.Model(&companies.Companies{}).Find(&companyList).Error
		if err != nil {
			return nil, errors.ServerErr(ctx, err)
		}
		for _, item := range companyList {
			newRes, err := steps.GetDataByCompanyId(ctx, item.ID)
			if err != nil {
				return nil, err
			}
			res.WeightTotal += newRes.WeightTotal
			res.YearWeightTotal += newRes.YearWeightTotal
			res.FeeTotal += newRes.FeeTotal
			res.YearFeeTotal += newRes.YearFeeTotal
		}

	} else {
		// 不是超管
		newRes , err := steps.GetDataByCompanyId(ctx, me.CompanyId)
		if err != nil {
			return nil, err
		}
		res = *newRes
	}

	return &res, nil
}

func (*GetSummarySteps) GetDataByCompanyId(ctx context.Context, companyId int64) (*graphModel.GetSummaryRes, error) {
	res := graphModel.GetSummaryRes{}
	var total int64
	var yearTotal int64
	specificationInfoTable := specificationinfo.SpecificationInfo{}.TableName()
	steelTable := steels.Steels{}.TableName()
	// 总重量
	totalModelIn := model.DB.Debug().Model(&steels.Steels{}).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.specification_id", specificationInfoTable, specificationInfoTable, steelTable))
	totalModelIn = totalModelIn.Where(fmt.Sprintf("%s.company_id = ?", steelTable), companyId)
	err := totalModelIn.Select("sum(weight) as WeightTotal").Scan(&res).Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	// 总量
	if err := totalModelIn.Count(&total).Error; err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	// 今年总重量
	startTime, endTime := helper.GetYearBetween(time.Now())
	totalModelIn = totalModelIn.Where(fmt.Sprintf("%s.created_at between ? AND ?", steelTable), endTime, startTime)
	err = totalModelIn.Select("sum(weight) as YearWeightTotal").Scan(&res).Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	// 今年总量
	if err := totalModelIn.Count(&yearTotal).Error; err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	configItem := configs.Configs{}
	err = model.DB.Model(&configItem).Where("name = ? ", configs.PRICE_NAME).
		Where("company_id = ?", companyId).
		First(&configItem).
		Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	price, _ := strconv.ParseFloat(configItem.Value, 64)
	res.FeeTotal = price * float64(total)
	res.YearFeeTotal = price * float64(yearTotal)

	return &res, nil
}
