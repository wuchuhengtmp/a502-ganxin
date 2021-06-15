/**
 * @Desc    获取仓库对应管理员模型
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/15
 * @Listen  MIT
 */
package repository_leader

import "gorm.io/gorm"

type RepositoryLeader struct {
	ID           int64 `json:"id" sql:"unique_index"`
	RepositoryId int64 `json:"repositoryId" gorm:"comment:仓库id"`
	Uid          int64 `json:"uid" gorm:"comment:用户id"`
	gorm.Model
}
