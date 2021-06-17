/**
 * @Desc    获取制造商家列表解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/7
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"http-api/app/http/graph/errors"
	"http-api/app/models/codeinfo"
)

func (*QueryResolver) GetManufacturers(ctx context.Context) ([]*codeinfo.CodeInfo, error) {
	c := codeinfo.CodeInfo{}
	cs, err := c.GetManufacturers(ctx)
	if err != nil {
		return cs, errors.ServerErr(ctx, err)
	}

	return cs, err
}

