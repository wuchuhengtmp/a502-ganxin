/**
 * @Desc    获取订单列表
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/16
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"http-api/app/http/graph/errors"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/http/graph/schema/services"
	"http-api/app/models/codeinfo"
	"http-api/app/models/order_express"
	"http-api/app/models/order_specification"
	"http-api/app/models/order_specification_steel"
	"http-api/app/models/orders"
	"http-api/app/models/projects"
	"http-api/app/models/repositories"
	"http-api/app/models/specificationinfo"
	"http-api/app/models/steels"
	"http-api/app/models/users"
	"http-api/pkg/model"
)

func (*QueryResolver) GetOrderList(ctx context.Context, input graphModel.GetOrderListInput) ([]*orders.Order, error) {
	if err := requests.ValidateGetOrderListRequest(ctx, input); err != nil {
		var tmp []*orders.Order
		return tmp, errors.ValidateErr(ctx, err)
	}

	return services.GetOrderList(ctx, input)
}

type OrderItemResolver struct{}

func (OrderItemResolver) Project(ctx context.Context, obj *orders.Order) (*projects.Projects, error) {
	orderService := services.OrderService{}
	orderService.OrderMoveIntoMe(obj)
	p, err := orderService.GetProject()

	return &p, err
}
func (OrderItemResolver) CreateUser(ctx context.Context, obj *orders.Order) (*users.Users, error) {
	u := users.Users{}
	err := model.DB.Unscoped().Model(&u).Where("id = ?", obj.CreateUid).First(&u).Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	return &u, nil
}

func (OrderItemResolver) Repository(ctx context.Context, obj *orders.Order) (*repositories.Repositories, error) {
	r := repositories.Repositories{ID: obj.RepositoryId}
	if err := r.GetSelf(); err != nil {
		return nil, err
	}

	return &r, nil
}

func (OrderItemResolver) ConfirmedUser(ctx context.Context, obj *orders.Order) (*users.Users, error) {
	// 有确认用户
	if obj.State >= orders.StateConfirmed {
		u := users.Users{}
		err := model.DB.Unscoped().Model(&u).Where("id = ?", obj.ConfirmedUid).First(&u).Error
		if  err != nil {
			return nil, err
		}
		return &u, nil
	} else {
		// 还没有确认用户
		return nil, nil
	}
}

/**
 * 订单数量字段解析
 */
func (OrderItemResolver) Total(ctx context.Context, obj *orders.Order) (int64, error) {
	res := struct {
		Tmp int64
	}{}
	err := model.DB.Model(&order_specification.OrderSpecification{}).
		Select("sum(total) as tmp").
		Where("order_id = ?", obj.Id).
		Scan(&res).
		Error
	if err != nil {
		return 0, err
	}

	return res.Tmp, nil
}

/**
 * 重量字段解析
 */
func (OrderItemResolver) Weight(ctx context.Context, obj *orders.Order) (float64, error) {
	var orderSpecificationList []*order_specification.OrderSpecification
	err := model.DB.Unscoped().Model(&order_specification.OrderSpecification{}).
		Where("order_id = ?", obj.Id).
		Find(&orderSpecificationList).Error
	if err != nil {
		return 0, err
	}
	var weight float64
	for _, item := range orderSpecificationList {
		s, err := item.GetSpecification()
		if err != nil {
			return 0, err
		}
		weight += float64(item.Total) * s.Weight
	}

	return weight, nil
}

func (OrderItemResolver) ExpressCompany(ctx context.Context, obj *orders.Order) (*codeinfo.CodeInfo, error) {
	if obj.State >= orders.StateSend {
		c := codeinfo.CodeInfo{}
		if err := model.DB.Model(&c).Where("id = ?", obj).First(&c).Error; err != nil {
			return nil, err
		}
		return &c, nil
	}

	return nil, nil
}

func (OrderItemResolver) OrderSpecificationList(ctx context.Context, obj *orders.Order) (list []*order_specification.OrderSpecification, err error) {
	oo := order_specification.OrderSpecification{}
	err = model.DB.Unscoped().Model(&oo).Where("order_id = ?", obj.Id).Find(&list).Error

	return
}

type OrderSpecificationItemResolver struct{}
func (OrderSpecificationItemResolver)Order(ctx context.Context, obj *order_specification.OrderSpecification) (*orders.Order, error) {
	orderItem := orders.Order{}
	err := model.DB.Model(&orderItem).Where("id = ?", obj.OrderId).First(&orderItem).Error
	if err != nil {
		return nil , err
	}

	return &orderItem, nil
}

func (OrderSpecificationItemResolver) Weight(ctx context.Context, obj *order_specification.OrderSpecification) (float64, error) {
	sp := specificationinfo.SpecificationInfo{ID: obj.SpecificationId}
	if err := sp.GetUnscopedSelf(); err != nil {
		return 0, err
	}

	return sp.Weight * float64(obj.Total), nil
}

func (OrderSpecificationItemResolver) TotalSend(ctx context.Context, obj *order_specification.OrderSpecification) (total int64, err error) {
	err = model.DB.Model(&order_specification_steel.OrderSpecificationSteel{}).
		Where("order_specification_id = ?", obj.Id).
		Count(&total).Error

	return
}

func (OrderSpecificationItemResolver) TotalToBeSend(ctx context.Context, obj *order_specification.OrderSpecification) (total int64, err error) {
	err = model.DB.Model(&order_specification_steel.OrderSpecificationSteel{}).
		Where("order_specification_id = ?", obj.Id).
		Count(&total).Error

	return obj.Total - total, err
}

func (OrderItemResolver) ExpressList(ctx context.Context, obj *orders.Order) (orderExpressList []*order_express.OrderExpress, err error) {
	err = model.DB.Model(&order_express.OrderExpress{}).
		Where("order_id = ?", obj.Id).
		Find(&orderExpressList).
		Error

	return orderExpressList, err
}
func (OrderSpecificationItemResolver) SpecificationInfo(ctx context.Context, obj *order_specification.OrderSpecification) (*specificationinfo.SpecificationInfo, error) {
	s := specificationinfo.SpecificationInfo{}
	err := model.DB.Unscoped().Model(&s).Where("id = ?", obj.SpecificationId).First(&s).Error

	return &s, err
}

// 已归库(出场并已保存到仓库中)
func (OrderSpecificationItemResolver) StoreTotal(ctx context.Context, obj *order_specification.OrderSpecification) (int64, error) {
	var total int64
	stateList := []int64{
		steels.StateInStore,              //【仓库】-在库
	}
	err := model.DB.Model(&order_specification_steel.OrderSpecificationSteel{}).
		Where("state in ?", stateList).
		Count(&total).
		Error

	return total, err
}

//（场地）已接收过的统计
func (OrderSpecificationItemResolver) WorkshopReceiveTotal(ctx context.Context, obj *order_specification.OrderSpecification) (int64, error) {
	var total int64
	stateList := steels.GetStateForProject()
	err := model.DB.Model(&order_specification_steel.OrderSpecificationSteel{}).
		Where("state in ?", stateList).
		Count(&total).
		Error

	return total, err
}
