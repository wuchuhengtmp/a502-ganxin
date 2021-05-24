/**
 * @Desc    角色的种子数据
 * @Author  wuchuheng<wuchuheng@163.com>
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
			return createRole(db,
				1,
				"管理员",
					"admin",
				)
		},
	},
}

func createRole(db *gorm.DB, id int64, name string, tag string)  error {
	return db.Create(&roles.Roles{
		ID: id,
		Name: name,
		Tag: tag,
	}).Error
}
