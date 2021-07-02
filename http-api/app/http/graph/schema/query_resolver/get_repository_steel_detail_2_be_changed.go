/**
 * @Desc    获取用于修改的仓库型钢
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

func (*QueryResolver)Get2BeChangedRepositorySteel(ctx context.Context, input graphModel.Get2BeChangedRepositorySteelInput) (*steels.Steels, error) {
	if err := requests.Get2BeChangedRepositorySteelRequest(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	var s steels.Steels
	me := auth.GetUser(ctx)
	err := model.DB.Model(&s).Where("company_id = ?", me.CompanyId).
		Where("identifier = ?", input.Identifier).
		First(&s).
		Error

	return &s, err
}
