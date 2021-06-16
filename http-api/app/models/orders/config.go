/**
 * @Desc    订单模型
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/26
 * @Listen  MIT
 */
package orders

import (
	"gorm.io/gorm"
	"http-api/pkg/model"
	"time"
)

type Order struct {
	Id               int64     `json:"id"`
	ProjectId        int64     `json:"ProjectId" gorm:"comment:项目id"`
	State            int       `json:"state" gorm:"仓库状态100待确认200已确认300已拒绝400已发货500待收货600已收货(部分)700已收货全部800已归库"`
	ReceiveState     int       `json:"receiveState" gorm:"comment:场地状态 "`
	ExpectedReturnAt time.Time `json:"expectedReturnAt" gorm:"comment:预计归还时间"`
	PartList         string    `json:"partList" gorm:"comment:配件清单"`
	CreateUid        int64     `json:"createUid" gorm:"comment:创建人"`
	ConfirmedAt      time.Time `json:"confirmAt" gorm:"comment:确认时间"`
	ReceiveUid       int64     `json:"receiveUid" gorm:"comment:收货人id"`
	ReceiveAt        time.Time `json:"receiveAt" gorm:"comment:收货时间"`
	ExpressCompanyId int64     `json:"expressCompanyId" gorm:"comment:快递公司(码表id)"`
	ExpressNo        string     `json:"expressNo" gorm:"comment:物流号"`
	Remark           string    `json:"remark" gorm:"comment:备注"`
	gorm.Model
}

/**
 * 根据物流公司获取订单列表
 */
func (*Order) GetOrdersByExpressId(expressId int64) (os []*Order, err error) {
	err = model.DB.Model(&Order{}).Where("express_company_id = ?", expressId).Find(&os).Error

	return
}
