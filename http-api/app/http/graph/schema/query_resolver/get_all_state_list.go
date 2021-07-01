/**
 * @Desc    获取全部状态列表解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/1
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"http-api/app/models/steels"
)

func (*QueryResolver)GetAllStateList(ctx context.Context) (res []*steels.StateItem, err error) {
	for _, i := range steels.GetAllStateList() {
		res = append(res, &steels.StateItem{
			State: i,
			Desc: steels.StateCodeMapDes[i],
		})

	}

	return
}
