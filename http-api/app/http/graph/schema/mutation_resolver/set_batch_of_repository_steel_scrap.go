/**
 * @Desc    仓库型钢批量报废
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

func (*MutationResolver) SetBatchOfRepositorySteelScrap(ctx context.Context, input graphModel.SetBatchOfRepositorySteelScrapInput) (res []*steels.Steels, err error) {
	if err := requests.ValidateSetBatchOfRepositorySteelScrapRequest(ctx, input); err != nil {
		return res, errors.ValidateErr(ctx, err)
	}
	me := auth.GetUser(ctx)
	err = model.DB.Transaction(func(tx *gorm.DB) error {
		for _, identifier := range input.IdentifierList {
			// 标记报废
			err := tx.Model(&steels.Steels{}).Where("identifier = ?", identifier).
				Where("company_id = ?", me.CompanyId).
				Update("state", steels.StateScrap).
				Error
			if err != nil {
				return err
			}
			i := steels.Steels{}
			err = tx.Model(&i).Where("identifier = ?", identifier).
				Where("company_id = ?", me.CompanyId).
				First(&i).Error
			if err != nil {
				return err
			}
			// 型钢日志
			s := steel_logs.SteelLog{
				Type:    steel_logs.CrapSteelType,
				SteelId: i.ID,
				Uid:     me.Id,
			}
			if err := tx.Create(&s).Error; err != nil {
				return err
			}

		}

		return nil
	})
	if err != nil {
		return res, errors.ServerErr(ctx, err)
	}
	err = model.DB.Model(&steels.Steels{}).Where("identifier IN ?", input.IdentifierList).
		Where("company_id = ?", me.CompanyId).
		Find(&res).Error
	if err != nil {
		return res, errors.ServerErr(ctx, err)
	}

	return
}
