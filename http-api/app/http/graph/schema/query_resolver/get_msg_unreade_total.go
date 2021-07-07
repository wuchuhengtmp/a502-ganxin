/**
 * @Desc    获取未读消息量
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/7
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/errors"
	"http-api/app/models/msg"
	"http-api/pkg/model"
)

func (*QueryResolver)GetMsgUnReadeTotal(ctx context.Context) (int64, error) {
	me := auth.GetUser(ctx)
	var total int64
	err := model.DB.Model(&msg.Msg{}).Where("uid = ?", me.Id).
		Where("is_read = ?", false).
		Count(&total).
		Error
	if err != nil {
		return total, errors.ServerErr(ctx, err)
	}

	return total, nil
}
