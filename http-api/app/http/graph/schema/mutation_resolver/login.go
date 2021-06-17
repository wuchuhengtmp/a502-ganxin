/**
 * @Desc    The mutation_resolver is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/22
 * @Listen  MIT
 */
package mutation_resolver

import (
	"context"
	"http-api/app/http/graph/errors"
	"http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/devices"
	"http-api/app/models/users"
	"http-api/pkg/helper"
	"http-api/pkg/jwt"
	sqlModel "http-api/pkg/model"
)

/**
 * 登录
 */
func (r *MutationResolver) Login(ctx context.Context, phone string, password string, mac *string) (*model.LoginRes, error) {
	if err := requests.ValidateLoginRequest(phone, password, mac); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	user := users.Users{}
	sqlDB := sqlModel.DB
	sqlDB.Where("phone=? AND password=?", phone, helper.GetHashByStr(password)).First(&user)
	var isDevice bool
	var macAddres string
	if mac != nil && len(*mac) > 0 {
		macAddres = *mac
		isDevice = true
		d := devices.Device{Uid: user.ID, Mac: *mac}
		_, err := d.GetDeviceSelf()
		if err != nil {
			d.IsAble = true
			_ = d.CreateSelf()
		}
	}
	accessToken, _ := jwt.GenerateTokenByUID(user.ID, isDevice, macAddres)
	expired := jwt.GetExpiredAt()
	role, _ := user.GetRole()
	return &model.LoginRes{
		AccessToken: accessToken,
		Expired:     expired,
		Role:        role.Tag,
		RoleName:    role.Name,
	}, nil
}
