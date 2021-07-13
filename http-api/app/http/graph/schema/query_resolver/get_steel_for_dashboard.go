/**
 * @Desc    获取型钢列表(用于仪表盘)
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/19
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/models/steels"
)

func (*QueryResolver) GetSteelForDashboard(ctx context.Context, input *graphModel.GetSteelForDashboardInput) (*steels.GetSteelListRes, error) {
	var res steels.GetSteelListRes

	return &res, nil
}
