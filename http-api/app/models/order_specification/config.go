/**
 * @Desc    订单详情
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/26
 * @Listen  MIT
 */
package order_specification

import (
	"gorm.io/gorm"
	"http-api/pkg/model"
)

type OrderSpecification struct {
	Id              int64  `json:"id"`
	SpecificationId int64  `json:"specificationId" gorm:"comment:规格id"`
	Total           int64  `json:"total" gorm:"总量"`
	Specification   string `json:"specification" gorm:"冗余规格"`
	OrderId         int64  `json:"orderId" gorm:"comment:订单id"`
	gorm.Model
}

func (OrderSpecification)TableName() string {
	return "order_specifications"
}

func (*OrderSpecification) GetOrderBySpecificationId(specificationId int64) (orderDetails []*OrderSpecification, err error) {
	db := model.DB
	if err := db.Model(&OrderSpecification{}).Where("specification_id = ?", specificationId).Find(&orderDetails).Error; err != nil {
		return orderDetails, err
	}

	return
}