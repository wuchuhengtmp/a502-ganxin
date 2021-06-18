/**
 * @Desc    仓库模型
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/22
 * @Listen  MIT
 */
package repositories

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"http-api/app/http/graph/auth"
	"http-api/app/models/logs"
	"http-api/app/models/repository_leader"
	"http-api/app/models/users"
	sqlModel "http-api/pkg/model"
)

type Repositories struct {
	ID        int64   `json:"id" sql:"unique_index"`
	Name      string  `json:"name" gorm:"comment:仓库名"`
	PinYin    string  `json:"pinYin" gorm:"comment:拼音"`
	City      string  `json:"city" gorm:"comment:城市"`
	Address   string  `json:"address" gorm:"comment:地址"`
	Total     int64   `json:"total" gorm:"comment:总量(根)"`
	Weight    float64 `json:"weight" gorm:"comment:重量(t/吨)"`
	Remark    string  `json:"remark" gorm:"comment:备注"`
	IsAble    bool    `json:"isAble" gorm:"comment:是否启用"`
	CompanyId int64   `json:"companyId" gorm:"comment:所属的公司id"`
	gorm.Model
}

func (Repositories) TableName() string {
	return "repositories"
}

/**
 * 获取公司下的仓库
 */
func (Repositories) GetAllRepositoryByCompanyId(CompanyId int64) ([]*Repositories, error) {
	db := sqlModel.DB
	var res []*Repositories
	if err := db.Model(&Repositories{}).Where("company_id = ?", CompanyId).Find(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

// 创建新的仓库
func (r *Repositories) CreatSelf(ctx context.Context) error {
	tx := sqlModel.DB.Begin()
	tx.Create(r)
	me := auth.GetUser(ctx)
	log := logs.Logos{
		Content: fmt.Sprintf("创建仓库: 仓库id为%d", r.ID),
		Uid:     me.Id,
		Type:    logs.CreateActionType,
	}
	tx.Create(&log)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *Repositories) GetSelf() error {
	db := sqlModel.DB

	return db.Model(r).Where("id = ?", r.ID).First(r).Error
}

/**
 * 删除一个仓库
 */
func DeleteById(ctx context.Context, id int64) error {
	tx := sqlModel.DB.Begin()
	tx.Where("id = ?", id).Delete(&Repositories{ID: id})
	me := auth.GetUser(ctx)
	l := logs.Logos{
		Uid:     me.Id,
		Content: fmt.Sprintf("删除仓库:仓库id为%d", id),
		Type:    logs.DeleteActionType,
	}
	tx.Create(&l)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *Repositories) IsExists(ctx context.Context) error {
	me := auth.GetUser(ctx)
	err := sqlModel.DB.Model(&Repositories{}).Where("id = ? AND company_id = ?", r.ID, me.CompanyId).First(&r).Error

	return err
}

/**
 * 获取仓库负责人列表
 */
func (r *Repositories) GetLeaders() (userList []*users.Users, err error) {
	userTable := users.Users{}.TableName()
	repositoryLeaderTable := repository_leader.RepositoryLeader{}.TableName()
	err = sqlModel.DB.Model(&users.Users{}).
		Select(fmt.Sprintf("%s.*", userTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.uid", repositoryLeaderTable, userTable, repositoryLeaderTable)).
		Where(fmt.Sprintf("%s.repository_id = %d", repositoryLeaderTable, r.ID)).
		Find(&userList).Error
	return
}

/**
 * 获取仓库负责人列表
 */
func GetLeaders(tx *gorm.DB, repositoryId int64) (userList []*users.Users, err error) {
	userTable := users.Users{}.TableName()
	repositoryLeaderTable := repository_leader.RepositoryLeader{}.TableName()
	err = tx.Model(&users.Users{}).
		Select(fmt.Sprintf("%s.*", userTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.uid", repositoryLeaderTable, userTable, repositoryLeaderTable)).
		Where(fmt.Sprintf("%s.repository_id = %d", repositoryLeaderTable, repositoryId)).
		Find(&userList).Error
	return
}
