/**
 * @Desc    获取维修出库型钢
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/3
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

func (*QueryResolver)Get2BeMaintainSteel(ctx context.Context, input graphModel.Get2BeMaintainSteelInput) (*steels.Steels, error) {
	if err := requests.ValidateGet2BeMaintainSteelRequest(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	steelItem := steels.Steels{}
	me := auth.GetUser(ctx)
	err :=  model.DB.Model(&steelItem).Where("company_id  = ? ", me.CompanyId).
		Where("identifier = ?", input.Identifier).
		First(&steelItem).
		Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	return &steelItem, nil
}


