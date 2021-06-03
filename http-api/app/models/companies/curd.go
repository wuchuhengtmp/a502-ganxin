/**
 * @Desc    公司模型curd操作
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/28
 * @Listen  MIT
 */
package companies

import (
	"context"
	"fmt"
	"http-api/app/http/graph/auth"
	graphQL "http-api/app/http/graph/model"
	"http-api/app/http/graph/util/helper"
	"http-api/app/models/logs"
	"http-api/app/models/roles"
	"http-api/app/models/users"
	helper2 "http-api/pkg/helper"
	"http-api/pkg/model"
	sqlModel "http-api/pkg/model"
	"time"
)

func (c *Companies) Create(ctx context.Context) error {
	//  todo 这里要加入会话 保证操作的一致性
	db := model.DB
	err := db.Model(c).Create(c).Error
	me := auth.GetUser(ctx)
	logsModel := logs.Logos{}
	logsModel.Type = logs.CreateActionType
	logsModel.Uid = me.ID
	logsModel.Content = fmt.Sprintf("添加公司,名称:%s,ID: %d", c.Name, c.ID)
	_ = logsModel.CreateSelf()

	return err
}

func GetAllByUid(uid int64) (cs []Companies) {
	db := model.DB
	me := users.Users{}
	_ = me.GetSelfById(uid)
	roleModel := roles.Role{}
	_ = roleModel.GetSelfById(me.RoleId)
	//  超级管理员能查看到的数据域--就是没限制
	if roleModel.Tag == roles.RoleAdmin {
		db.Model(&Companies{}).Find(&cs)
	} else {
		// 不是超级管员查看数据进行限制
		db.Model(&Companies{}).Where("id = ?", me.CompanyId).Find(&cs)
	}
	// 公司状态根据有郊期期限修正
	for i, company := range cs {
		if company.IsAble == true && company.EndedAt.Unix() >= time.Now().Unix() {
			company.IsAble = false
			db.Model(&company).Where("id = ?", company.ID).Updates(company)
			cs[i].IsAble = false
		}
	}

	return cs
}

/**
 * 获取公司的管理员
 */
func (c *Companies) GetAdmin() (user users.Users, err error) {
	db := model.DB
	err = db.Model(&user).
		Where("role_id = ? AND company_id = ?", roles.RoleCompanyAdminId, c.ID).
		First(&user).Error
	return user, err
}

/**
 *  有没有这家公司
 */
func (Companies) HasCompanyId(cid int64) (*Companies, error) {
	db := model.DB
	cp := Companies{}
	err := db.Model(&Companies{}).Where("id = ?", cid).First(&cp).Error
	return &cp, err
}

/**
 *  更新公司
 */
func (Companies) Update(ctx context.Context, input graphQL.EditCompanyInput) bool {
	// todo  这里涉及3个表的一致性操作，要用会话解决
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
	err := db.Model(&Companies{}).Where("id = ?", input.ID).Updates(c).Error
	if err != nil {

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
	// 添加一个操作日志
	me := auth.GetUser(ctx)
	logsModel := logs.Logos{}
	logsModel.Type = logs.UpdateActionType
	logsModel.Uid = me.ID
	logsModel.Content = fmt.Sprintf("修改公司，名称: %s, ID:%d", c.Name, input.ID)
	_ = logsModel.CreateSelf()

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

