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

const (
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
			Uid: me.ID,
			Content: fmt.Sprintf("修改价格: %f -> %s", oldPrice, c.Name),
			Type: logs.UpdateActionType,
		}
		if err := tx.Create(&l).Error; err != nil {
			return err
		}

		return nil
	})
}
