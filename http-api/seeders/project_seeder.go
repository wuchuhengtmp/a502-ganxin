/**
 * @Desc    The seeders is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/21
 * @Listen  MIT
 */
package seeders

import (
	"gorm.io/gorm"
	"http-api/app/models/project_leader"
	"http-api/app/models/projects"
	"http-api/pkg/seed"
	"time"
)

var projectSeeder = []seed.Seed{
	seed.Seed{
		Name: "create project",
		Run: func(db *gorm.DB) error {
			return createProject(db,
				1,
				"测试项目1",
				"测试城市1",
				"测试地址1",
				time.Unix(time.Now().Unix()+7*24*60*60, 0),
				"测试项目1的备注",
				CompanyId,
			)
		},
	},
}

func createProject(
	db *gorm.DB,
	id int64,
	name string,
	city string,
	address string,
	startedAt time.Time,
	remark string,
	companyId int64,
) error {
	return db.Transaction(func(tx *gorm.DB) error {
		r := projects.Projects{
			ID:        id,
			Name:      name,
			City:      city,
			Address:   address,
			StartedAt: startedAt,
			Remark:    remark,
			CompanyId: companyId,
		}
		err := tx.Create(&r).Error
		if err != nil {
			return err
		}
		tx.Create(&project_leader.ProjectLeader{
			ProjectId: r.ID,
			Uid:       ProjectAdmin.Id,
		})

		return nil
	})
}
