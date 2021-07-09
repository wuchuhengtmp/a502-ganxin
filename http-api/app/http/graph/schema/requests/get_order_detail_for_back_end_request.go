/**
 * @Desc    获取订单详情验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/9
 * @Listen  MIT
 */
package requests

import (
	"context"
	graphModel "http-api/app/http/graph/model"
)

func ValidateGetOrderDetailForBackEndRequest(ctx context.Context, input graphModel.GetOrderDetailForBackEndInput) error  {
	steps := StepsForOrder{}
	// 检验分页
	if err := steps.CheckPagination(input.IsShowAll, input.PageSize, input.Page); err != nil {
		return err
	}

	return nil
}
