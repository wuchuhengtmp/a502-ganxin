/**
 * @Desc    获取型钢列表解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/11
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"fmt"
	"http-api/app/http/graph/errors"
	grpahModel "http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/codeinfo"
	"http-api/app/models/maintenance_record"
	"http-api/app/models/order_specification"
	"http-api/app/models/order_specification_steel"
	"http-api/app/models/orders"
	"http-api/app/models/repositories"
	"http-api/app/models/specificationinfo"
	"http-api/app/models/steels"
	"http-api/app/models/users"
	"http-api/pkg/model"
)


func (*QueryResolver)GetSteelList(ctx context.Context, input grpahModel.PaginationInput) (*steels.GetSteelListRes, error) {
	if err := requests.ValidateGetSteelListRequest(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	s := steels.Steels{}
	list, err := s.GetPagination(ctx, input.Page, input.PageSize, input.RepositoryID, input.SpecificationID)
	if  err != nil {
		return nil, err
	}

	res := steels.GetSteelListRes {
		List: list,
		Total: s.GetTotal(ctx, input.RepositoryID, input.SpecificationID),
	}

	return &res, nil
}

type SteelItemResolver struct { }

func (SteelItemResolver) Specifcation(ctx context.Context, obj *steels.Steels) (*specificationinfo.SpecificationInfo, error) {
	return obj.GetSpecification()
}

func (SteelItemResolver)MaterialManufacturer(ctx context.Context, obj *steels.Steels) (*codeinfo.CodeInfo, error) {
	return obj.GetMaterialManufacturer()
}
func (SteelItemResolver)Manufacturer(ctx context.Context, obj *steels.Steels) (*codeinfo.CodeInfo, error) {
	return obj.GetManufacturer()
}
func (SteelItemResolver)Repository(ctx context.Context, obj *steels.Steels) (*repositories.Repositories, error) {
	return obj.GetRepository()
}
func (SteelItemResolver)CreateUser(ctx context.Context, obj *steels.Steels) (*users.Users, error) {
	u := users.Users{}
	err := model.DB.Model(&users.Users{}).Where("id = ?", obj.CreatedUid).First(&u).Error
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (SteelItemResolver)SteelInMaintenance(ctx context.Context, obj *steels.Steels) ([]*maintenance_record.MaintenanceRecord, error) {
	var recordList []*maintenance_record.MaintenanceRecord
	m := maintenance_record.MaintenanceRecord{}
	err := model.DB.Model(&m).Where("steel_id = ?", obj.ID).Find(&recordList).Error
	if err != nil {
		return recordList, nil
	}

	return recordList, nil
}
func (SteelItemResolver)SteelInProject(ctx context.Context, obj *steels.Steels) ([]*order_specification_steel.OrderSpecificationSteel, error) {
	o := order_specification_steel.OrderSpecificationSteel{}
	var orderList []*order_specification_steel.OrderSpecificationSteel
	err := model.DB.Model(&o).Where("steel_id = ?", obj.ID).Find(&orderList).Error

	return orderList, err
}
type SteelInProjectResolver struct {}
func (SteelInProjectResolver) Steel(ctx context.Context, obj *order_specification_steel.OrderSpecificationSteel) (*steels.Steels, error) {
	s := steels.Steels{}
	err := model.DB.Model(&s).Where("id = ?", obj.SteelId).First(&s).Error

	return &s, err
}

func (SteelInProjectResolver)Order(ctx context.Context, obj *order_specification_steel.OrderSpecificationSteel) (*orders.Order, error) {
	o := orders.Order{}
	orderTable := o.TableName()
	orderSpecificationTable := order_specification.OrderSpecification{}.TableName()
	err := model.DB.Model(&o).
		Select(fmt.Sprintf("%s.*", orderTable)).
		Joins(fmt.Sprintf("join %s ON %s.order_id = %s.id", orderSpecificationTable, orderSpecificationTable, orderTable)).
		Where(fmt.Sprintf("%s.id = ?", orderTable), obj.OrderSpecificationId).
		First(&o).
		Error

	return &o, err
}
func (SteelInProjectResolver)UseDays(ctx context.Context, obj *order_specification_steel.OrderSpecificationSteel) (*int64, error) {
	var days int64
	// todo 使用天数
	return &days, nil
}
