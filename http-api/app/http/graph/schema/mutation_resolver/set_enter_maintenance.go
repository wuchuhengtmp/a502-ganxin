/**
 * @Desc    型钢入库解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/5
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
	"http-api/app/models/maintenance"
	"http-api/app/models/maintenance_leader"
	"http-api/app/models/maintenance_record"
	"http-api/app/models/steel_logs"
	"http-api/app/models/steels"
	"http-api/pkg/model"
	"time"
)

func (*MutationResolver) SetEnterMaintenance(ctx context.Context, input graphModel.SetMaintenanceInput) (res []*maintenance_record.MaintenanceRecord, err error) {
	if err := requests.ValidateSetMaintenanceRequest(ctx, input); err != nil {
		return res, errors.ValidateErr(ctx, err)
	}
	err = model.DB.Transaction(func(tx *gorm.DB) error {
		steps := SetMaintenanceSteps{}
		for _, identifier := range input.IdentifierList {
			// 标记型钢状态
			if err := steps.FlagSteel(ctx, identifier, tx); err != nil {
				return err
			}
			// 标记维修厂型钢
			newRecord, err := steps.FlagMaintenanceRecord(ctx, identifier, tx);
			if  err != nil {
				return err
			}
			// 型钢日志
			if err := steps.CreateSteelLog(ctx, identifier, tx); err != nil {
				return err
			}
			res = append(res, newRecord)
		}
		return nil
	})
	if err != nil {
		return res, errors.ServerErr(ctx, err)
	}

	return

}

type SetMaintenanceSteps struct{}

/**
 *  获取详情
 */
func (*SetMaintenanceSteps) GetMaintenanceRecord(ctx context.Context, identifier string, tx *gorm.DB) (maintenance_record.MaintenanceRecord, error) {
	me := auth.GetUser(ctx)
	maintenanceTable := maintenance.Maintenance{}.TableName()
	maintenanceRecordTable := maintenance_record.MaintenanceRecord{}.TableName()
	maintenanceLeaderTable := maintenance_leader.MaintenanceLeader{}.TableName()
	steelTable := steels.Steels{}.TableName()
	maintenanceRecordItem := maintenance_record.MaintenanceRecord{}
	err := tx.Debug().Model(&maintenanceRecordItem).
		Select(fmt.Sprintf("%s.*", maintenanceRecordTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.steel_id", steelTable, steelTable, maintenanceRecordTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.maintenance_id", maintenanceTable, maintenanceTable, maintenanceRecordTable)).
		Joins(fmt.Sprintf("join %s ON %s.maintenance_id = %s.id", maintenanceLeaderTable, maintenanceLeaderTable, maintenanceTable)).
		Where(fmt.Sprintf("%s.identifier = ?", steelTable), identifier).
		Where(fmt.Sprintf("%s.uid = ? ", maintenanceLeaderTable), me.Id).
		Where(fmt.Sprintf("%s.state = ?", maintenanceRecordTable), steels.StateRepository2Maintainer).
		First(&maintenanceRecordItem).
		Error

	return maintenanceRecordItem, err
}

/**
 * 标记型钢状态
 */
func (*SetMaintenanceSteps) FlagSteel(ctx context.Context, identifier string, tx *gorm.DB) error {
	me := auth.GetUser(ctx)
	steelItem := steels.Steels{}
	err := tx.Model(&steelItem).Where("identifier = ?", identifier).
		Where("company_id = ?", me.CompanyId).
		Update("state", steels.StateMaintainerWillBeMaintained).
		Error

	return err
}

/**
 * 更新维修厂的型钢详情
 */
func (s *SetMaintenanceSteps) FlagMaintenanceRecord(ctx context.Context, identifier string, tx *gorm.DB) (*maintenance_record.MaintenanceRecord, error) {
	me := auth.GetUser(ctx)
	maintenanceRecordItem, err := s.GetMaintenanceRecord(ctx, identifier, tx)
	if err != nil {
		return nil, err
	}
	err = tx.Model(&maintenanceRecordItem).Where("id = ?", maintenanceRecordItem.Id).
		Update("state", steels.StateMaintainerWillBeMaintained).
		Update("entered_at", time.Now()).
		Update("entered_uid", me.Id).
		Error

	return &maintenanceRecordItem ,err
}

/**
 * 型钢日志
 */
func (s *SetMaintenanceSteps) CreateSteelLog(ctx context.Context, identifier string, tx *gorm.DB) error {
	me := auth.GetUser(ctx)
	steelItem := steels.Steels{}
	err := tx.Model(&steelItem).Where("company_id = ?", me.CompanyId).
		Where("identifier = ?", identifier).
		First(&steelItem).
		Error

	if err != nil {
		return err
	}
	steelLogItem := steel_logs.SteelLog{
		Type:    steel_logs.ToBeMaintenanceType,
		SteelId: steelItem.ID,
		Uid:     me.Id,
	}

	return tx.Create(&steelLogItem).Error
}
