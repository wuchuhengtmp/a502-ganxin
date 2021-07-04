/**
 * @Desc    获取型钢入厂的型钢信息解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/3
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

func (*QueryResolver)GetEnterMaintenanceSteel(ctx context.Context, input graphModel.EnterMaintenanceInput) (*maintenance_record.MaintenanceRecord, error) {
	if err := requests.ValidateGetEnterMaintenanceRequest(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	me := auth.GetUser(ctx)
	recorderItem := maintenance_record.MaintenanceRecord{}
	steelsTable := steels.Steels{}.TableName()
	err := model.DB.Model(&recorderItem).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.steel_id", steelsTable, steelsTable, recorderItem.TableName())).
		Where(fmt.Sprintf("%s.identifier = ?", steelsTable), input.Identifier).
		Where(fmt.Sprintf("%s.company_id = ?", steelsTable), me.CompanyId).
		First(&recorderItem).Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)

	}

	return &recorderItem, nil
}


