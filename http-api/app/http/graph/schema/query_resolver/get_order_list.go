/**
 * @Desc    The query_resolver is part of http-api
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/16
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"http-api/app/models/companies"
	"http-api/app/models/orders"
	"http-api/app/models/projects"
	"http-api/app/models/users"
)

type OrderItemResolver struct { }
func (OrderItemResolver)Project(ctx context.Context, obj *orders.Order) (*projects.Projects, error) {
	// todo
	var p projects.Projects

	return &p, nil
}
func (OrderItemResolver)State(ctx context.Context, obj *orders.Order) (int64, error) {
	// todo

	return 1, nil
}

func (OrderItemResolver)CreateUser(ctx context.Context, obj *orders.Order) (*users.Users, error) {
	// todo
	var u users.Users

	return &u, nil
}

func (OrderItemResolver)ReceiveUser(ctx context.Context, obj *orders.Order) (*users.Users, error) {
	// todo
	var u users.Users

	return &u, nil
}


func (OrderItemResolver)ExpressCompany(ctx context.Context, obj *orders.Order) (*companies.Companies, error) {
	// todo
	var c companies.Companies

	return &c, nil
}
