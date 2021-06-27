/**
 * @Desc 	待修改武钢详细信息解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/27
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
	"http-api/app/models/projects"
	"http-api/app/models/specificationinfo"
	"http-api/app/models/steels"
	"http-api/pkg/model"
)

func (*QueryResolver) GetProjectSteel2BeChangeDetail(ctx context.Context, input graphModel.ProjectSteel2BeChangeInput) (*projects.GetProjectSteel2BeChangeDetailRes,  error) {
	var res  projects.GetProjectSteel2BeChangeDetailRes
	if err := requests.GetProjectSteel2BeChangeDetailRequest(ctx, input); err != nil {
		return &res, errors.ValidateErr(ctx, err)
	}
	orderSpecificationSteelTable := order_specification_steel.OrderSpecificationSteel{}.TableName()
	steelTable := steels.Steels{}.TableName()
	specificationInfoTable := specificationinfo.SpecificationInfo{}.TableName()
	queryInstance := model.DB.Model(&order_specification_steel.OrderSpecificationSteel{}).
		Joins(fmt.Sprintf("join %s ON %s.order_specification_steel_id = %s.id", steelTable, steelTable, orderSpecificationSteelTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.specification_id", specificationInfoTable, specificationInfoTable, steelTable)).
		Where(fmt.Sprintf("%s.identifier in ?", steelTable), input.IdentifierList)
	// 规格过滤
	if input.SpecificationID != nil {
		queryInstance = queryInstance.Where(fmt.Sprintf("%s.id = ?", specificationInfoTable), *input.SpecificationID)
	}
	// 状态过滤
	if input.State != nil {
		queryInstance = queryInstance.Where(fmt.Sprintf("%s.state = ?", orderSpecificationSteelTable), *input.State)
	}
	err := queryInstance.
		Select(fmt.Sprintf("%s.*", orderSpecificationSteelTable)).
		Find(&res.List).Error
	if err != nil {
		return &res, errors.ServerErr(ctx, err)
	}
	// 数量
	res.Total = int64(len(res.List))
	// 重量
	var weightInfo struct{ WeightTotal float64 }
	err = queryInstance.
		Select("sum(weight) as WeightTotal").
		Scan(&weightInfo).Error
	if err != nil {
		return &res, errors.ServerErr(ctx, err)
	}
	res.WeightTotal = weightInfo.WeightTotal

	return &res, err
}
