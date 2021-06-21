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
	"http-api/app/models/steels"
	"http-api/pkg/model"
)

func (*QueryResolver)GetProjectOrder2WorkshopDetail(ctx context.Context, input graphModel.ProjectOrder2WorkshopDetailInput) (steelList []*steels.Steels, err error) {
	if err = requests.ValidateGetProject2WorkshopDetailRequest(ctx, input); err != nil {
		return steelList, errors.ValidateErr(ctx, err)
	}
	if input.SpecificationID != nil {
		if err := model.DB.
			Model(&steels.Steels{}).
			Where("identifier in ?", input.IdentifierList).
			Where("specification_id = ?", *input.SpecificationID).
			Find(&steelList).
			Error; err != nil {
			return steelList, err
		}
	} else {
		if err := model.DB.Model(&steels.Steels{}).Where("identifier in ?", input.IdentifierList).Find(&steelList).Error; err != nil {
			return steelList, err
		}
	}

	return
}
