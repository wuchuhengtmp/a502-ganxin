/**
 * @Desc    The query_resolver is part of http-api
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
	graphModel"http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/steels"
	"http-api/pkg/model"
)

func (*QueryResolver)GetOneSteelDetail(ctx context.Context, input graphModel.GetOneSteelDetailInput) (*steels.Steels, error) {
	me := auth.GetUser(ctx)
	if err := requests.ValidateGetOneSteelDetailRequest(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	s := steels.Steels{}
	model.DB.Model(&s).Where("company_id = ? AND identifier = ?", me.CompanyId, input.Identifier).First(&s)

	return &s, nil
}
