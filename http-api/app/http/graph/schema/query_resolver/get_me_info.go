/**
 * @Desc    获取我的信息
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/15
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"http-api/app/http/graph/auth"
	"http-api/app/models/users"
)

func (QueryResolver)GetMyInfo(ctx context.Context) (*users.Users, error) {
	return auth.GetUser(ctx), nil
}
