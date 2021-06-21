/**
 * @Desc    型钢填充数据
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/21
 * @Listen  MIT
 */
package seeders

import (
	"gorm.io/gorm"
	"http-api/app/models/steel_logs"
	"http-api/app/models/steels"
	"http-api/pkg/seed"
	"time"
)

var steelsSeeds = []seed.Seed{
	seed.Seed{
		Name: "create steel",
		Run: func(db *gorm.DB) error {
			return createSteel(db,
				1,
				"8",
				3,
				100,
				1,
				2,
				1,
				1,
				5,
				0,
				0,
				0,
				"GSM1-SJS01-000001",
				time.Now(),
			)
		},
	},
	seed.Seed{
		Name: "create steel",
		Run: func(db *gorm.DB) error {
			return createSteel(db,
				2,
				"9",
				3,
				100,
				1,
				2,
				1,
				1,
				4,
				0,
				0,
				0,
				"GSM1-SJS01-000002",
				time.Now(),
			)
		},
	},
	seed.Seed{
		Name: "create steel",
		Run: func(db *gorm.DB) error {
			return createSteel(db,
				3,
				"11",
				3,
				100,
				1,
				2,
				1,
				1,
				6,
				0,
				0,
				0,
				"GSM1-SJS01-000003",
				time.Now(),
			)
		},
	},
}

func createSteel(db *gorm.DB, id int64,
	identifier string,
	createdUid int64,
	state int64,
	specificationId int64,
	companyId int64,
	repositoryId int64,
	materialManufacturerId int64,
	manufacturerId int64,
	turnover int64,
	usageYearRate float64,
	totalUsageRate float64,
	code string,
	producedDate time.Time,
) error {
	return db.Transaction(func(tx *gorm.DB) error {
		err := db.Create(&steels.Steels{
			ID:                     id,
			Identifier:             identifier,
			CreatedUid:             createdUid,
			State:                  state,
			SpecificationId:        specificationId,
			CompanyId:              companyId,
			RepositoryId:           repositoryId,
			MaterialManufacturerId: materialManufacturerId,
			ManufacturerId:         manufacturerId,
			Turnover:               turnover,
			UsageYearRate:          usageYearRate,
			TotalUsageRate:         totalUsageRate,
			Code:                   code,
			ProducedDate:           producedDate,
		}).Error
		if err != nil {
			return err
		}
		err = db.Create(&steel_logs.SteelLog{
			Type:    steel_logs.CreateType,
			SteelId: id,
			Uid: RepositoryAdmin.Id,
		}).Error

		return nil
	})
}
