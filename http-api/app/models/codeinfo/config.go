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
	me := auth.GetUser(ctx)
	return model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(c).Error; err != nil {
			return err
		}
		// 更新默认配置
		if c.IsDefault {
			err := tx.Model(&CodeInfo{}).
				Where("id != ? AND type = ? AND company_id = ?", c.ID, MaterialManufacturer, me.CompanyId).
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
			Uid:     me.ID,
			Content: fmt.Sprintf("添加材料商：id为%d", c.ID),
			Type:    logs.CreateActionType,
		}
		if err := tx.Create(&l).Error; err != nil {
			return err
		}

		return nil
	})
}

/**
 * 检验是否有默认材料商，否则尝试指定一个材料商作为默认选项
 */
func TryManufactureDefault(ctx context.Context, tx *gorm.DB) (err error) {
	var cs []*CodeInfo
	me := auth.GetUser(ctx)
	if err = tx.Model(&CodeInfo{}).Where("type = ? AND company_id = ? AND is_default = ?", MaterialManufacturer, me.CompanyId, false).Find(&cs).Error; err != nil {
		return err
	}
	if len(cs) > 0 {
		var c CodeInfo
		err = tx.Model(&CodeInfo{}).Where("type = ? AND is_default = ? AND company_id = ?", MaterialManufacturer, true, me.CompanyId).First(&c).Error
		if err != nil {
			if err = tx.Model(&CodeInfo{}).Where("id = ?", cs[0].ID).Update("is_default", true).Error; err != nil {
				return err
			}
		}
	}

	return
}

/**
 * 获取材料商列表
 */
func (CodeInfo) GetMaterialManufacturers(ctx context.Context) (cs []*CodeInfo, err error) {
	db := model.DB
	me := auth.GetUser(ctx)
	err = db.Model(&CodeInfo{}).Where("type = ? AND company_id = ?", MaterialManufacturer, me.CompanyId).Find(&cs).Error

	return
}

func (c *CodeInfo) GetSelf() error {
	db := model.DB

	return db.Model(&CodeInfo{}).Where("id = ? ", c.ID).First(c).Error
}

func (c *CodeInfo) EditMaterialManufacturer(ctx context.Context) error {
	me := auth.GetUser(ctx)
	c.CompanyId = me.CompanyId
	return model.DB.Transaction(func(tx *gorm.DB) error {
		fmt.Println(c)
		if err := tx.Model(c).Where("id = ?", c.ID).Updates(c).Error; err != nil {
			return err
		}
		if err := tx.Model(c).Where("id = ?", c.ID).Update("is_default", c.IsDefault).Error; err != nil {
			return err
		}
		if c.IsDefault {
			err := tx.Model(&CodeInfo{}).Where("company_id = ? AND id != ? AND type = ?", me.CompanyId, c.ID, MaterialManufacturer).
				Update("is_default", false).Error
			if err != nil {
				return err
			}
		} else {
			if err := TryManufactureDefault(ctx, tx); err != nil {
				return err
			}
		}
		l := logs.Logos{
			Uid:     me.ID,
			Content: fmt.Sprintf("编辑材料商:被修改id为%d", c.ID),
			Type:    logs.UpdateActionType,
		}
		if err := tx.Create(&l).Error; err != nil {
			return err
		}

		return nil
	})
}

// 删除材料商
func (c *CodeInfo) DeleteMaterial(ctx context.Context) error {
	if err := c.GetSelf(); err != nil {
		return err
	}

	return model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(c).Where("id = ?", c.ID).Delete(c).Error; err != nil {
			return err
		}
		me := auth.GetUser(ctx)
		l := logs.Logos{
			Uid:     me.ID,
			Content: fmt.Sprintf("删除材料商:被删除id为 %d", c.ID),
			Type:    logs.DeleteActionType,
		}
		if err := tx.Create(&l).Error; err != nil {
			return err
		}
		// 尝试设置一个默认项
		if c.IsDefault {
			if err := TryManufactureDefault(ctx, tx); err != nil {
				return err
			}
		}

		return nil
	})
}