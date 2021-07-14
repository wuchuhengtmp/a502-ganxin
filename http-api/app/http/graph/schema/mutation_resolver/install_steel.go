/**
 * @Desc    安装型钢
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/26
 * @Listen  MIT
 */
package mutation_resolver

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/errors"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/order_specification_steel"
	"http-api/app/models/steel_logs"
	"http-api/app/models/steels"
	"http-api/pkg/model"
	"time"
)

func (*MutationResolver) InstallSteel(ctx context.Context, input graphModel.InstallLocationInput) (bool, error) {
	if err := requests.ValidateInstallSteelRequest(ctx, input); err != nil {
		return false, errors.ValidateErr(ctx, err)
	}
	err := model.DB.Transaction(func(tx *gorm.DB) error {
		steps := InstallSteelSteps{}
		// 更新订单中的型钢
		if err := steps.UpdateOrderSpecificationSteel(tx, ctx, input); err != nil {
			return err
		}
		// 创建型钢日志
		if err := steps.CreateSteelLog(tx, ctx, input); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return false, errors.ServerErr(ctx, err)
	}

	return true, nil
}

type InstallSteelSteps struct { }

/**
 * 更新订单中的型钢
 */
func (*InstallSteelSteps)UpdateOrderSpecificationSteel(tx *gorm.DB, ctx context.Context, input graphModel.InstallLocationInput) error {
	me := auth.GetUser(ctx)
	// 标记订单中的型钢为安装状态
	orderSpecificationSteelItem := order_specification_steel.OrderSpecificationSteel{}
	orderSpecificationSteelTable := orderSpecificationSteelItem.TableName()
	steelTable := steels.Steels{}.TableName()
	err := tx.Model(&orderSpecificationSteelItem).
		Select(fmt.Sprintf("%s.*", orderSpecificationSteelTable)).
		Joins(fmt.Sprintf("join %s ON %s.order_specification_steel_id = %s.id", steelTable, steelTable, orderSpecificationSteelTable)).
		Where(fmt.Sprintf("%s.identifier = ?", steelTable), input.Identifier).
		First(&orderSpecificationSteelItem).
		Error
	if err != nil {
		return err
	}
	err = tx.Model(&orderSpecificationSteelItem).Where("id = ?", orderSpecificationSteelItem.Id).
		// 型钢安装码
		Update("location_code", input.LocationCode).
		// 安装用户
		Update("installation_uid", me.Id).
		// 订单中型钢状态标记为使用中...
		Update("state", steels.StateProjectInUse).
		// 安装时间
		Update("installation_at", time.Now()).
		Error
	if err != nil {
		return err
	}

	return nil
}

/**
 * 创建型钢日志
 */
func (*InstallSteelSteps)CreateSteelLog(tx *gorm.DB, ctx context.Context, input graphModel.InstallLocationInput) error {
	me := auth.GetUser(ctx)
	//型钢日志
	steelItem := steels.Steels{}
	err := tx.Model(&steelItem).Where("identifier = ?", input.Identifier).First(&steelItem).Error
	if err != nil {
		return err
	}
	steelLogItem := steel_logs.SteelLog{
		Type:    steel_logs.InstallationType,
		SteelId: steelItem.ID,
		Uid:     me.Id,
	}
	if err := tx.Create(&steelLogItem).Error; err != nil {
		return err
	}

	return nil
}
