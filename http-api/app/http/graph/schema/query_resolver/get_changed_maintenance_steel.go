/**
 * @Desc  	查询用于维修型钢的信息解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/5
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"fmt"
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/errors"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/maintenance"
	"http-api/app/models/maintenance_leader"
	"http-api/app/models/maintenance_record"
	"http-api/app/models/steels"
	"http-api/pkg/model"
)

func (*QueryResolver) GetChangedMaintenanceSteel(ctx context.Context, input graphModel.GetChangedMaintenanceSteelInput) (*maintenance_record.MaintenanceRecord, error) {
	if err := requests.ValidateGetChangedMaintenanceSteelRequest(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	recordItem := maintenance_record.MaintenanceRecord{}
	leaderTable := maintenance_leader.MaintenanceLeader{}.TableName()
	maintenanceTable := maintenance.Maintenance{}.TableName()
	steelTable := steels.Steels{}.TableName()
	me := auth.GetUser(ctx)
	err := model.DB.Model(&recordItem).
		Select(fmt.Sprintf("%s.*", recordItem.TableName())).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.steel_id", steelTable, steelTable, recordItem.TableName())).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.maintenance_id", maintenanceTable, maintenanceTable, recordItem.TableName())).
		Joins(fmt.Sprintf("join %s ON %s.maintenance_id = %s.id", leaderTable, leaderTable, maintenanceTable)).
		Where(fmt.Sprintf("%s.identifier = ?", steelTable), input.Identifier).
		Where(fmt.Sprintf("%s.uid = ?", leaderTable), me.Id).
		First(&recordItem).
		Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	return &recordItem, nil
}
