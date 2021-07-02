/**
 * @Desc    型钢待报废查询解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/2
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/errors"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/steels"
	"http-api/pkg/model"
)

func (*QueryResolver)Get2BeScrapRepositorySteel(ctx context.Context, input graphModel.Get2BeScrapRepositorySteelInput) (*steels.Steels, error) {
	if err := requests.ValidateGet2BeScrapRepositorySteelRequest(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	s := steels.Steels{}
	me := auth.GetUser(ctx)
	err := model.DB.Model(&s).Where("company_id = ?", me.CompanyId).
		Where("identifier = ?", input.Identifier).
		First(&s).Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	return &s, nil
}