/**
 * @Desc    The devices is part of http-api
 * @Author  wuchuheng<wuchuheng@163.com>
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
	IsAble bool   `json:"is_abl" gorm:"comment:是否启用"`
	gorm.Model
}

/**
 * 获取用户信息
 */
func (d *Device) GetUser() (*users.Users, error) {
	u := users.Users{ID: d.Uid}
	err := u.GetSelfById(u.ID)

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
	err = model.DB.Raw(`
		SELECT
			devices.* 
		FROM
			devices
			JOIN users ON users.id = devices.uid 
		WHERE
			users.company_id = ?
	`, me.CompanyId).Scan(&ds).Error

	return
}

/**
 * 编辑设备
 */
func (d *Device) EditSelf(ctx context.Context) error {
	return model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Device{}).Where("id = ?", d.ID).Updates(&Device{
			Mac:    d.Mac,
			Uid:    d.Uid,
			IsAble: d.IsAble,
		}).Error; err != nil {
			return err
		}
		me := auth.GetUser(ctx)
		l := logs.Logos{
			Uid:     me.ID,
			Content: fmt.Sprintf("编辑设备: 设备id为%d", d.ID),
			Type:    logs.UpdateActionType,
		}
		if err := tx.Create(&l).Error; err != nil {
			return err
		}

		return nil
	})
}
