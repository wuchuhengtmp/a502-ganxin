/**
 * @Desc    获取仓库对应管理员模型
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/15
 * @Listen  MIT
 */
package repository_leader

import (
	"fmt"
	"gorm.io/gorm"
	"http-api/app/models/users"
)

type RepositoryLeader struct {
	ID           int64 `json:"id" sql:"unique_index"`
	RepositoryId int64 `json:"repositoryId" gorm:"comment:仓库id"`
	Uid          int64 `json:"uid" gorm:"comment:用户id"`
	gorm.Model
}

func (RepositoryLeader) TableName() string {
	return "repository_leader"
}

/**
 * 获取负责人列表
 */
func (RepositoryLeader) GetLeaders(tx *gorm.DB) (userList []*users.Users, err error) {
	rt := RepositoryLeader{}.TableName()
	ut := users.Users{}.TableName()
	err = tx.Model(&users.Users{}).Select(fmt.Sprintf("%s.*", ut)).
		Joins(fmt.Sprintf("join %s ON %s.uid = %s.id", rt, rt, ut)).
		Scan(&userList).Error

	return
}
