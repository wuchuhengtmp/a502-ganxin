/**
 * @Desc    The mutation_resolver is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/23
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
	"http-api/app/models/logs"
	"http-api/app/models/msg"
	"http-api/app/models/orders"
	"http-api/app/models/repositories"
	"http-api/app/models/repository_leader"
	"http-api/pkg/model"
)

func (*MutationResolver) DeleteOrder(ctx context.Context, input graphModel.DeleteOrderInput) (bool, error) {
	if err := requests.ValidateDeleteOrder(ctx, input); err != nil {
		return false, errors.ValidateErr(ctx, err)
	}
	me := auth.GetUser(ctx)
	myRole, _ := me.GetRole()
	err := model.DB.Transaction(func(tx *gorm.DB) error {
		orderItem := orders.Order{}
		if err := tx.Model(&orderItem).Where("id = ?", input.ID).First(&orderItem).Error; err != nil {
			return err
		}
		// 删除
		if err := tx.Where("id = ?", input.ID).Delete(&orders.Order{}).Error; err != nil {
			return err
		}
		// 通知仓库管理员
		leaderItem := repository_leader.RepositoryLeader{}
		var leaders []*repository_leader.RepositoryLeader
		repositoryTable := repositories.Repositories{}.TableName()
		err := tx.Model(&leaderItem).
			Select(fmt.Sprintf("%s.*", leaderItem.TableName())).
			Joins(fmt.Sprintf("join %s ON %s.id = %s.repository_id", repositoryTable, repositoryTable, leaderItem.TableName())).
			Where(fmt.Sprintf("%s.id = ?", repositoryTable), orderItem.RepositoryId).
			Find(&leaders).
			Error
		if err != nil {
			return err
		}
		for _, i := range leaders {
			msgItem := msg.Msg{
				Content: fmt.Sprintf(
					"订单%s 已被%s %s 删除, 更多详情请拨打对方电话:%s",
					orderItem.OrderNo,
					myRole.Name,
					me.Name,
					me.Phone,
				),
				Uid:  i.Uid,
				Type: msg.DeleteOrder,
			}
			if err := msgItem.CreateSelf(tx); err != nil {
				return err
			}
		}
		// 添加操作日志
		logItem := logs.Logos{
			Type:    logs.DeleteActionType,
			Content: fmt.Sprintf("删除订单: %s", orderItem.OrderNo),
			Uid:     me.Id,
		}
		if err := tx.Create(&logItem).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return false, errors.ServerErr(ctx, err)
	}

	return true, nil
}
