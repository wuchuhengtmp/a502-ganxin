/**
 * @Desc    The configs is part of http-api
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @DATE    2021/4/27
 * @Listen  MIT
 */

package configs

import (
	"context"
	"gorm.io/gorm"
	"http-api/app/http/graph/auth"
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

const (
	APP_ICON   = "APP_ICON"
	APP_NAME   = "APP_NAME"
	PRICE_NAME = "PRICE"
)

func getVal(key string, ctx context.Context) string {
	var about Configs
	me := auth.GetUser(ctx)
	err := model.DB.Model(&Configs{}).Where("name = ? AND company_id  = ?", key, me.CompanyId).First(&about).Error
	if err != nil {
		logger.LogError(err)
	}
	return about.Value
}

func (c *Configs) GetPrice(ctx context.Context) float64 {
	v := getVal(PRICE_NAME, ctx)
	s, _ := strconv.ParseFloat(v, 64)

	return s
}
