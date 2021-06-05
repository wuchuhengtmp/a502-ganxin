/**
 * @Desc    The query_resolver is part of http-api
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/5
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"fmt"
	"http-api/app/models/specificationinfo"
)

type SpecificationItemResolver struct { }

func (SpecificationItemResolver)Specification(ctx context.Context, obj *specificationinfo.SpecificationInfo) (string, error) {
	return fmt.Sprintf("%sx%.2fx%.2f", obj.Type, obj.Length, obj.Weight), nil
}
