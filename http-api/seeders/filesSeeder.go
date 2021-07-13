/**
 * @Desc    文件种子数据
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/24
 * @Listen  MIT
 */
package seeders

import (
	"gorm.io/gorm"
	"http-api/app/models/files"
	"http-api/pkg/seed"
)

var fileSeeders =  []seed.Seed{
	seed.Seed{
		Name: "create files",
		Run: func(db *gorm.DB) error {
			return CreateFile( db, 1, "test.png", "local")
		},
	},
	seed.Seed{
		Name: "create files",
		Run: func(db *gorm.DB) error {
			return CreateFile( db, 2, "logo.png", "local")
		},
	},
	seed.Seed{
		Name: "create files",
		Run: func(db *gorm.DB) error {
			return CreateFile( db, 3, "background.png", "local")
		},
	},
	seed.Seed{
		Name: "create files",
		Run: func(db *gorm.DB) error {
			return CreateFile( db, 4, "2021-5-28/1622171038-1622171009253.jpg", "local")
		},
	},
	seed.Seed{
		Name: "create files",
		Run: func(db *gorm.DB) error {
			return CreateFile( db, 5, "2021-5-28/1622171414-companyBackground.png", "local")
		},
	},
	seed.Seed{
		Name: "create files",
		Run: func(db *gorm.DB) error {
			return CreateFile( db, 6, "2021-5-28/1622171422-companyLogo.png", "local")
		},
	},
	seed.Seed{
		Name: "create files",
		Run: func(db *gorm.DB) error {
			return CreateFile( db, 7, "2021-7-13/1626155017-tutor.mp4", "local")
		},
	},
}

func CreateFile(
	db *gorm.DB,
	id int64,
	path string,
	disk string,
)  error {
	return db.Create(&files.File{
		ID: id,
		Path: path,
		Disk: disk,
	}).Error
}
