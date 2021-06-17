/**
 * @Desc    编辑规格求验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/5
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

func ValidateEditSpecificationRequest(ctx context.Context, input model.EditSpecificationInput) error {
	rules := govalidator.MapData{
		"id":     []string{"isSpecificationId"},
		"length": []string{"isGreaterZero"},
		"weight": []string{"isGreaterZero"},
		"type":   []string{"min:6"},
	}
	opt := govalidator.Options{
		Data:            &input,
		Rules:           rules,
		TagIdentifier:   "json",
	}
	res := govalidator.New(opt).ValidateStruct()
	if len(res) > 0 {
		for _, fieldErrors := range res {
			for _, err := range fieldErrors {
				return fmt.Errorf("%s", err)
			}
		}
	}
	me := auth.GetUser(ctx)
	s := specificationinfo.SpecificationInfo{
		ID: input.ID,
	}
	_ = s.GetSelf()
	if me.CompanyId != s.CompanyId {
		return fmt.Errorf("当前规格记录与您不是归属于同一家公司，您无权操作")
	}

	return nil
}
