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
