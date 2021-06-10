/**
 * @Desc    The query_resolver is part of http-api
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/10
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"http-api/app/models/devices"
	"http-api/app/models/users"
)

func (*QueryResolver) GetDeviceList(ctx context.Context) ([]*devices.Device, error) {
	d := devices.Device{}

	return d.GetAll(ctx)
}

type DeviceItemResolver struct{}

/**
 * 用户字段解析
 */
func (DeviceItemResolver) UserInfo(ctx context.Context, obj *devices.Device) (*users.Users, error) {
	return obj.GetUser()
}
