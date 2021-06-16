/**
 * @Desc    获取仓库概览请求验证器
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/16
 * @Listen  MIT
 */
package requests

import (
	"context"
	"fmt"
	"github.com/thedevsaddam/govalidator"
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/model"
	"http-api/app/models/specificationinfo"
)

func ValidateGetOverviewRequest(ctx context.Context, input model.GetRepositoryOverviewInput) error {
	// 规格验证
	if input.SpecificationID != nil {
		s := specificationinfo.SpecificationInfo{
			ID: *input.SpecificationID,
		}
		if err := s.IsExist(ctx); err != nil {
			return fmt.Errorf("这家公司没有这个规格")
		}
	}
	me := auth.GetUser(ctx)
	idStr := fmt.Sprintf("%d", me.ID)
	// 仓库验证
	rules := govalidator.MapData{
		"id": []string{"required", "isCompanyRepository:" + idStr},
	}
	opts := govalidator.Options{
		Data:          &input,
		Rules:         rules,
		TagIdentifier: "json",
	}

	return Validate(opts)
}
