/**
 * @Desc    日志模型
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/25
 * @Listen  MIT
 */
package logs

import (
	"gorm.io/gorm"
	"http-api/pkg/model"
)

type Logos struct {
	ID        int64      `json:"id"`
	Type      ActionType `json:"type" gorm:"comment:操作类型增删改"`
	Content   string     `json:"content" gorm:"comment:操作内容"`
	Uid       int64      `json:"uid" gorm:"comment:用户id"`
	gorm.Model
}

func (Logos) TableName() string {
	return "logos"
}

type ActionType string

const (
	DeleteActionType ActionType = "DELETE"
	UpdateActionType ActionType = "UPDATE"
	CreateActionType ActionType = "CREATE"
)

/**
 * 添加条新的记录
 */
func (l *Logos) CreateSelf() error {
	db := model.DB
	err := db.Create(l).Error

	return err
}

// 获取日志列表响应结果
type GetLogListRes struct {
	List  []*Logos `json:"list"`
	Total int64    `json:"total"`
}
