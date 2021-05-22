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
	Name 		 string `json:"name" gorm:"comment:用户名"`
	Password     string `json:"password" gorm:"comment:密码"`
	Phone 		 string `json:"phone" gorm:"comment:手机号"`
	Mac 		 string `json:"MAC" gorm:"comment:设备mac地址"`
	RoteId		 int8   `json:"roteId" gorm:"comment:角色id"`
	Company 	 string `json:"company" gorm:"comment:公司名"`
	DeviceState  int8 	`json:"deviceState" gorm:"手持机状态"`
	gorm.Model
}
