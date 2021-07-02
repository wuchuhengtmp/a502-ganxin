/**
 * @Desc    批量修改仓库型钢请求验证器
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

func ValidateSetBatchOfRepositorySteelRequest(ctx context.Context, input graphModel.SetBatchOfRepositorySteelInput) error {
	steps := StepsForRepository{}
	for _, identifier := range input.IdentifierList {
		// 检验有没有这根型钢
		if err := steps.CheckHasSteel(ctx, identifier); err != nil {
			return err
		}
	}
	// 检验有没有这个规格
	if err := steps.CheckHasSpecification(ctx, input.SpecificationID); err != nil {
		return err
	}
	// 检验有没有这个材料商
	if err := steps.CheckHasMaterialManufacturer(ctx, input.MaterialManufacturersID); err != nil {
		return err
	}
	// 检验有没有这个制造商
	if err := steps.CheckHasManufacturer(ctx, input.ManufacturerID); err != nil {
		return err
	}

	return nil
}