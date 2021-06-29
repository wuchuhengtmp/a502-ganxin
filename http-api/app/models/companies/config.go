/**
 * @Desc    公司模型
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/24
 * @Listen  MIT
 */
package companies

import (
	"gorm.io/gorm"
	"http-api/pkg/model"
	"time"
)

type Companies struct {
	ID               int64     `json:"id" sql:"unique_index"`
	Name             string    `json:"name" gorm:"comment:公司名"`
	PinYin           string    `json:"pinyin" gorm:"comment:用于型钢编码生成"`
	Symbol           string    `json:"symbol" gorm:"comment:APP 企业宗旨"`
	LogoFileId       int64     `json:"logoFileId" gorm:"comment:文件id"`
	BackgroundFileId int64     `json:"backgroundFileId" gorm:"comment:app背景文件id"`
	IsAble           bool      `json:"state" gorm:"comment:账号状态"`
	Phone            string    `json:"phone" gorm:"comment:公司的电话"`
	Wechat           string    `json:"wechat" gorm:"comment:公司的微信"`
	StartedAt        time.Time `json:"startedAt" gorm:"comment:开始时间"`
	EndedAt          time.Time `json:"ended_at" gorm:"comment:结束时间"`
	gorm.Model
}
func(Companies)TableName() string {
	return "companies"
}

func (Companies)GetAll() (cs []*Companies, err error) {
	err = model.DB.Model(&Companies{}).Find(&cs).Error

	return
}
