/**
 * @Desc    删除仓库解析器
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/4
 * @Listen  MIT
 */
package mutation_resolver

import (
	"context"
	"http-api/app/http/graph/errors"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/repositories"
)

func (*MutationResolver)DeleteRepository(ctx context.Context, repositoryID int64) (bool, error) {
	if err := requests.ValidateDeleteRepository(ctx, repositoryID);err != nil {
		return false, errors.ValidateErr(ctx, err)
	}
	if err := repositories.DeleteById(ctx, repositoryID); err != nil {
		return false, errors.ServerErr(ctx, err)
	}

	return true, nil
}
