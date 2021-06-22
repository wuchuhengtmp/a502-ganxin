/**
 * @Desc    获取型钢详情请求验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/22
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

func ValidateGetOneSteelDetailRequest(ctx context.Context, input graphModel.GetOneSteelDetailInput) error {
	me := auth.GetUser(ctx)
	err := model.DB.Model(&steels.Steels{}).
		Where("company_id = ? AND identifier = ?", me.CompanyId, input.Identifier).
		First(&steels.Steels{}).Error
	if err != nil {
		return fmt.Errorf("没有识别码为：%s的型钢", input.Identifier)
	}

	return nil
}
