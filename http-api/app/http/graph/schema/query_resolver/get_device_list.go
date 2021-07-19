/**
 * @Desc    The query_resolver is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/10
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"fmt"
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/errors"
	"http-api/app/models/devices"
	"http-api/app/models/users"
	"http-api/pkg/model"
)

func (*QueryResolver) GetDeviceList(ctx context.Context) ([]*devices.Device, error) {
	deviceItem := devices.Device{}
	deviceTable := deviceItem.TableName()
	userTable := users.Users{}.TableName()
	var res []*devices.Device
	me := auth.GetUser(ctx)
	err := model.DB.Model(&deviceItem).
		Select(fmt.Sprintf("%s.*", deviceTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.uid", userTable, userTable, deviceTable)).
		Where(fmt.Sprintf("%s.deleted_at is NULL", userTable)).
		Where(fmt.Sprintf("%s.company_id = ?", userTable), me.CompanyId).
		Find(&res).
		Error
	if err != nil {
		return res, errors.ServerErr(ctx, err)
	}

	return res, nil
}

type DeviceItemResolver struct{}

/**
 * 用户字段解析
 */
func (DeviceItemResolver) UserInfo(ctx context.Context, obj *devices.Device) (*users.Users, error) {
	return obj.GetUser()
}
