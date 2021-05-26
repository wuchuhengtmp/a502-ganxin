/**
 * @Desc    The seeders is part of http-api
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/25
 * @Listen  MIT
 */
package seeders

import (
	"gorm.io/gorm"
	"http-api/app/models/companies"
	"http-api/pkg/seed"
	"time"
)

var companySeeders =  []seed.Seed{
	seed.Seed{
		Name: "create files",
		Run: func(db *gorm.DB) error {
			return CreateCompany( db, 1,
				"上海北星建筑工程有限公司",
				"BXJZ",
						"悦人达己 创新卓越",
						2,
						3,
						true,
						"13427969604",
						"wc20030318",
						getAfterDayTime(0),
						getAfterDayTime(365),
			)
		},
	},
}

func getAfterDayTime(day int) time.Time {
	now := time.Now()
	yyyy, mm, dd := now.Date()
	h := now.Hour()
	s := now.Second()
	i := now.Minute()
	tomorrow := time.Date(yyyy, mm, dd + day, h, i, s, 0, now.Location())
	return tomorrow
}

func CreateCompany(
	db *gorm.DB,
	id int64,
	name             string,
	pinYin           string,
	symbol           string,
	logoFileId       int64,
	backgroundFileId int64,
	state            bool,
	phone 			 string,
	wechat 			 string,
	startedAt        time.Time,
	endedAt          time.Time,
)  error {
	return db.Create(&companies.Companies{
		ID: id,
		Name: name,
		PinYin: pinYin,
		Symbol: symbol,
		LogoFileId: logoFileId,
		BackgroundFileId: backgroundFileId,
		State: state,
		Phone: phone,
		Wechat: wechat,
		StartedAt: startedAt,
		EndedAt: endedAt,
	}).Error
}
