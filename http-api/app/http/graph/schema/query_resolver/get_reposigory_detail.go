/**
 * @Desc    The query_resolver is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/12
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/errors"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/models/repositories"
	"http-api/pkg/model"
)

func (*QueryResolver)GetRepositoryDetail(ctx context.Context, input graphModel.GetRepositoryDetailInput) ([]*repositories.Repositories, error)  {
	me := auth.GetUser(ctx)
	repositoryItem := repositories.Repositories{}
	modelIn := model.DB.Model(&repositoryItem).
		Where("company_id = ?", me.CompanyId)
	if input.RepositoryID != nil {
		modelIn = modelIn.Where( "id = ?", input.RepositoryID)
	}
	var res []*repositories.Repositories
	err := modelIn.Find(&res).Error
	if err != nil {
		return res, errors.ServerErr(ctx, err)
	}

	return res, nil
}
