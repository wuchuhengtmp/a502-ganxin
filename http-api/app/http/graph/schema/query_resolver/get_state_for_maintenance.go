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
	"http-api/app/models/steels"
)

func (*QueryResolver)GetStateForMaintenance(ctx context.Context) ([]*steels.StateItem, error) {
	var res []*steels.StateItem
	for _, state := range steels.GetMaintenanceStateList() {
		res = append(res, &steels.StateItem{
			State: state,
			Desc: steels.StateCodeMapDes[state],
		})
	}

	return res, nil
}
