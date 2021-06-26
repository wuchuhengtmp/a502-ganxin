/**
 * @Desc    安装码是否可用解析验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/26
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"fmt"
	"http-api/app/http/graph/errors"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/order_specification_steel"
	"http-api/app/models/steels"
	"http-api/pkg/model"
)

func (*QueryResolver) IsAccessLocationCode(ctx context.Context, input graphModel.IsAccessLocationCodeInput) (bool, error) {
	if err := requests.ValidateIsAccessLocationCodeRequest(ctx, input); err != nil {
		return false, errors.ValidateErr(ctx, err)
	}
	orderSpecificationSteelTable := order_specification_steel.OrderSpecificationSteel{}.TableName()
	orderSpecificationSteelItem := order_specification_steel.OrderSpecificationSteel{}
	steelTalbe := steels.Steels{ }.TableName()
	err := model.DB.Model(&orderSpecificationSteelItem).
		Select(fmt.Sprintf("%s.*", orderSpecificationSteelTable)).
		Joins(fmt.Sprintf("join %s ON %s.order_specification_steel_id = %s.id", steelTalbe, steelTalbe, orderSpecificationSteelTable)).
		Where(fmt.Sprintf("%s.location_code = ?", orderSpecificationSteelTable), input.LocationCode).
		Where(fmt.Sprintf("%s.identifier = ?", steelTalbe), input.Identifier).
		First(&orderSpecificationSteelItem).
		Error
	if err != nil && err.Error() == "record not found" {
		return true, nil
	}
	if err != nil {
		return false, errors.ServerErr(ctx, err)
	}

	return false, nil
}
