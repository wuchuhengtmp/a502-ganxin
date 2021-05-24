/**
 * @Desc    The users is part of http-api
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @DATE    2021/4/27
 * @Listen  MIT
 */
package users

import (
	"gorm.io/gorm"
)

type Users struct {
	ID           int64  `json:"id"`
	Name         string `json:"name" gorm:"comment:用户名"`
	Password     string `json:"password" gorm:"comment:密码"`
	Phone        string `json:"phone" gorm:"comment:手机号"`
	Mac          string `json:"MAC" gorm:"comment:设备mac地址"`
	RoteId       int8   `json:"roteId" gorm:"comment:角色id"`
	DeviceState  int8   `json:"deviceState" gorm:"comment:手持机状态"`
	Wechat       string `json:"wechat" gorm:"comment:微信"`
	CompanyId    int64  `json:"CompanyId" gorm:"comment:所属公司id"`
	State        bool   `json:"state" gorm:"state:启用状态"`
	AvatarFileId int64  `json:"avatar" gorm:"comment:头像文件id"`
	gorm.Model
}

func (Users) TableName() string {
	return "users"
}
