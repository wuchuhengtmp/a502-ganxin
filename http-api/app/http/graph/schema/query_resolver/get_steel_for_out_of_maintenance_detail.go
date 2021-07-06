/**
 * @Desc    待出厂的型钢详情
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
	"http-api/app/http/graph/errors"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/maintenance"
	"http-api/app/models/maintenance_record"
	"http-api/app/models/specificationinfo"
	"http-api/app/models/steels"
	"http-api/pkg/model"
)

func (*QueryResolver) GetSteelForOutOfMaintenanceDetail(ctx context.Context, input graphModel.GetSteelForOutOfMaintenanceDetailInput) (*maintenance.GetSteelForOutOfMaintenanceDetailRes, error) {
	var res maintenance.GetSteelForOutOfMaintenanceDetailRes
	if err := requests.ValidateGetSteelForOutOfMaintenanceDetailRequest(ctx, input); err != nil {
		return &res, err
	}
	steelTable := steels.Steels{}.TableName()
	record := maintenance_record.MaintenanceRecord{}
	modelIn := model.DB.Model(&record).
		Joins(fmt.Sprintf("join %s ON %s.maintenance_record_steel_id = %s.id", steelTable, steelTable, record.TableName())).
		Where(fmt.Sprintf("%s.identifier IN ?", steelTable), input.IdentifierList)
	if input.SpecificationID != nil {
		modelIn = modelIn.Where(fmt.Sprintf("%s.specification_id = ?", steelTable), *input.SpecificationID)
	}
	if err := modelIn.Select(fmt.Sprintf("%s.*", record.TableName())).Scan(&res.List).Error; err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	specificationTable := specificationinfo.SpecificationInfo{}.TableName()
	var weightInfo struct{ Weight float64 }
	err := modelIn.Select("sum(weight) as Weight").
		Joins(fmt.Sprintf("join %s ON %s.id = %s.specification_id", specificationTable, specificationTable, steelTable)).
		Scan(&weightInfo).
		Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	res.Weight = weightInfo.Weight
	res.Total = int64(len(res.List))

	return &res, nil
}

//func (*QueryResolver) GetSteelForOutOfMaintenanceDetail(ctx context.Context, input graphModel.GetSteelForOutOfMaintenanceDetailInput) (*maintenance.GetSteelForOutOfMaintenanceDetailRes, error) {
//	var res maintenance.GetSteelForOutOfMaintenanceDetailRes
//	if err := requests.ValidateGetSteelForOutOfMaintenanceDetailRequest(ctx, input); err != nil {
//		return &res, err
//	}
//	err := model.DB.Transaction(func(tx *gorm.DB) error {
//		steps := GetSteelForOutOfMaintenanceDetailSteps{}
//		for _, identifier := range input.IdentifierList {
//			// 标记型钢状态
//			if err := steps.FlagSteel(ctx, identifier, tx); err != nil {
//				return err
//			}
//			// 标记维修型钢详情
//			if err := steps.FlagMaintenanceSteel(ctx, identifier, tx); err != nil {
//				return err
//			}
//			// 型钢日志
//			if err := steps.CreateLog(ctx, identifier, tx); err != nil {
//				return err
//			}
//		}
//
//		return nil
//	})
//	if err != nil {
//		return nil, errors.ServerErr(ctx, err)
//	}
//
//	return &res, nil
//}
//
//type GetSteelForOutOfMaintenanceDetailSteps struct{}
//
///**
// * 标记型钢状态（标记为归库途中）
// */
//func (*GetSteelForOutOfMaintenanceDetailSteps) FlagSteel(ctx context.Context, identifier string, tx *gorm.DB) error {
//	steelsItem := steels.Steels{}
//	me := auth.GetUser(ctx)
//	err := tx.Model(&steelsItem).Where("identifier = ?", identifier).
//		Where("company_id = ?", me.CompanyId).
//		Update("state", steels.StateMaintainerOnTheStoreWay).
//		Error
//
//	return err
//}
//
///**
// * 标记维修型钢状态（标记为归库途中）
// */
//func (*GetSteelForOutOfMaintenanceDetailSteps) FlagMaintenanceSteel(ctx context.Context, identifier string, tx *gorm.DB) error {
//	steelsTable := steels.Steels{}.TableName()
//	me := auth.GetUser(ctx)
//	record := maintenance_record.MaintenanceRecord{}
//	err := tx.Model(&record).
//		Select(fmt.Sprintf("%s.*", record.TableName())).
//		Joins(fmt.Sprintf("join %s ON %s.maintenance_record_steel_id = %s.id", steelsTable, steelsTable, record.TableName())).
//		Where(fmt.Sprintf("%s.identifier = ?", steelsTable), identifier).
//		Where(fmt.Sprintf("%s.company_id = ?", steelsTable), me.CompanyId).
//		First(&record).
//		Error
//	if err != nil {
//		return err
//	}
//	err = tx.Model(&record).Where("id = ?", record.Id).
//		Update("state", steels.StateMaintainerOnTheStoreWay).
//		Update("outed_uid", me.Id).
//		Update("outed_at", time.Now()).
//		Error
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
///**
// * 添加日志
// */
//func (*GetSteelForOutOfMaintenanceDetailSteps) CreateLog(ctx context.Context, identifier string, tx *gorm.DB) error {
//	me := auth.GetUser(ctx)
//	steelItem := steels.Steels{}
//	err := tx.Model(&steelItem).
//		Where("identifier = ?", identifier).
//		Where("company_id = ?", me.CompanyId).
//		First(&steelItem).
//		Error
//	if err != nil {
//		return err
//	}
//	l := steel_logs.SteelLog{
//		Type:    steel_logs.OutOfMaintenance,
//		SteelId: steelItem.ID,
//		Uid:     me.Id,
//	}
//	if err := tx.Create(&l).Error; err != nil {
//		return err
//	}
//
//	return nil
//}
