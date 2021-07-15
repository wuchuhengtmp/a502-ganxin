/**
 * @Desc    维修厂扩展解析
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/1
 * @Listen  MIT
 */
package extends_resolver

import (
	"context"
	"fmt"
	"http-api/app/models/maintenance"
	"http-api/app/models/maintenance_leader"
	"http-api/app/models/maintenance_record"
	"http-api/app/models/specificationinfo"
	"http-api/app/models/steels"
	"http-api/app/models/users"
	"http-api/pkg/model"
)

type MaintenanceItemResolver struct{}

func (MaintenanceItemResolver) Admin(ctx context.Context, obj *maintenance.Maintenance) (res []*users.Users, err error) {
	maintenanceLeaderTable := maintenance_leader.MaintenanceLeader{}.TableName()
	userTable := users.Users{}.TableName()
	err = model.DB.Model(&users.Users{}).
		Select(fmt.Sprintf("%s.*", userTable)).
		Joins(fmt.Sprintf("join %s ON %s.uid = %s.id", maintenanceLeaderTable, maintenanceLeaderTable, userTable)).
		Where(fmt.Sprintf("%s.maintenance_id = ?", maintenanceLeaderTable), obj.Id).
		Find(&res).
		Error
	return

}

// 维修厂数量
func (MaintenanceItemResolver) Total(ctx context.Context, obj *maintenance.Maintenance) (int64, error) {
	recordItem := maintenance_record.MaintenanceRecord{}
	var total int64
	err := model.DB.Model(&recordItem).Where("maintenance_id = ?", obj.Id).Count(&total).Error
	return total, err
}

//  维修厂重量
func (MaintenanceItemResolver) WeightTotal(ctx context.Context, obj *maintenance.Maintenance) (float64, error) {
	recordItem := maintenance_record.MaintenanceRecord{}
	var weightInfo struct {
		Weight float64
	}
	steelsTable := steels.Steels{}.TableName()
	specificationTable := specificationinfo.SpecificationInfo{}.TableName()
	recordTable := recordItem.TableName()
	err := model.DB.Model(&recordItem).
		Select("sum(weight) as Weight").
		Joins(fmt.Sprintf("join %s on %s.id = %s.steel_id", steelsTable, steelsTable, recordTable)).
		Joins(fmt.Sprintf("join %s on %s.id = %s.specification_id", specificationTable, specificationTable, steelsTable)).
		Where(fmt.Sprintf("%s.maintenance_id = ?", recordTable), obj.Id).
		Scan(&weightInfo).
		Error

	return weightInfo.Weight, err
}
