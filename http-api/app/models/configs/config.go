/**
 * @Desc    The configs is part of http-api
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @DATE    2021/4/27
 * @Listen  MIT
 */

package configs

import (
	"gorm.io/gorm"
	"http-api/pkg/logger"
	"http-api/pkg/model"
)

type Configs struct {
	ID     int64  `json:"id"`
	Name   string `json:"name" gorm:"comment:参数名"`
	Value  string `json:"value" gorm:"comment:参数值"`
	Remark string `json:"remark" gorm:"comment:配置备注"`
	gorm.Model
}

var AboutKey = "about"

const APP_NAME = "APP_NAME"

const APP_ICON = "APP_ICON"

/**
 * 小程序 appId配置key
 */
const MINI_WECHAT_APP_ID = "MINI_WECHAT_APP_ID"

/**
 * 小程序 appSecret配置key
 */
const MINI_WECHAT_APP_SECRET = "MINI_WECHAT_APP_SECRET"

func getVal(key string) string  {
	var about Configs
	err := model.DB.Model(&Configs{}).Where("name", key).First(&about).Error
	if err != nil { logger.LogError(err) }
	return about.Value
}
