/**
 * @Desc    获取维修型钢的状态列表
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/6
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"http-api/app/models/steels"
)

func (*QueryResolver) GetStateListForMaintenanceSteelDetail(ctx context.Context) ([]*steels.StateItem, error) {
	var res []*steels.StateItem
	for _, state := range steels.GetMaintenanceStateListForDetail() {
		res = append(res, &steels.StateItem{
			State: state,
			// 状态说明
			Desc: steels.StateCodeMapDes[state],
		})
	}

	return res, nil
}
