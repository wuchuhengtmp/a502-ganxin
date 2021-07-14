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
	"gorm.io/gorm"
	"http-api/app/http/graph/errors"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/models/configs"
	"http-api/pkg/model"
)

func (*MutationResolver)SetSMSConfig(ctx context.Context, input graphModel.SetSMSConfigInput) (*graphModel.GetSMSConfigRes, error) {
	setConfig := func(keyName string, value string, tx *gorm.DB) error {
		err := tx.Model(&configs.Configs{}).Where("name = ?", keyName).
			Where("company_id = ?", 0).
			Update("value", value).Error

		return err
	}
	err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 设置短信模板
		if err := setConfig(configs.SMS_TEMPLATECODE, input.Template, tx); err != nil {
			return err
		}
		// 设置签名
		if err := setConfig(configs.SMS_SIGN, input.Sign, tx); err != nil {
			return err
		}
		// 设置accessKey
		if err := setConfig(configs.SMS_ACCESS_KEY, input.AccessKey, tx); err != nil {
			return err
		}
		// 设置accessSecretKey
		if err := setConfig(configs.SMS_ACCESS_SECRET_KEY, input.AccessKey, tx); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	res := graphModel.GetSMSConfigRes{
		AccessKey: configs.GetGlobalVal(configs.SMS_ACCESS_KEY),
		AccessSecretKey: configs.GetGlobalVal(configs.SMS_ACCESS_SECRET_KEY),
		Sign: configs.GetGlobalVal(configs.SMS_SIGN),
		Template: configs.GetGlobalVal(configs.SMS_TEMPLATECODE),
	}

	return &res, nil
}

