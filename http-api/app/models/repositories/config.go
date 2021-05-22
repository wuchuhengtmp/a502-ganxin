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
	ID           int64  `json:"id"`
	Name string `json:"name" gorm:"comment:仓库名"`
	gorm.Model
}
