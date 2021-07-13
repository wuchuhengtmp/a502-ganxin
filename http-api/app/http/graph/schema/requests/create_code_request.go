/**
 * @Desc    The requests is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/13
 * @Listen  MIT
 */
package requests

import (
	"context"
	"fmt"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/models/users"
	"http-api/pkg/model"
	"regexp"
)

func ValidateCreateCodeRequest(ctx context.Context, input graphModel.GetCodeForForgetPasswordInput) error {
	_, err := regexp.Match("^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199)\\d{8}$", []byte(input.Phone) )
	if err != nil {
		return fmt.Errorf("手机号为: %s 不是正确的手机号", input.Phone)
	}
	userItem := users.Users{}
	err = model.DB.Model(&userItem).Where("phone = ?", &input.Phone).Find(&userItem).Error
	if err != nil {
		return fmt.Errorf("没有手机号为: %s 的用户", input.Phone)
	}

	return nil
}