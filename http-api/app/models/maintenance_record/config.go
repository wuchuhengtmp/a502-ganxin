/**
 * @Desc    维修详情表
 * @Author  wuchuheng<wuchuheng@163.com>
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
	Id        int64     `json:"id"`
	FixId     int64     `json:"fix_id" gorm:"comment:维修厂id"`
	SteelId   int64     `json:"steel_id" gorm:"comment:型钢id"`
	OutedAt   time.Time `json:"outed_at" gorm:"comment:出厂时间"`
	EnteredAt time.Time `json:"entered_at" gorm:"comment:入厂时间"`
	gorm.Model
}
