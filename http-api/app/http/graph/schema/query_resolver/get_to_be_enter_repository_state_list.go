/**
 * @Desc    The query_resolver is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/30
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"http-api/app/models/steels"
)

func (*QueryResolver)GetToBeEnterRepositoryStateList(ctx context.Context) (res []*steels.StateItem, err error) {
	for _, state := range steels.GetStateListForEnterRepository() {
		res = append(res, &steels.StateItem{
			Desc: steels.StateCodeMapDes[state],
			State: state,
		})
	}

	return
}
