/**
 * @Desc    The roles is part of http-api
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/24
 * @Listen  MIT
 */
package roles

import "gorm.io/gorm"

type Roles struct {
	ID int64	`json:"id"`
	Name string `json:"name" gorm:"comment:角色名"`
	Tag string 	`json:"tag" gorm:"comment:角色标识"`
	gorm.Model
}
