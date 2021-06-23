/**
 * @Desc    快速查询多个型钢解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/22
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/errors"
	graphqlModel "http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/steels"
	"http-api/pkg/model"
)

func (*QueryResolver) GetMultipleSteelDetail(ctx context.Context, input *graphqlModel.GetMultipleSteelDetailInput) ([]*steels.Steels, error) {
	var ss []*steels.Steels
	if err := requests.ValidateGetMultipleSteelDetailRequest(ctx, input); err != nil {
		return ss, errors.ValidateErr(ctx, err)
	}
	me := auth.GetUser(ctx)
	err := model.DB.
		Model(&steels.Steels{}).
		Where("identifier in ? AND company_id = ?", input.IdentifierList, me.CompanyId).
		Find(&ss).
		Error
	if err != nil {
		return ss, errors.ServerErr(ctx, err)
	}

	return ss, nil
}
