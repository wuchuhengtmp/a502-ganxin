/**
 * @Desc    The seeders is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @DATE    2021/4/27
 * @Listen  MIT
 */
package seeders

import (
	"gorm.io/gorm"
	"http-api/app/models/configs"
	"http-api/pkg/seed"
)

var configsSeeders =  []seed.Seed{
	seed.Seed{
		Name: "create config",
		Run: func(db *gorm.DB) error {
			return CreateConfig(db, 1, configs.PRICE_NAME, "1", "型钢单价", CompanyId)
		},
	},
	seed.Seed{
		Name: "create config",
		Run: func(db *gorm.DB) error {
			return CreateConfig(db, 2, configs.TUTOR_FILE_NAME, "7", "教学视频文件", CompanyId)
		},
	},
	seed.Seed{
		Name: "create config",
		Run: func(db *gorm.DB) error {
			return CreateConfig(db, 3, configs.WECHAT_NAME, "12345678", "微信号", CompanyId)
		},
	},
	seed.Seed{
		Name: "create config",
		Run: func(db *gorm.DB) error {
			return CreateConfig(db, 4, configs.PHONE_NAME, "12345678901", "电话号", CompanyId)
		},
	},
	seed.Seed{
		Name: "create config",
		Run: func(db *gorm.DB) error {
			return CreateConfig(db, 5, configs.SMS_SIGN, "惠州市蚁人科技有限公司", "短信签名", 0)
		},
	},


	seed.Seed{
		Name: "create config",
		Run: func(db *gorm.DB) error {
			return CreateConfig(db, 6, configs.SMS_ACCESS_KEY, "LTAI4GEZUaCwsv8J6TxRCc55", "短信密key", 0)
		},
	},
	seed.Seed{
		Name: "create config",
		Run: func(db *gorm.DB) error {
			return CreateConfig(db, 7, configs.SMS_ACCESS_SECRET_KEY, "e6tAMTcN____EuQwyP1nDF5d2xn8IwfSwU____", "短信accessKey", 0)
		},
	},
	seed.Seed{
		Name: "create config",
		Run: func(db *gorm.DB) error {
			return CreateConfig(db, 8, configs.SMS_TEMPLATECODE, "SMS_215455029", "短信模板", 0)
		},
	},
	seed.Seed{
		Name: "create config",
		Run: func(db *gorm.DB) error {
			return CreateConfig(db, 9, configs.GLOBAL_PRICE_NAME, "1", "用于初始化型钢单价", 0)
		},
	},
	seed.Seed{
		Name: "create config",
		Run: func(db *gorm.DB) error {
			return CreateConfig(db, 10, configs.GLOBAL_TUTOR_FILE_NAME, "7", "用于初始化教学视频文件",0)
		},
	},
	seed.Seed{
		Name: "create config",
		Run: func(db *gorm.DB) error {
			return CreateConfig(db, 11, configs.GLOBAL_WECHAT_NAME, "12345678", "初始化微信号", 0)
		},
	},
	seed.Seed{
		Name: "create config",
		Run: func(db *gorm.DB) error {
			return CreateConfig(db, 12, configs.GLOBAL_PHONE_NAME, "12345678901", "初始化电话号", 0)
		},
	},

}

func CreateConfig(db *gorm.DB, id int64,  name string, value string, remark string, companyId int64)  error {
	return db.Create(&configs.Configs{ID: id, Name: name, Value: value, Remark: remark, CompanyId: companyId}).Error
}
