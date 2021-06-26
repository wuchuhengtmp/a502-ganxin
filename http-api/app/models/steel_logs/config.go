/**
 * @Desc    型钢流转操作记录模型
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/15
 * @Listen  MIT
 */
package steel_logs

import "gorm.io/gorm"

type SteelLog struct {
	ID      int64        `json:"id"`
	Type    SteelLogType `json:"type" gorm:"comment:操作类型"`
	SteelId int64        `json:"steelId" gorm:"comment:型钢id"`
	Uid     int64        `json:"uid" gorm:"comment:操作人的id"`
	gorm.Model
}
type SteelLogType string

const (
	CreateType SteelLogType = "create" // 创建
	OutSteelType SteelLogType = "outSteel" //  出库了
	EnterWorkshopType SteelLogType = "enterWorkshop"  // 入场了
	InstallationType SteelLogType = "InstallationType" // 安装
)

// 操作类型映射说明
var TypeMapName = map[SteelLogType]string {
	CreateType: "型钢入库",
	OutSteelType: "型钢出库",
}

/*
 * 定义表名，用于那些联表查询需要直接使用表名等情况
 */
func (SteelLog) TableName() string {
	return "steel_logs"
}
