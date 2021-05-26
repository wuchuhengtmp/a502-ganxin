/**
 * @Desc    The steels is part of http-api
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/25
 * @Listen  MIT
 */
package steels

import (
	"gorm.io/gorm"
	"time"
)

type Steels struct {
	ID                     int64     `json:"id"`
	Identifier             string    `json:"identifier" gorm:"comment:识别码"`
	CreatedUid             int64     `json:"createdUid" gorm:"comment:首次入库用户id"`
	State                  int64     `json:"state" gorm:"comment:0已归库100项目-待使用101项目-使用中102项目-异常103项目—准备归库104项目—准备归库201维修-待维修202维修-维修中203维修-准备归库204维修-归库途中301废弃"`
	SpecificationInfoId    int64     `json:"specificationInfoId" gorm:"comment:规格表id"`
	CompanyId              int64     `json:"companyId" gorm:"comment:所属的公司id"`
	RepositoryId           int64     `json:"repositoryId" gorm:"comment:当前存放的仓库id"`
	MaterialManufacturerId int64     `json:"materialManufacturerId" gorm:"comment:code表的材料商类型id"`
	ManufacturerId         int64     `json:"manufacturerId" gorm:"comment:code表的制造商类型id"`
	Turnover               int64     `json:"turnover" gorm:"comment:周转次数"`
	UsageYearRate          float64   `json:"usageYearRate" gorm:"comment:年使用率"`
	TotalUsageRate         float64   `json:"totalUsageRate" gorm:"comment:总使用率"`
	ProducedDate           time.Time `json:"producedDate" gorm:"comment:生产时间"`
	gorm.Model
}
