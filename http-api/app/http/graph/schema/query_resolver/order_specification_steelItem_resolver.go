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
	// "errors"
	// "fmt"
	// "runtime/debug"
	"http-api/app/models/codeinfo"
	"http-api/app/models/order_specification"
	"http-api/app/models/order_specification_steel"
	"http-api/app/models/steels"
	"http-api/app/models/users"
	"http-api/pkg/model"
)

type OrderSpecificationSteelItemResolver struct{}

func (OrderSpecificationSteelItemResolver) UseDays(ctx context.Context, obj *order_specification_steel.OrderSpecificationSteel) (*int64, error) {
	var days int64 = 0
	if obj.OutWorkshopAt.Unix() > obj.EnterWorkshopAt.Unix() {
		days = (obj.OutWorkshopAt.Unix() - obj.EnterWorkshopAt.Unix()) / (60 * 60 * 24)
	}
	return &days, nil

}

func (OrderSpecificationSteelItemResolver) Steel(ctx context.Context, obj *order_specification_steel.OrderSpecificationSteel) (*steels.Steels, error) {
	s := steels.Steels{}
	err := model.DB.Model(&s).Where("id = ?", obj.SteelId).First(&s).Error

	return &s, err
}

func (OrderSpecificationSteelItemResolver) StateInfo(ctx context.Context, obj *order_specification_steel.OrderSpecificationSteel) (*steels.StateItem, error) {
	item := steels.StateItem{
		State: obj.State,
		Desc:  steels.StateCodeMapDes[obj.State],
	}

	return &item, nil
}

func (OrderSpecificationSteelItemResolver) OrderSpecification(ctx context.Context, obj *order_specification_steel.OrderSpecificationSteel) (*order_specification.OrderSpecification, error) {
	o := order_specification.OrderSpecification{}
	err := model.DB.Model(&o).Where("id = ?", obj.OrderSpecificationId).First(&o).Error
	return &o, err
}
func (OrderSpecificationSteelItemResolver) ToWorkshopExpress(ctx context.Context, obj *order_specification_steel.OrderSpecificationSteel) (*codeinfo.CodeInfo, error) {
	c := codeinfo.CodeInfo{}
	err := model.DB.Model(&c).Where("id = ?", obj.ToWorkshopExpressId).First(&c).Error

	return &c, err
}

func (OrderSpecificationSteelItemResolver) ToRepositoryExpress(ctx context.Context, obj *order_specification_steel.OrderSpecificationSteel) (*codeinfo.CodeInfo, error) {
	c := codeinfo.CodeInfo{}
	err := model.DB.Model(&c).Where("id = ?", obj.ToRepositoryExpressId).First(&c).Error

	return &c, err
}

func (OrderSpecificationSteelItemResolver) EnterRepositoryUser(ctx context.Context, obj *order_specification_steel.OrderSpecificationSteel) (*users.Users, error) {
	u := users.Users{}
	if obj.EnterRepositoryUid != 0 {
		err := model.DB.Model(&u).Where("id = ?", obj.EnterRepositoryUid).First(&u).Error
		return &u, err
	}

	return nil, nil

}
