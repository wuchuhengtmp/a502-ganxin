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
	"http-api/app/models/roles"
	"http-api/app/models/users"
	"http-api/pkg/helper"
	"http-api/pkg/seed"
)

type AccountType struct {
	Password string
	Username string
}

var SuperAdmin = AccountType{
	Password: "12345678",
	Username: "13427969604",
}
var CompanyAdmin = AccountType{
	Password: "12345678",
	Username: "13427969605",
}
var RepositoryAdmin = AccountType{
	Password: "12345678",
	Username: "13427969606",
}
var ProjectAdmin = AccountType{
	Password: "12345678",
	Username: "13427969607",
}
var MaintenanceAdmin = AccountType{
	Password: "12345678",
	Username: "13427969608",
}

var UsersSeeders = []seed.Seed{
	seed.Seed{
		Name: "create user",
		Run: func(db *gorm.DB) error {
			return CreateUser(
				db,
				1,
				"吴楚衡",
				helper.GetHashByStr(SuperAdmin.Password),
				SuperAdmin.Username,
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
				helper.GetHashByStr(CompanyAdmin.Password),
				CompanyAdmin.Username,
				2,
				4,
				2,
			)
		},
	},
	seed.Seed{
		Name: "create user",
		Run: func(db *gorm.DB) error {
			return CreateUser(
				db,
				3,
				"仓库管理员1",
				helper.GetHashByStr(RepositoryAdmin.Password),
				RepositoryAdmin.Username,
				roles.RoleRepositoryAdminId,
				4,
				2,
			)
		},
	},
	seed.Seed{
		Name: "create user",
		Run: func(db *gorm.DB) error {
			return CreateUser(
				db,
				4,
				"项目管理员1",
				helper.GetHashByStr(ProjectAdmin.Password),
				ProjectAdmin.Username,
				roles.RoleProjectAdminId,
				4,
				2,
			)
		},
	},
	seed.Seed{
		Name: "create user",
		Run: func(db *gorm.DB) error {
			return CreateUser(
				db,
				5,
				"维修管理员1",
				helper.GetHashByStr(MaintenanceAdmin.Password),
				MaintenanceAdmin.Username,
				roles.RoleMaintenanceAdminId,
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
	roleId int64,
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
		IsAble:       true,
	}).Error
}
