/**
 * @Desc    获取消息列表
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/26
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"http-api/app/http/graph/auth"
	"http-api/app/models/msg"
	"http-api/pkg/model"
)

func (*QueryResolver)GetMsgList(ctx context.Context) (res []*msg.Msg, err error) {
	me := auth.GetUser(ctx)
	err = model.DB.Model(&msg.Msg{}).Where("uid = ?", me.Id).Order("id desc").Find(&res).Error

	return
}
