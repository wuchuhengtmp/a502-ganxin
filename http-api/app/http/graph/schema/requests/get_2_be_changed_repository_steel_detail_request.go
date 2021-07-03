/**
 * @Desc    The requests is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/3
 * @Listen  MIT
 */
package requests

import (
	"context"
	graphModel "http-api/app/http/graph/model"
)

func Get2BeChangedRepositorySteelDetail(ctx context.Context, input graphModel.Get2BeChangedRepositorySteelDetailInput) error {
	steps := StepsForRepository{}
	for _, identifier := range input.IdentifierList {
		// 检验有没有这根型钢
		if err := steps.CheckHasSteel(ctx, identifier); err != nil {
			return err
		}
		// 检验是不是归属我
		if err := steps.CheckIsSteelBeLongMe(ctx, identifier); err != nil {
			return err
		}
		// 检验能否修改
		if err := steps.CheckIsChangeAccess(ctx, identifier); err != nil {
			return err
		}
	}
	// 检验规格
	if input.SpecificationID != nil {
		if err := steps.CheckHasSpecification(ctx, *input.SpecificationID); err != nil {
			return err
		}
	}

	return nil
}
