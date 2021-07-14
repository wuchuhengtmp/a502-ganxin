/**
 * @Desc    角色鉴权指令
 * @Author  wuchuheng<root@wuchuheng.com>
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
	if err := auth.ValidateToken(ctx); err != nil {
		return errors.AccessDenied(ctx, err.Error())
	}
	myRole, _ := me.GetRole()
	if allRoles.isContain(myRole.Tag) {
		return next(ctx)
	} else {
		var errMsg string
		roleStr := ""
		for _, role := range roles {
			roleStr += fmt.Sprintf("%s ", role)
		}
		errMsg = fmt.Sprintf("拒绝访问:需要任一的 %s %s 权限", errMsg, roleStr)

		return errors.AccessDenied(ctx, errMsg)
	}
}




