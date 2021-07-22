/**
 * @Desc    创建仓库服务
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/16
 * @Listen  MIT
 */
package services

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"http-api/app/http/graph/auth"
	grlModel "http-api/app/http/graph/model"
	"http-api/app/models/logs"
	"http-api/app/models/repositories"
	"http-api/app/models/repository_leader"
	"http-api/pkg/model"
)

func CreateRepository(ctx context.Context, input grlModel.CreateRepositoryInput) (*repositories.Repositories, error) {
	me := auth.GetUser(ctx)
	r := repositories.Repositories{
		Name:      input.Name,
		Address:   input.Address,
		PinYin:    input.PinYin,
		CompanyId: me.CompanyId,
		IsAble:    true,
	}
	if input.Remark != nil {
		r.Remark = *input.Remark
	}
	err := model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&r).Error; err != nil {
			return err
		}
		// 添加管理员
		for _, adminUid := range  input.RepositoryAdminID {
			rl := repository_leader.RepositoryLeader{
				RepositoryId: r.ID,
				Uid: adminUid,
			}
			if err := tx.Create(&rl).Error; err != nil {
				return err
			}
		}

		// 操作日志
		l := logs.Logos{
			Type:    logs.CreateActionType,
			Content: fmt.Sprintf("创建新的仓库: 仓库id为%d", r.ID),
			Uid:     me.Id,
		}
		if err := tx.Create(&l).Error; err != nil {
			return err
		}
		return nil
	})

	return &r, err
}
