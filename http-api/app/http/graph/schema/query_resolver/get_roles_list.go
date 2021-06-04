/**
 * @Desc    获取角色列表解析器
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/4
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"http-api/app/http/graph/errors"
	"http-api/app/models/roles"
)

func (*QueryResolver)GetRoleList(ctx context.Context) ([]*roles.RoleItem, error) {
	res, err := roles.GetRolesGraphRes()
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	return res, nil
}
