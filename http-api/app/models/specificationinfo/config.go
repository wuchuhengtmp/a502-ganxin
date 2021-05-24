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
	"gorm.io/gorm"
)

type SpecificationInfo struct {
	ID        int64   `json:"id"`
	Type      string  `json:"type" gorm:"comment:类型"`
	Length    float64 `json:"length" gorm:"comment:长度(m/米)"`
	Weight    float64 `json:"weight" gorm:"comment:重量(t/吨)"`
	IsDefault bool    `json:"isDefault" gorm:"comment:是否默认"`
	gorm.Model
}
