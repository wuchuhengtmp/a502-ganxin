/**
 * @Desc    The devices is part of http-api
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/25
 * @Listen  MIT
 */
package devices

type Devices struct {
	ID    int64  `json:"id"`
	Mac   string `json:"mac" gorm:"comment:mac地址"`
	Uid   int64  `json:"uid" gorm:"comment:用户id"`
	State bool   `json:"state" gorm:"comment:是否启用"`
}