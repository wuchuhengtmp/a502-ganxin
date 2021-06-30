/**
 * @Desc    获取待归库详情请求验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/30
 * @Listen  MIT
 */
package requests

import (
	"context"
	"fmt"
	"http-api/app/http/graph/auth"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/models/specificationinfo"
	"http-api/pkg/model"
)

func ValidateGetToBeEnterRepositoryDetailRequest(ctx context.Context, input graphModel.GetToBeEnterRepositoryDetailInput) error {
	steps := StepsForProject{}
	if input.State != nil {
		if err := steps.CheckIsEnterRepositoryState(*input.State); err != nil {
			return err
		}
	}
	// 检验有没有这个规格id
	me := auth.GetUser(ctx)
	if input.SpecificationID != nil {
		s := specificationinfo.SpecificationInfo{}
		err := model.DB.Model(&s).Where("company_id = ?", me.CompanyId).
			Where("id = ?", *input.SpecificationID).
			First(&s).Error
		if err != nil && err.Error() == "record not found" {
			return fmt.Errorf("没有这个规格")
		} else  if err != nil {
			return  err
		}
	}
	// 有没有这个项目
	if err := steps.CheckHasProject(ctx, input.ProjectID); err != nil {
		return err
	}

	return nil
}
