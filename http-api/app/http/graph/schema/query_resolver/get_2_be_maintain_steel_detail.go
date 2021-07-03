/**
 * @Desc    获取待维修型钢详情
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
	"http-api/app/http/graph/errors"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/specificationinfo"
	"http-api/app/models/steels"
	"http-api/pkg/model"
)

func (*QueryResolver)Get2BeMaintainSteelDetail(ctx context.Context, input graphModel.Get2BeMaintainSteelDetailInput) (*steels.Get2BeScrapRepositorySteelDetailRes, error) {
	res := steels.Get2BeScrapRepositorySteelDetailRes{}
	if err := requests.ValidateGet2BeMaintainSteelDetailRequest(ctx, input); err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	steelsItem := steels.Steels{}
	steelsTable := steels.Steels{}.TableName()
	modelInstance := model.DB.Model(&steelsItem).
		Where(fmt.Sprintf("%s.identifier IN ?", steelsTable), input.IdentifierList)
	if input.SpecificationID != nil {
		modelInstance = modelInstance.Where(fmt.Sprintf("%s.specification_id = ?", steelsTable), *input.SpecificationID)
	}
	if err := modelInstance.Select( fmt.Sprintf("%s.*", steelsTable)).Find(&res.List).Error; err != nil {

		return nil, errors.ServerErr(ctx, err)
	}
	specificationInfoTable := specificationinfo.SpecificationInfo{}.TableName()
	var weightInfo struct{
		Weight float64
	}
	err := modelInstance.Select("sum(weight) as Weight").
		Joins(fmt.Sprintf("join %s ON %s.id = %s.specification_id", specificationInfoTable, specificationInfoTable, steelsTable)).
		Scan(&weightInfo).
		Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	res.Total = int64(len(res.List))
	res.Weight = weightInfo.Weight

	return &res, nil
}