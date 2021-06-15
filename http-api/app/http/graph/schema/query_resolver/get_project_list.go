/**
 * @Desc    The query_resolver is part of http-api
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/15
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"http-api/app/models/companies"
	"http-api/app/models/projects"
	"http-api/app/models/users"
)

type ProjectItemResolver struct { }

func (ProjectItemResolver)Company(ctx context.Context, obj *projects.Projects) (*companies.Companies, error) {
	cm, err := obj.GetCompany()

	return &cm, err
}

func (ProjectItemResolver)LeaderList(ctx context.Context, obj *projects.Projects) ([]*users.Users, error) {
	return obj.GetLeaderList()
}
