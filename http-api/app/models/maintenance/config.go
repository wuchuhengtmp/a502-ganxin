/**
 * @Desc    维修厂模型
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/26
 * @Listen  MIT
 */
package maintenance

import "gorm.io/gorm"

type Maintenance struct {
	Id        int64  `json:"id"`
	Name      string `json:"name" gorm:"comment:维修厂名"`
	Address   string `json:"address" gorm:"comment:地址"`
	Remark    string `json:"remark" gorm:"comment:备注"`
	IsAble    bool   `json:"isAble" gorm:"comment:是否启动"`
	CompanyId int64  `json:"companyId" gorm:"comment:公司id"`
	gorm.Model
}
func(Maintenance)TableName() string {
	return "maintenance"
}



