/**
 * @Desc    The mutation_resolver is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/14
 * @Listen  MIT
 */
package mutation_resolver

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/errors"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/http/graph/schema/services"
	"http-api/app/models/configs"
	"http-api/pkg/model"
)

func (*MutationResolver) SetCompanyInfo(ctx context.Context, input graphModel.SetCompanyInfoInput) (*graphModel.GetCompnayInfoRes, error) {
	if err := requests.ValidateSetCompanyInfoRequest(input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	me := auth.GetUser(ctx)
	err := model.DB.Transaction(func(tx *gorm.DB) error {
		setConfig := func(key string, value string) error {
			err := tx.Model(&configs.Configs{}).Where("company_id = ?", me.CompanyId).
				Where("name = ?", key).
				Update("value", value).Error
			return err
		}
		// 设置教学
		if err := setConfig(configs.TUTOR_FILE_NAME, fmt.Sprintf("%d", input.TutorFileID)); err != nil {
			return err
		}
		// 设置微信
		if err := setConfig(configs.WECHAT_NAME, input.Wechat); err != nil {
			return err
		}
		// 设置客服
		if err := setConfig(configs.WECHAT_NAME, input.Wechat); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	res, err := services.GetCompanyInfo(ctx)
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	return res, nil
}
