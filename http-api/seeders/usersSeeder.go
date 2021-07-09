/**
 * @Desc    The seeders is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
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
	Id       int64
}

var SuperAdmin = AccountType{
	Password: "12345678",
	Username: "13427969604",
	Id:       1,
}
var CompanyAdmin = AccountType{
	Password: "12345678",
	Username: "13427969605",
	Id:       2,
}
var RepositoryAdmin = AccountType{
	Password: "12345678",
	Username: "13427969606",
	Id:       3,
}
var ProjectAdmin = AccountType{
	Password: "12345678",
	Username: "13427969607",
	Id:       4,
}
var MaintenanceAdmin = AccountType{
	Password: "12345678",
	Username: "13427969608",
	Id:       5,
}

var UsersSeeders = []seed.Seed{
	seed.Seed{
		Name: "create user",
		Run: func(db *gorm.DB) error {
			return CreateUser(
				db,
				SuperAdmin.Id,
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
				CompanyAdmin.Id,
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
				RepositoryAdmin.Id,
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
				ProjectAdmin.Id,
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
				MaintenanceAdmin.Id,
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
	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&users.Users{
			Id:           id,
			Name:         name,
			Password:     password,
			Phone:        phone,
			RoleId:       roleId,
			AvatarFileId: avatarFileId,
			CompanyId:    companyId,
			IsAble:       true,
		}).Error
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
