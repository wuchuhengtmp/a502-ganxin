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
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/errors"
	"http-api/app/models/roles"
)

type Roles []roles.GraphqlRole

/**
 *  是否包含这个角色
 */
func (r *Roles)isContain(role roles.GraphqlRole) bool {
	for _, e := range *r {
		if e == role  {
			return true
		}
	}

	return false
}

func HasRole (ctx context.Context, obj interface{}, next graphql.Resolver, roles []roles.GraphqlRole) (interface{}, error) {
	var allRoles Roles = roles
	me := auth.GetUser(ctx)
	if me == nil {
		return errors.InvalidToken(ctx)
	}
	myRole, _ := me.GetRole()
	if allRoles.isContain(myRole.Tag) {
		return next(ctx)
	} else {
		var errMsg string
		for _, role := range roles {
			errMsg = fmt.Sprintf("%s %s", errMsg, role)
		}

		return errors.AccessDenied(ctx, errMsg)
	}
}




