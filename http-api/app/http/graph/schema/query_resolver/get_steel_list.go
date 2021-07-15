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
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/errors"
	grpahModel "http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/http/graph/util/helper"
	"http-api/app/models/codeinfo"
	"http-api/app/models/maintenance_record"
	"http-api/app/models/order_specification"
	"http-api/app/models/order_specification_steel"
	"http-api/app/models/orders"
	"http-api/app/models/projects"
	"http-api/app/models/repositories"
	"http-api/app/models/specificationinfo"
	"http-api/app/models/steels"
	"http-api/app/models/users"
	"http-api/pkg/model"
	"sync"
)

type GetSteelListSteps struct{}

func (*QueryResolver) GetSteelList(ctx context.Context, input grpahModel.PaginationInput) (*steels.GetSteelListRes, error) {
	if err := requests.ValidateGetSteelListRequest(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	steps := GetSummarySteps{}
	offset := 0
	if input.Page > 1 {
		offset = int((input.Page - 1) * input.PageSize)
	}
	me := auth.GetUser(ctx)
	whereMap := fmt.Sprintf("company_id = %d", me.CompanyId)
	// 仓库
	if input.RepositoryID != nil {
		whereMap = fmt.Sprintf("%s AND repository_id = %d", whereMap, *input.RepositoryID)
	}
	// 规格
	if input.SpecificationID != nil {
		whereMap = fmt.Sprintf("%s AND specification_id = %d", whereMap, *input.SpecificationID)
	}
	// 识别码
	if input.Identifier != nil {
		whereMap = fmt.Sprintf("%s AND identifier like '%s'", whereMap, "%"+*input.Identifier+"%")
	}
	// 编码
	if input.Code != nil {
		whereMap = fmt.Sprintf("%s AND code like '%s'", whereMap, "%"+*input.Code+"%")
	}
	// 状态过滤
	if input.State != nil {
		whereMap = fmt.Sprintf("%s AND state = %d", whereMap, *input.State)
	}
	// 材料商过滤
	if input.MaterialManufacturerID != nil {
		whereMap = fmt.Sprintf("%s AND material_manufacturer_id = %d", whereMap, *input.MaterialManufacturerID)
	}
	// 制造商过滤
	if input.ManufacturerID != nil {
		whereMap = fmt.Sprintf("%s AND manufacturer_id = %d", whereMap, *input.ManufacturerID)
	}
	// 首次入库时间
	if input.CreatedAt != nil {
		s, e := helper.GetSecondBetween(*input.CreatedAt)
		whereMap = fmt.Sprintf("%s AND (created_at between '%s' AND '%s' )", whereMap, s, e)
	}
	// 生产时间过滤
	if input.ProduceAt != nil {
		s, e := helper.GetSecondBetween(*input.ProduceAt)
		whereMap = fmt.Sprintf("%s AND (produced_date between '%s' AND '%s' )", whereMap, s, e)
	}
	res := steels.GetSteelListRes{}
	if err := model.DB.Model(&steels.Steels{}).Where(whereMap).Count(&res.Total).Error; err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	err := model.DB.Model(&steels.Steels{}).Where(whereMap).Offset(offset).Limit(int(input.PageSize)).Find(&res.List).Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	wg := &sync.WaitGroup{}
	resChan := make(chan ChanItemRes, len(res.List))
	resWg := &sync.WaitGroup{}
	go func() {
		resWg.Add(1)
		defer resWg.Done()
		for item := range resChan {
			res.List[item.Index].Turnover = item.Res
		}
	}()
	limiter := make(chan bool, 20)
	for i, item := range res.List {
		limiter <- true
		wg.Add(1)
		go steps.GetTurnoverById(i, item.ID, wg, &limiter, &resChan)
	}
	wg.Wait()
	close(resChan)
	resWg.Wait()

	return &res, nil
}

type ChanItemRes struct {
	Index int
	Res   int64
}

func (GetSummarySteps) GetTurnoverById(index int, id int64, wg *sync.WaitGroup, limiter *chan bool, resChan *chan ChanItemRes) {
	defer func() {
		wg.Done()
		<-*limiter
	}()
	res := ChanItemRes{
		Index: index,
	}
	recordItem := order_specification_steel.OrderSpecificationSteel{}
	model.DB.Model(&recordItem).Where("steel_id = ?", id).Count(&res.Res)
	*resChan <- res
}

type SteelItemResolver struct{}

func (SteelItemResolver) StateInfo(ctx context.Context, obj *steels.Steels) (*steels.StateItem, error) {
	return &steels.StateItem{
		State: obj.State,
		Desc:  steels.StateCodeMapDes[obj.State],
	}, nil
}

func (SteelItemResolver) Specifcation(ctx context.Context, obj *steels.Steels) (*specificationinfo.SpecificationInfo, error) {
	return obj.GetSpecification()
}

func (SteelItemResolver) MaterialManufacturer(ctx context.Context, obj *steels.Steels) (*codeinfo.CodeInfo, error) {
	return obj.GetMaterialManufacturer()
}
func (SteelItemResolver) Manufacturer(ctx context.Context, obj *steels.Steels) (*codeinfo.CodeInfo, error) {
	return obj.GetManufacturer()
}
func (SteelItemResolver) Repository(ctx context.Context, obj *steels.Steels) (*repositories.Repositories, error) {
	return obj.GetRepository()
}
func (SteelItemResolver) CreateUser(ctx context.Context, obj *steels.Steels) (*users.Users, error) {
	u := users.Users{}
	err := model.DB.Model(&users.Users{}).Where("id = ?", obj.CreatedUid).First(&u).Error
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (SteelItemResolver) SteelInMaintenance(ctx context.Context, obj *steels.Steels) ([]*maintenance_record.MaintenanceRecord, error) {
	var recordList []*maintenance_record.MaintenanceRecord
	m := maintenance_record.MaintenanceRecord{}
	err := model.DB.Model(&m).Where("steel_id = ?", obj.ID).Find(&recordList).Error
	if err != nil {
		return recordList, nil
	}

	return recordList, nil
}
func (SteelItemResolver) SteelInProject(ctx context.Context, obj *steels.Steels) ([]*order_specification_steel.OrderSpecificationSteel, error) {
	o := order_specification_steel.OrderSpecificationSteel{}
	var orderList []*order_specification_steel.OrderSpecificationSteel
	err := model.DB.Model(&o).Where("steel_id = ?", obj.ID).Find(&orderList).Error

	return orderList, err
}

type SteelInProjectResolver struct{}

func (SteelInProjectResolver) Steel(ctx context.Context, obj *order_specification_steel.OrderSpecificationSteel) (*steels.Steels, error) {
	s := steels.Steels{}
	err := model.DB.Model(&s).Where("id = ?", obj.SteelId).First(&s).Error

	return &s, err
}

func (SteelInProjectResolver) Order(ctx context.Context, obj *order_specification_steel.OrderSpecificationSteel) (*orders.Order, error) {
	o := orders.Order{}
	orderTable := o.TableName()
	orderSpecificationTable := order_specification.OrderSpecification{}.TableName()
	err := model.DB.Model(&o).
		Select(fmt.Sprintf("%s.*", orderTable)).
		Joins(fmt.Sprintf("join %s ON %s.order_id = %s.id", orderSpecificationTable, orderSpecificationTable, orderTable)).
		Where(fmt.Sprintf("%s.id = ?", orderSpecificationTable), obj.OrderSpecificationId).
		First(&o).
		Error

	return &o, err
}
func (SteelInProjectResolver) UseDays(ctx context.Context, obj *order_specification_steel.OrderSpecificationSteel) (*int64, error) {
	var days int64
	recordItem := order_specification_steel.OrderSpecificationSteel{}
	if err := model.DB.Model(&recordItem).Where("id = ?", obj.Id).First(&recordItem).Error; err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	if recordItem.State == steels.StateInStore || recordItem.State == steels.StateProjectOnTheStoreWay {
		timeLen := recordItem.OutWorkshopAt.Unix() - recordItem.EnterWorkshopAt.Unix()
		days = timeLen / (60 * 60 * 24)
	}

	return &days, nil
}
func (SteelInProjectResolver) ProjectName(ctx context.Context, obj *order_specification_steel.OrderSpecificationSteel) (string, error) {
	item := projects.Projects{}
	projectsTable := projects.Projects{}.TableName()
	orderTable := orders.Order{}.TableName()
	orderSpecificationTable := order_specification.OrderSpecification{}.TableName()
	orderSpecificationSteelTable := order_specification_steel.OrderSpecificationSteel{}.TableName()

	err := model.DB.Model(&item).
		Select(fmt.Sprintf("%s.*", projectsTable)).
		Joins(fmt.Sprintf("join %s ON %s.project_id = %s.id", orderTable, orderTable, projectsTable)).
		Joins(fmt.Sprintf("join %s ON %s.order_id = %s.id", orderSpecificationTable, orderSpecificationTable, orderTable)).
		Joins(fmt.Sprintf("join %s ON %s.order_specification_id = %s.id", orderSpecificationSteelTable, orderSpecificationSteelTable, orderSpecificationTable)).
		Where(fmt.Sprintf("%s.id = ?", orderSpecificationSteelTable), obj.Id).
		First(&item).
		Error
	return item.Name, err
}
