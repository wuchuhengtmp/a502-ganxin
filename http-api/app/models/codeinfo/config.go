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
	"gorm.io/gorm"
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
	return model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(c).Error; err != nil {
			return err
		}
		// 更新默认配置
		if c.IsDefault {
			err := tx.Model(&CodeInfo{}).
				Where("id != ? AND type = ?", c.ID, Manufacturer).
				Update("is_default", false).Error
			if err != nil {
				return err
			}
		} else {
			// 尝试指定默认选项
			if err := TryManufactureDefault(tx); err != nil {
				return err
			}
		}

		return nil
	})
}

/**
 * 尝试设置材料商默认选项
 */
func  TryManufactureDefault(tx *gorm.DB) (err error) {
	var cs []*CodeInfo
	if err = tx.Model(&CodeInfo{}).Where("type = ?", Manufacturer).Find(&cs).Error; err != nil {
		return err
	}
	if len(cs) > 0 {
		var c CodeInfo
		tx.Model(&CodeInfo{}).Where("type = ? AND is_default = ?", Manufacturer, true).Find(&c)
		if c.ID == 0 {
			if err = tx.Model(&CodeInfo{}).Where("id = ?", c.ID).Update("is_default", true).Error; err != nil {
				return err
			}
		}
	}

	return
}