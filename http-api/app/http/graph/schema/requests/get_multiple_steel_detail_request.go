/**
 * @Desc    快速查询多个型钢请求验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/23
 * @Listen  MIT
 */
package requests

import (
	"context"
	"fmt"
	"http-api/app/http/graph/auth"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/models/steels"
	"http-api/pkg/model"
)

func ValidateGetMultipleSteelDetailRequest(ctx context.Context, input *graphModel.GetMultipleSteelDetailInput) error {
	var ss []*steels.Steels
	// 检验识别码是否为冗余
	steps := ValidateGetProject2WorkshopDetailRequestSteps{}
	if err := steps.CheckRedundancyIdentification(input.IdentifierList); err != nil {
		return err
	}
	// 检验是否有识别码不存在
	err := model.DB.Model(&steels.Steels{}).Where("identifier in ?", input.IdentifierList).
		Find(&ss).
		Error
	if err != nil {
		return err
	}
	if len(ss) != len(input.IdentifierList) {
		me := auth.GetUser(ctx)
		for _, indentifier := range input.IdentifierList{
			err := model.DB.Model(&steels.Steels{}).
				Where("identifier = ? AND company_id = ?", indentifier, me.CompanyId).
				First(&steels.Steels{}).
				Error
			if err != nil {
				return fmt.Errorf("没有识别码为:%s的型钢", indentifier)
			}
		}
	}

	return nil
}
