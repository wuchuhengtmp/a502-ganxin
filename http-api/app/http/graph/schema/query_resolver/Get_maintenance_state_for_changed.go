/**
 * @Desc    获取用于修改型钢状态的状态列表
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/5
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"http-api/app/models/steels"
)

func (*QueryResolver)GetMaintenanceStateForChanged(ctx context.Context) (res []*steels.StateItem, err error) {
	for _, state := range steels.GetMaintenanceStateListForChanged() {
		res = append(res, &steels.StateItem{
			State: state,
			Desc: steels.StateCodeMapDes[state],
		})
	}

	return
}
