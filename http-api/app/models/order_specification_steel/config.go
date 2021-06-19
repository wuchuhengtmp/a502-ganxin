/**
 * @Desc    订单规格型钢模型
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/16
 * @Listen  MIT
 */
package order_specification_steel

import "gorm.io/gorm"

type OrderSpecificationSteel struct {
	Id                    int64 `json:"id"`
	SteelId               int64 `json:"steelId" gorm:"comment:型钢id"`
	OrderSpecificationId  int64 `json:"orderSpecificationId" gorm:"comment: 订单规格id"`
	ToWorkshopExpressId   int64 `json:"toWorkshopExpressId" gorm:"comment:去出库场地的物流单id"`
	ToRepositoryExpressId int64 `json:"toRepositoryExpressId" gorm:"comment:去场地归库的物流单id"`
	gorm.Model
}

func (OrderSpecificationSteel) TableName() string {
	return "order_specification_steels"
}
