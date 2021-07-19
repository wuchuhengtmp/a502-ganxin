/**
 * @Desc    The msg is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/26
 * @Listen  MIT
 */
package msg

import (
	"gorm.io/gorm"
)

type Msg struct {
	Id      int64  `json:"id"`
	IsRead  bool   `json:"isRead" gorm:"comment:是否已读"`
	Content string `json:"content" gorm:"comment:内容"`
	Uid     int64  `json:"uid" gorm:"comment:用户id"`
	Type    string `json:"type" gorm:"comment:消息类型"`
	Extends string `json:"extends" gorm:"comment:扩展参数json格式,参数用于点击消息能识别类型并自动跳到对应的页面"`
	gorm.Model
}

const (
	CreateOrderType       string = "createOrder"                // 创建订单类型
	ConfirmOrderType      string = "confirmOrder"               // 确认订单类型
	RejectOrderType       string = "rejectOrderType"            // 拒绝订单类型
	OutProject2Workshop   string = "outOfProjectToWorkshopType" // 出库到项目场地
	EnterProject2Workshop string = "enter project to Workshop"  // 项目入场
	OutOfWorkshop         string = "OutOfWorkshop"              // 型钢出场
	DelMaintenance        string = "DelMaintenance"             // 删除维修厂（下岗了呗）
	ToBeMaintained        string = "ToBeMaintained"             // 待维修
	OutOfMaintenance      string = "OutOfMaintenance"           // 出厂
	EnterMaintenance      string = "EnterMaintenance"           // 入厂
	DeleteProject         string = "DeleteProject"              // 删除项目
)

func (Msg) TableName() string {
	return "msg"
}

/**
 * 推送消息
 */
func (*Msg) Push() error {
	// todo 推送消息逻辑
	return nil
}

/**
 * 创建消息
 */
func (m *Msg) CreateSelf(tx *gorm.DB) error {
	if err := tx.Create(m).Error; err != nil {
		return err
	}
	_ = m.Push()

	return nil
}
