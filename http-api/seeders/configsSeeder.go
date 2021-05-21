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
			return CreateConfig(db, 1, "MINI_WECHAT_APP_ID",  "wxb0e720433455d030")
		},
	},
	seed.Seed{
		Name: "create config",
		Run: func(db *gorm.DB) error {
			return CreateConfig(db, 2, "MINI_WECHAT_APP_SECRET",  "df9dfaf44cc589ff51a886f0b288e9ba")
		},
	},
	seed.Seed{
		Name: "create config",
		Run: func(db *gorm.DB) error {
			return CreateConfig(db, 3, "FREE_DAYS",  "3")
		},
	},
	seed.Seed{
		Name: "create config",
		Run: func(db *gorm.DB) error {
			return CreateConfig(db, 4, "ADDRESS",  "公司地址")
		},
	},
	seed.Seed{
		Name: "create config",
		Run: func(db *gorm.DB) error {
			return CreateConfig(db, 5, "CONTACT",  "13420000000000")
		},
	},
	seed.Seed{
		Name: "create config",
		Run: func(db *gorm.DB) error {
			return CreateConfig(db, 6, "CONTACT",  "13420000000000")
		},
	},
	seed.Seed{
		Name: "create config",
		Run: func(db *gorm.DB) error {
			return CreateConfig(db, 7, "USER_AGREEMENT",  "用户协议")
		},
	},
	seed.Seed{
		Name: "create config",
		Run: func(db *gorm.DB) error {
			return CreateConfig(db, 8, "PRIVACY_AGREEMENT",  "隐私协议")
		},
	},
	seed.Seed{
		Name: "create config",
		Run: func(db *gorm.DB) error {
			return CreateConfig(db, 9, "APP_ICON",  "http://storage.360buyimg.com/mtd/home/32443566_635798770100444_2113947400891531264_n1533825816008.jpg")
		},
	},
	seed.Seed{
		Name: "create config",
		Run: func(db *gorm.DB) error {
			return CreateConfig(db, 10, "APP_NAME",  "泥头车小程序")
		},
	},
}

func CreateConfig(db *gorm.DB, id int64,  name string, value string)  error {
	return db.Create(&configs.Configs{ID: id, Name: name, Value: value}).Error
}
