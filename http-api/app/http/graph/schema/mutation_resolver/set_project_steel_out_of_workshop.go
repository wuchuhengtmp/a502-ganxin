/**
 * @Desc    型钢出场解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/28
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
	"http-api/app/http/graph/util/helper"
	"http-api/app/models/msg"
	"http-api/app/models/order_specification_steel"
	"http-api/app/models/project_leader"
	"http-api/app/models/projects"
	"http-api/app/models/specificationinfo"
	"http-api/app/models/steel_logs"
	"http-api/app/models/steels"
	"http-api/pkg/model"
	"time"
)

func (*MutationResolver) SetProjectSteelOutOfWorkshop(ctx context.Context, input graphModel.SetProjectSteelOutOfWorkshopInput) (bool, error) {
	if err := requests.ValidateSetProjectSteelOutOfWorkshopRequest(ctx, input); err != nil {
		return false, errors.ValidateErr(ctx, err)
	}
	me := auth.GetUser(ctx)
	var weithTotal float64
	var total int64
	err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 型钢标记为出库
		for _, identifier := range input.IdentifierList {
			steelItem := steels.Steels{}
			err := tx.Model(&steelItem).Where("identifier = ? ", identifier).
				Where("company_id = ?", me.CompanyId).
				Update("state", steels.StateProjectOnTheStoreWay).Error
			if err != nil {
				return err
			}
			// 订单型钢标记为出库
			tx.Model(&steelItem).Where("company_id = ?", me.CompanyId).
				Where("identifier = ?", identifier).
				First(&steelItem)
			orderSpecificationSteelItem := order_specification_steel.OrderSpecificationSteel{}
			err = tx.Model(&orderSpecificationSteelItem).
				Where("id = ?", steelItem.OrderSpecificationSteelId).
				Update("state", steels.StateProjectOnTheStoreWay).Error
			if err != nil {
				return err
			}
			// 添加型钢操作记录
			err = tx.Create(&steel_logs.SteelLog{
				Type:    steel_logs.OutOfWorkshop,
				SteelId: steelItem.ID,
				Uid:     me.Id,
			}).Error
			if err != nil {
				return err
			}
			// 统计
			total++
			specificationInfoItem := specificationinfo.SpecificationInfo{}
			err = tx.Model(&specificationInfoItem).Where("id = ?", steelItem.ID).First(&specificationInfoItem).Error
			if err != nil {
				return err
			}
			weithTotal += specificationInfoItem.Weight
		}

		// 通知仓库管理员
		projectItem := projects.Projects{}
		tx.Model(&projectItem).Where("id = ?", input.ProjectID).First(&projectItem)
		contente := fmt.Sprintf(
			"%s 项目于 %s 从场地往仓库发一批型钢,总数: %d根，重量: %.2f, 请注意查收!",
			projectItem.Name,
			helper.Time2Str(time.Now()),
			total,
			weithTotal,
		)
		var leaders []*project_leader.ProjectLeader
		err := tx.Model(&project_leader.ProjectLeader{}).Where("project_id = ?", input.ProjectID).Find(&leaders).Error
		if err != nil {
			return err
		}
		for _, leaderItem := range leaders {
			msgItem := msg.Msg{
				Content: contente,
				Uid:     leaderItem.Uid,
				Type:    msg.OutOfWorkshop,
			}
			err = msgItem.CreateSelf(tx)
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
