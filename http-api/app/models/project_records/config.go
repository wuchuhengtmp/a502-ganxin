/**
 * @Desc    项目详情模型
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/26
 * @Listen  MIT
 */
package project_records

import (
	"gorm.io/gorm"
	"time"
)

type ProjectRecord struct {
	Id                 int64     `json:"id"`
	ProjectId          int64     `json:"projectId" gorm:"comment:项目id"`
	OrderDetailId      int64     `json:"orderDetailId" gorm:"comment:订单详情id"`
	SteelId            int64     `json:"steelId" gorm:"comment:型钢id"`
	OutRepositoryAt    time.Time `json:"outRepositoryAt" gorm:"comment:出库时间"`
	ReturnRepositoryAt time.Time `json:"returnRepositoryAt" gorm:"comment:入库时间"`
	OutSiteAt          time.Time `json:"outSiteAt" gorm:"comment:出库时间"`
	ReturnSiteAt       time.Time `json:"returnSiteAt" gorm:"comment:入库时间"`
	InstallNo          int64     `json:"installNo" gorm:"comment:安装编号"`
	InstalledAt        time.Time `json:"installedAt" gorm:"comment:安装时间"`
	gorm.Model
}
