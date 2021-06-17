/**
 * @Desc    角色的种子数据
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/24
 * @Listen  MIT
 */
package seeders

import (
	"gorm.io/gorm"
	"http-api/app/models/roles"
	"http-api/pkg/seed"
)

var rolesSeeders =  []seed.Seed{
	seed.Seed{
		Name: "create role",
		Run: func(db *gorm.DB) error {
			return createRole(db, 1, "超级管理员", "admin")
		},
	},
	seed.Seed{
		Name: "create role",
		Run: func(db *gorm.DB) error {
			return createRole(db, 2, "公司管理员", "companyAdmin")
		},
	},
	seed.Seed{
		Name: "create role",
		Run: func(db *gorm.DB) error {
			return createRole(db, 3, "仓库管理员", "repositoryAdmin")
		},
	},
	seed.Seed{
		Name: "create role",
		Run: func(db *gorm.DB) error {
			return createRole(db, 4, "项目管理员", "projectAdmin")
		},
	},
	seed.Seed{
		Name: "create role",
		Run: func(db *gorm.DB) error {
			return createRole(db, 5, "维修管理员", "maintenanceAdmin")
		},
	},
}

func createRole(db *gorm.DB, id int64, name string, tag roles.GraphqlRole)  error {
	return db.Create(&roles.Role{
		ID: id,
		Name: name,
		Tag: tag,
	}).Error
}
