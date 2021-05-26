/**
 * @Desc    订单详情
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/26
 * @Listen  MIT
 */
package order_details

import "gorm.io/gorm"

type OrderDetail struct {
	Id              int64 `json:"id"`
	SpecificationId int64 `json:"specificationId" gorm:"comment:规格id"`
	Total           int64 `json:"total" gorm:"总量"`
	OrderId         int64 `json:"orderId" gorm:"comment:订单id"`
	gorm.Model
}
