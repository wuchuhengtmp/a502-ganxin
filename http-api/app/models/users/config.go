/**
 * @Desc    The users is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @DATE    2021/4/27
 * @Listen  MIT
 */
package users

import (
	"gorm.io/gorm"
	"http-api/app/models"
	"http-api/app/models/roles"
	helper2 "http-api/pkg/helper"
	sqlModel "http-api/pkg/model"
)

type Users struct {
	ID           int64  `json:"id" sql:"unique_index"`
	Name         string `json:"name" gorm:"comment:用户名"`
	Password     string `json:"password" gorm:"comment:密码"`
	Phone        string `json:"phone" gorm:"comment:手机号"`
	RoleId       int64  `json:"roleId" gorm:"comment:角色id"`
	Wechat       string `json:"wechat" gorm:"comment:微信"`
	CompanyId    int64  `json:"CompanyId" gorm:"comment:所属公司id"`
	IsAble       bool   `json:"is_able" gorm:"comment:启用状态"`
	AvatarFileId int64  `json:"avatar" gorm:"comment:头像文件id"`
	models.Base
	gorm.Model
}

/*
 * 定义表名，用于那些联表查询需要直接使用表名等情况
 */
func (Users)TableName() string {
	return "users"
}

/**
 * 用户手机号是否存在
 */
func (Users) IsPhoneExists(phone string) bool {
	db := sqlModel.DB
	u := Users{}
	err := db.Model(&Users{}).Where("phone = ?", phone).First(&u).Error
	if err == nil {
		return true
	} else {
		return false
	}
}

/**
 * 有没有这个用户
 */
func (Users) HasUserById(id int64) (*Users, error) {
	db := sqlModel.DB
	u := Users{}
	err := db.Model(&u).Where("id = ?", id).First(&u).Error
	return &u, err
}

/**
 * 获取关联的角色
 */
func (u Users) GetRole() (roles.Role, error) {
	role := roles.Role{}
	sqlDB := sqlModel.DB
	err := sqlDB.Model(&role).Where("id = ?", u.RoleId).First(&role).Error
	return role, err
}

/**
 * 公司管理员电话新的电话比较原有的， 是否已经更改了
 */
func (Users) IsChangeCompanyAdminPhone(companyId int64, phone string) bool {
	db := sqlModel.DB
	u := Users{}
	db.Model(&Users{}).Where(
		"role_id = ? AND company_id = ?",
		roles.RoleCompanyAdminId,
		companyId,
	).First(&u)
	if u.Phone != phone {
		return false
	} else {
		return false
	}
}

/**
 * 是否是唯一的手机号
 */
func (Users) IsUniPhone(phone string) bool {
	db := sqlModel.DB
	u := Users{}
	err := db.Model(Users{}).Where("phone = ?", phone).First(&u).Error
	if err != nil {
		return false
	} else {
		return true
	}
}

func (u *Users) GetSelfById(uid int64) error {
	db := sqlModel.DB
	return db.Model(u).Where("id = ?", uid).First(u).Error
}

func (u *Users)SetSelfPassword(password string) error {
	err := sqlModel.DB.Model(&Users{}).Where("id = ?", u.ID).Update("password", helper2.GetHashByStr(password)).Error
	return err
}
