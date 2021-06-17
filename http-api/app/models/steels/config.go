/**
 * @Desc    The steels is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/25
 * @Listen  MIT
 */
package steels

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"http-api/app/http/graph/auth"
	"http-api/app/models/codeinfo"
	"http-api/app/models/companies"
	"http-api/app/models/repositories"
	"http-api/app/models/specificationinfo"
	"http-api/app/models/steel_logs"
	"http-api/pkg/model"
	"time"
)

type Steels struct {
	ID         int64  `json:"id"`
	Identifier string `json:"identifier" gorm:"comment:识别码"`
	CreatedUid int64  `json:"createdUid" gorm:"comment:首次入库用户id"`
	State int64 `json:"state" gorm:"comment:
100【仓库】-在库
101【仓库】-运送至项目途中
102【仓库】-运送至维修厂途中
200【项目】-待使用
201【项目】-使用中
202【项目】-异常
203【项目】-闲置
204【项目】-准备归库
205【项目】-归库途中
300【维修】-待维修
301【维修】-维修中
302【维修】-准备归库
303【维修】-归库途中
400丢失
500报废
`
	SpecificationId        int64     `json:"specificationId" gorm:"comment:规格表id"`
	CompanyId              int64     `json:"companyId" gorm:"comment:所属的公司id"`
	RepositoryId           int64     `json:"repositoryId" gorm:"comment:当前存放的仓库id"`
	MaterialManufacturerId int64     `json:"materialManufacturerId" gorm:"comment:code表的材料商类型id"`
	ManufacturerId         int64     `json:"manufacturerId" gorm:"comment:code表的制造商id"`
	Turnover               int64     `json:"turnover" gorm:"comment:周转次数"`
	UsageYearRate          float64   `json:"usageYearRate" gorm:"comment:年使用率"`
	TotalUsageRate         float64   `json:"totalUsageRate" gorm:"comment:总使用率"`
	Code                   string    `json:"code gorm:comment:编码"`
	ProducedDate           time.Time `json:"producedDate" gorm:"comment:生产时间"`
	gorm.Model
}

func (Steels) TableName() string {
	return  "steels"
}

// 获取型钢的响应格式
type GetSteelListRes struct {
	Total int64     `json:"total"`
	List  []*Steels `json:"list"`
}

// 状态码声明
const (
	StateInStore                    = 100 //【仓库】-在库
	StateRepository2Project         = 101 //【仓库】-运送至项目途中
	StateRepository2Maintainer      = 102 //【仓库】-运送至维修厂途中
	StateProjectWillBeUsed          = 200 //【项目】-待使用
	StateProjectInUse               = 201 //【项目】-使用中
	StateProjectException           = 202 //【项目】-异常
	StateProjectIdle                = 203 //【项目】-闲置
	StateProjectWillBeStore         = 204 //【项目】-准备归库
	StateProjectOnTheStoreWay       = 205 //【项目】-归库途中
	StateMaintainerWillBeMaintained = 300 //【维修】-待维修
	StateMaintainerBeMaintaining    = 301 //【维修】-维修中
	StateMaintainerWillBeStore      = 302 //【维修】-准备归库
	StateMaintainerOnTheStoreWay    = 303 //【维修】-归库途中
	StateLost                       = 400 //丢失
	StateScrap                      = 500 //报废
)

/**
 * 根据规格id获取型钢
 */
func (*Steels) GetSteelsBySpecificationId(specificationId int64) (res []*Steels, err error) {
	db := model.DB
	err = db.Model(&Steels{}).Where("specification_id = ?", specificationId).Find(&res).Error
	if err != nil {
		return nil, err
	}

	return
}

/**
 * 通过材料商家id获取数据
 */
func (s *Steels) GetListByMMID(MMID int64) (ss []*Steels, err error) {
	db := model.DB
	err = db.Model(s).Where("material_manufacturer_id = ?", MMID).Find(&ss).Error

	return ss, err
}

/**
 * 通过制造商家id获取数据
 */
func (s *Steels) GetListByManufacturerId(manufacturerId int64) (ss []*Steels, err error) {
	db := model.DB
	err = db.Model(s).Where("manufacturer_id = ?", manufacturerId).Find(&ss).Error

	return ss, err
}

/**
 * 批量入库
 */
func (s *Steels) CreateMultipleSteel(ctx context.Context, steels []*Steels) error {
	me := auth.GetUser(ctx)
	c := companies.Companies{}
	if err := c.GetSelfById(me.CompanyId); err != nil {
		return err
	}
	r := repositories.Repositories{ID: steels[0].RepositoryId}
	if err := r.GetSelf(); err != nil {
		return err
	}
	return model.DB.Transaction(func(tx *gorm.DB) error {
		for _, steel := range steels {
			if err := tx.Create(&steel).Error; err != nil {
				return nil
			}
			code := fmt.Sprintf("%s-%s%.2d-%.6d", c.PinYin, r.PinYin, r.ID, steel.ID)
			if err := tx.Model(&Steels{}).Where("id = ?", steel.ID).Update("code", code).Error; err != nil {
				return err
			}
			// 型钢入库日志
			ll := steel_logs.SteelLog{
				Type:    steel_logs.CreateType,
				SteelId: steel.ID,
				Uid:     me.ID,
			}
			if err := tx.Create(&ll).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func IsExistIdentifier(ctx context.Context, identifier string) bool {
	me := auth.GetUser(ctx)
	var ss Steels
	model.DB.Model(&Steels{}).Where("identifier = ? AND company_id = ? ", identifier, me.CompanyId).Find(&ss)
	if ss.ID == 0 {
		return false
	} else {
		return true
	}
}

func (Steels) GetPagination(ctx context.Context, page int64, pageSize int64, repositoryId *int64, specificationInfoId *int64) (ss []*Steels, err error) {
	offset := 0
	if pageSize > 1 {
		offset = int((page - 1) * pageSize)
	}
	me := auth.GetUser(ctx)
	whereMap := fmt.Sprintf("company_id = %d", me.CompanyId)
	if repositoryId != nil {
		whereMap = fmt.Sprintf("%s AND repository_id = %d", whereMap, *repositoryId)
	}
	if specificationInfoId != nil {
		whereMap = fmt.Sprintf("%s AND specification_id = %d", whereMap, *specificationInfoId)
	}
	// todo 总使用率 年使用率 Turnover周围次数
	err = model.DB.Model(&Steels{}).Where(whereMap).Offset(offset).Limit(int(pageSize)).Find(&ss).Error

	return
}

func (Steels) GetTotal(ctx context.Context, repositoryId *int64, specificationInfoId *int64) (total int64) {
	me := auth.GetUser(ctx)
	whereMap := fmt.Sprintf("company_id = %d", me.CompanyId)
	if repositoryId != nil {
		whereMap = fmt.Sprintf("%s AND repository_id = %d", whereMap, *repositoryId)
	}
	if specificationInfoId != nil {
		whereMap = fmt.Sprintf("%s AND specification_id = %d", whereMap, *specificationInfoId)
	}
	model.DB.Model(&Steels{}).Where(whereMap).Count(&total)

	return
}
func (s *Steels) GetSpecification() (*specificationinfo.SpecificationInfo, error) {
	sp := specificationinfo.SpecificationInfo{}
	err := model.DB.
		Model(&specificationinfo.SpecificationInfo{}).
		Where("id = ?", s.SpecificationId).
		First(&sp).Error
	if err != nil {
		return nil, err
	}

	return &sp, nil
}
func (s *Steels) GetMaterialManufacturer() (*codeinfo.CodeInfo, error) {
	c := codeinfo.CodeInfo{}
	err := model.DB.Model(&codeinfo.CodeInfo{}).
		Where("type = ? AND id = ?", codeinfo.MaterialManufacturer, s.MaterialManufacturerId).
		First(&c).Error

	return &c, err
}

func (s *Steels) GetManufacturer() (*codeinfo.CodeInfo, error) {
	c := codeinfo.CodeInfo{}
	err := model.DB.Model(&codeinfo.CodeInfo{}).
		Where("type = ? AND id = ?", codeinfo.Manufacturer, s.ManufacturerId).
		First(&c).Error

	return &c, err
}

func (s Steels) GetRepository() (*repositories.Repositories, error) {
	r := repositories.Repositories{}
	err := model.DB.Model(repositories.Repositories{}).Where("id = ?", s.RepositoryId).
		First(&r).Error

	return &r, err
}
