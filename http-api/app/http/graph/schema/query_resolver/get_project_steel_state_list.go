/**
 * @Desc    The query_resolver is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/25
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"http-api/app/models/steels"
)

func (*QueryResolver)GetProjectSteelStateList(ctx context.Context) ([]*steels.StateItem, error) {
	var res []*steels.StateItem
	for _, state := range steels.GetStateForProject() {
		item := steels.StateItem{
			Desc: steels.StateCodeMapDes[state],
			State: state,
		}
		res = append(res, &item)
	}

	return res, nil
}
