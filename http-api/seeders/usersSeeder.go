/**
 * @Desc    The seeders is part of http-api
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @DATE    2021/4/27
 * @Listen  MIT
 */
package seeders

import (
	"gorm.io/gorm"
	"http-api/app/models/users"
	"http-api/pkg/helper"
	"http-api/pkg/seed"
)

var UsersSeeders = []seed.Seed{
	seed.Seed{
		Name: "create user",
		Run: func(db *gorm.DB) error {
			return CreateUser(
				db,
				1,
				"吴楚衡",
				helper.GetHashByStr("12345678"),
				"13427969604",
				1,
				1,
				0,
			)
		},
	},
	seed.Seed{
		Name: "create user",
		Run: func(db *gorm.DB) error {
			return CreateUser(
				db,
				2,
				"公司管理员1",
				helper.GetHashByStr("12345678"),
				"13427969600",
				2,
				4,
				2,
			)
		},
	},
}

func CreateUser(
	db *gorm.DB,
	id int64,
	name string,
	password string,
	phone string,
	roleId int8,
	avatarFileId int64,
	companyId int64,
) error {
	return db.Create(&users.Users{
		ID:           id,
		Name:         name,
		Password:     password,
		Phone:        phone,
		RoleId:       roleId,
		AvatarFileId: avatarFileId,
		CompanyId:    companyId,
	}).Error
}
