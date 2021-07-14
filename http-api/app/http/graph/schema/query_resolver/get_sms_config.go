/**
 * @Desc    The query_resolver is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/14
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"http-api/app/http/graph/model"
	"http-api/app/models/configs"
)

func (*QueryResolver) GetSMSConfig(ctx context.Context) (*model.GetSMSConfigRes, error) {
	res := model.GetSMSConfigRes{
		AccessKey:       configs.GetGlobalVal(configs.SMS_ACCESS_KEY),
		AccessSecretKey: configs.GetGlobalVal(configs.SMS_ACCESS_SECRET_KEY),
		Sign:            configs.GetGlobalVal(configs.SMS_SIGN),
		Template:        configs.GetGlobalVal(configs.SMS_TEMPLATECODE),
	}

	return &res, nil
}
