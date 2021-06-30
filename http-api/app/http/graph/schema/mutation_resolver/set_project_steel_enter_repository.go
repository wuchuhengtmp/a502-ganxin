/**
 * @Desc    型钢归库解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/30
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

func (*MutationResolver) SetProjectSteelEnterRepository(ctx context.Context, input graphModel.SetProjectSteelEnterRepositoryInput) (bool, error) {
	if err := requests.ValidateSetProjectSteelEnterRepositoryRequest(ctx, input); err != nil {
		return false, errors.ValidateErr(ctx, err)
	}
	me := auth.GetUser(ctx)
	err := model.DB.Transaction(func(tx *gorm.DB) error {
		for _, identifier := range input.IdentifierListr {
			// 修改型钢状态
			steelItem := steels.Steels{}
			err := tx.Model(&steelItem).Where("identifier = ?", identifier).Update("state", steels.StateInStore).Error
			if err != nil {
				return err
			}
			//  修改订单型钢
			orderSpecificationSteelItem := order_specification_steel.OrderSpecificationSteel{}
			err = tx.Model(&orderSpecificationSteelItem).
				Joins(fmt.Sprintf("join %s ON %s.order_specification_steel_id = %s.id", steelItem.TableName(), steelItem.TableName(), orderSpecificationSteelItem.TableName())).
				Where(fmt.Sprintf("%s.identifier = ?", steelItem.TableName()), identifier).
				First(&orderSpecificationSteelItem).
				Error
			if err != nil {
				return err
			}
			err = tx.Model(&orderSpecificationSteelItem).Where("id = ?", orderSpecificationSteelItem.Id).
				Update("enter_repository_uid", me.Id).
				Update("enter_repository_at", time.Now()).
				Update("state", steels.StateInStore).
				Error
			if err != nil {
				return err
			}
			// 添加型钢日志
			err = tx.Model(&steelItem).Where("identifier = ?", identifier).First(&steelItem).Error
			if err != nil {
				return err
			}
			l := steel_logs.SteelLog{
				Type:    steel_logs.EnterRepositoryType,
				SteelId: steelItem.ID,
				Uid:     me.Id,
			}
			if err := tx.Create(&l).Error; err != nil {
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
