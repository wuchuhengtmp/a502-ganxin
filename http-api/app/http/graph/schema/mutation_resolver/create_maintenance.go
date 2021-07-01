/**
 * @Desc    创建维修厂解析器
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
	"http-api/pkg/model"
)

func (*MutationResolver) CreateMaintenance(ctx context.Context, input graphModel.CreateMaintenanceInput) (*maintenance.Maintenance, error) {
	var res maintenance.Maintenance
	if err := requests.ValidateCreateMaintenanceRequest(input); err != nil {
		return &res, errors.ValidateErr(ctx, err)
	}
	me := auth.GetUser(ctx)
	m := maintenance.Maintenance{
		Name:      input.Name,
		Address:   input.Address,
		IsAble:    true,
		CompanyId: me.CompanyId,
	}
	if input.Remark != nil {
		m.Remark = *input.Remark
	}
	err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 创建维修厂
		if err := tx.Create(&m).Error; err != nil {
			return err
		}
		tx.Create(&maintenance_leader.MaintenanceLeader{
			MaintenanceId: m.Id,
			Uid:           input.UID,
		})
		// 操作日志
		l := logs.Logos{
			Type:    logs.CreateActionType,
			Content: fmt.Sprintf("创建维修厂,id为:%d", m.Id),
			Uid:     me.Id,
		}
		if err := tx.Create(&l).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	return &m, nil
}
