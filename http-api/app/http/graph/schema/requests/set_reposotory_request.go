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

func ValidateSetRepositoryRequest(ctx context.Context, input graphModel.SetRepositoryInput) error {
	steps := StepsForRepository{}
	if err := steps.CheckHasRepository(ctx, input.ID); err != nil {
		return err
	}
	for _, uid := range input.LeaderIDList {
		if err := steps.CheckHasUser(ctx, uid); err != nil {
			return err
		}
		if err := steps.CheckIsRepositoryLeader(ctx, uid); err != nil {
			return err
		}
	}

	return nil
}
