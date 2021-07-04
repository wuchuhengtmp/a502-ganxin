/**
 * @Desc    The requests is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/22
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"http-api/app/http/graph/auth"
	"http-api/app/models/maintenance"
	"http-api/app/models/maintenance_record"
	"http-api/app/models/steels"
	"http-api/pkg/model"
)



func (*QueryResolver) GetMaintenanceList(ctx context.Context) (res []*maintenance.Maintenance, err error) {
	me := auth.GetUser(ctx)
	err = model.DB.Model(&maintenance.Maintenance{}).
		Where("company_id = ?", me.CompanyId).
		Find(&res).
		Error

	return
}

type MaintenanceRecordItemResolver struct { }

func (MaintenanceRecordItemResolver)StateInfo(ctx context.Context, obj *maintenance_record.MaintenanceRecord) (*steels.StateItem, error) {
	return &steels.StateItem{
		State: obj.State,
		Desc: steels.StateCodeMapDes[obj.State],
	}, nil
}

func (MaintenanceRecordItemResolver) Maintenance(ctx context.Context, obj *maintenance_record.MaintenanceRecord) (*maintenance.Maintenance, error) {
	me := auth.GetUser(ctx)
	m := maintenance.Maintenance{}
	err := model.DB.Model(&m).Where("company_id = ? AND id = ?", me.CompanyId, obj.MaintenanceId).First(&m).Error

	return &m, err
}

func (MaintenanceRecordItemResolver)Steel(ctx context.Context, obj *maintenance_record.MaintenanceRecord) (*steels.Steels, error) {
	me := auth.GetUser(ctx)
	s := steels.Steels{}
	err := model.DB.Model(&s).Where("company_id = ? AND id = ?",me.CompanyId, obj.SteelId).First(&s).Error

	return &s, err
}
func (MaintenanceRecordItemResolver)UseDays(ctx context.Context, obj *maintenance_record.MaintenanceRecord) (*int64, error) {
	// todo 维修天数
	var useDays int64

	return &useDays, nil
}
