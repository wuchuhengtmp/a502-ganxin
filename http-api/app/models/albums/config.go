/**
 * @Desc    The albums is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @DATE    2021/4/27
 * @Listen  MIT
 */
package albums

import "gorm.io/gorm"

type Albums struct {
	ID   int64  `json:"id"`
	Path string `json:"path" gorm:"comment:路径"`
	Disk string `json:"disk" gorm:"comment:硬盘"`
	gorm.Model
}
