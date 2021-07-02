/**
 * @Desc    批量修改仓库型钢
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/2
 * @Listen  MIT
 */
package mutation_resolver

import (
	"context"
	"http-api/app/http/graph/errors"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/steels"
	"http-api/pkg/model"
)

func (*MutationResolver) SetBatchOfRepositorySteel(ctx context.Context, input graphModel.SetBatchOfRepositorySteelInput) (res []*steels.Steels, err error) {
	if err := requests.ValidateSetBatchOfRepositorySteelRequest(ctx, input); err != nil {
		return res, errors.ValidateErr(ctx, err)
	}
	err = model.DB.Model(&steels.Steels{}).
		Where("identifier IN ?", input.IdentiferList).
		Update("manufacturer_id", input.ManufacturerID).
		Update("material_manufacturer_id", input.MaterialManufacturersID).
		Update("produced_date", input.ProducedAt).
		Update("specification_id", input.SpecificationID).
		Error
	if err != nil {
		return res, errors.ServerErr(ctx, err)
	}
	err = model.DB.Model(&steels.Steels{}).
		Where("identifier in ?", input.IdentiferList).
		Find(&res).Error
	if err != nil {
		return res, errors.ServerErr(ctx, err)
	}

	return
}
