/**
 * @Desc    The mutation_resolver is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/28
 * @Listen  MIT
 */
package mutation_resolver

import (
	"context"
	"gorm.io/gorm"
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/errors"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/order_specification_steel"
	"http-api/app/models/steel_logs"
	"http-api/app/models/steels"
	"http-api/pkg/model"
)

func (*MutationResolver) SetProjectSteelState(ctx context.Context, input graphModel.SetProjectSteelInput) (bool, error) {
	if err := requests.ValidateSetProjectSteelStateRequest(ctx, input); err != nil {
		return false, errors.ValidateErr(ctx, err)
	}
	me := auth.GetUser(ctx)
	err := model.DB.Transaction(func(tx *gorm.DB) error {
		for _, identifier := range input.IdentifierList {
			// 修改型钢状态
			steelItem := steels.Steels{}
			err := tx.Model(&steelItem).Where("identifier = ?", identifier).First(&steelItem).Error
			if err != nil {
				return err
			}
			if err := tx.Model(&steelItem).Update("state", input.State).Error; err != nil {
				return err
			}
			// 修改订单规格型钢状态
			orderSpecificationSteelItem := order_specification_steel.OrderSpecificationSteel{}
			err = tx.Model(&orderSpecificationSteelItem).Where("id = ?", steelItem.OrderSpecificationSteelId).
				First(&orderSpecificationSteelItem).
				Error
			if err != nil {
				return err
			}
			err = tx.Model(&orderSpecificationSteelItem).Where("id = ?", orderSpecificationSteelItem.Id).
				Update("state", input.State).
				Error
			if err != nil {
				return err
			}
			// 型钢操作日志
			err = tx.Create(&steel_logs.SteelLog{
				Type:    steel_logs.ChangeType,
				SteelId: steelItem.ID,
				Uid:     me.Id,
			}).Error
			if err != nil {
				return err
			}

		}

		return nil
	})

	if err != nil {
		return false, errors.ServerErr(ctx, err)
	}

	return true, nil
}
