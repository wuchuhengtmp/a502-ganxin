/**
 * @Desc 待出库详情信息解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/19
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"http-api/app/http/graph/errors"
	"http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/steels"
)

func (*QueryResolver)GetProjectOrder2WorkshopDetail(ctx context.Context, input model.ProjectOrder2WorkshopDetail) (steelList []*steels.Steels, err error) {
	if err = requests.ValidateGetProject2WorkshopDetailRequest(ctx, input); err != nil {
		return steelList, errors.ValidateErr(ctx, err)
	}
	// todo 待出库详情信息解析器
	var steelsList []*steels.Steels

	return steelsList, nil
}
