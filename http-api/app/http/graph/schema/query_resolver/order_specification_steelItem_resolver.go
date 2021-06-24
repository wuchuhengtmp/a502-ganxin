/**
 * @Desc    订单中的规格列表中的一根型钢的解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/23
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"http-api/app/models/codeinfo"
	"http-api/app/models/order_specification"
	"http-api/app/models/order_specification_steel"
	"http-api/app/models/steels"
	"http-api/app/models/users"
	"http-api/pkg/model"
)

type OrderSpecificationSteelItemResolver struct { }


func (OrderSpecificationSteelItemResolver)Steel(ctx context.Context, obj *order_specification_steel.OrderSpecificationSteel) (*steels.Steels, error) {
	s := steels.Steels{}
	err := model.DB.Model(&s).Where("id = ?", obj.SteelId).First(&s).Error

	return &s, err
}

func (OrderSpecificationSteelItemResolver)OrderSpecification(ctx context.Context, obj *order_specification_steel.OrderSpecificationSteel) (*order_specification.OrderSpecification, error) {
	o := order_specification.OrderSpecification{}
	err := model.DB.Model(&o).Where("id = ?", o.SpecificationId).First(&o).Error

	return &o, err
}
func (OrderSpecificationSteelItemResolver)ToWorkshopExpress(ctx context.Context, obj *order_specification_steel.OrderSpecificationSteel) (*codeinfo.CodeInfo, error) {
	c := codeinfo.CodeInfo{}
	err := model.DB.Model(&c).Where( "id = ?", obj.ToWorkshopExpressId).First(&c).Error

	return &c, err
}

func (OrderSpecificationSteelItemResolver)ToRepositoryExpress(ctx context.Context, obj *order_specification_steel.OrderSpecificationSteel) (*codeinfo.CodeInfo, error) {
	c := codeinfo.CodeInfo{}
	err := model.DB.Model(&c).Where( "id = ?", obj.ToRepositoryExpressId).First(&c).Error

	return &c, err
}

func (OrderSpecificationSteelItemResolver)EnterRepositoryUser(ctx context.Context, obj *order_specification_steel.OrderSpecificationSteel) (*users.Users, error) {
	u := users.Users{}
	err := model.DB.Model(&u).Where("id = ?", obj.EnterRepositoryUid).First(&u).Error

	return &u, err
}
