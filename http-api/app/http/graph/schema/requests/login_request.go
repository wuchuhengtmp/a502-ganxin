/**
 * @Desc    登录验证器
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/9
 * @Listen  MIT
 */
package requests

import (
	"fmt"
	"http-api/app/models/devices"
	"http-api/app/models/roles"
	"http-api/app/models/users"
	"http-api/pkg/helper"
	sqlModel "http-api/pkg/model"
)

func  ValidateLoginRequest(phone string, password string, mac *string) error  {
	sqlDB := sqlModel.DB
	user := users.Users{}
	err := sqlDB.Where("phone=? AND password=?", phone, helper.GetHashByStr(password)).First(&user).Error
	if err != nil {
		return fmt.Errorf("没有这个账号或密码错误")
	}
	if !user.IsAble {
		return fmt.Errorf("您的账号已被禁用，请联系管理员解禁")
	}
	// 设备端登录检验
	if mac != nil {
		d := devices.Device{Uid: user.ID,Mac: *mac}
		recordDevice, err := d.GetDeviceSelf()
		if err == nil && recordDevice.IsAble == false {
			return fmt.Errorf("您的账号已禁止在当前设备登录，请联系管理员解禁")
		}
		// 设备端检验禁止超级管理员和公司管理员登录
		role, _ := user.GetRole()
		if role.Tag == roles.RoleAdmin {
			return fmt.Errorf("超级管理员无法在手持设备登录")
		} else if role.Tag == roles.RoleCompanyAdmin {
			return fmt.Errorf("公司管理员无法在手持设备登录")
		}
	}

	return nil
}