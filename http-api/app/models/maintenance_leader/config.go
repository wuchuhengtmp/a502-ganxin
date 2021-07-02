/**
 * @Desc    维修管理员模型
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/15
 * @Listen  MIT
 */
package maintenance_leader

type MaintenanceLeader struct {
	ID            int64 `json:"id" sql:"unique_index"`
	MaintenanceId int64 `json:"maintenanceId" gorm:"comment:维修厂id"`
	Uid           int64 `json:"uid" gorm:"comment:用户id"`
}
func(MaintenanceLeader)TableName() string {
	return "maintenance_leaders"
}
