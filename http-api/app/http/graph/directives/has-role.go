/**
 * @Desc    角色鉴权指令
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/27
 * @Listen  MIT
 */
package directives

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"http-api/app/models/roles"
)

func HasRole (ctx context.Context, obj interface{}, next graphql.Resolver, role roles.Role) (interface{}, error) {
//if !getCurrentUser(ctx).HasRole(role) {
//// block calling the next resolver
//return nil, fmt.Errorf("Access denied")
//}

// or let it pass through
return next(ctx)
}




