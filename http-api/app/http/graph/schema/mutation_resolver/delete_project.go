/**
 * @Desc    The mutation_resolver is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/19
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
	"http-api/app/models/logs"
	"http-api/app/models/msg"
	"http-api/app/models/project_leader"
	"http-api/app/models/projects"
	"http-api/pkg/model"
	"time"
)

func (*MutationResolver) DeleteProject(ctx context.Context, input graphModel.DeleteProjectInput) (bool, error) {
	if err := requests.ValidateDeleteProjectRequest(ctx, input); err != nil {
		return false, errors.ValidateErr(ctx, err)
	}
	me := auth.GetUser(ctx)
	myRole, _ := me.GetRole()
	err := model.DB.Transaction(func(tx *gorm.DB) error {
		leaderItem := project_leader.ProjectLeader{}
		projectItem := projects.Projects{}
		// 管理员通知管理员 并删除
		var leaders []project_leader.ProjectLeader
		if err := tx.Model(&leaderItem).Where("project_id = ?", input.ID).Find(&leaders).Error; err != nil {
			return err
		}
		if err := tx.Model(&projectItem).Where("id = ?", input.ID).First(&projectItem).Error; err != nil {
			return err
		}
		for _, item := range leaders {
			msgItem := msg.Msg{
				Content: fmt.Sprintf("您项目: %s 于 %s, 被 %s %s 删除了, 对方电话为: %s",
					projectItem.Name,
					helper.Time2Str(time.Now()),
					myRole.Name,
					me.Name,
					me.Phone,
				),
				Uid:  item.Uid,
				Type: msg.DeleteProject,
			}
			if err := msgItem.CreateSelf(tx); err != nil {
				return err
			}
			// 删除管理员
			if err := tx.Where("id = ?", item.ID).Delete(&leaderItem).Error; err != nil {
				return err
			}
		}
		// 删除项目
		if err := tx.Where("id = ?", input.ID).Delete(&projectItem).Error; err != nil {
			return err
		}
		// 操作日志
		logsItem := logs.Logos{
			Type:    logs.DeleteActionType,
			Content: fmt.Sprintf("删除 %s 项目， id: %d", projectItem.Name, projectItem.ID),
			Uid:     me.Id,
		}
		if err := tx.Create(&logsItem).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return false, errors.ServerErr(ctx, err)
	}

	return true, nil
}
