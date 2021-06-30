/**
 * @Desc    获取规格解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/5
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"http-api/app/http/graph/auth"
	"http-api/app/models/specificationinfo"
	"http-api/pkg/model"
)

type SpecificationItemResolver struct { }

func (*QueryResolver)GetSpecification(ctx context.Context) ([]*specificationinfo.SpecificationInfo, error) {
	var ss  []*specificationinfo.SpecificationInfo
	me := auth.GetUser(ctx)
	model.DB.Model(&specificationinfo.SpecificationInfo{}).
		Where("company_id = ?", me.CompanyId).
		Find(&ss)

	return ss, nil
}

func (SpecificationItemResolver)Specification(ctx context.Context, obj *specificationinfo.SpecificationInfo) (string, error) {
	return obj.GetSelfSpecification(), nil
}
