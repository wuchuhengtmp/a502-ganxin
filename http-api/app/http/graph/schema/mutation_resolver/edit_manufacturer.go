/**
 * @Desc    编辑制造商解析器
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
	"http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/codeinfo"
)

func (*MutationResolver) EditManufacturer(ctx context.Context, input model.EditManufacturerInput) (*codeinfo.CodeInfo, error) {
	if err := requests.ValidateEditManufacturerRequest(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	c := codeinfo.CodeInfo{ID: input.ID}
	_ = c.GetSelf()
	c.Name = input.Name
	c.Remark = input.Remark
	c.IsDefault = input.IsDefault
	if err := c.EditManufacturer(ctx); err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	return &c, nil
}
