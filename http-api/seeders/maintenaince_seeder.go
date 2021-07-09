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
	"http-api/app/models/maintenance_leader"
	"http-api/app/models/roles"
	"http-api/app/models/users"
	"http-api/pkg/seed"
)

var MaintenanceSeeder = []seed.Seed{
	seed.Seed{
		Name: "create maintenance",
		Run: func(db *gorm.DB) error {
			return createMaintenance(db,
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
	Name string,
	Address string,
	Remark string,
	IsAble bool,
	CompanyId int64,
) error {
	return db.Transaction(func(tx *gorm.DB) error {
		r := maintenance.Maintenance{
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
		u := users.Users{}
		err = tx.Model(&u).Where("role_id = ?", roles.RoleMaintenanceAdminId).First(&u).Error
		if err != nil {
			return err
		}
		err = tx.Create(&maintenance_leader.MaintenanceLeader{
			MaintenanceId: r.Id,
			Uid:           u.Id,
		}).Error
		if err != nil {
			return err
		}

		return nil
	})
}
