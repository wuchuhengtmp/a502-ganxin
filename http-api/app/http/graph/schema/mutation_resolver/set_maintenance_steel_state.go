/**
 * @Desc    修改维修型钢状态解析器
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
	graphModel "http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/maintenance_record"
	"http-api/app/models/steel_logs"
	"http-api/app/models/steels"
	"http-api/pkg/model"
)

func (*MutationResolver) SetMaintenanceSteelState(ctx context.Context, input graphModel.SetMaintenanceSteelStateInput) (res []*maintenance_record.MaintenanceRecord, err error) {
	if err := requests.ValidateSetMaintenanceSteelStateRequest(ctx, input); err != nil {
		return res, err
	}
	steelTable := steels.Steels{}.TableName()
	me := auth.GetUser(ctx)

	err = model.DB.Transaction(func(tx *gorm.DB) error {
		for _, identifier := range input.IdentifierList {
			recordItem := maintenance_record.MaintenanceRecord{}
			// 修改状态
			err := tx.Model(&recordItem).
				Joins(fmt.Sprintf("join %s ON %s.maintenance_record_steel_id = %s.id", steelTable, steelTable, recordItem.TableName())).
				Where(fmt.Sprintf("%s.identifier = ?", steelTable), identifier).
				Where(fmt.Sprintf("%s.company_id = ?", steelTable), me.CompanyId).
				First(&recordItem).
				Error
			if err != nil {
				return err
			}
			err = tx.Model(&recordItem).Where("id = ?", recordItem.Id).
				Update("state", input.State).
				Error
			if err != nil {
				return err
			}
			res = append(res, &recordItem)
			// 型钢操作日志
			steelItem := steels.Steels{}
			err = tx.Model(&steelItem).Where("identifier = ?", identifier).
				Where("company_id = ?", me.CompanyId).
				First(&steelItem).
				Error
			if err != nil {
				return err
			}
			logItem := steel_logs.SteelLog{
				Uid: me.Id,
				SteelId: steelItem.ID,
				Type:  steel_logs.ChangedMaintenanceSteel,
			}
			if err := tx.Create(&logItem).Error; err != nil {
				return err
			}
		}

		return nil
	})

	return
}
