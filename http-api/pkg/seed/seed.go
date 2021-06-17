/**
 * @Desc    The seed is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @DATE    2021/4/27
 * @Listen  MIT
 */
package seed

import "gorm.io/gorm"

type Seed struct {
	Name string
	Run func(db *gorm.DB) error
}