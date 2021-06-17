/**
 * @Desc    项目管理员模型
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/15
 * @Listen  MIT
 */
package project_leader

import "gorm.io/gorm"

type ProjectLeader struct {
	ID           int64 `json:"id" sql:"unique_index"`
	ProjectId int64 `json:"projectId" gorm:"comment:项目id"`
	Uid          int64 `json:"uid" gorm:"comment:用户id"`
	gorm.Model
}

/*
 * 定义表名，用于那些联表查询需要直接使用表名等情况
 */
func (ProjectLeader)TableName() string {
	return "project_leader"
}
