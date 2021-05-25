/**
 * @Desc    The seeders is part of http-api
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/25
 * @Listen  MIT
 */
package seeders

import (
	"gorm.io/gorm"
	"http-api/app/models/devices"
	"http-api/pkg/seed"
)

var deviceSeeders =  []seed.Seed{
	seed.Seed{
		Name: "create device",
		Run: func(db *gorm.DB) error {
			return createDevice(db, 1, "20:82:c0:2d:a5:d6", 1, true)
		},
	},
	seed.Seed{
		Name: "create device",
		Run: func(db *gorm.DB) error {
			return createDevice(db, 2, "20:82:c0:2d:a5:d6", 1, false)
		},
	},
}

func createDevice(db *gorm.DB, id int64, mac string, uid int64, state bool)  error {
	return db.Create(&devices.Devices{
		ID: id,
		Mac: mac,
		Uid: uid,
		State: state,
	}).Error
}
