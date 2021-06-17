/**
 * @Desc    获取物流公司列表解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/8
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"http-api/app/http/graph/errors"
	"http-api/app/models/codeinfo"
)

func (*QueryResolver)GetExpressList(ctx context.Context) ([]*codeinfo.CodeInfo, error) {
	c := codeinfo.CodeInfo{}
	var cs []*codeinfo.CodeInfo
	cs, err := c.GetExpressList(ctx)
	if  err != nil {
		return cs, errors.ServerErr(ctx, err)
	}

	return cs, nil
}
