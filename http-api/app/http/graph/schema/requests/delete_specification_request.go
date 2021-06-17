/**
 * @Desc    删除规格请求验证
 * @Author  wuchuheng<wuchuheng@163.com>
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
	order_details "http-api/app/models/order_specification"
	"http-api/app/models/specificationinfo"
	"http-api/app/models/steels"
)

func ValidateDeleteSpecificationRequest(ctx context.Context, id int64) error {
	rules := govalidator.MapData{
		"id": []string{"isSpecificationId"},
	}
	opt := govalidator.Options{
		Data: &struct {
			ID int64 `json:"id"`
		}{id},
		Rules:         rules,
		TagIdentifier: "json",
	}
	res := govalidator.New(opt).ValidateStruct()
	if len(res) > 0 {
		for _, fieldErrors := range res {
			for _, err := range fieldErrors {
				return fmt.Errorf("%s", err)
			}
		}
	}
	s := specificationinfo.SpecificationInfo{ID: id}
	_ = s.GetSelf()
	me := auth.GetUser(ctx)
	if me.CompanyId != s.CompanyId {
		return fmt.Errorf("当前规格记录与您不是归属于同一家公司，您无权操作")
	}
	steelsModel := steels.Steels{}
	if _, err := steelsModel.GetSteelsBySpecificationId(id); err != nil {
		return fmt.Errorf("该规格正被型钢使用中，无法删除")
	}
	OrderDetailModel := order_details.OrderSpecification{}
	if _, err := OrderDetailModel.GetOrderBySpecificationId(id); err != nil {

	}

	return nil
}
