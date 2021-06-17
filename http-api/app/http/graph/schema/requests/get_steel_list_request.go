/**
 * @Desc    获取型钢列表验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/16
 * @Listen  MIT
 */
package requests

import (
	"context"
	"fmt"
	"http-api/app/http/graph/model"
	"http-api/app/models/repositories"
	"http-api/app/models/specificationinfo"
)

func ValidateGetSteelListRequest(ctx context.Context, input model.PaginationInput) error {
	// 有没有这个仓库
	if input.RepositoryID != nil {
		r := repositories.Repositories{
			ID: *input.RepositoryID,
		}
		if err := r.IsExists(ctx); err != nil {
			return fmt.Errorf("没有这个仓库")
		}
	}
	// 有没有这个规格
	if input.SpecificationID != nil {
		sp := specificationinfo.SpecificationInfo{ID: *input.SpecificationID}
		if err := sp.IsExist(ctx); err != nil {
			return fmt.Errorf("没有这个格式")
		}
	}

	return nil
}
