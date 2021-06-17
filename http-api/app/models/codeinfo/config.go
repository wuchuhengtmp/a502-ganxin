/**
 * @Desc    其它的码表信息模型
 * @Author  wuchuheng<root@wuchuheng.com>
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
			if err := TryMaterialManufactureDefault(ctx, tx); err != nil {
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
func TryMaterialManufactureDefault(ctx context.Context, tx *gorm.DB) (err error) {
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
			if err := TryMaterialManufactureDefault(ctx, tx); err != nil {
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
			if err := TryMaterialManufactureDefault(ctx, tx); err != nil {
				return err
			}
		}

		return nil
	})
}

// 创建制造商
func (c *CodeInfo) CreateManufacturerSelf(ctx context.Context) error {
	return model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(c).Error; err != nil {
			return err
		}
		if err := c.SetManufactureDefault(tx); err != nil {
			return err
		}
		me := auth.GetUser(ctx)
		l := logs.Logos{
			Uid:     me.ID,
			Content: fmt.Sprintf("添加制造商:id为%d", c.ID),
			Type:    logs.CreateActionType,
		}
		if err := tx.Create(&l).Error; err != nil {
			return err
		}

		return nil
	})
}

// 设置制造商默认选项
func (c *CodeInfo) SetManufactureDefault(tx *gorm.DB) error {
	// 设置默认选项
	if c.IsDefault {
		err := tx.Model(&CodeInfo{}).Where("company_id = ? AND id != ? AND type = ?", c.CompanyId, c.ID, Manufacturer).Update("is_default", false).Error
		if err != nil {
			return err
		}
	} else {
		hasDefault := CodeInfo{}
		// 没有默认选项并且有多条数据，则指定一条为默认
		if err := tx.Model(&CodeInfo{}).Where("company_id = ? AND type = ? AND is_default = ?", c.CompanyId, Manufacturer, true).First(&hasDefault).Error; err != nil {
			var ms []*CodeInfo
			tx.Model(&CodeInfo{}).Where("type = ? AND is_default = ? AND company_id = ?", Manufacturer, false, c.CompanyId).Find(&ms)
			if len(ms) > 0 {
				m := ms[0]
				if m.ID == c.ID {
					m = c
				}
				m.IsDefault = true
				if err := tx.Model(&CodeInfo{}).Where("id = ?", m.ID).Update("is_default", m.IsDefault).Error; err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (c *CodeInfo) GetManufacturers(ctx context.Context) (cs []*CodeInfo, err error) {
	me := auth.GetUser(ctx)
	err = model.DB.Model(&CodeInfo{}).Where("type = ? AND company_id = ?", Manufacturer, me.CompanyId).Find(&cs).Error

	return cs, err
}

/**
 *  编辑制造商
 */
func (c *CodeInfo) EditManufacturer(ctx context.Context) error {
	return model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(c).Where("id = ?", c.ID).Updates(c).Error; err != nil {
			return err
		}
		if err := c.SetManufactureDefault(tx); err != nil {
			return err
		}
		me := auth.GetUser(ctx)
		l := logs.Logos{
			Uid:     me.ID,
			Content: fmt.Sprintf("编辑制作商:修改的id为%d", c.ID),
			Type:    logs.UpdateActionType,
		}
		if err := tx.Create(&l).Error; err != nil {
			return err
		}

		return nil
	})
}

/**
 * 删除制造商
 */
func (c *CodeInfo) DeleteManufacturer(ctx context.Context) error {
	_ = c.GetSelf()
	return model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&CodeInfo{}).Where("id = ?", c.ID).Delete(&CodeInfo{}).Error; err != nil {
			return err
		}
		// 尝试指定一个新的默认项
		if c.IsDefault {
			var cs []*CodeInfo
			tx.Model(&CodeInfo{}).Where("type = ? AND company_id = ?", Manufacturer, c.CompanyId).Find(&cs)
			if len(cs) > 0 {
				if err := tx.Model(&CodeInfo{}).Where("id = ?", cs[0].ID).Update("is_default", true).Error; err != nil {
					return err
				}
			}
		}
		me := auth.GetUser(ctx)
		l := logs.Logos{
			Uid:     me.ID,
			Content: fmt.Sprintf("删除制造商: id为%d", c.ID),
			Type:    logs.CreateActionType,
		}
		if err := tx.Create(&l).Error; err != nil {
			return err
		}

		return nil
	})
}

/**
 * 创建新的物流商
 */
func (c *CodeInfo) CreateExpress(ctx context.Context) error {
	return model.DB.Transaction(func(tx *gorm.DB) error {
		me := auth.GetUser(ctx)
		c.CompanyId = me.CompanyId
		c.Type = ExpressCompany
		if err := tx.Create(c).Error; err != nil {
			return err
		}
		if err := c.SetDefaultExpress(tx, ctx); err != nil {
			return err
		}
		l := logs.Logos{
			Uid:     me.ID,
			Content: fmt.Sprintf("创建物流公司:id为%d", c.ID),
			Type:    logs.CreateActionType,
		}
		if err := tx.Create(&l).Error; err != nil {
			return err
		}

		return nil
	})
}

func (c *CodeInfo) SetDefaultExpress(tx *gorm.DB, ctx context.Context) error {
	me := auth.GetUser(ctx)
	if c.IsDefault {
		tx.Model(&CodeInfo{}).
			Where("type = ? AND company_id = ? AND id != ?", ExpressCompany, me.CompanyId, c.ID).
			Update("is_default", false)
	} else {
		// 没有默认就指定一个默认的
		hasDefault := CodeInfo{}
		err := tx.Model(&CodeInfo{}).
			Where("type = ? AND company_id = ? AND is_default = ?", ExpressCompany, me.CompanyId, true).
			First(&hasDefault).Error
		if err != nil {
			var cs []*CodeInfo
			tx.Model(&CodeInfo{}).Where("type = ? AND company_id = ?", ExpressCompany, me.CompanyId).Find(&cs)
			if len(cs) > 0 {
				if err := tx.Model(&CodeInfo{}).Where("id = ?", cs[0].ID).Update("is_default", true).Error; err != nil {
					return err
				}
			}
		}
	}

	return nil
}

/**
 * 获取物流列表
 */
func (c *CodeInfo) GetExpressList(ctx context.Context) (cs []*CodeInfo, err error) {
	me := auth.GetUser(ctx)
	model.DB.Model(&CodeInfo{}).
		Where("company_id = ? AND type = ?", me.CompanyId, ExpressCompany).
		Find(&cs)

	return
}

/**
 * 编辑物流公司
 */
func (c *CodeInfo) EditExpress(ctx context.Context) error {
	return model.DB.Transaction(func(tx *gorm.DB) error {
		me := auth.GetUser(ctx)
		c.CompanyId = me.CompanyId
		copayCodeInfo := CodeInfo{
			Type:      c.Type,
			Name:      c.Name,
			IsDefault: c.IsDefault,
			Remark:    c.Remark,
			CompanyId: c.CompanyId,
		}
		if err := tx.Model(&CodeInfo{}).Where("id = ?", c.ID).Updates(&copayCodeInfo).Error; err != nil {
			return err
		}
		if err := c.SetDefaultExpress(tx, ctx); err != nil {
			return nil
		}
		l := logs.Logos{
			Uid:     me.ID,
			Content: fmt.Sprintf("编辑物流公司:id为%d", c.ID),
			Type:    logs.UpdateActionType,
		}
		if err := tx.Create(&l).Error; err != nil {
			return err
		}

		return nil
	})
}

func (c *CodeInfo) DeleteExpress(ctx context.Context) error {
	return model.DB.Transaction(func(tx *gorm.DB) error {
		_ = c.GetSelf()
		if err := tx.Model(&CodeInfo{}).Where("id = ?", c.ID).Delete(c).Error; err != nil {
			return err
		}
		me := auth.GetUser(ctx)
		// 指定一个新的默认
		if c.IsDefault {
			var cs []*CodeInfo
			tx.Model(&CodeInfo{}).
				Where("company_id = ? AND type = ? AND id != ?", me.CompanyId, ExpressCompany, c.ID).
				Find(&cs)
			if len(cs) > 0 {
				err := tx.Model(&CodeInfo{}).Where("id = ?", cs[0].ID).Update("is_default", true).Error
				if err != nil {
					return err
				}
			}
		}
		l := logs.Logos{
			Uid:     me.ID,
			Content: fmt.Sprintf("删除物流商:id为:%d", c.ID),
			Type:    logs.DeleteActionType,
		}
		if err := tx.Create(&l).Error; err != nil {
			return err
		}

		return nil
	})
}
