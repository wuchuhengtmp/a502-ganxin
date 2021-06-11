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
	"http-api/app/models/specificationinfo"
	"http-api/pkg/seed"
)

var specificationinfoSeeds =  []seed.Seed{
	seed.Seed{
		Name: "create specificationinfo",
		Run: func(db *gorm.DB) error {
			return createSpecificationInfo(db, 1, "2H500×200×10×16", 0.1, 0.016, true, CompanyId)
		},
	},
	seed.Seed{
		Name: "create specificationinfo",
		Run: func(db *gorm.DB) error {
			return createSpecificationInfo(db, 2, "2H500×200×10×16", 0.1, 0.016, false, CompanyId)
		},
	},
}

func createSpecificationInfo(db *gorm.DB, id int64, codeType string, weight float64, length float64, isDefault bool, companyId int64)  error {
	return db.Create(&specificationinfo.SpecificationInfo{
		ID:        id,
		Type:      codeType,
		Length:    length,
		Weight:    weight,
		IsDefault: isDefault,
		CompanyId: companyId,
	}).Error
}
