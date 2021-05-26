/**
 * @Desc    维修模型
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/26
 * @Listen  MIT
 */
package fix

import "gorm.io/gorm"

type Fix struct {
	Id        int64  `json:"id"`
	name      string `json:"name" gorm:"comment:维修厂名"`
	Address   string `json:"address" gorm:"comment:地址"`
	Remark    string `json:"remark" gorm:"comment:备注"`
	IsAble    bool   `json:"isAble" gorm:"comment:是否启动"`
	CompanyId int64  `json:"companyId" gorm:"comment:公司id"`
	gorm.Model
}
