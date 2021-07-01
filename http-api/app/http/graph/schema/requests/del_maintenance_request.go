/**
 * @Desc    删除维修厂请求验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/1
 * @Listen  MIT
 */
package requests

import (
	"context"
	"fmt"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/models/maintenance_record"
	"http-api/pkg/model"
)

func ValidateDelMaintenanceRequest(ctx context.Context, input graphModel.DelMaintenanceInput) error {
	steps := StepsForMaintenance{}
	// 检验有没有这维修厂
	if err := steps.CheckHashMaintenance(ctx, input.ID) ; err != nil {
		return err
	}
	// 有没有型钢在维修
	err := model.DB.Model(&maintenance_record.MaintenanceRecord{}).
		Where("maintenance_id", input.ID).
		First(&maintenance_record.MaintenanceRecord{}).
		Error
	if err != nil {
		if err.Error() == "record not found" {
			return nil
		}
		return err
	} else {
		return fmt.Errorf("维修有相关的维修型钢数据，不能删除")
	}
}
