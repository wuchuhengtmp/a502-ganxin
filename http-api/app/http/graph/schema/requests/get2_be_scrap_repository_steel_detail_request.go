/**
 * @Desc    获取待报废型钢详情请求验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/2
 * @Listen  MIT
 */
package requests

import (
	"context"
	graphModel "http-api/app/http/graph/model"
)

func ValidateGet2BeScrapRepositorySteelDetailRequest(ctx context.Context, input graphModel.Get2BeScrapRepositorySteelDetailInput) error {
	steps := StepsForRepository{}
	for _, identifier := range input.IdentifierList {
		// 检验有没有这个型钢
		if err := steps.CheckHasSteel(ctx, identifier); err != nil {
			return err
		}
		// 检验型钢能否报废
		if err := steps.CheckIsScrapAccess(ctx, identifier); err != nil {
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