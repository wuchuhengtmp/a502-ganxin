/**
 * @Desc    仓库填充数据
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/25
 * @Listen  MIT
 */
package seeders

import (
	"gorm.io/gorm"
	"http-api/app/models/repositories"
	"http-api/app/models/repository_leader"
	"http-api/pkg/seed"
)

var repositorySeeder = []seed.Seed{
	seed.Seed{
		Name: "create repository",
		Run: func(db *gorm.DB) error {
			return createRepository(db, 1,
				"石景山仓库",
				"SJS",
				"北京",
				"石景山城通街26号院",
				1,
				700,
				5322.6,
				"",
				true,
				CompanyId,
			)
		},
	},
	seed.Seed{
		Name: "create repository",
		Run: func(db *gorm.DB) error {
			return createRepository(db, 2,
				"大兴仓库",
				"DX",
				"北京",
				"大兴生物医药基地",
				1,
				390,
				355.3,
				"",
				false,
				1,
			)
		},
	},
}

func createRepository(
	db *gorm.DB,
	id int64,
	name string,
	pinYin string,
	city string,
	address string,
	uid int64,
	total int64,
	weight float64,
	remark string,
	isAble bool,
	companyId int64,
) error {
	return db.Transaction(func(tx *gorm.DB) error {
		r := repositories.Repositories{
				ID:        id,
				Name:      name,
				PinYin:    pinYin,
				City:      city,
				Address:   address,
				Total:     total,
				Weight:    weight,
				Remark:    remark,
				IsAble:    isAble,
				CompanyId: companyId,
			}
		err := tx.Create(&r).Error
		if err != nil { return err }
		tx.Create(&repository_leader.RepositoryLeader{
			Uid: uid,
			RepositoryId: r.ID,
		})

		return nil
	})
}
