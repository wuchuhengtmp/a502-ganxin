/**
 * @Desc    The requests is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/16
 * @Listen  MIT
 */
package requests

import (
	"context"
	graphModel "http-api/app/http/graph/model"
)

func ValidateSetProjectRequest(ctx context.Context, input graphModel.SetProjectInput) error {
	steps := StepsForProject{}
	if err := steps.CheckHasProject(ctx, input.ID); err != nil {
		return err
	}
	for _, uid := range input.LeaderIDList {
		if err := steps.CheckHasUser(ctx, uid); err != nil {
			return err
		}
		if err := steps.CheckIsProjectRole(ctx, uid); err != nil {
			return err
		}
	}

	return nil
}
