/**
 * @Desc    文件模型
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/24
 * @Listen  MIT
 */
package files

type Files struct {
	ID int64 `json:"id"`
	Path string `json:"path" gorm:"comment:文件路径"`
	Disk string `json:"disk" gorm:"comment:硬盘,default:local"`
}
