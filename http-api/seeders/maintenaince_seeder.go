/**
 * @Desc
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/3
 * @Listen  MIT
 */
package seeders

import (
	"gorm.io/gorm"
	"http-api/app/models/maintenance"
	"http-api/pkg/seed"
)

var MaintenanceSeeder = []seed.Seed{
	seed.Seed{
		Name: "create maintenance",
		Run: func(db *gorm.DB) error {
			return createMaintenance(db,
				2,
				"edit name",
				"edit address",
				"",
				true,
				2,
			)
		},
	},
}

func createMaintenance(
	db *gorm.DB,
	Id int64,
	Name string,
	Address string,
	Remark string,
	IsAble bool,
	CompanyId int64,
) error {
	return db.Transaction(func(tx *gorm.DB) error {
		r := maintenance.Maintenance{
			Id:        Id,
			Name:      Name,
			Address:   Address,
			Remark:    Remark,
			IsAble:    IsAble,
			CompanyId: CompanyId,
		}
		err := tx.Create(&r).Error
		if err != nil {
			return err
		}

		return nil
	})
}
