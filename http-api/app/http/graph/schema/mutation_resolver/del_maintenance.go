/**
 * @Desc    The mutation_resolver is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/1
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
	"http-api/app/models/logs"
	"http-api/app/models/maintenance"
	"http-api/app/models/maintenance_leader"
	"http-api/app/models/msg"
	"http-api/pkg/model"
)

func (*MutationResolver) DelMaintenance(ctx context.Context, input graphModel.DelMaintenanceInput) (bool, error) {
	if err := requests.ValidateDelMaintenanceRequest(ctx, input); err != nil {
		return false, errors.ValidateErr(ctx, err)
	}
	err := model.DB.Transaction(func(tx *gorm.DB) error {
		m := maintenance.Maintenance{}
		if err := tx.Model(&m).Where("id = ?", input.ID).First(&m).Error; err != nil {
			return err
		}
		// 删除
		if err := tx.Where("id = ?", input.ID).Delete(&maintenance.Maintenance{}).Error; err != nil {
			return err
		}
		// 通知管理理员，你下岗了
		var leaders []*maintenance_leader.MaintenanceLeader
		if err := tx.Model(&leaders).Where("maintenance_id = ?", input.ID).Find(&leaders).Error; err != nil {
			return err
		}
		me := auth.GetUser(ctx)
		for _, l := range leaders {
			msgInstance := msg.Msg{
				IsRead:  false,
				Content: fmt.Sprintf("公司管理员:%s, 删除%s维修厂同时也删除了你在该厂的管理权，有什么不明白的，找他去问去", me.Name, m.Name),
				Uid:     l.Uid,
				Type:    msg.DelMaintenance,
			}
			if err := msgInstance.CreateSelf(tx); err != nil {
				return err
			}
		}
		// 删除对应管理员
		err := tx.Where("maintenance_id = ?", input.ID).Delete(&maintenance_leader.MaintenanceLeader{}).Error
		if err != nil {
			return err
		}
		// 操作日志
		l := logs.Logos{
			Type:    logs.CreateActionType,
			Content: fmt.Sprintf("删除维修厂"),
			Uid:     me.Id,
		}
		if err := tx.Create(&l).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return false, errors.ServerErr(ctx, err)
	}

	return true, nil
}
