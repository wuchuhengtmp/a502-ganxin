/**
 * @Desc    The logs is part of http-api
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/25
 * @Listen  MIT
 */
package logs

import "gorm.io/gorm"

type Logos struct {
	ID      int64  `json:"id"`
	Type    string `json:"type" gorm:"comment:操作类型增删改"`
	Content string `json:"content" gorm:"comment:操作内容"`
	Uid     int64  `json:"uid" gorm:"comment:用户id"`
	gorm.Model
}
