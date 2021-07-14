/**
 * @Desc    The configs is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @DATE    2021/4/27
 * @Listen  MIT
 */

package configs

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"http-api/app/http/graph/auth"
	"http-api/app/models/logs"
	"http-api/pkg/logger"
	"http-api/pkg/model"
	"strconv"
)

type Configs struct {
	ID        int64  `json:"id"`
	Name      string `json:"name" gorm:"comment:参数名"`
	Value     string `json:"value" gorm:"comment:参数值"`
	Remark    string `json:"remark" gorm:"comment:配置备注"`
	CompanyId int64  `json:"company_id" gorm:"comment:归属公司"`
	gorm.Model
}

func (Configs) TableName() string {
	return "configs"
}

const (
	PRICE_NAME             = "PRICE"                  // 价格字段名
	TUTOR_FILE_NAME        = "TUTOR"                  // 教学文件字段名
	WECHAT_NAME            = "WECHAT"                 // 微信
	PHONE_NAME             = "PHONE"                  // 电话名
	SMS_SIGN               = "SMS_SIGN"               // 短信签名
	SMS_TEMPLATECODE       = "SMS_TEMPLATECODE"       // 短信模板
	SMS_ACCESS_KEY         = "SMS_ACCESS_KEY"         // 短信密钥
	SMS_ACCESS_SECRET_KEY  = "SMS_ACCESS_SECRET_KEY"  // 短信密钥
	GLOBAL_PRICE_NAME      = "Global_PRICE_NAME"      // 用于初始化价格
	GLOBAL_TUTOR_FILE_NAME = "GLOVAL_TUTOR_FILE_NAME" // 用于初始化教学文件
	GLOBAL_WECHAT_NAME     = "GLOBAL_WECHAT_NAME"     // 用于初始化微信
	GLOBAL_PHONE_NAME      = "GLOBAL_PHONE_NAME"      // 用于初始化公司的客户电话
)

func GetVal(key string, ctx context.Context) string {
	var about Configs
	me := auth.GetUser(ctx)
	err := model.DB.Model(&Configs{}).Where("name = ? AND company_id  = ?", key, me.CompanyId).First(&about).Error
	if err != nil {
		logger.LogError(err)
	}
	return about.Value
}

func GetGlobalVal(key string) string {
	var about Configs
	err := model.DB.Model(&Configs{}).Where("name = ? ", key).First(&about).Error
	if err != nil {
		logger.LogError(err)
	}
	return about.Value
}

func (c *Configs) GetPrice(ctx context.Context) float64 {
	v := GetVal(PRICE_NAME, ctx)
	s, _ := strconv.ParseFloat(v, 64)

	return s
}

/**
 * 编辑价格
 */
func (c *Configs) EditPrice(ctx context.Context) error {
	return model.DB.Transaction(func(tx *gorm.DB) error {
		oldPrice := c.GetPrice(ctx)
		me := auth.GetUser(ctx)
		if err := tx.Model(&Configs{}).
			Where("name = ? AND company_id = ?", PRICE_NAME, me.CompanyId).
			Update("value", c.Value).Error; err != nil {
			return err
		}
		l := logs.Logos{
			Uid:     me.Id,
			Content: fmt.Sprintf("修改价格: %.4f -> %s", oldPrice, c.Value),
			Type:    logs.UpdateActionType,
		}
		if err := tx.Create(&l).Error; err != nil {
			return err
		}

		return nil
	})
}
