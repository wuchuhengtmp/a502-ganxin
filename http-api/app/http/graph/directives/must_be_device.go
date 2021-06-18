/**
 * @Desc    必须是设备指令
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/18
 * @Listen  MIT
 */
package directives

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/errors"
)

func MustBeDevice (ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	if !auth.IsDevice(ctx) {
		errMsg := fmt.Sprintf("必须是设备才能访问")
		return errors.AccessDenied(ctx, errMsg)
	}

	return next(ctx)
}
