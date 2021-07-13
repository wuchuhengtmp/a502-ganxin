/**
 * @Desc    获取型钢列表(用于仪表盘)
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/19
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"fmt"
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/errors"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/models/order_specification_steel"
	"http-api/app/models/projects"
	"http-api/app/models/roles"
	"http-api/app/models/steels"
	"http-api/pkg/model"
)

func (*QueryResolver) GetSteelForDashboard(ctx context.Context, input *graphModel.GetSteelForDashboardInput) (*projects.GetProjectSteelDetailRes, error) {
	var res projects.GetProjectSteelDetailRes
	recordItem := order_specification_steel.OrderSpecificationSteel{}
	recordTable := recordItem.TableName()
	steelTable := steels.Steels{}.TableName()
	modelIn := model.DB.Model(&recordItem).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.steel_id", steelTable, steelTable, recordTable))
	me := auth.GetUser(ctx)
	role, _ := me.GetRole()
	if role.ID != roles.RoleAdminId {
		modelIn.Where(fmt.Sprintf("%s.company_id = ?", steelTable), me.CompanyId)
	}
	if input.RepositoryID != nil {
		modelIn.Where(fmt.Sprintf("%s.repository_id = ?", steelTable), *input.RepositoryID)
	}
	if err := modelIn.Count(&res.Total).Error; err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	o := input.PageSize * (input.Page - 1)
	if err := modelIn.Select(fmt.Sprintf("%s.*", recordTable)).Offset(int(o)).Limit(int(input.PageSize)).Find(&res.List).Error ; err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	return &res, nil
}
