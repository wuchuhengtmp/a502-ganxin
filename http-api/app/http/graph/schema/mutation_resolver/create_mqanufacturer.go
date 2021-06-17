/**
 * @Desc    创建制造商解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/7
 * @Listen  MIT
 */
package mutation_resolver

import (
	"context"
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/errors"
	"http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/codeinfo"
)

func (*MutationResolver) CreateManufacturer(ctx context.Context, input model.CreateManufacturerInput) (*codeinfo.CodeInfo, error) {
	if err := requests.ValidateCreateManufacturerRequest(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	me := auth.GetUser(ctx)
	c := codeinfo.CodeInfo{
		Name:      input.Name,
		IsDefault: input.IsDefault,
		Remark:    input.Remark,
		Type:      codeinfo.Manufacturer,
		CompanyId: me.CompanyId,
	}
	if err := c.CreateManufacturerSelf(ctx); err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	return &c, nil
}
