/**
 * @Desc    The mutation_resolver is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/14
 * @Listen  MIT
 */
package mutation_resolver

import (
	"context"
	"fmt"
	"http-api/app/http/graph/errors"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/models/users"
	"http-api/pkg/helper"
	"http-api/pkg/model"
	"time"
)

func (*MutationResolver)ResetPassword(ctx context.Context, input graphModel.ResetPasswordInput) (bool, error) {
	if err := ValidateResetPasswordRequest(input); err != nil {
		return false, errors.ValidateErr(ctx, err)
	}
	smsInfo, _ := SmsKeyMapCode[input.Key]
	err := model.DB.Model(&users.Users{}).
		Where("phone = ?", smsInfo.Phone).
		Update("password", helper.GetHashByStr(input.NewPassword)).
		Error
	if err != nil {
		return false, errors.ServerErr(ctx, err)
	}
	delete(SmsKeyMapCode, input.Key)

	return true, nil
}

func ValidateResetPasswordRequest(input graphModel.ResetPasswordInput) error {
	item, ok := SmsKeyMapCode[input.Key]
	if !ok {
		return fmt.Errorf("没有这个验证码key")
	}
	if item.Code != input.Code {
		return fmt.Errorf("验证码不正确")
	}
	timeLen := time.Now().Unix() - item.CreatedAt.Unix()
	if timeLen > 60*30 {
		return fmt.Errorf("验证码已过期")
	}
	passwordLen := 6
	if len(input.NewPassword) < passwordLen {
		return fmt.Errorf("密码不能小于%d位", passwordLen)
	}
	return nil
}
