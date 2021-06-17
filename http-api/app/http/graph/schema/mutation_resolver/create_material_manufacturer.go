/**
 * @Desc    添加材料商解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/5
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

func  (*MutationResolver)CreateMaterialManufacturer(ctx context.Context, input model.CreateMaterialManufacturerInput) (*codeinfo.CodeInfo, error) {
	if err := requests.ValidateCreateMaterialManufacturerRequest(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	me := auth.GetUser(ctx)
	c := codeinfo.CodeInfo{
		Name: input.Name,
		Type: codeinfo.MaterialManufacturer,
		IsDefault: input.IsDefault,
		Remark: input.Remark,
		CompanyId: me.CompanyId,
	}
	if err := c.CreateMaterialManufacturer(ctx); err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	return &c, nil
}
