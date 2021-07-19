/**
 * @Desc    The requests is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/19
 * @Listen  MIT
 */
package requests

import (
	"context"
	"fmt"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/models/orders"
	"http-api/pkg/model"
)

func ValidateDeleteProjectRequest(ctx context.Context, input graphModel.DeleteProjectInput) error {
	steps := StepsForProject{}
	if err := steps.CheckHasProject(ctx, input.ID); err != nil {
		return err
	}
	// 有没有订单
	orderItem := orders.Order{}
	var total int64
	if err := model.DB.Model(&orderItem).Where( "project_id = ?", input.ID).Count(&total).Error; err != nil {
		return err
	}
	if total != 0 {
		return fmt.Errorf("该项目已有订单，不能删除")
	}

	return nil
}
