/**
 * @Desc    编辑材料商解析器
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/5
 * @Listen  MIT
 */
package mutation_resolver

import (
	"context"
	"http-api/app/http/graph/errors"
	"http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/codeinfo"
)

func (*MutationResolver) EditMaterialManufacturer(ctx context.Context, input model.EditMaterialManufacturerInput) (*codeinfo.CodeInfo, error) {
	var c codeinfo.CodeInfo
	if err := requests.ValidateEditMaterialManufacturerRequest(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	c.Name = input.Name
	c.IsDefault = input.IsDefault
	c.Remark = input.Remark
	c.ID = input.ID
	if err := c.EditMaterialManufacturer(ctx); err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	return &c, nil
}
