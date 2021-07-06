/**
 * @Desc    获取维修厂型钢
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/6
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
	"sync"
)

func (*QueryResolver) GetMaintenanceSteel(ctx context.Context, input graphModel.GetMaintenanceSteelInput) (*maintenance.GetMaintenanceSteelRes, error) {
	if err := requests.ValidateGetMaintenanceSteelRequest(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	steps := GetMaintenanceSteelSteps{}
	recordList, err := steps.GetRecordList(input)
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	var res maintenance.GetMaintenanceSteelRes
	specificationMapListItem, total, weight, err := steps.getRes(ctx, recordList)
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	res.Total = total
	res.Weight = weight
	for _, item := range specificationMapListItem {
		// :xxx 重在问题 这里如果不使用中间值，则map出来的值会是同一个指针，也就是说，直接用会造成只有一个元素的slice 必须基于一个中间值才行
		// 我这想当然是错误的，理由是就算是遍历出来的是的值的指针是同一个， 也不会影响slice列表的数量， 所以如果有一天，你发现了答案，请
		// 邮件告知我，谢谢!!! 😁
		tmp := item
		res.List = append(res.List, &tmp)
	}

	return &res, nil
}

type GetMaintenanceSteelSteps struct{}

/**
 * 遍历查询响应结果需要的数据
 */
type SpecificationItemChanType struct {
	RecordItem            *maintenance_record.MaintenanceRecord
	SpecificationInfoItem *specificationinfo.SpecificationInfo
}
/**
 * 获取响应结果需要的数量
 * :xxx 这里采用并发方式来查询数据，在获得性能4.5倍的提升的同时， 并发的代码量大大增加和可读性能也大大降低，当你看到这时里，我想作业就是提高代码的可以性了
 */
func (steps *GetMaintenanceSteelSteps) getRes(ctx context.Context, recordList []*maintenance_record.MaintenanceRecord) (specificationMapListItem map[int64]maintenance.GetMaintenanceSteelResItem, total int64, weight float64, err error) {
	limit := 20
	limiter := make(chan bool, limit)
	specificationItemChan := make(chan *SpecificationItemChanType, limit)
	isResultWGOk := &sync.WaitGroup{} // 收集结果信号
	specificationMapListItem = make(map[int64]maintenance.GetMaintenanceSteelResItem)
	// 订阅结果
	go func() {
		isResultWGOk.Add(1)
		for specificationItem := range specificationItemChan {
			i := steps.GetSteelResItem(*specificationItem.SpecificationInfoItem, specificationItem.RecordItem)
			specificationMapListItem = steps.Merge(specificationMapListItem, i)
			total++
			weight += specificationItem.SpecificationInfoItem.Weight
		}
		isResultWGOk.Done()
	}()
	// 订阅错误
	errChan := make(chan error, limit)
	isErrWgOk := &sync.WaitGroup{}
	go func() {
		isErrWgOk.Add(1)
		for newErr := range errChan {
			err = newErr
		}
		isErrWgOk.Done()
	}()
	isTaskWGOk := &sync.WaitGroup{}
	for _, recordItem := range recordList {
		limiter <- true
		isTaskWGOk.Add(1)
		go steps.GetSpecificationItem(
			*recordItem,
			&limiter,
			&specificationItemChan,
			&errChan,
			isTaskWGOk,
		)
	}
	isTaskWGOk.Wait()
	close(specificationItemChan)
	isResultWGOk.Wait()
	close(errChan)
	isErrWgOk.Wait()

	return
}

/**
 * 是否归库了
 */
func (*GetMaintenanceSteelSteps) IsStored(state int64) bool {
	if state == steels.StateInStore {
		return true
	} else {
		return false
	}
}

/**
 * 合并统计
 */
func (*GetMaintenanceSteelSteps) Merge(group map[int64]maintenance.GetMaintenanceSteelResItem, newItem maintenance.GetMaintenanceSteelResItem) map[int64]maintenance.GetMaintenanceSteelResItem {
	if oldItem, ok := group[newItem.Id]; ok {
		oldItem.StoredWeight += newItem.StoredWeight
		oldItem.StoredTotal += newItem.StoredTotal
		oldItem.ReceivedWeight += newItem.ReceivedWeight
		oldItem.ReceivedTotal += newItem.ReceivedTotal
		group[newItem.Id] = oldItem
	} else {
		group[newItem.Id] = newItem
	}

	return group
}

/**
 * 获取规格
 */
func (*GetMaintenanceSteelSteps) GetSpecificationItem(
	recordItem maintenance_record.MaintenanceRecord,
	limiter *chan bool,
	specificationItemChan *chan *SpecificationItemChanType,
	errChan *chan error,
	isTaskWGOK *sync.WaitGroup,
) {
	defer func() {
		isTaskWGOK.Done() // 任务搞定
		<-*limiter // 释放并发量-1
	}()

	specificationItem := specificationinfo.SpecificationInfo{}
	specificationInfoTable := specificationinfo.SpecificationInfo{}.TableName()
	steelTable := steels.Steels{}.TableName()
	err := model.DB.Model(&specificationItem).
		Select(fmt.Sprintf("%s.*", specificationItem.TableName())).
		Joins(fmt.Sprintf("join %s ON %s.specification_id = %s.id", steelTable, steelTable, specificationInfoTable)).
		Where(fmt.Sprintf("%s.id = ?", steelTable), recordItem.SteelId).
		First(&specificationItem).
		Error
	if err != nil {
		*errChan <- err // 把错误传递出去
	} else {
		i := SpecificationItemChanType{
			RecordItem:            &recordItem,
			SpecificationInfoItem: &specificationItem,
		}
		*specificationItemChan <- &i // 把查询结果传递出去
	}
}

/**
 * 获取维修的型钢列表
 */
func (*GetMaintenanceSteelSteps) GetRecordList(input graphModel.GetMaintenanceSteelInput) ([]*maintenance_record.MaintenanceRecord, error) {
	recordItem := maintenance_record.MaintenanceRecord{}
	var recordList []*maintenance_record.MaintenanceRecord
	err := model.DB.Model(&recordItem).Where("maintenance_id = ?", input.MaintenanceID).
		Find(&recordList).
		Error

	return recordList, err
}

/**
 * 获取响应数据中的列表的一项数据
 */
func (steps *GetMaintenanceSteelSteps) GetSteelResItem(specificationItem specificationinfo.SpecificationInfo, recordItem *maintenance_record.MaintenanceRecord) maintenance.GetMaintenanceSteelResItem {
	i := maintenance.GetMaintenanceSteelResItem{
		Id:             specificationItem.ID,
		Specification:  specificationItem.GetSelfSpecification(),
		ReceivedTotal:  1,
		ReceivedWeight: specificationItem.Weight,
	}
	if steps.IsStored(recordItem.State) {
		i.StoredTotal = 1
		i.StoredWeight = specificationItem.Weight
	}
	return i
}
