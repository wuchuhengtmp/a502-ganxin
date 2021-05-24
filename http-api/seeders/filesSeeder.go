/**
 * @Desc    文件种子数据
 * @Author  wuchuheng<wuchuheng@163.com>
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
}

func CreateFile(
	db *gorm.DB,
	id int64,
	path string,
	disk string,
)  error {
	return db.Create(&files.Files{
		ID: id,
		Path: path,
		Disk: disk,
	}).Error
}
