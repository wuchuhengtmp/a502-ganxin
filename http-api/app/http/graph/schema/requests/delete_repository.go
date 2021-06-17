/**
 * @Desc    删除仓库请求验证
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/4
 * @Listen  MIT
 */
package requests

import (
	"context"
	"fmt"
	"http-api/app/http/graph/auth"
	"http-api/app/models/repositories"
)

func  ValidateDeleteRepository(ctx context.Context, repositoryId int64) error {
	r := repositories.Repositories{ ID:  repositoryId, }
	if err := r.GetSelf(); err != nil {
		return fmt.Errorf("没有这个仓库")
	}
	me := auth.GetUser(ctx)
	if r.CompanyId != me.CompanyId {
		return fmt.Errorf("您要删除的仓库,与您不是同一家公司的，您无权删除")
	}

	return nil
 }

