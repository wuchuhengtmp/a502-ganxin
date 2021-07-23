/**
 * @Desc    The requests is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/23
 * @Listen  MIT
 */
package requests

import (
	"context"
	graphModel "http-api/app/http/graph/model"
)

func ValidateDeleteOrder(ctx context.Context, input graphModel.DeleteOrderInput) error {
	 steps := StepsForProject{}
	 if err := steps.CheckHasOrder(ctx, input.ID); err != nil {
	 	return err
	 }

	return nil
}


