/**
 * @Desc 待出库详情信息解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/19
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"http-api/app/http/graph/errors"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/projects"
	"http-api/app/models/steels"
	"http-api/pkg/model"
)

func (*QueryResolver) GetProjectOrder2WorkshopDetail(ctx context.Context, input graphModel.ProjectOrder2WorkshopDetailInput) (*projects.GetProjectOrder2WorkshopDetailRes, error) {
	res := projects.GetProjectOrder2WorkshopDetailRes{}
	if err := requests.ValidateGetProject2WorkshopDetailRequest(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	modelInstance := model.DB.
		Model(&steels.Steels{}).
		Where("identifier in ?", input.IdentifierList)
	if input.SpecificationID != nil {
		modelInstance = modelInstance.Where("specification_id = ?", *input.SpecificationID)
	}
	if err := modelInstance.Find(&res.List).Error; err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	res.Total = int64(len(res.List))
	for _, item := range res.List {
		s, err := item.GetSpecification()
		if err != nil {
			return nil, errors.ServerErr(ctx, err)
		}
		res.TotalWeight += s.Weight
	}

	return &res, nil
}
