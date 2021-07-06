/**
 * @Desc    型钢出厂
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/6
 * @Listen  MIT
 */
package mutation_resolver

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/errors"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/maintenance_record"
	"http-api/app/models/steel_logs"
	"http-api/app/models/steels"
	"http-api/pkg/model"
	"time"
)

type SetSteelForOutOfMaintenanceSteps struct { }

func (*MutationResolver) SetSteelForOutOfMaintenance(ctx context.Context, input graphModel.SetSteelForOutOfMaintenanceInput) ([]*maintenance_record.MaintenanceRecord, error) {
	var res  []*maintenance_record.MaintenanceRecord
	if err := requests.ValidateSetSteelForOutOfMaintenanceRequest(ctx, input); err != nil {
		return res, errors.ValidateErr(ctx, err)
	}
	err := model.DB.Transaction(func(tx *gorm.DB) error {
		steps := SetSteelForOutOfMaintenanceSteps{}
		for _, identifier := range input.IdentifierList {
			// 标记型钢状态
			if err := steps.FlagSteel(ctx, identifier, tx); err != nil {
				return err
			}
			// 标记维修型钢详情
			r, err := steps.FlagMaintenanceSteel(ctx, identifier, tx)
			if err != nil {
				return err
			}
			res = append(res, r)
			// 型钢日志
			if err := steps.CreateLog(ctx, identifier, tx); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return res, errors.ServerErr(ctx, err)
	}

	return res, nil
}



/**
* 标记型钢状态（标记为归库途中）
*/
func (*SetSteelForOutOfMaintenanceSteps) FlagSteel(ctx context.Context, identifier string, tx *gorm.DB) error {
	steelsItem := steels.Steels{}
	me := auth.GetUser(ctx)
	err := tx.Model(&steelsItem).Where("identifier = ?", identifier).
		Where("company_id = ?", me.CompanyId).
		Update("state", steels.StateMaintainerOnTheStoreWay).
		Error

	return err
}

/**
* 标记维修型钢状态（标记为归库途中）
*/
func (*SetSteelForOutOfMaintenanceSteps) FlagMaintenanceSteel(ctx context.Context, identifier string, tx *gorm.DB) (*maintenance_record.MaintenanceRecord, error) {
	steelsTable := steels.Steels{}.TableName()
	me := auth.GetUser(ctx)
	record := maintenance_record.MaintenanceRecord{}
	err := tx.Model(&record).
		Select(fmt.Sprintf("%s.*", record.TableName())).
		Joins(fmt.Sprintf("join %s ON %s.maintenance_record_steel_id = %s.id", steelsTable, steelsTable, record.TableName())).
		Where(fmt.Sprintf("%s.identifier = ?", steelsTable), identifier).
		Where(fmt.Sprintf("%s.company_id = ?", steelsTable), me.CompanyId).
		First(&record).
		Error
	if err != nil {
		return nil, err
	}
	err = tx.Model(&record).Where("id = ?", record.Id).
		Update("state", steels.StateMaintainerOnTheStoreWay).
		Update("outed_uid", me.Id).
		Update("outed_at", time.Now()).
		Error
	if err != nil {
		return nil, err
	}

	return &record, nil
}

/**
* 添加日志
*/
func (*SetSteelForOutOfMaintenanceSteps) CreateLog(ctx context.Context, identifier string, tx *gorm.DB) error {
	me := auth.GetUser(ctx)
	steelItem := steels.Steels{}
	err := tx.Model(&steelItem).
		Where("identifier = ?", identifier).
		Where("company_id = ?", me.CompanyId).
		First(&steelItem).
		Error
	if err != nil {
		return err
	}
	l := steel_logs.SteelLog{
		Type:    steel_logs.OutOfMaintenance,
		SteelId: steelItem.ID,
		Uid:     me.Id,
	}
	if err := tx.Create(&l).Error; err != nil {
		return err
	}

	return nil
}
