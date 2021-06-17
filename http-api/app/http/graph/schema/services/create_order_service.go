/**
 * @Desc    创建订单服务
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/16
 * @Listen  MIT
 */
package services

import (
	"context"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"http-api/app/http/graph/auth"
	gqlModel "http-api/app/http/graph/model"
	"http-api/app/http/graph/util/helper"
	msg2 "http-api/app/models/msg"
	"http-api/app/models/order_specification"
	"http-api/app/models/orders"
	"http-api/app/models/projects"
	"http-api/app/models/repositories"
	"http-api/app/models/repository_leader"
	"http-api/app/models/specificationinfo"
	"http-api/app/models/users"
	"http-api/pkg/model"
	"time"
)

func CreateOrder(ctx context.Context, input gqlModel.CreateOrderInput) (*orders.Order, error) {
	o := orders.Order{}
	return &o, model.DB.Transaction(func(tx *gorm.DB) error {
		me := auth.GetUser(ctx)
		o.CreateUid = me.ID
		o.State = orders.StateToBeConfirmed
		o.ProjectId = input.ProjectID
		o.RepositoryId = input.RepositoryID
		o.ExpectedReturnAt = input.ExpectedReturnAt
		o.PartList = input.PartList
		o.Remark = input.Remark
		o.ProjectId = input.ProjectID
		t := time.Now()
		year, month, day := t.Date()
		h := t.Hour()
		i := t.Minute()
		s := t.Second()
		n := t.Nanosecond()
		o.OrderNo = fmt.Sprintf("%d%d%d%d%d%d%d", year, month, day, h, i, s, n)

		if err := tx.Create(&o).Error; err != nil {
			return err
		}
		// 创建订单规格单
		for _, sp := range input.SteelList {
			spc := specificationinfo.SpecificationInfo{}
			if err := tx.Model(&spc).Where("id = ?", sp.SpecificationID).First(&spc).Error; err != nil {
				return err
			}
			oo := order_specification.OrderSpecification{
				SpecificationId: spc.ID,
				Total:           sp.Total,
				Specification:   spc.GetSelfSpecification(),
				OrderId:         o.Id,
			}
			if err := tx.Create(&oo).Error; err != nil {
				return err
			}
			//  消息通知
			if err := _createOrderMsg(tx, &o); err != nil {
				return err
			}
		}

		return nil
	})
}

/**
 * 创建新订单消息
 */
func _createOrderMsg(tx *gorm.DB, o *orders.Order) error {
	userList, err := repository_leader.RepositoryLeader{}.GetLeaders(tx)
	if err != nil { return err }
	r := repositories.Repositories{ID: o.RepositoryId}
	if err := r.GetSelf(); err != nil {
		return err
	}
	p := projects.Projects{}
	if err := tx.Model(&p).Where("id = ?", o.ProjectId).First(&p).Error; err != nil {
		return err
	}
	timeStr := helper.Time2Str(time.Now())
	projectUser := users.Users{}
	if err := tx.Model(&projectUser).Where("id = ?", o.CreateUid).First(&projectUser).Error; err != nil {
		return err
	}
	var totalSteels struct {
		TotalSteels int64 `json:"totalSteels"`
		TotalWeight float64
	}
	// 获取数量
	err = tx.Model(&order_specification.OrderSpecification{}).
		Debug().
		Select("sum(total) as TotalSteels").
		Where("order_id = ?", o.Id).
		Scan(&totalSteels).Error
	if err != nil { return err }
	// 获取重量
	var oss []*order_specification.OrderSpecification
	tx.Model(&order_specification.OrderSpecification{}).Where("order_id = ?", o.Id).Find(&oss)
	for _, specificationItem := range oss {
		spec := specificationinfo.SpecificationInfo{}
		if err := tx.Model(&spec).Where("id = ?", specificationItem.SpecificationId).First(&spec).Error; err != nil {
			return err
		}
		totalSteels.TotalWeight += spec.Weight * float64(specificationItem.Total)
	}

	extends, _ := json.Marshal(o)
	for _, user := range userList {
		msg := msg2.Msg{
			Content: fmt.Sprintf(
				"%s 项目的管理员%s于%s, 创建需求单:%s, 总数:%d根, %.2f吨,请点击查看详情",
				p.Name,
				projectUser.Name,
				timeStr,
				o.OrderNo,
				totalSteels.TotalSteels,
				totalSteels.TotalWeight,
			),
			Type:    msg2.CreateOrderType,
			Uid:     user.ID,
			Extends: string(extends),
			IsRead:  false,
		}
		// 创建消息
		if err := tx.Create(&msg).Error; err != nil { return err }
		// 推送消息
		_: msg.Push()
	}

	return nil
}
