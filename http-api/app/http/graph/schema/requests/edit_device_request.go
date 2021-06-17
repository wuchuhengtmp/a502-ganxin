/**
 * @Desc    The requests is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/11
 * @Listen  MIT
 */
package requests

import (
	"context"
	"fmt"
	"github.com/thedevsaddam/govalidator"
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/model"
	"http-api/app/models/devices"
)

func ValidateEditDevice(ctx context.Context, input model.EditDeviceInput) error {
	rules := govalidator.MapData{
		"id": []string{"isDeviceId"},
	}
	opts := govalidator.Options{
		Data:          &input,
		Rules:         rules,
		TagIdentifier: "json",
	}
	res := govalidator.New(opts).ValidateStruct()
	if len(res) > 0 {
		for _, fieldErrors := range res {
			for _, err := range fieldErrors {
				return fmt.Errorf("%s", err)
			}
		}
	}
	d := devices.Device{}
	_ = d.GetDeviceSelfById(input.ID)
	me := auth.GetUser(ctx)
	u, _ := d.GetUser()
	if u.CompanyId != me.CompanyId {
		return fmt.Errorf("该设备与您不归属于同一家公司名下，您无权操作")
	}

	return nil
}
