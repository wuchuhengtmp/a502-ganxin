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
	"http-api/app/models/order_specification"
	"http-api/app/models/orders"
	"http-api/app/models/projects"
	"http-api/app/models/repositories"
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
	if err := u.GetSelfById(obj.CreateUid); err != nil {
		return nil, err
	}

	return &u, nil
}

func (OrderItemResolver) ReceiveUser(ctx context.Context, obj *orders.Order) (*users.Users, error) {
	// 有收货人
	if obj.State >= orders.StateReceipted {
		u := users.Users{}
		if err := u.GetSelfById(obj.ReceiveUid); err != nil {
			return nil, err
		} else {
			return &u, nil
		}
	}

	return nil, nil
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
		if err := u.GetSelfById(obj.ConfirmedUid); err != nil {
			return nil, err
		}
		return &u, nil
	} else {
		// 还没有确认用户
		return nil, nil
	}
}

func (OrderItemResolver) Sender(ctx context.Context, obj *orders.Order) (*users.Users, error) {
	// 有发货
	if obj.State >= orders.StateSend {
		u := users.Users{}
		if err := u.GetSelfById(obj.SenderUid); err != nil {
			return nil, err
		} else {
			return &u, nil
		}
	}

	return nil, nil
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
	err := model.DB.Model(&order_specification.OrderSpecification{}).
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
