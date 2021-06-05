/**
 * @Desc    订单详情
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/26
 * @Listen  MIT
 */
package order_details

import (
	"gorm.io/gorm"
	"http-api/pkg/model"
)

type OrderDetail struct {
	Id              int64 `json:"id"`
	SpecificationId int64 `json:"specificationId" gorm:"comment:规格id"`
	Total           int64 `json:"total" gorm:"总量"`
	OrderId         int64 `json:"orderId" gorm:"comment:订单id"`
	gorm.Model
}

func (*OrderDetail) GetOrderBySpecificationId(specificationId int64) (orderDetails []*OrderDetail, err error)  {
	db := model.DB
	if err := db.Model(&OrderDetail{}).Where("specification_id = ?", specificationId).Find(&orderDetails).Error; err != nil {
		return orderDetails, err
	}

	return
}
