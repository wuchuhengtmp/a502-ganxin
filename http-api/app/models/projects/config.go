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
	"gorm.io/gorm"
	"time"
)

type Projects struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name" gorm:"comment:项目名"`
	City      string    `json:"city" gorm:"comment:城市"`
	LeaderUid int64     `json:"leaderUid" gorm:"comment:负责人id"`
	Address   string    `json:"address" gorm:"comment:地址"`
	StartedAt time.Time `json:"statedAt" gorm:"comment:项目开始时间"`
	EndedAt   time.Time `json:"endedAt" gorm:"comment:线束时间"`
	Remark    string    `json:"remark" gorm:"comment:备注"`
	CompanyId int64     `json:"companyId" gorm:"comment:所属公司id"`
	gorm.Model
}