/**
 * @Desc    获取日志类型列表
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/8
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"fmt"
	"http-api/app/models/logs"
)

func (*QueryResolver)GetLogTypeList(ctx context.Context) ([]*logs.LogTypeItem, error) {
	var res  []*logs.LogTypeItem
	for _, flag := range logs.GetAllType() {
		l := logs.LogTypeItem{
			Flag: fmt.Sprintf("%s", flag),
			Desc: logs.ActionTypeMapDes[flag],
		}
		res = append(res, &l)
	}

	return res, nil
}
