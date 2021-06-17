/**
 * @Desc    设置密码解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/15
 * @Listen  MIT
 */
package mutation_resolver

import (
	"context"
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/errors"
	"http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
)

func (*MutationResolver)SetPassword(ctx context.Context, input *model.SetPasswordInput) (bool, error) {
	if err := requests.ValidateSetPasswordRequest(ctx, input); err != nil {
		return false, errors.ValidateErr(ctx, err)
	}
	me := auth.GetUser(ctx)
	if err := me.SetSelfPassword(input.Password); err != nil {
		return false, err
	}

	return true, nil
}

