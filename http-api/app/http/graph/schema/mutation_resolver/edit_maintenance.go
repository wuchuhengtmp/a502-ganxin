/**
 * @Desc    编辑维修厂解析器
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
	"http-api/app/http/graph/util/helper"
	"http-api/app/models/logs"
	"http-api/app/models/maintenance"
	"http-api/app/models/maintenance_leader"
	"http-api/pkg/model"
)

func (*MutationResolver) EditMaintenance(ctx context.Context, input graphModel.EditMaintenanceInput) (*maintenance.Maintenance, error) {
	if err := requests.ValidateEditMaintenanceRequest(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	res := maintenance.Maintenance{}
	me := auth.GetUser(ctx)
	err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 修改维修厂
		err := tx.Model(&res).
			Where("id = ?", input.ID).
			Update("name", input.Name).
			Update("address", input.Address).
			Update("is_able", input.IsAble).
			Error
		if err != nil {
			return err
		}
		if input.Remark != nil {
			if err := tx.Model(&res).Where("id = ?", input.ID).Update("remark", *input.Remark).Error; err != nil {
				return err
			}
		}
		// 修改管理员
		var oldLeaderList []*maintenance_leader.MaintenanceLeader // 原来的管理员
		err = tx.Model(&maintenance_leader.MaintenanceLeader{}).Where("maintenance_id = ?", input.ID).Find(&oldLeaderList).Error
		if err != nil {
			return err
		}
		var oldUidList []int64
		for _, leader := range oldLeaderList {
			oldUidList = append(oldUidList, leader.Uid)
		}
		addUidList, deleteUidList := helper.CompareCollect(input.AdminIDList, oldUidList)
		// 删除
		err = tx.Debug().Where("uid in ?", deleteUidList).
			Where("maintenance_id = ?", input.ID).
			Delete(&maintenance_leader.MaintenanceLeader{}).Error
		if err != nil {
			return err
		}
		// 新增
		for _, uid := range addUidList {
			if err := tx.Create(&maintenance_leader.MaintenanceLeader{Uid: uid, MaintenanceId: input.ID}).Error; err != nil {
				return err
			}
		}
		// 操作日志
		err = tx.Create(
			&logs.Logos{
				Type:    logs.EditActionType,
				Content: fmt.Sprintf("编辑维修厂"),
				Uid:     me.Id,
			},
		).Error
		if err != nil {
			return err
		}
		if err := tx.Model(&res).Where("id = ?", input.ID).First(&res).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	return &res, nil
}
