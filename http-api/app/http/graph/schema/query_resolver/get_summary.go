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
	"http-api/app/models/maintenance_record"
	"http-api/app/models/order_specification_steel"
	"http-api/app/models/projects"
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
			res.LeaseWeightTotal += newRes.LeaseWeightTotal
			res.IdleWeightTotal += newRes.IdleWeightTotal
			res.ScrapWeightTotal += newRes.ScrapWeightTotal
			res.MaintenanceWeightTotal += newRes.MaintenanceWeightTotal
			res.ProjectTotal += newRes.ProjectTotal
			res.LeaseTotal += newRes.LeaseTotal
			res.MaintenanceTotal += newRes.MaintenanceTotal
			res.LossTotal += newRes.LossTotal
		}

	} else {
		// 不是超管
		newRes, err := steps.GetDataByCompanyId(ctx, me.CompanyId)
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
	res.FeeTotal = price * res.WeightTotal
	res.YearFeeTotal = price * res.YearWeightTotal
	// 租赁数量(吨)
	orderSteelItem := order_specification_steel.OrderSpecificationSteel{}
	orderSteelTable := orderSteelItem.TableName()
	err = model.DB.Model(&orderSteelItem).
		Select("sum(weight) as LeaseWeightTotal").
		Joins(fmt.Sprintf("join %s ON %s.id = %s.steel_id", steelTable, steelTable, orderSteelTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.specification_id", specificationInfoTable, specificationInfoTable, steelTable)).
		Where(fmt.Sprintf("%s.company_id = ?", steelTable), companyId).Scan(&res).Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	// 闲置数量(吨)
	steelsItem := steels.Steels{}
	err = model.DB.Model(&steelsItem).
		Select("sum(weight) as IdleWeightTotal").
		Joins(fmt.Sprintf("join %s ON %s.id = %s.specification_id", specificationInfoTable, specificationInfoTable, steelTable)).
		Joins(fmt.Sprintf("join %s ON %s.steel_id = %s.id", orderSteelTable, orderSteelTable, steelTable)).
		Where(fmt.Sprintf("%s.company_id = ?", steelTable), companyId).
		Where(fmt.Sprintf("%s.id != %s.steel_id", steelTable, orderSteelTable)).
		Scan(&res).
		Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	// 报废数量(吨)
	err = model.DB.Model(&steelsItem).
		Select(fmt.Sprintf("sum(weight) as ScrapWeightTotal")).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.specification_id", specificationInfoTable, specificationInfoTable, steelTable)).
		Where(fmt.Sprintf("%s.company_id = ?", steelTable), companyId).
		Where(fmt.Sprintf("%s.state = ?", steelTable), steels.StateScrap).
		Scan(&res).
		Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	// 维修数量(吨)
	recordItem := maintenance_record.MaintenanceRecord{}
	recordTable := recordItem.TableName()
	err = model.DB.Model(&recordItem).
		Select(fmt.Sprintf("sum(weight) as MaintenanceWeightTotal")).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.steel_id", steelTable, steelTable, recordTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.specification_id", specificationInfoTable, specificationInfoTable, steelTable)).
		Where(fmt.Sprintf("%s.company_id = ?", steelTable), companyId).
		Scan(&res).
		Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	// 项目数
	projectItem := projects.Projects{}
	err = model.DB.Model(&projectItem).Where("company_id = ?", companyId).Count(&res.ProjectTotal).Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	// 总体租出
	err = model.DB.Model(&orderSteelItem).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.steel_id", steelTable, steelTable, orderSteelTable)).
		Where(fmt.Sprintf("%s.company_id = ?", steelTable), companyId).
		Count(&res.LeaseTotal).
		Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	// 维修数量
	err = model.DB.Model(&recordItem).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.steel_id", steelTable, steelTable, recordTable)).
		Where(fmt.Sprintf("%s.company_id = ?", steelTable), companyId).
		Count(&res.MaintenanceTotal).Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	// 丢失数量
	err = model.DB.Model(&steelsItem).
		Where("company_id = ?", companyId).
		Where("state = ?", steels.StateLost).
		Count(&res.LossTotal).
		Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	return &res, nil
}
