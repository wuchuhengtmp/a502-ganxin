/**
 * @Desc    项目模型
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/26
 * @Listen  MIT
 */
package projects

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"http-api/app/http/graph/auth"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/models/companies"
	"http-api/app/models/logs"
	"http-api/app/models/project_leader"
	"http-api/app/models/users"
	"http-api/pkg/model"
	"time"
)

type Projects struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name" gorm:"comment:项目名"`
	City      string    `json:"city" gorm:"comment:城市"`
	Address   string    `json:"address" gorm:"comment:地址"`
	StartedAt time.Time `json:"statedAt" gorm:"comment:项目开始时间"`
	EndedAt   time.Time `json:"endedAt" gorm:"comment:线束时间"`
	Remark    string    `json:"remark" gorm:"comment:备注"`
	CompanyId int64     `json:"companyId" gorm:"comment:所属公司id"`
	gorm.Model
}
/*
 * 定义表名，用于那些联表查询需要直接使用表名等情况
 */
func (Projects)TableName() string {
	return "projects"
}

/**
 * 创建新的项目
 */
func (p *Projects) CreateProject(ctx context.Context, input graphModel.CreateProjectInput) error {
	me := auth.GetUser(ctx)
	p.StartedAt = input.StartAt
	p.Name = input.Name
	p.City = input.City
	p.CompanyId = me.CompanyId
	p.Remark = input.Remark
	p.Address = input.Address
	return model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(p).Error; err != nil {
			return err
		}
		// 批量添加项目的管理员
		for _, uid := range input.LeaderIDS {
			user := users.Users{}
			if err := tx.Model(&users.Users{}).Where("id = ?", uid).First(&user).Error; err != nil {
				return err
			}
			leader := project_leader.ProjectLeader{
				Uid:       user.ID,
				ProjectId: p.ID,
			}
			if err := tx.Create(&leader).Error; err != nil {
				return err
			}
		}
		// 添加操作日志
		l := logs.Logos{
			Type:    logs.CreateActionType,
			Content: fmt.Sprintf("新增新项目:项目id为:%d", p.ID),
			Uid:     me.ID,
		}
		if err := tx.Create(&l).Error; err != nil {
			return err
		}

		return nil
	})
}

func (c *Projects)GetLeaderList() (userList []*users.Users, err error) {
	projectLeaderTableName := project_leader.ProjectLeader{}.TableName()
	userTableName := users.Users{}.TableName()
	err = model.DB.Debug().Model(&users.Users{}).
		Select(fmt.Sprintf("%s.*", userTableName)).
		Joins(fmt.Sprintf("join %s on %s.id = %s.uid", projectLeaderTableName, userTableName, projectLeaderTableName )).
		Where("project_id = ?", c.ID).
		Scan(&userList).Error

	return
}

func (c *Projects)GetCompany() (cm companies.Companies, err error) {
	err = model.DB.Model(&companies.Companies{}).Where("id = ?", c.CompanyId).First(&cm).Error
	return
}

