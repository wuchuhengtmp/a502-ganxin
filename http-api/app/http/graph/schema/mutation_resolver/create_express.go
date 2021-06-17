/**
 * @Desc    创建物流商解析器
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

func (*MutationResolver) CreateExpress(ctx context.Context, input model.CreateExpressInput) (*codeinfo.CodeInfo, error) {
	if err := requests.ValidateCreateExpressRequest(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	c := codeinfo.CodeInfo{
		Name:      input.Name,
		Remark:    input.Remark,
		IsDefault: input.IsDefault,
	}
	if err := c.CreateExpress(ctx); err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	return &c, nil
}
