/**
 * @Desc
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/14
 * @Listen  MIT
 */
package services

import (
	"fmt"
	"http-api/app/models/order_specification_steel"
	"http-api/app/models/steels"
	"http-api/pkg/config"
	"http-api/pkg/logger"
	"http-api/pkg/model"
	"time"
)

/**
 *  启动时间
 */
var launchAt string = "00:00:00"

func init()  {
	 launchAt = config.GetString("LAUNCH_AT_FOR_CALCULATE_RAGE")
	// 定时计算型钢年利用率
	go CalculateSteelUsageYearRate()
}


/**
 * 定时计算型钢年利用率
 */
func CalculateSteelUsageYearRate()  {
	ticker := time.NewTicker(time.Second)
	for {
		<- ticker.C
		now := time.Now()
		ctime := fmt.Sprintf("%02d:%02d:%02d", now.Hour(), now.Minute(), now.Second())
		if ctime == launchAt {
			var steelList []steels.Steels
			err := model.DB.Model(&steels.Steels{}).Find(&steelList).Error
			if err != nil {
				logger.LogError(err)
			} else {
				limiter := make(chan bool, 20)
				for _, item := range steelList {
					limiter <- true
					go CalculateSteelUsageYearRateById(item, &limiter)
					go CalculateSteelUsageRateById(item, &limiter)
				}
			}
		}
	}
}

/**
 * 计算每根型钢的年使用率
 */
func CalculateSteelUsageYearRateById(item steels.Steels, limiter *chan bool) {
	defer func() {
		<- *limiter
	}()
	var createAt time.Time
	if	item.CreatedAt.Unix() > time.Now().AddDate(-1, 0, 0).Unix() {
		createAt = item.CreatedAt
	} else {
		createAt = item.CreatedAt.AddDate(-1, 0, 0)
	}
	// 一年中的时间
	totalTimeLen := time.Now().Unix() - createAt.Unix()
	recordItem := order_specification_steel.OrderSpecificationSteel{}
	var items []order_specification_steel.OrderSpecificationSteel
	err := model.DB.
		Model(&recordItem).
		Where("steel_id = ?", item.ID).
		Where("out_workshop_at IS NOT NULL").
		Find(&items).
		Error
	if err != nil {
		logger.LogError(err)
		return
	}
	// 场地使用的时间
	var usageTimeLen int64
	for _, record := range items{
		usageTimeLen += record.OutWorkshopAt.Unix() - record.EnterWorkshopAt.Unix()
	}
	usageRage := usageTimeLen / totalTimeLen
	err = model.DB.Model(&steels.Steels{}).
		Where("id = ?", item.ID ).
		Update("usage_year_rate", usageRage).
		Error
	if err != nil {
		logger.LogError(err)
	}
}

/**
 * 计算型钢总使用率
 */
func CalculateSteelUsageRateById(item steels.Steels, limiter *chan bool) {
	defer func() {
		<- *limiter
	}()
	totalTimeLen := time.Now().Unix() - item.CreatedAt.Unix()
	recordItem := order_specification_steel.OrderSpecificationSteel{}
	var items []order_specification_steel.OrderSpecificationSteel
	err := model.DB.
		Model(&recordItem).
		Where("steel_id = ?", item.ID).
		Where("out_workshop_at IS NOT NULL").
		Find(&items).
		Error
	if err != nil {
		logger.LogError(err)
		return
	}
	// 场地使用的时间
	var usageTimeLen int64
	for _, record := range items{
		usageTimeLen += record.OutWorkshopAt.Unix() - record.EnterWorkshopAt.Unix()
	}
	usageRage := float64(usageTimeLen) / float64(totalTimeLen)
	err = model.DB.Model(&steels.Steels{}).
		Where("id = ?", item.ID ).
		Update("total_usage_rate", usageRage).
		Error
	if err != nil {
		logger.LogError(err)
	}
}
