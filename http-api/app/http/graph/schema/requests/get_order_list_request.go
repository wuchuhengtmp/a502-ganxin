/**
 * @Desc    获取订单列表验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/18
 * @Listen  MIT
 */
package requests

import (
	"context"
	"fmt"
	"http-api/app/http/graph/auth"
	graphModel "http-api/app/http/graph/model"
)

func ValidateGetOrderListRequest(ctx context.Context, input graphModel.GetOrderListInput) error {
	isDevice := auth.IsDevice(ctx)
	if isDevice {
		if input.QueryType == nil {
			return fmt.Errorf("设备端获取订单列表需要要选择确认订单类型或未确认订单类型")
		}
	}

	return nil
}
