/**
 * @Desc    更新项目
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/16
 * @Listen  MIT
 */
package mutation_resolver

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"http-api/app/http/graph/errors"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/project_leader"
	"http-api/app/models/projects"
	"http-api/pkg/model"
)

func (*MutationResolver) SetProject(ctx context.Context, input graphModel.SetProjectInput) (*projects.Projects, error) {
	if err := requests.ValidateSetProjectRequest(ctx, input); err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	leaderTable := project_leader.ProjectLeader{}.TableName()
	projectItem := projects.Projects{}
	err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 修改管理员列表
		err := tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE project_id = %d", leaderTable, input.ID)).Error
		if err != nil {
			return err
		}
		for _, uid := range input.LeaderIDList {
			err = tx.Create(&project_leader.ProjectLeader{
				ProjectId: input.ID,
				Uid:       uid,
			}).Error
			if err != nil {
				return err
			}
		}
		// 修改项目
		err = tx.Model(&projectItem).
			Where("id =  ?", input.ID).
			Update("name", input.Name).
			Update("city", input.City).
			Update("address", input.Address).
			Update("started_at", input.StartedAt).
			Update("remark", input.Remark).Error
		if err != nil {
			return err
		}
		if input.EndedAt != nil {
			err = tx.Model(&projectItem).
				Where("id =  ?", input.ID).
				Update("ended_at", *input.EndedAt).Error
		}
		if err != nil {
			return err
		}
		if err := tx.Model(&projectItem).Where("id = ?", input.ID).First(&projectItem).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	return &projectItem, nil
}
