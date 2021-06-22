/**
 * @Desc    订单规格型钢模型
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/16
 * @Listen  MIT
 */
package order_specification_steel

import (
	"gorm.io/gorm"
	"time"
)

type OrderSpecificationSteel struct {
	Id                    int64     `json:"id"`
	SteelId               int64     `json:"steelId" gorm:"comment:型钢id"`
	OrderSpecificationId  int64     `json:"orderSpecificationId" gorm:"comment: 订单规格id"`
	ToWorkshopExpressId   int64     `json:"toWorkshopExpressId" gorm:"comment:去出库场地的物流单id"`
	ToRepositoryExpressId int64     `json:"toRepositoryExpressId" gorm:"comment:去场地归库的物流单id"`
	LocationCode          int64     `json:"locationCode" gorm:"comment:安装位置编码"`
	EnterRepositoryUid    int64     `json:"enterRepositoryUid" gorm:"comment:出库用户id"`
	InstallationAt        time.Time `json:"installationAt" gorm:"comment:安装时间"`
	EnterWorkshopAt       time.Time `json:"intoWorkshopAt" gorm:"comment:入场时间"`
	OutWorkshopAt         time.Time `json:"outWorkshopAt" gorm:"comment:出场时间"`
	OutRepositoryAt       time.Time `json:"outRepositoryAt" gorm:"comment:出库时间"`
	EnterRepositoryAt     time.Time `json:"intoRepositoryAt" gorm:"comment:归库时间"`
	gorm.Model
}

func (OrderSpecificationSteel) TableName() string {
	return "order_specification_steels"
}
