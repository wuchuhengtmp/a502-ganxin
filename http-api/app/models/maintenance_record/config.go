/**
 * @Desc    维修详情表
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/26
 * @Listen  MIT
 */
package maintenance_record

import (
	"gorm.io/gorm"
	"time"
)

type MaintenanceRecord struct {
	Id                int64     `json:"id"`
	State             int64     `json:"state" gorm:"comment:维修周期的状态"`
	MaintenanceId     int64     `json:"maintenance_id" gorm:"comment:维修厂id"`
	SteelId           int64     `json:"steel_id" gorm:"comment:型钢id"`
	OutedAt           time.Time `json:"outed_at" gorm:"comment:出厂时间"`
	EnteredAt         time.Time `json:"entered_at" gorm:"comment:入厂时间"`
	EnteredUid        int64 `json:"entered_uid" gorm:"comment:入厂用户id"`
	OutRepositoryAt   time.Time `json:"outRepository" gorm:"comment:出库时间"`
	EnterRepositoryAt time.Time `json:"enterRepositoryAt" gorm:"comment:入库时间"`
	gorm.Model
}

func (MaintenanceRecord) TableName() string {
	return "maintenance_records"
}
