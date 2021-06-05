/**
 * @Desc    其它的码表信息模型
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/24
 * @Listen  MIT
 */
package codeinfo

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"http-api/app/http/graph/auth"
	"http-api/app/models/logs"
	"http-api/pkg/model"
)

// 材料厂商类型
const MaterialManufacturer string = "MaterialManufacturer"

// 制造厂商类型
const Manufacturer string = "Manufacturer"

// 运输公司类型
const ExpressCompany string = "ExpressCompany"

type CodeInfo struct {
	ID        int64  `json:"id"`
	Type      string `json:"type" gorm:"comment:类型"`
	Name      string `json:"length" gorm:"comment:厂商名称"`
	IsDefault bool   `json:"isDefault" gorm:"comment:是否默认"`
	Remark    string `json:"remark" gorm:"comment:备注"`
	CompanyId int64  `json:"companyId" gorm:"comment:公司id"`
	gorm.Model
}

/**
 * 添加一个新的材料商家
 */
func (c *CodeInfo) CreateMaterialManufacturer(ctx context.Context) error {
	c.Type = Manufacturer
	me := auth.GetUser(ctx)
	return model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(c).Error; err != nil {
			return err
		}
		// 更新默认配置
		if c.IsDefault {
			err := tx.Model(&CodeInfo{}).
				Where("id != ? AND type = ? AND company_id = ?", c.ID, Manufacturer, me.CompanyId).
				Update("is_default", false).Error
			if err != nil {
				return err
			}
		} else {
			// 尝试指定默认选项
			if err := TryManufactureDefault(ctx, tx); err != nil {
				return err
			}
		}
		l := logs.Logos{
			Uid:  me.ID,
			Content: fmt.Sprintf("添加材料商：id为%d", c.ID),
			Type: logs.CreateActionType,
		}
		if err := tx.Create(&l).Error; err != nil {
			return err
		}

		return nil
	})
}

/**
 * 尝试设置材料商默认选项
 */
func  TryManufactureDefault(ctx context.Context, tx *gorm.DB) (err error) {
	var cs []*CodeInfo
	me := auth.GetUser(ctx)
	if err = tx.Model(&CodeInfo{}).Where("type = ? AND company_id = ?", Manufacturer, me.CompanyId).Find(&cs).Error; err != nil {
		return err
	}
	if len(cs) > 0 {
		var c CodeInfo
		tx.Model(&CodeInfo{}).Where("type = ? AND is_default = ? AND company_id = ?", Manufacturer, true, me.CompanyId).Find(&c)
		if c.ID == 0 {
			if err = tx.Model(&CodeInfo{}).Where("id = ?", c.ID).Update("is_default", true).Error; err != nil {
				return err
			}
		}
	}

	return
}

/**
 * 获取材料商列表
 */
func (CodeInfo)GetMaterialManufacturers(ctx context.Context) (cs []*CodeInfo, err error) {
	db := model.DB
	me := auth.GetUser(ctx)
	err = db.Model(&CodeInfo{}).Where("type = ? AND company_id = ?", Manufacturer, me.CompanyId).Find(&cs).Error;

	return
}