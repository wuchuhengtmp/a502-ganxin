/**
 * @Desc    编辑仓库
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
	"http-api/app/models/repositories"
	"http-api/app/models/repository_leader"
	"http-api/pkg/model"
)

func (*MutationResolver) SetRepository(ctx context.Context, input graphModel.SetRepositoryInput) (*repositories.Repositories, error) {
	if err := requests.ValidateSetRepositoryRequest(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	repositoryItem := repositories.Repositories{}
	leaderTable := repository_leader.RepositoryLeader{}.TableName()
	err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 修改管理员列表
		err := tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE  repository_id = %d", leaderTable, input.ID)).Error
		if err != nil {
			return err
		}
		for _, uid := range input.LeaderIDList {
			err = tx.Create(&repository_leader.RepositoryLeader{
				RepositoryId: input.ID,
				Uid:          uid,
			}).Error
			if err != nil {
				return err
			}
		}
		// 修改仓库
		err = tx.Model(&repositoryItem).Where("id = ?", input.ID).
			Update("name", input.Name).
			Update("address", input.Address).
			Update("remark", input.Remark).
			Update("is_able", input.IsAble).Error
		if err !=  nil {
			return err
		}
		if err := tx.Model(&repositoryItem).Where("id = ?", input.ID).First(&repositoryItem).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	return &repositoryItem, nil
}

