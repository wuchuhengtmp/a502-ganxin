/**
 * @Desc    规格尺寸模型
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/24
 * @Listen  MIT
 */
package specificationinfo

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/model"
	"http-api/app/models/logs"
	sqlModel "http-api/pkg/model"
)

type SpecificationInfo struct {
	ID        int64   `json:"id" sql:"unique_index"`
	Type      string  `json:"type" gorm:"comment:类型"`
	Length    float64 `json:"length" gorm:"comment:长度(m/米)"`
	Weight    float64 `json:"weight" gorm:"comment:重量(t/吨)"`
	IsDefault bool    `json:"isDefault" gorm:"comment:是否默认"`
	CompanyId int64   `json:"companyId" gorm:"comment:所属的公司id"`
	gorm.Model
}

/**
 * 添加一条新的规格记录
 */
func (s *SpecificationInfo) CreateSelf(ctx context.Context) error {
	return sqlModel.DB.Transaction(func(tx *gorm.DB) error {
		me := auth.GetUser(ctx)
		s.CompanyId = me.CompanyId
		if err := tx.Create(s).Error; err != nil {
			return err
		}
		if s.IsDefault {
			err := tx.Model(&SpecificationInfo{}).Where("company_id = ? AND id != ?", s.CompanyId, s.ID).Update("is_default", false).Error
			if err != nil {
				return err
			}
		} else {
			// 可能没有默认选项，尝试指定一条为默认选项
			if err := recoverDefaultByCompanyId(tx, s.CompanyId); err != nil {
				return err
			}
		}
		l := logs.Logos{
			Content: fmt.Sprintf("添加一条新的规格记录:id为%d", s.ID),
			Type:    logs.CreateActionType,
			Uid:     me.ID,
		}
		if err := tx.Create(&l).Error; err != nil {
			return err
		}

		return nil
	})
}

func (s *SpecificationInfo) GetSelf() error {
	db := sqlModel.DB

	return db.Model(s).Where("id = ?", s.ID).First(s).Error
}

// 如果没有默认规格，尝试指定一条为默认
func recoverDefaultByCompanyId(tx *gorm.DB, companyId int64) error {
	var cs []*SpecificationInfo
	tx.Model(&SpecificationInfo{}).Where("company_id = ?", companyId).Find(&cs)
	if len(cs) > 0 {
		var c SpecificationInfo
		tx.Model(&SpecificationInfo{}).Where("company_id = ? AND is_default = ? ", companyId, true).First(&c)
		// 没有公司的规格没有默认选项，则指定第一个为默认选项
		if c.ID == 0 {
			defaultCs := cs[0]
			err := tx.Model(&SpecificationInfo{}).
				Where("company_id = ? AND id = ?", defaultCs.CompanyId, defaultCs.ID).
				Update("is_default", true).
				Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// 编辑规格记录
func (s *SpecificationInfo) Edit(ctx context.Context, input model.EditSpecificationInput) error {
	return sqlModel.DB.Transaction(func(tx *gorm.DB) error {
		s.ID = input.ID
		_ = s.GetSelf()
		if err := tx.Model(&s).Where("id = ?", input.ID).First(&s).Error; err != nil {
			return err
		}
		if err := tx.Model(&s).Updates(SpecificationInfo{
			Type:      input.Type,
			Length:    input.Length,
			Weight:    input.Weight,
			IsDefault: input.IsDefault,
			CompanyId: s.CompanyId,
		}).Error; err != nil {
			return err
		}
		// 把其它标记为非默认
		if input.IsDefault {
			err := tx.Model(&SpecificationInfo{}).Where("company_id = ? AND id != ?", s.CompanyId, s.ID).Update("is_default", false).Error
			if err != nil {
				return err
			}
		} else {
			// 尝试指定一个默认
			if err := recoverDefaultByCompanyId(tx, s.CompanyId); err != nil {
				return err
			}
		}
		me := auth.GetUser(ctx)
		l := logs.Logos{
			Type:    logs.UpdateActionType,
			Content: fmt.Sprintf("编辑规格: 被修改的id为:%d", s.ID),
			Uid:     me.ID,
		}
		if err := tx.Create(&l).Error; err != nil {
			return err
		}

		return nil
	})
}

// 删除一条规格记录
func (s *SpecificationInfo)DeleteSelf(ctx context.Context) error {
	return sqlModel.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(s).Where("id = ?", s.ID).Delete(&s).Error; err != nil {
			return err
		}
		me := auth.GetUser(ctx)
		l := logs.Logos{
			Uid: me.ID,
			Type: logs.CreateActionType,
			Content: fmt.Sprintf("删除一条规格记录:id为%d", s.ID),
		}
		if err := tx.Create(&l).Error; err != nil {
			return err
		}

		return nil
	})
}

