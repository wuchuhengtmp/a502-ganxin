/**
 * @Desc    创建项目解析器
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/15
 * @Listen  MIT
 */
package mutation_resolver

import (
	"context"
	"http-api/app/http/graph/errors"
	"http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/projects"
)

func (*MutationResolver)CreateProject(ctx context.Context, input model.CreateProjectInput) (*projects.Projects, error) {
	if err := requests.ValidateCreateProjectRequest(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	p := projects.Projects{}
	if err := p.CreateProject(ctx, input); err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	return &p, nil
}
