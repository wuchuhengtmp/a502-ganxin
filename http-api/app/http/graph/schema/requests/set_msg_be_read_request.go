/**
 * @Desc    标记消息为已读请求验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/7
 * @Listen  MIT
 */
package requests

import (
	"context"
	"fmt"
	"http-api/app/http/graph/auth"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/models/msg"
	"http-api/pkg/model"
)

func ValidateSetMsgBeReadRequest(ctx context.Context, input graphModel.SetMsgReadedInput) error {
	me := auth.GetUser(ctx)
	// 检验有没有这条消息
	err := model.DB.Model(&msg.Msg{}).Where("uid = ?", me.Id).
		Where("id = ?", input.ID).
		First(&msg.Msg{}).
		Error
	if err != nil && err.Error() == "record not found" {
		return fmt.Errorf("消息id为: %d 不存在", input.ID)
	}

	return nil
}