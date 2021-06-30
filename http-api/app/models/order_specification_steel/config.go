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
	Id      int64 `json:"id"`
	SteelId int64 `json:"steelId" gorm:"comment:型钢id"`
	// 型钢在仓库和场地的流转周期状态 在库-->场地-->在库
	State int64 `json:"state" gorm:"comment:
100【仓库】-在库
101【仓库】-运送至项目途中
102【仓库】-运送至维修厂途中
200【项目】-待使用
201【项目】-使用中
202【项目】-异常
203【项目】-闲置
204【项目】-准备归库
205【项目】-归库途中
"`
	OrderSpecificationId  int64     `json:"orderSpecificationId" gorm:"comment: 订单规格id"`
	ToWorkshopExpressId   int64     `json:"toWorkshopExpressId" gorm:"comment:从仓库出库到场地的物流单id"`
	ToRepositoryExpressId int64     `json:"toRepositoryExpressId" gorm:"comment:从场地归库的物流单id"`
	LocationCode          int64     `json:"locationCode" gorm:"comment:安装位置编码"`
	EnterRepositoryUid    int64     `json:"enterRepositoryUid" gorm:"comment:入库用户id"`
	OutOfRepositoryUid    int64     `json:"outOfRepositoryUid" gorm:"comment:出库用户id"`
	InstallationUid       int64     `json:"installationUid" gorm:"comment:安装用户id"`
	InstallationAt        time.Time `json:"installationAt" gorm:"comment:安装时间"`
	EnterWorkshopAt       time.Time `json:"intoWorkshopAt" gorm:"comment:入场时间"`
	EnterWorkshopUid      int64     `json:"enterWorkshopUid" gorm:"comment:入场用户"`
	OutWorkshopAt         time.Time `json:"outWorkshopAt" gorm:"comment:出场时间"`
	OutRepositoryAt       time.Time `json:"outRepositoryAt" gorm:"comment:出库时间"`
	EnterRepositoryAt     time.Time `json:"intoRepositoryAt" gorm:"comment:归库时间"`
	gorm.Model
}

func (OrderSpecificationSteel) TableName() string {
	return "order_specification_steels"
}
