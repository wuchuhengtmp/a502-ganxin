/**
 * @Desc    仓库模型
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/22
 * @Listen  MIT
 */
package repositories

import "gorm.io/gorm"

type Repositories struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name" gorm:"comment:仓库名"`
	PinYin    string  `json:"pinYin" gorm:"comment:拼音"`
	City      string  `json:"city" gorm:"comment:城市"`
	Address   string  `json:"address" gorm:"comment:地址"`
	Uid       int64   `json:"uid" gorm:"comment:管理员id"`
	Total     int64   `json:"total" gorm:"comment:总量(根)"`
	Weight    float64 `json:"weight" gorm:"comment:重量(t/吨)"`
	Remark    string  `json:"remark" gorm:"comment:备注"`
	State     bool    `json:"state" gorm:"comment:是否启用"`
	CompanyId int64   `json:"companyId" gorm:"comment:所属的公司id"`
	gorm.Model
}

