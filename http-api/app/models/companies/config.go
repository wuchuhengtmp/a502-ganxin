/**
 * @Desc    公司表
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/24
 * @Listen  MIT
 */
package companies

import (
	"gorm.io/gorm"
	"http-api/app/http/graph/model"
	"http-api/app/http/graph/util/helper"
	"http-api/app/models/roles"
	"http-api/app/models/users"
	helper2 "http-api/pkg/helper"
	sqlModel "http-api/pkg/model"
	"time"
)

type Companies struct {
	ID               int64     `json:"id"`
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

/**
 *  更新公司
 */
func (Companies) Update(input model.EditCompanyInput) bool {
	// todo  这里涉及2个表的一致性操作，要用会话解决
	startAt, _ := helper.Str2Time(input.StartedAt)
	endAt, _ := helper.Str2Time(input.EndedAt)
	c := Companies{
		Name:             input.Name,
		PinYin:           input.PinYin,
		Symbol:           input.Symbol,
		LogoFileId:       int64(input.LogoFileID),
		BackgroundFileId: int64(input.BackgroundFileID),
		IsAble:           input.IsAble,
		Phone:            input.Phone,
		Wechat:           input.Wechat,
		StartedAt:        startAt,
		EndedAt:          endAt,
	}
	db := sqlModel.DB
	err := db.Model(&Companies{}).Where("id = ?", c.ID).Updates(c).Error
	if err != nil {
		return false
	}
	// 更新用户信息
	u := users.Users{}
	db.Model(&u).Where("company_id = ? AND role_id = ?", input.ID, roles.RoleCompanyAdminId).First(&u)
	u.Name = input.AdminName
	u.Phone = input.Phone
	u.Wechat = input.AdminWechat
	u.AvatarFileId = int64(input.AdminAvatarFileID)
	if input.AdminPassword != nil {
		u.Password = helper2.GetHashByStr(*input.AdminPassword)
	}
	err = db.Model(&u).Updates(users.Users{
		Name:         u.Name,
		Password:     u.Password,
		AvatarFileId: u.AvatarFileId,
		Phone:        u.Phone,
		Wechat:       u.Wechat,
	}).Error
	if err != nil {
		return false
	}

	return true
}

/**
 * 获取本公司
 */
func (c *Companies) GetSelfById(id int64) (err error) {
	db := sqlModel.DB
	err = db.Model(c).Where("id = ?", id).First(c).Error
	return err
}

/**
 * 删除一家公司
 */
func (c *Companies) DeleteById(id int64) error {
	db := sqlModel.DB
	err := db.Delete(c, id).Error
	return err
}

