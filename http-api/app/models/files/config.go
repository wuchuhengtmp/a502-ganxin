/**
 * @Desc    文件模型
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/24
 * @Listen  MIT
 */
package files

import (
	"gorm.io/gorm"
	"http-api/pkg/filesystem"
	"http-api/pkg/model"
)

type File struct {
	ID   int64  `json:"id"`
	Path string `json:"path" gorm:"comment:文件路径"`
	Disk string `json:"disk" gorm:"comment:硬盘,default:local"`
	gorm.Model
}

func (file *File)GetSelfById(id int64) error {
	db := model.DB
	return db.Model(file).Where("id = ?", id).First(file).Error
}

/**
 * 保存文件
 */
func (file *File) CreateFile() error {
	db := model.DB
	file.Disk = filesystem.GetDefaultDisk()
	err := db.Model(file).Create(file).Error
	return err
}

func (file *File) IsExist() bool {
	db := model.DB
	f := File{}
	err := db.Model(file).Where("id = ?", file.ID).First(&f).Error
	if err != nil {
		return false
	} else {
		return true
	}
}

func (file *File) GetUrl() string {
	return filesystem.Disk(file.Disk).Url(file.Path)
}
