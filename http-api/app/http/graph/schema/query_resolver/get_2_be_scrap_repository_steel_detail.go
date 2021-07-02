/**
 * @Desc    获取待报废型钢详情
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

func (*QueryResolver)Get2BeScrapRepositorySteelDetail(ctx context.Context, input graphModel.Get2BeScrapRepositorySteelDetailInput) (*steels.Get2BeScrapRepositorySteelDetailRes, error) {
	if err := requests.ValidateGet2BeScrapRepositorySteelDetailRequest(ctx, input); err != nil {
		return nil, err
	}
	res := steels.Get2BeScrapRepositorySteelDetailRes{}
	s := steels.Steels{}
	i := model.DB.Model(&s).Where("identifier IN ?", input.IdentifierList)
	if input.SpecificationID != nil {
		i = i.Where("specification_id = ?", *input.SpecificationID)
	}
	if err := i.Find(&res.List).Error; err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	res.Total = int64(len(res.List))
	// 重量
	specificationInfoTable := specificationinfo.SpecificationInfo{}.TableName()
	steelTable := steels.Steels{}.TableName()
	m := model.DB.Model(&specificationinfo.SpecificationInfo{}).
		Select("sum(weight) as Weight").
		Joins(fmt.Sprintf("join %s ON %s.specification_id = %s.id", steelTable, steelTable, specificationInfoTable)).
		Where(fmt.Sprintf("%s.identifier IN ?", steelTable), input.IdentifierList)
	if input.SpecificationID != nil {
		m = m.Where(fmt.Sprintf("%s.id = ?", specificationInfoTable), *input.SpecificationID)
	}
	var weightInfo struct{Weight float64}
	if err := m.Scan(&weightInfo).Error; err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	 res.Weight = weightInfo.Weight

	return &res, nil
}
