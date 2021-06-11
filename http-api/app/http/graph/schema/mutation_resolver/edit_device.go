/**
 * @Desc    编辑设备解析器
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/11
 * @Listen  MIT
 */
package mutation_resolver

import (
	"context"
	"http-api/app/http/graph/errors"
	"http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/devices"
)

func (*MutationResolver) EditDevice(ctx context.Context, input model.EditDeviceInput) (bool, error) {
	if err := requests.ValidateEditDevice(ctx, input); err != nil {
		return false, errors.ValidateErr(ctx, err)
	}
	d := devices.Device{}
	_ = d.GetDeviceSelfById(input.ID)
	if err := d.EditSelf(ctx); err != nil {
		return false, errors.ServerErr(ctx, err)
	}

	return true, nil
}
