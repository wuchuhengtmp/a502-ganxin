/**
 * @Desc    创建公司解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/28
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
	"http-api/app/models/companies"
	"http-api/app/models/configs"
	"http-api/app/models/logs"
	"http-api/app/models/roles"
	"http-api/app/models/users"
	helper2 "http-api/pkg/helper"
	"http-api/pkg/model"
)

type CreateCompanySteps companies.Companies

func (m *MutationResolver) CreateCompany(ctx context.Context, input graphModel.CreateCompanyInput) (*companies.Companies, error) {
	CreateCompanyRequest := requests.CreateCompanyRequest{}
	err := CreateCompanyRequest.ValidateCreateCompanyRequest(input)
	if err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	c := companies.Companies{}
	steps := CreateCompanySteps{}
	err = model.DB.Transaction(func(tx *gorm.DB) error {
		//添加公司
		newCompany, err := steps.CreateSelf(ctx, input, tx)
		if err != nil {
			return err
		}
		c = *newCompany
		// 初始始化配置
		if err := steps.CreateConfig(tx); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	return &c, nil
}

/**
 * 创建公司
 */
func (steps *CreateCompanySteps) CreateSelf(ctx context.Context, input graphModel.CreateCompanyInput, tx *gorm.DB) (*companies.Companies, error) {
	startedAt, _ := helper.Str2Time(input.StartedAt)
	endedAt, _ := helper.Str2Time(input.EndedAt)
	c := companies.Companies{}
	c.Name = input.Name
	c.PinYin = input.PinYin
	c.Symbol = input.Symbol
	c.LogoFileId = int64(input.LogoFileID)
	c.BackgroundFileId = int64(input.BackgroundFileID)
	c.IsAble = input.IsAble
	c.Phone = input.Phone
	c.Wechat = input.Wechat
	c.StartedAt = startedAt
	c.EndedAt = endedAt
	if err := tx.Create(&c).Error; err != nil {
		return nil, err
	}
	me := auth.GetUser(ctx)
	logsModel := logs.Logos{}
	logsModel.Type = logs.CreateActionType
	logsModel.Uid = me.Id
	user := users.Users{
		Name:         input.AdminName,
		Password:     helper2.GetHashByStr(input.AdminPassword),
		Phone:        input.AdminPhone,
		RoleId:       roles.RoleCompanyAdminId,
		Wechat:       input.AdminWechat,
		CompanyId:    c.ID,
		IsAble:       true,
		AvatarFileId: input.AdminAvatarFileID,
	}
	logsModel.Content = fmt.Sprintf("添加公司,名称:%s,ID: %d", c.Name, c.ID)
	if err := tx.Create(&logsModel).Error; err != nil {
		return nil, err
	}
	if err := tx.Create(&user).Error; err != nil {
		return nil, err
	}
	steps.ID = c.ID

	return &c, nil
}

/**
 * 创建公司配置
 */
func (c CreateCompanySteps) CreateConfig(tx *gorm.DB) error {
	// 通过键名来初始化
	initConfigByKeyName := func(key string, newKey string, remark string) error {
		configItem := configs.Configs{}
		err := tx.Debug().Model(&configItem).Where("name = ?", key).
			Where("company_id = ?", 0).
			First(&configItem).
			Error
		if err != nil {
			return err
		}
		con := configs.Configs{
			Name:      newKey,
			Value:     configItem.Value,
			CompanyId: c.ID,
			Remark:    remark,
		}
		if err := tx.Create(&con).Error; err != nil {
			return err
		}

		return nil
	}
	// 初始化公司客服
	if err := initConfigByKeyName(configs.GLOBAL_PHONE_NAME, configs.PHONE_NAME, "公司客服电话"); err != nil {
		return err
	}
	// 初始化公司微信
	if err := initConfigByKeyName(configs.GLOBAL_WECHAT_NAME, configs.WECHAT_NAME, "公司微信"); err != nil {
		return err
	}
	// 初始化公司视频
	if err := initConfigByKeyName(configs.GLOBAL_TUTOR_FILE_NAME, configs.TUTOR_FILE_NAME, "型钢教学"); err != nil {
		return err
	}
	// 初始化公司价格
	if err := initConfigByKeyName(configs.GLOBAL_PRICE_NAME, configs.PRICE_NAME, "型钢价格"); err != nil {
		return err
	}

	return nil
}
