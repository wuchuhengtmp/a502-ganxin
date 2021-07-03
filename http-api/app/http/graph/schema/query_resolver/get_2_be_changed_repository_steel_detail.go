/**
 * @Desc    The query_resolver is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/3
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
	"http-api/app/models/specificationinfo"
	"http-api/app/models/steels"
	"http-api/pkg/model"
)

func (*QueryResolver) Get2BeChangedRepositorySteelDetail(ctx context.Context, input graphModel.Get2BeChangedRepositorySteelDetailInput) (*steels.Get2BeChangedRepositorySteelDetailRes, error) {
	if err := requests.Get2BeChangedRepositorySteelDetail(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	me := auth.GetUser(ctx)
	steelTable := steels.Steels{}.TableName()
	var res steels.Get2BeChangedRepositorySteelDetailRes
	// 列表
	modelInstance := model.DB.Model(&steels.Steels{}).
		Where(fmt.Sprintf("%s.company_id = ?", steelTable), me.CompanyId).
		Where(fmt.Sprintf("%s.identifier IN ?", steelTable), input.IdentifierList)
	if input.SpecificationID != nil {
		modelInstance = modelInstance.Where(fmt.Sprintf("%s.specification_id = ?", steelTable), *input.SpecificationID)
	}
	if err := modelInstance.Scan(&res.List).Error; err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	specificationInfoTable := specificationinfo.SpecificationInfo{}.TableName()
	var weightInfo struct{ Weight float64 }
	err := modelInstance.Select("sum(weight) as Weight").
		Joins(fmt.Sprintf("join %s ON %s.id = %s.specification_id", specificationInfoTable, specificationInfoTable, steelTable)).
		Scan(&weightInfo).
		Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	res.Total = int64(len(res.List))
	res.Weight = weightInfo.Weight

	return &res, nil
}
