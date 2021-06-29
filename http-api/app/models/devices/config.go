/**
 * @Desc    The devices is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/25
 * @Listen  MIT
 */
package devices

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"http-api/app/http/graph/auth"
	"http-api/app/models/logs"
	"http-api/app/models/users"
	"http-api/pkg/model"
)

type Device struct {
	ID     int64  `json:"id"`
	Mac    string `json:"mac" gorm:"comment:mac地址"`
	Uid    int64  `json:"uid" gorm:"comment:用户id"`
	IsAble bool   `json:"is_able" gorm:"comment:是否启用"`
	gorm.Model
}
func(Device)TableName() string {
	return "device"
}

/**
 * 获取用户信息
 */
func (d *Device) GetUser() (*users.Users, error) {
	u := users.Users{Id: d.Uid}
	err := u.GetSelfById(u.Id)

	return &u, err
}

func (d *Device) GetDeviceSelfById(id int64) error {
	return model.DB.Model(&Device{}).Where("id = ?", id).First(d).Error
}

func (d *Device) GetDeviceSelf() (*Device, error) {
	err := model.DB.Model(&Device{}).
		Where("uid = ? AND mac = ?", d.Uid, d.Mac).
		First(d).Error
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (d *Device) CreateSelf() error {
	return model.DB.Create(d).Error
}

/**
 * 获取公司的手持设备列表
 */
func (Device) GetAll(ctx context.Context) (ds []*Device, err error) {
	me := auth.GetUser(ctx)
	deviceTable := Device{}.TableName()
	ut := users.Users{}.TableName()
	err = model.DB.Model(&Device{}).
		Select(fmt.Sprintf("%s.*", deviceTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.uid", ut, ut, deviceTable)).
		Where(fmt.Sprintf("%s.company_id = ?", ut), me.CompanyId).
		Find(&ds).Error

	return
}

/**
 * 编辑设备
 */
func (d *Device) EditSelf(ctx context.Context) error {
	return model.DB.Transaction(func(tx *gorm.DB) error {
		// :xxx 不知道为什么这里把isAble 修改为false就是不行,也不报错
		if err := tx.Updates(d).Error; err != nil {
			return err
		}
		if err := tx.Model(&Device{}).Where("id = ?", d.ID).Update("is_able", d.IsAble).Error; err != nil {
			return err
		}

		me := auth.GetUser(ctx)
		l := logs.Logos{
			Uid:     me.Id,
			Content: fmt.Sprintf("编辑设备: 设备id为%d", d.ID),
			Type:    logs.UpdateActionType,
		}
		if err := tx.Create(&l).Error; err != nil {
			return err
		}

		return nil
	})
}
