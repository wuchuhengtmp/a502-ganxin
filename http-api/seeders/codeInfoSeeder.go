/**
 * @Desc    The seeders is part of http-api
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/24
 * @Listen  MIT
 */
package seeders

import (
	"gorm.io/gorm"
	"http-api/app/models/codeinfo"
	"http-api/pkg/seed"
)

var codeInfoSeeds = []seed.Seed{
	// 材料厂商
	seed.Seed{
		Name: "create codeInfo",
		Run: func(db *gorm.DB) error {
			return createCodeInfo( db, 1, codeinfo.MaterialManufacturer, "兴达工业", true, "这是备注1")
		},
	},
	seed.Seed{
		Name: "create codeInfo",
		Run: func(db *gorm.DB) error {
			return createCodeInfo( db, 2, codeinfo.MaterialManufacturer, "长洲工业", false, "这是备注2")
		},
	},
	seed.Seed{
		Name: "create codeInfo",
		Run: func(db *gorm.DB) error {
			return createCodeInfo( db, 3, codeinfo.MaterialManufacturer, " 北建工业", false, "这是备注3")
		},
	},
	// 制造厂商
	seed.Seed{
		Name: "create codeInfo",
		Run: func(db *gorm.DB) error {
			return createCodeInfo( db, 4, codeinfo.Manufacturer, "制作厂商1", true, "这是备注1")
		},
	},
	seed.Seed{
		Name: "create codeInfo",
		Run: func(db *gorm.DB) error {
			return createCodeInfo( db, 5, codeinfo.Manufacturer, "制作厂商2", false, "这是备注2")
		},
	},
	seed.Seed{
		Name: "create codeInfo",
		Run: func(db *gorm.DB) error {
			return createCodeInfo( db, 6, codeinfo.Manufacturer, "制作厂商3", false, "这是备注3")
		},
	},
	// 运输公司
	seed.Seed{
		Name: "create codeInfo",
		Run: func(db *gorm.DB) error {
			return createCodeInfo( db, 7, codeinfo.ExpressCompany, "运输公司xxx", true, "这是备注1")
		},
	},

}

func createCodeInfo(db *gorm.DB, id int64, codeInfoType string, name string, isDefault bool, remark string)  error {
	return db.Create(&codeinfo.CodeInfo{
		ID:        id,
		Type:      codeInfoType,
		Name:      name,
		IsDefault: isDefault,
		Remark:    remark,
	}).Error
}
