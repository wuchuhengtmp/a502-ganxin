/**
 * @Desc    The query_resolver is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/2
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"fmt"
	"http-api/app/http/graph/errors"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/repositories"
	"http-api/app/models/specificationinfo"
	"http-api/app/models/steels"
	"http-api/pkg/model"
)

func (*QueryResolver) GetRepositorySteel(ctx context.Context, input graphModel.GetRepositorySteelInput) (*repositories.GetRepositorySteelRes, error) {
	var res repositories.GetRepositorySteelRes
	if err := requests.ValidateGetRepositorySteelRequest(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	steelsItem := steels.Steels{}
	modelInstance := model.DB.Debug().Model(&steelsItem).
		Select(fmt.Sprintf("count(*) as total, specification_id")).
		Where("repository_id = ?", input.ReposirotyID)
	if input.SpecificationID != nil {
		modelInstance = modelInstance.Where("specification_id = ?", *input.SpecificationID)
	}
	if input.State != nil {
		modelInstance = modelInstance.Where("state = ?", *input.State)
	}
	var totalInfo  []struct {
		Total           int64 `json:"total"`
		SpecificationId int64 `json:"specification_id"`
	}
	if err := modelInstance.Group("specification_id").Scan(&totalInfo).Error; err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	for _, i := range totalInfo {
		var listItem  repositories.GetRepositorySteelListItemRes
		listItem.Total = i.Total
		var weightInfo struct{
			WeightTotal float64 `json:"weight_total"`
		}
		err := model.DB.Model(&specificationinfo.SpecificationInfo{}).
			Select("SUM(weight) as weight_total").
			Where("id = ?", i.SpecificationId).
			Scan(&weightInfo).Error
		if err != nil {
			return nil, errors.ServerErr(ctx, err)
		}
		specificationInfoItem := specificationinfo.SpecificationInfo{}
		err = model.DB.Model(&specificationInfoItem).
			Where("id = ?", i.SpecificationId).
			First(&specificationInfoItem).
			Error
		if err != nil {
			return nil, err
		}
		listItem.SpecificationInfo = specificationInfoItem
		listItem.Weight = weightInfo.WeightTotal * float64( listItem.Total)
		res.List = append(res.List, &listItem)
		res.Total += listItem.Total
		res.Weight += listItem.Weight
	}

	return &res, nil
}
