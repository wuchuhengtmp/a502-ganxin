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
	"http-api/app/models/configs"
	"http-api/pkg/seed"
)

var configsSeeders =  []seed.Seed{
	seed.Seed{
		Name: "create config",
		Run: func(db *gorm.DB) error {
			return CreateConfig(db, 1, "PRICE", "1", "型钢单价", CompanyId)
		},
	},
}

func CreateConfig(db *gorm.DB, id int64,  name string, value string, remark string, companyId int64)  error {
	return db.Create(&configs.Configs{ID: id, Name: name, Value: value, Remark: remark, CompanyId: companyId}).Error
}
