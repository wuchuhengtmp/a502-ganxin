/**
 * @Desc    获取规格解析器
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
	"http-api/pkg/model"
)

type SpecificationItemResolver struct { }

func (*QueryResolver)GetSpecification(ctx context.Context) ([]*specificationinfo.SpecificationInfo, error) {
	var ss  []*specificationinfo.SpecificationInfo
	model.DB.Model(&specificationinfo.SpecificationInfo{}).Find(&ss)

	return ss, nil
}

func (SpecificationItemResolver)Specification(ctx context.Context, obj *specificationinfo.SpecificationInfo) (string, error) {
	return fmt.Sprintf("%sx%.2fx%.2f", obj.Type, obj.Length, obj.Weight), nil
}
