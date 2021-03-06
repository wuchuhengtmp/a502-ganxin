/**
 * @Desc    获取仓库概览
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/16
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"fmt"
	"http-api/app/http/graph/errors"
	"http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/orders"
	"http-api/app/models/repositories"
	"http-api/app/models/specificationinfo"
	"http-api/app/models/steels"
	sqlModel "http-api/pkg/model"
)

func (*QueryResolver) GetRepositoryOverview(ctx context.Context, input model.GetRepositoryOverviewInput) (*repositories.GetRepositoryOverviewRes, error) {
	if err := requests.ValidateGetOverviewRequest(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	var g repositories.GetRepositoryOverviewRes
	mapWhere := fmt.Sprintf("id = %d", input.ID)
	if input.SpecificationID != nil {
		mapWhere = fmt.Sprintf("%s AND specification_id = %d", mapWhere, *input.SpecificationID)
	}
	sTable := specificationinfo.SpecificationInfo{}.TableName()
	steelsTable := steels.Steels{}.TableName()
	err := sqlModel.DB.Debug().Model(&steels.Steels{}).
		Select(fmt.Sprintf("count(*) as total, sum(%s.weight) as weight", sTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.specification_id", sTable, sTable, steelsTable)).
		Where(fmt.Sprintf("%s.state = %d ", steelsTable, steels.StateInStore)).
		Where(fmt.Sprintf("%s.specification_id = ?", steelsTable), input.SpecificationID).
		Scan(&g).Error
	if err != nil {
		return nil, err
	}
	//  要减去还没发货但已经确认的数量和重量
	cTotal, _ := orders.GetConfirmSteelTotalBySpecificationId(*input.SpecificationID)
	g.Total -= cTotal
	w, _ := orders.GetConfirmSteelTotalWeightBySpecificationId(*input.SpecificationID)
	g.Weight -= w

	return &g, nil
}
