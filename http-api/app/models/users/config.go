/**
 * @Desc    The users is part of http-api
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @DATE    2021/4/27
 * @Listen  MIT
 */
package users

import "gorm.io/gorm"


type Users struct {
	ID           int64  `json:"id"`
	Username     string `json:"username" gorm:"comment:账号"`
	Password     string `json:"password" gorm:"comment:密码"`
	Nickname     string `json:"nickname" gorm:"comment:昵称"`
	Gender       int64  `json:"gender" gorm:"comment:性别 1男2女0不知"`
	AvatarUrl    string `json:"avatarUrl" gorm:"comment:头像"`
	WechatModel
	gorm.Model
}
