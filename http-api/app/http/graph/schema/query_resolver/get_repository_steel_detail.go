/**
 * @Desc    获取仓库型钢详情
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
	"http-api/app/models/specificationinfo"
	"http-api/app/models/steels"
	"http-api/pkg/model"
)

func (*QueryResolver) GetRepositorySteelDetail(ctx context.Context, input graphModel.GetRepositorySteelInput) (*steels.GetRepositorySteelDetailRes, error) {
	if err := requests.ValidateGetRepositorySteelDetailRequest(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	var res steels.GetRepositorySteelDetailRes
	var list []*steels.Steels
	steelInstance := model.DB.Debug().Model(&steels.Steels{}).
		Where("repository_id = ?", input.ReposirotyID)
	if input.State != nil {
		steelInstance = steelInstance.Where("state = ?", *input.State)
	}
	if input.SpecificationID != nil {
		steelInstance = steelInstance.Where("specification_id = ?", *input.SpecificationID)
	}
	if err := steelInstance.Find(&list).Error; err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	var weightInfo struct{
		WeightTotal float64 `json:"weight_total"`
	}
	specificationInfoTable := specificationinfo.SpecificationInfo{}.TableName()
	err := steelInstance.
		Select("sum(weight) as weight_total").
		Joins(fmt.Sprintf("join %s ON %s.id = %s.specification_id", specificationInfoTable, specificationInfoTable, steels.Steels{}.TableName())).
		Scan(&weightInfo).Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	res.Weight = weightInfo.WeightTotal
	res.List = list
	res.Total = int64(len(res.List))

	return &res, nil
}
