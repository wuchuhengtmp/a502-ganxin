/**
 * @Desc    订单规格型钢模型
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/16
 * @Listen  MIT
 */
package order_specification_steel

import "gorm.io/gorm"

type OrderSpecificationSteel struct {
	Id                   int64 `json:"id"`
	SteelId              int64 `json:"steelId" gorm:"comment:型钢id"`
	OrderSpecificationId int64 `json:"orderSpecificationId" gorm:"comment: 订单规格id"`
	gorm.Model
}

func (OrderSpecificationSteel) TableName () string {
	return "order_specification_steels"
}
