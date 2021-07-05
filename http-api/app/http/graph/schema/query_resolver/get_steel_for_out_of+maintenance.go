/**
 * @Desc    获取可出厂的型钢解析器
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
	"http-api/app/models/maintenance_record"
	"http-api/app/models/steels"
	"http-api/pkg/model"
)

func (*QueryResolver) GetSteelForOutOfMaintenance(ctx context.Context, input graphModel.GetSteelForOutOfMaintenanceInput) (*maintenance_record.MaintenanceRecord, error) {
	var res maintenance_record.MaintenanceRecord
	if err := requests.ValidateGetSteelForOutOfMaintenanceRequest(ctx, input); err != nil {
		return &res, err
	}
	recordItem := maintenance_record.MaintenanceRecord{}
	steelTable := steels.Steels{}.TableName()
	me := auth.GetUser(ctx)
	err := model.DB.Model(&recordItem).
		Select(fmt.Sprintf("%s.*", recordItem.TableName())).
		Joins(fmt.Sprintf("join %s ON %s.maintenance_record_steel_id = %s.id", steelTable, steelTable, recordItem.TableName())).
		Where(fmt.Sprintf("%s.identifier = ?", steelTable), input.Identifier).
		Where(fmt.Sprintf("%s.company_id = ?", steelTable), me.CompanyId).
		First(&recordItem).
		Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	return &recordItem, nil
}
