/**
 * @Desc    The mutation_resolver is part of http-api
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/22
 * @Listen  MIT
 */
package mutation_resolver

import (
	"context"
	"errors"
	"http-api/app/http/graph/model"
	"http-api/app/models/users"
	"http-api/pkg/helper"
	"http-api/pkg/jwt"
	sqlModel "http-api/pkg/model"
)

/**
 * 登录
 */
func (r *MutationResolver) Login (ctx context.Context, phone string, password string, mac *string) (*model.LoginRes, error)   {
	// todo 这里要验证当前用户是否有能使用这个mac地址对应的设备
	sqlDB := sqlModel.DB
	user := users.Users{}
	err := sqlDB.Where("phone=? AND password=?", phone, helper.GetHashByStr(password)).First(&user).Error
	if err != nil {
		err = errors.New("没有这个账号或密码错误")
		return &model.LoginRes{ }, err
	} else {
		accessToken, _ := jwt.GenerateTokenByUID(user.ID)
		expired := jwt.GetExpiredAt()
		role, _ := user.GetRole()
		return &model.LoginRes{
			AccessToken: accessToken,
			Expired: expired,
			Role: role.Tag,
		}, nil
	}
}
