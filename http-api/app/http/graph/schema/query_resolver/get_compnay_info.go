/**
 * @Desc    获取公司人员解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/3
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"http-api/app/http/graph/errors"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/services"
)

func (*QueryResolver)GetCompanyInfo(ctx context.Context) (*graphModel.GetCompnayInfoRes, error) {
	res, err := services.GetCompanyInfo(ctx)
	if  err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	return res, nil
}
