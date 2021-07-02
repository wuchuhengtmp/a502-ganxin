/**
 * @Desc    订单物流模型
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/19
 * @Listen  MIT
 */
package order_express

import (
	"gorm.io/gorm"
	"time"
)

type OrderExpress struct {
	Id               int64                 `json:"id"`
	OrderId          int64                 `json:"orderId"`
	ExpressCompanyId int64                 `json:"codeInfoId" gorm:"comment:码表id上"`
	ExpressNo        string                `json:"expressNo" gorm:"comment物流号"`
	SenderUid        int64                 `json:"senderUid" gorm:"comment:发货人"`
	CompanyId        int64                 `json:"companyId" gorm:"comment:公司id"`
	ReceiveUid       int64                 `json:"receiveUid" gorm:"comment:收货人id"`
	Direction        OrderExpressDirection `json:"direction" gorm:"comment:物流方向toWorkshop去工场,toRepository去仓库"`
	ReceiveAt        time.Time             `json:"receiveAt" gorm:"comment:收货时间"`
	gorm.Model
}
func(OrderExpress)TableName() string {
	return "order_expresses"
}

//  物流方向
type OrderExpressDirection string

const (
	//  去工场方向
	OrderExpressDirectionToWorkshop OrderExpressDirection = "toWorkshop"
	//  归库方向
	OrderExpressDirectionToRepository OrderExpressDirection = "toRepository"
)
