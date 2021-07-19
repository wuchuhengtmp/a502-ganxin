/**
 * @Desc    必须是设备指令
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/18
 * @Listen  MIT
 */
package directives

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/errors"
	"http-api/app/models/devices"
	"http-api/pkg/model"
)

func MustBeDevice (ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	if err := auth.ValidateToken(ctx); err != nil {
		return errors.AccessDenied(ctx, err.Error())
	}
	if !auth.IsDevice(ctx) {
		errMsg := fmt.Sprintf("必须是设备才能访问")
		return errors.AccessDenied(ctx, errMsg)
	}
	// 检验设备是否禁用
	me := auth.GetUser(ctx)
	mac := auth.GetMac(ctx)
	deviceItem := devices.Device{}
	err = model.DB.Model(&deviceItem).Where("uid = ?", me.Id).
		Where("mac = ?", mac).
		First(&deviceItem).
		Error
	if err == nil && deviceItem.IsAble == false{
		return errors.DeviceDenied(ctx, "当前设备禁止使用这个账号，请找公司管理员解禁")
	}

	return next(ctx)
}
