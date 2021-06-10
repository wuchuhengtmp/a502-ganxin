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
)

func (*QueryResolver) GetDeviceList(ctx context.Context) ([]*devices.Device, error) {
	var ds  []*devices.Device

	return ds, nil
}

type DeviceItemResolver struct { }

//func (DeviceItemResolver)UserInfo(ctx context.Context, obj *devices.Device) (*model.UserItem, error) {
//
//}


