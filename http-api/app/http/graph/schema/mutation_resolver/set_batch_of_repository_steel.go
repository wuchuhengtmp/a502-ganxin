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
	"gorm.io/gorm"
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/errors"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/steel_logs"
	"http-api/app/models/steels"
	"http-api/pkg/model"
)

func (*MutationResolver) SetBatchOfRepositorySteel(ctx context.Context, input graphModel.SetBatchOfRepositorySteelInput) (res []*steels.Steels, err error) {
	if err := requests.ValidateSetBatchOfRepositorySteelRequest(ctx, input); err != nil {
		return res, errors.ValidateErr(ctx, err)
	}
	me := auth.GetUser(ctx)
	err = model.DB.Transaction(func(tx *gorm.DB) error {
		for _, identifier := range input.IdentifierList {
			// 修改
			err = model.DB.Model(&steels.Steels{}).
				Where("identifier = ?", identifier).
				Update("manufacturer_id", input.ManufacturerID).
				Update("material_manufacturer_id", input.MaterialManufacturersID).
				Update("produced_date", input.ProducedAt).
				Update("specification_id", input.SpecificationID).
				Error
			if err != nil {
				return err
			}
			// 型钢日志
			s := steels.Steels{}
			err = tx.Model(&s).Where("identifier = ?", identifier).
				Where("company_id = ?", me.CompanyId).
				First(&s).Error
			if err != nil {
				return err
			}
			err = tx.Create(&steel_logs.SteelLog{
				Type:    steel_logs.EditType,
				SteelId: s.ID,
				Uid:     me.Id,
			}).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return res, errors.ServerErr(ctx, err)
	}

	err = model.DB.Model(&steels.Steels{}).
		Where("identifier in ?", input.IdentifierList).
		Find(&res).Error
	if err != nil {
		return res, errors.ServerErr(ctx, err)
	}

	return
}
