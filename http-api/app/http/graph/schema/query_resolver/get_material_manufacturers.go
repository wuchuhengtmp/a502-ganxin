/**
 * @Desc    获取材料商列表解析器
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/5
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"http-api/app/http/graph/errors"
	"http-api/app/models/codeinfo"
)

func (*QueryResolver)GetMaterialManufacturers(ctx context.Context) ([]*codeinfo.CodeInfo, error) {
	c := codeinfo.CodeInfo{}

	cs, err := c.GetMaterialManufacturers(ctx)
	if err != nil {
		return cs, errors.ServerErr(ctx, err)
	}

	return cs, nil
}