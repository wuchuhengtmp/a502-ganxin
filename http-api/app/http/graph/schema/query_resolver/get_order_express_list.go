/**
 * @Desc    The query_resolver is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/19
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"http-api/app/http/graph/auth"
	"http-api/app/models/codeinfo"
	"http-api/app/models/order_express"
	"http-api/app/models/orders"
	"http-api/app/models/users"
	"http-api/pkg/model"
)

type OrderExpressItemResolver struct { }

func (OrderExpressItemResolver) Sender(ctx context.Context, obj *order_express.OrderExpress) (u *users.Users, err error) {
	err = model.DB.Model(&users.Users{}).Where("id = ?", obj.SenderUid).First(&u).Error

	return
}

func (OrderExpressItemResolver)Receiver(ctx context.Context, obj *order_express.OrderExpress) (u *users.Users, err error) {
	err = model.DB.Model(&users.Users{}).Where("id = ?", obj.ReceiveUid).First(&u).Error

	return
}
func (OrderExpressItemResolver)ExpressList(ctx context.Context, obj *orders.Order) (expressList []*order_express.OrderExpress, err error) {
	me := auth.GetUser(ctx)
	err = model.DB.Model(&order_express.OrderExpress{}).
		Where("order_id = ?", obj.Id).
		Where("company_id = ?", me.CompanyId).
		Find(&expressList).Error

	return
}
func (OrderExpressItemResolver) ExpressCompany(ctx context.Context, obj *order_express.OrderExpress) (*codeinfo.CodeInfo, error) {
	c := codeinfo.CodeInfo{}
	if err := model.DB.Model(&codeinfo.CodeInfo{}).Where("id = ?", obj.CompanyId).First(&c).Error; err != nil {
		return nil, err
	}

	return &c, nil
}
