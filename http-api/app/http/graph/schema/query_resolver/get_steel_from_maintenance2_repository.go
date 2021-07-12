/**
 * @Desc    型钢归库查询
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/12
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

func (*QueryResolver)GetSteelFromMaintenance2Repository(ctx context.Context, input graphModel.GetSteelFromMaintenance2RepositoryInput) (*steels.Steels, error) {
	if err := requests.ValidateGetSteelFromMaintenance2RepositoryRequest(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	steelItem := steels.Steels{}
	me := auth.GetUser(ctx)
	err := model.DB.Model(&steelItem).Where("identifier = ?", input.Identifer).
		Where("company_id = ?", me.CompanyId).
		First(&steelItem).Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	return &steelItem, nil
}
