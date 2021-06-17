/**
 * @Desc    获取项目列表解析器
 * @Author  wuchuheng<root@wuchuheng.com>
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

func (QueryResolver)GetProjectLis(ctx context.Context) ([]*projects.Projects, error) {
	return projects.Projects{}.GetProjectList(ctx)
}

type ProjectItemResolver struct { }

func (ProjectItemResolver)Company(ctx context.Context, obj *projects.Projects) (*companies.Companies, error) {
	cm, err := obj.GetCompany()

	return &cm, err
}

func (ProjectItemResolver)LeaderList(ctx context.Context, obj *projects.Projects) ([]*users.Users, error) {
	return obj.GetLeaderList()
}
