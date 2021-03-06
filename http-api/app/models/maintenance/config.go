/**
 * @Desc    维修厂模型
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/26
 * @Listen  MIT
 */
package maintenance

import (
	"gorm.io/gorm"
	"http-api/app/models/maintenance_record"
)

type Maintenance struct {
	Id        int64  `json:"id"`
	Name      string `json:"name" gorm:"comment:维修厂名"`
	Address   string `json:"address" gorm:"comment:地址"`
	Remark    string `json:"remark" gorm:"comment:备注"`
	IsAble    bool   `json:"isAble" gorm:"comment:是否启动"`
	CompanyId int64  `json:"companyId" gorm:"comment:公司id"`
	gorm.Model
}

func (Maintenance) TableName() string {
	return "maintenance"
}

// 获取待入厂详细信息参数
type GetEnterMaintenanceSteelDetailRes struct {
	// 入厂型钢列表
	List []*maintenance_record.MaintenanceRecord
	// 数量
	Total int64
	// 重量
	Weight float64
}

/**
 * 待维修型钢详情响应
 */
type GetChangedMaintenanceSteelDetailRes struct {
	// 维修型钢列表
	List []*maintenance_record.MaintenanceRecord
	// 数量
	Total int64
	// 重量
	Weight float64
}

/**
 * 获取可出厂的型钢详情响应数据
 */
type GetSteelForOutOfMaintenanceDetailRes struct {
	//  维修型钢详情
	List []*maintenance_record.MaintenanceRecord
	// 数量
	Total int64
	// 重量
	Weight float64
}

/**
 * 获取维修厂维修记录响应格式
 */
type GetMaintenanceSteelResItem struct {
	Id int64
	// 规格
	Specification string
	// 已接收
	ReceivedTotal int64
	// 已接收重量
	ReceivedWeight float64
	// 已归库数量
	StoredTotal int64
	// 已归库重量
	StoredWeight float64
}

/**
 *  获取维修厂维修记录响应格式
 */
type GetMaintenanceSteelRes struct {
	// 规格列表
	List []*GetMaintenanceSteelResItem
	// 数量
	Total int64
	// 重量
	Weight float64
}
