/**
 * @Desc    创建需求单验证器
 * @Author  wuchuheng<root@wuchuheng.com>
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
	grapModel "http-api/app/http/graph/model"
)

func ValidateCreateOrderValidate(ctx context.Context, input grapModel.CreateOrderInput) error {
	steps := StepsForOrder{}
	// 检验开始时间
	if err := steps.CheckExpectedAt(ctx, input.ExpectedReturnAt); err != nil {
		return err
	}
	// 型钢规格需求验证
	if err := steps.CheckSteelSpecification(ctx, input.SteelList ); err != nil {
		return err
	}

	strId := fmt.Sprintf("%d", auth.GetUser(ctx).Id)
	rules := govalidator.MapData{
		"projectId":  []string{"required", "isCompanyProject:" + strId},
		"repositoryId": []string{"required", "isCompanyRepository:" + strId},
	}
	opts := govalidator.Options{
		Data:          &input,
		Rules:         rules,
		TagIdentifier: "json",
	}

	return Validate(opts)
}
