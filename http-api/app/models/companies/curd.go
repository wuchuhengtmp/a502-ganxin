/**
 * @Desc    公司模型curd操作
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/28
 * @Listen  MIT
 */
package companies

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"http-api/app/http/graph/auth"
	graphQL "http-api/app/http/graph/model"
	"http-api/app/http/graph/util/helper"
	"http-api/app/models/files"
	"http-api/app/models/logs"
	"http-api/app/models/roles"
	"http-api/app/models/users"
	helper2 "http-api/pkg/helper"
	"http-api/pkg/model"
	sqlModel "http-api/pkg/model"
)

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
func (c *Companies) Update(ctx context.Context, input graphQL.EditCompanyInput) error {
	startAt, _ := helper.Str2Time(input.StartedAt)
	endAt, _ := helper.Str2Time(input.EndedAt)
	c.ID = input.ID
	return model.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&Companies{}).Where("id = ?", input.ID).Updates(&Companies{
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
		}).Error
		if err != nil {
			return err
		}
		if err := tx.Model(&Companies{}).Where("id = ?", c.ID).First(&c).Error; err != nil {
			return err
		}
		// 更新用户信息
		u := users.Users{}
		tx.Model(&u).Where("company_id = ? AND role_id = ?", input.ID, roles.RoleCompanyAdminId).First(&u)
		u.Name = input.AdminName
		u.Phone = input.Phone
		u.Wechat = input.AdminWechat
		u.AvatarFileId = int64(input.AdminAvatarFileID)
		if input.AdminPassword != nil {
			u.Password = helper2.GetHashByStr(*input.AdminPassword)
		}
		err = tx.Model(&u).Updates(users.Users{
			Name:         u.Name,
			Password:     u.Password,
			AvatarFileId: u.AvatarFileId,
			Phone:        u.Phone,
			Wechat:       u.Wechat,
		}).Error
		if err != nil {
			return err
		}
		// 添加一个操作日志
		me := auth.GetUser(ctx)
		logsModel := logs.Logos{}
		logsModel.Type = logs.UpdateActionType
		logsModel.Uid = me.Id
		logsModel.Content = fmt.Sprintf("修改公司，名称: %s, ID:%d", c.Name, input.ID)
		if tx.Create(&logsModel).Error != nil {
			return err
		}

		return nil
	})
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

func GetCompanyAdminUserById(id int64) (*users.Users, error) {
	db := sqlModel.DB
	user := users.Users{}
	err := db.Model(&user).Where("company_id = ? AND role_id = ? ", id, roles.RoleCompanyAdminId).First(&user).Error

	return &user, err
}


/**
 * 获取对应解析器的公司下的员工数据
 */
func GetCompanyItems(companyId int64, input *graphQL.GetCompanyUserInput) ([]*users.Users, error) {
	var c []*users.Users
	db := sqlModel.DB
	if input != nil && input.RoleID != nil {
		db.Model(&users.Users{}).Where("company_id = ? AND role_id = ?", companyId, input.RoleID).Find(&c)
	} else {
		db.Model(&users.Users{}).Where("company_id = ?", companyId).Find(&c)
	}

	return c, nil
}

func UpdateCompanyUser(ctx context.Context, input *graphQL.EditCompanyUserInput) (*users.Users, error) {
	tx := model.DB.Begin()
	user := users.Users{}
	_ = user.GetSelfById(input.ID)
	user.Phone = input.Phone
	user.Name = input.Name
	user.RoleId = input.RoleID
	user.IsAble = input.IsAble
	if err := tx.Model(&user).Updates(&user).Error; err != nil {
		return nil, fmt.Errorf("update user failed")
	}
	log := logs.Logos{}
	log.Type = logs.UpdateActionType
	log.Content = fmt.Sprintf("编辑公司人员: 被更新的用户id为%d", user.ID)
	me := auth.GetUser(ctx)
	log.Uid = me.Id
	tx.Create(&log)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("execu the UpdateCompanyUser method was failed")
	}
	roleInfo := roles.Role{}
	roleInfo.GetSelfById(user.RoleId)
	avatarInfo := files.File{}
	_ = avatarInfo.GetSelfById(user.AvatarFileId)
	res := users.Users{
		Id: user.Id,
		Phone:  user.Phone,
		Wechat: user.Wechat,
		IsAble: user.IsAble,
	}

	return &res, nil
}

/**
 *删除公司员工
 */
func DeleteCompanyUserByUid(ctx context.Context, uid int64) error {
	tx := sqlModel.DB.Begin()
	user := users.Users{}
	_ = user.GetSelfById(uid)
	tx.Model(&users.Users{}).Where("id = ?", uid).Delete(&users.Users{})
	me := auth.GetUser(ctx)
	log := logs.Logos{}
	log.Uid = me.Id
	log.Content = fmt.Sprintf("删除用户:用户id为 %d;用户名为 %s", user.ID, user.Name)
	log.Type = logs.DeleteActionType
	tx.Create(&log)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
