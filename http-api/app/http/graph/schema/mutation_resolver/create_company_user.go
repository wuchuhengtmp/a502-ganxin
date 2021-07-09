/**
 * @Desc    添加公司成员解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/3
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
	"http-api/app/models/companies"
	"http-api/app/models/logs"
	"http-api/app/models/roles"
	"http-api/app/models/users"
	helper2 "http-api/pkg/helper"
	sqlModel "http-api/pkg/model"
)

type CreateCompanyUserSteps companies.Companies
func (m *MutationResolver) CreateCompanyUser(ctx context.Context, input graphModel.CreateCompanyUserInput) (*users.Users, error) {
	validate := requests.CreateCompanyUserRequest{}
	err := validate.ValidateCreateCompanyUserRequest(input)
	if err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	var steps  CreateCompanyUserSteps
	newUser, err := steps.CreateUser(ctx, input)
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	return newUser, nil
}


/**
 * 添加公司归属下的人员
 */
func (CreateCompanyUserSteps) CreateUser(ctx context.Context, input graphModel.CreateCompanyUserInput) (*users.Users, error) {
	me := auth.GetUser(ctx)
	user := users.Users{}
	err := sqlModel.DB.Transaction(func(tx *gorm.DB) error {
		// 添加人员
		user.Name = input.Name
		user.Phone = input.Phone
		user.Wechat = input.Wechat
		user.Password = helper2.GetHashByStr(input.Password)
		user.AvatarFileId = input.AvatarID
		user.RoleId = roles.RoleTagMapId[input.Role.String()]
		user.CompanyId = me.CompanyId
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		// 操作日志
		log := logs.Logos{}
		log.Uid = me.Id
		log.Content = fmt.Sprintf("添加 %s", roles.RoleTagMapName[input.Role.String()])
		log.Type = logs.CreateActionType
		log.Uid = me.Id
		if err := tx.Create(&log).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	return &user, nil
}
