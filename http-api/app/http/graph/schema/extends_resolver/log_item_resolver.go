/**
 * @Desc    The extends_resolver is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/8
 * @Listen  MIT
 */
package extends_resolver

import (
	"context"
	"fmt"
	"http-api/app/models/logs"
	"http-api/app/models/users"
	"http-api/pkg/model"
)

type LogItemResolver struct{}

func (LogItemResolver) User(ctx context.Context, obj *logs.Logos) (*users.Users, error) {
	userItem := users.Users{}
	err := model.DB.Model(&userItem).Where("id = ?", obj.Uid).First(&userItem).Error

	return &userItem, err
}
func (LogItemResolver) TypeInfo(ctx context.Context, obj *logs.Logos) (*logs.LogTypeItem, error) {
	l := logs.LogTypeItem{
		Flag: fmt.Sprintf("%s", obj.Type),
		Desc: logs.ActionTypeMapDes[obj.Type],
	}
	return &l, nil
}
