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
)

type OrderSpecificationSteelItemResolver struct { }


func (OrderSpecificationSteelItemResolver)Steel(ctx context.Context, obj *order_specification_steel.OrderSpecificationSteel) (*steels.Steels, error) {
	s := steels.Steels{}
	// todo

	return &s, nil
}

func (OrderSpecificationSteelItemResolver)OrderSpecification(ctx context.Context, obj *order_specification_steel.OrderSpecificationSteel) (*order_specification.OrderSpecification, error) {
	o := order_specification.OrderSpecification{}

	return &o, nil
}
func (OrderSpecificationSteelItemResolver)ToWorkshopExpress(ctx context.Context, obj *order_specification_steel.OrderSpecificationSteel) (*codeinfo.CodeInfo, error) {
	o := codeinfo.CodeInfo{}

	return &o, nil
}

func (OrderSpecificationSteelItemResolver)ToRepositoryExpress(ctx context.Context, obj *order_specification_steel.OrderSpecificationSteel) (*codeinfo.CodeInfo, error) {
	// todo
	o := codeinfo.CodeInfo{}

	return &o, nil
}

func (OrderSpecificationSteelItemResolver)EnterRepositoryUser(ctx context.Context, obj *order_specification_steel.OrderSpecificationSteel) (*users.Users, error) {
	// todo
	u := users.Users{}

	return &u, nil
}
