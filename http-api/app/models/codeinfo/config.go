/**
 * @Desc    其它的码表信息模型
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/24
 * @Listen  MIT
 */
package codeinfo

import "gorm.io/gorm"

// 材料厂商类型
const MeterialFacturer string = "MeterialFacturer"
// 制造厂商类型
const Producer string = "Producer"
// 运输公司类型
const ExpressCompany string = "ExpressCompany"

type CodeInfo struct {
	ID        int64  `json:"id"`
	Type      string `json:"type" gorm:"comment:类型"`
	Name      string `json:"length" gorm:"comment:厂商名称"`
	IsDefault bool   `json:"isDefault" gorm:"comment:是否默认"`
	Remark    string `json:"remark" gorm:"comment:备注"`
	gorm.Model
}
