/**
 * @Desc 待出库详情信息解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/19
 * @Listen  MIT
 */
package mutation_resolver

import (
	"context"
	"http-api/app/http/graph/model"
	"http-api/app/models/steels"
)

func (*MutationResolver)ProjectOrder2WorkshopDetail(ctx context.Context, input model.ProjectOrder2WorkshopDetail) ([]*steels.Steels, error) {
	// todo 待出库详情信息解析器
	var steelsList []*steels.Steels

	return steelsList, nil
}
