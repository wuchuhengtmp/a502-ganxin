/**
 * @Desc    The msg is part of http-api
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/26
 * @Listen  MIT
 */
package msg

import "gorm.io/gorm"

type Msg struct {
	Id      int64  `json:"id"`
	IsRead  bool   `json:"isRead" gorm:"comment:是否已读"`
	Content string `json:"content" gorm:"comment:内容"`
	Uid     int64  `json:"uid" gorm:"comment:用户id"`
	gorm.Model
}
