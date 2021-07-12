/**
 * @Desc    型钢入库
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/12
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
	"http-api/app/models/maintenance_record"
	"http-api/app/models/steel_logs"
	"http-api/app/models/steels"
	"http-api/pkg/model"
)

func (*MutationResolver) EnterMaintenanceSteelToRepository(ctx context.Context, input graphModel.EnterMaintenanceSteelToRepositoryInput) (bool, error) {
	if err := requests.ValidateEnterMaintenanceSteelToRepositoryRequest(ctx, input); err != nil {
		return false, errors.ValidateErr(ctx, err)
	}
	me := auth.GetUser(ctx)
	err := model.DB.Transaction(func(tx *gorm.DB) error {
		for _, identifier := range input.IdentifierList {
			// 标记型钢为入库
			steelItem := steels.Steels{}
			err := tx.Model(&steelItem).Where("identifier = ?", identifier).
				Where("company_id = ?", me.CompanyId).
				First(&steelItem).
				Error
			if err != nil {
				return err
			}
			err = tx.Model(&steelItem).
				Where("id = ?", steelItem.ID).
				Update("state", steels.StateInStore).
				Error
			if err != nil {
				return err
			}
			// 标记维修记录为已归库了
			recordItem := maintenance_record.MaintenanceRecord{}
			err = tx.Model(&recordItem).
				Where("id = ?", steelItem.MaintenanceRecordSteelId).
				Update("state", steels.StateInStore).
				Error
			if err != nil {
				return err
			}
			// 型钢日志
			log := steel_logs.SteelLog{
				Type:    steel_logs.EnterRepositoryFromMaintenance,
				Uid:     me.Id,
				SteelId: steelItem.ID,
			}
			 if err := tx.Create(&log).Error; err != nil {
			 	return err
			 }
		}
		return nil
	})

	if err != nil {
		return false, errors.ServerErr(ctx, err)
	}

	return true, nil
}
