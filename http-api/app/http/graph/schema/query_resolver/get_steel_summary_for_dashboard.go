/**
 * @Desc    型钢概览-饼图
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/13
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"gorm.io/gorm"
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/errors"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/steels"
	"http-api/pkg/model"
)

func (*QueryResolver) GetSteelSummaryForDashboard(ctx context.Context, input graphModel.GetSteelSummaryForDashboardInput) (*graphModel.GetSteelSummaryForDashboardRes, error) {
	if err := requests.ValidateGetSteelSummaryForDashboardRequest(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	getModelIn := func() *gorm.DB {
		modelIn := model.DB.Model(&steels.Steels{})
		// 指定仓库
		if input.RepositoryID != nil {
			modelIn = modelIn.Where("repository_id = ?", *input.RepositoryID)
		}
		me := auth.GetUser(ctx)
		role, _ := me.GetRole()
		if role.ID != role.ID {
			modelIn = modelIn.Where("company_id = ?", me.CompanyId)
		}
		return modelIn
	}
	var total int64
	if err := getModelIn().Count(&total).Error; err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	var res graphModel.GetSteelSummaryForDashboardRes
	// 项目中
	var usingPercent int64
	err := getModelIn().Where("state IN ?",
		[]int64{
			steels.StateRepository2Project,   //【仓库】-运送至项目途中
			steels.StateProjectWillBeUsed,    //【项目】-待使用
			steels.StateProjectInUse,         //【项目】-使用中
			steels.StateProjectException,     //【项目】-异常
			steels.StateProjectIdle,          //【项目】-闲置
			steels.StateProjectWillBeStore,   //【项目】-准备归库
			steels.StateProjectOnTheStoreWay, //【项目】-归库途中
		},
	).Count(&usingPercent).Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	res.UsingPercent = float64(usingPercent) / float64(total)
	// 维修中
	var maintaning int64
	err =getModelIn().Where("state IN ?", []int64{
		steels.StateRepository2Maintainer,      //"【仓库】-运送至维修厂途中",
		steels.StateMaintainerWillBeMaintained, //"【维修】-待维修",
		steels.StateMaintainerBeMaintaining,    //"【维修】-维修中",
		steels.StateMaintainerWillBeStore,      //"【维修】-准备归库",
		steels.StateMaintainerOnTheStoreWay,    //"【维修】-归库途中",
	}).Count(&maintaning).Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	res.MaintainingPercent = float64(maintaning) / float64(total)
	// 报废
	var scraping int64
	err = getModelIn().Where("state IN ?", []int64{
		steels.StateScrap,
	}).Count(&scraping).Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	res.CrappedPercent = float64(scraping) / float64(total)
	// 丢失
	var losting int64
	err =getModelIn().Where("state IN ?", []int64{
		steels.StateLost,
	}).Count(&losting).Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	res.LostPercent = float64(losting) / float64(total)
	// 在库
	var storing int64
	err =getModelIn().Where("state IN ?", []int64{
		steels.StateInStore,
	}).Count(&storing).Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	res.StoredPercent = float64(storing) / float64(total)

	return &res, nil
}
