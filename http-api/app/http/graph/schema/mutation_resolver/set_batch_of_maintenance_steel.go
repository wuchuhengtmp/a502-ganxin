/**
 * @Desc    型钢维修出库
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/3
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
	"http-api/app/http/graph/util/helper"
	"http-api/app/models/maintenance"
	"http-api/app/models/maintenance_leader"
	"http-api/app/models/maintenance_record"
	"http-api/app/models/msg"
	"http-api/app/models/specificationinfo"
	"http-api/app/models/steel_logs"
	"http-api/app/models/steels"
	"http-api/pkg/model"
	"time"
)

// 型钢维修出库操作步骤
type SetBatchOfMaintenanceSteelSteps struct{}

func (*MutationResolver) SetBatchOfMaintenanceSteel(ctx context.Context, input graphModel.SetBatchOfMaintenanceSteelInput) (res []*steels.Steels, err error) {
	if err := requests.ValidateSetBatchOfMaintenanceSteelRequest(ctx, input); err != nil {
		return res, err
	}
	steps := SetBatchOfMaintenanceSteelSteps{}
	err = model.DB.Transaction(func(tx *gorm.DB) error {
		for _, identifier := range input.IdentifierList {
			// 标记型钢为待维修状态
			if err := steps.FlagMaintenance(ctx, tx, identifier); err != nil {
				return err
			}
			// 添加维修详情记录
			if err := steps.CreateDetail(ctx, tx, identifier, input.MaintenanceID); err != nil {
				return err
			}
			// 添加型钢日志
			if err := steps.CreateSteelLog(ctx, tx, identifier); err != nil {
				return err
			}

		}
		// 添加消息通知
		if err := steps.createMsg(ctx, tx, input); err != nil {
			return err
		}
		// 获取响应数据
		newRes, err := steps.GetRes(ctx, tx, input)
		if err != nil {
			return err
		}
		res = newRes
		return nil
	})

	return
}

func (*SetBatchOfMaintenanceSteelSteps) GetSteelItem(ctx context.Context, tx *gorm.DB, identifier string) (steelItem steels.Steels, err error) {
	me := auth.GetUser(ctx)
	err = tx.Model(&steelItem).Where("identifier = ?", identifier).
		Where("company_id = ?", me.CompanyId).
		First(&steelItem).
		Error

	return
}

/**
 * 获取响应数据
 */
func (*SetBatchOfMaintenanceSteelSteps) GetRes(ctx context.Context, tx *gorm.DB, input graphModel.SetBatchOfMaintenanceSteelInput) (res []*steels.Steels, err error) {
	steelItem := steels.Steels{}
	me := auth.GetUser(ctx)
	err = tx.Model(&steelItem).Where("identifier IN ?", input.IdentifierList).
		Where("company_id = ?", me.CompanyId).
		Find(&res).
		Error

	return
}

/**
 * 添加详情记录
 */
func (s *SetBatchOfMaintenanceSteelSteps) CreateDetail(ctx context.Context, tx *gorm.DB, identifier string, maintenanceId int64) error {
	steelItem, err := s.GetSteelItem(ctx, tx, identifier)
	if err != nil {
		return err
	}
	i := maintenance_record.MaintenanceRecord{
		State:           steels.StateRepository2Maintainer,
		MaintenanceId:   maintenanceId,
		SteelId:         steelItem.ID,
		OutRepositoryAt: time.Now(),
	}
	if err := tx.Create(&i).Error; err != nil {
		return err
	}
	// 标记当前维修的型钢详情id 对应到 型钢表中
	me := auth.GetUser(ctx)
	steelItem = steels.Steels{}
	err = tx.Model(&steelItem).Where("identifier = ?", identifier).
		Where("company_id = ?", me.CompanyId).
		Update("maintenance_record_steel_id", i.Id).
		Error
	if err != nil {
		return  err
	}

	return nil
}

/**
 * 添加型钢日志
 */
func (s *SetBatchOfMaintenanceSteelSteps) CreateSteelLog(ctx context.Context, tx *gorm.DB, identifier string) error {
	me := auth.GetUser(ctx)
	steelItem, err := s.GetSteelItem(ctx, tx, identifier)
	if err != nil {
		return err
	}
	logItem := steel_logs.SteelLog{
		Type:    steel_logs.OutOfRepositoryForMaintenance,
		SteelId: steelItem.ID,
		Uid:     me.Id,
	}
	if err := tx.Create(&logItem).Error; err != nil {
		return err
	}

	return nil
}

/**
 *  型钢标记维修
 */
func (s *SetBatchOfMaintenanceSteelSteps) FlagMaintenance(ctx context.Context, tx *gorm.DB, identifier string) error {
	steelItem, err := s.GetSteelItem(ctx, tx, identifier)
	if err != nil {
		return err
	}
	// 标记为维修
	err = tx.Model(&steels.Steels{}).Where("id = ?", steelItem.ID).
		Update("state", steels.StateRepository2Maintainer).
		Error
	return nil
}

/**
 * 发送消息
 */
func (*SetBatchOfMaintenanceSteelSteps) createMsg(ctx context.Context, tx *gorm.DB, input graphModel.SetBatchOfMaintenanceSteelInput) error {
	me := auth.GetUser(ctx)
	specificationInfoTable := specificationinfo.SpecificationInfo{}.TableName()
	steelTable := steels.Steels{}.TableName()
	var weightInfo struct {
		Weight float64
	}
	err := tx.Model(&steels.Steels{}).
		Select(fmt.Sprintf("sum(weight) as Weight")).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.specification_id", specificationInfoTable, specificationInfoTable, steelTable)).
		Where(fmt.Sprintf("%s.identifier IN ?", steelTable), input.IdentifierList).
		Where(fmt.Sprintf("%s.company_id = ?", steelTable), me.CompanyId).
		Scan(&weightInfo).Error
	if err != nil {
		return err
	}
	var leaders []*maintenance_leader.MaintenanceLeader
	maintenanceItem := maintenance.Maintenance{}
	leaderTable := maintenance_leader.MaintenanceLeader{}.TableName()
	err = tx.Model(maintenanceItem).
		Select(fmt.Sprintf("%s.*", leaderTable)).
		Joins(fmt.Sprintf("join %s ON %s.maintenance_id = %s.id", leaderTable, leaderTable, maintenanceItem.TableName())).
		Where(fmt.Sprintf("%s.id = ?", maintenanceItem.TableName()), input.MaintenanceID).
		Scan(&leaders).
		Error
	if err != nil {
		return err
	}
	for _, leaderItem := range leaders {
		content := fmt.Sprintf(
			"仓库管理员:%s 于%s 发一批型钢需要维修，总数:%d根 %.2f吨, 请注意查收.",
			me.Name,
			helper.Time2Str(time.Now()),
			len(input.IdentifierList),
			weightInfo.Weight,
		)
		m := msg.Msg{
			Content: content,
			Type:    msg.ToBeMaintained,
			Uid:     leaderItem.Uid,
		}
		if err := m.CreateSelf(tx); err != nil {
			return err
		}
	}

	return nil
}
