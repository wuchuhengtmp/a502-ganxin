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
	"http-api/app/models"
	"http-api/app/models/roles"
	sqlModel "http-api/pkg/model"
)

type Users struct {
	ID           int64       `json:"id"`
	Name         string      `json:"name" gorm:"comment:用户名"`
	Password     string      `json:"password" gorm:"comment:密码"`
	Phone        string      `json:"phone" gorm:"comment:手机号"`
	RoteId       int8        `json:"roteId" gorm:"comment:角色id"`
	Wechat       string      `json:"wechat" gorm:"comment:微信"`
	CompanyId    int64       `json:"CompanyId" gorm:"comment:所属公司id"`
	IsAble       bool        `json:"is_able" gorm:"comment:启用状态"`
	AvatarFileId int64       `json:"avatar" gorm:"comment:头像文件id"`
	models.Base
	gorm.Model
}

func (Users) TableName() string {
	return "users"
}

/**
 * 获取关联的角色
 */
func (u Users) GetRole() (roles.Role, error) {
	role := roles.Role{}
	sqlDB := sqlModel.DB
	err := sqlDB.Model(&role).Where("id = ?", u.RoteId).First(&role).Error
	return role, err
}
