/**
 * @Desc    删除制造商解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/8
 * @Listen  MIT
 */
package mutation_resolver

import (
	"context"
	"http-api/app/http/graph/errors"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/codeinfo"
)

func (*MutationResolver)DeleteManufacturer(ctx context.Context, id int64) (bool, error) {
	if err := requests.ValidateDeleteManufacturerRequest(ctx, id); err != nil {
		return false, errors.ValidateErr(ctx, err)
	}

	c := codeinfo.CodeInfo{ID:id}
	if err := c.DeleteManufacturer(ctx); err != nil {
		return false, errors.ServerErr(ctx, err)
	}
	return true, nil
}
