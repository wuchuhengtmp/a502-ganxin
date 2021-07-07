/**
 * @Desc    标记消息为已读
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/7
 * @Listen  MIT
 */
package mutation_resolver

import (
	"context"
	"http-api/app/http/graph/errors"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/msg"
	"http-api/pkg/model"
)

func (*MutationResolver)SetMsgBeRead(ctx context.Context, input graphModel.SetMsgReadedInput) (bool, error) {
	if err := requests.ValidateSetMsgBeReadRequest(ctx, input); err != nil {
		return false, errors.ValidateErr(ctx, err)
	}

	err := model.DB.Model(&msg.Msg{}).Where("id = ?", input.ID).
		Update("is_read", true).
		Error
	if err != nil {
		return false, errors.ServerErr(ctx, err)
	}

	return true, nil
}
