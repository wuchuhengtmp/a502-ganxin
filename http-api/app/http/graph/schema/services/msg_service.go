/**
 * @Desc    消息服务
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/18
 * @Listen  MIT
 */
package services

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"http-api/app/http/graph/util/helper"
	"http-api/app/models/msg"
	"http-api/app/models/orders"
	"http-api/app/models/repositories"
	"http-api/app/models/users"
	"time"
)

/**
 * 创建或拒绝订单消息
 */
func CreateConfirmOrRejectOrderMsg(tx *gorm.DB, o *orders.Order) error {
	r := repositories.Repositories{}
	if err := tx.Model(&r).Where("id = ?", o.RepositoryId).First(&r).Error; err != nil {
		return err
	}
	confirmAdmin := users.Users{}
	if err := tx.Model(&confirmAdmin).Where("id = ?", o.ConfirmedUid).First(&confirmAdmin ).Error; err != nil {
		return err
	}
	state := fmt.Sprintf("仓库管理员 %s", confirmAdmin.Name)
	msgType := ""
	if o.State != orders.StateRejected && o.State != orders.StateConfirmed {
		return fmt.Errorf("系统错误，订单状态不正确")
	} else if o.State == orders.StateConfirmed {
		msgType = msg.ConfirmOrderType
		state += "确认订单"
	} else {
		msgType = msg.RejectOrderType
		state += fmt.Sprintf("拒绝订单,拒绝原因为: %s", o.RejectReason)
	}
	total, err := orders.GetTotal(tx, o)
	if err != nil {
		return err
	}
	weight, err := orders.GetWeight(tx, o)
	if err != nil {
		return err
	}

	t := helper.Time2Str(time.Now())
	c := fmt.Sprintf("%s 仓库于%s, %s, 订单编号为: %s, 总数: %d根，%.2f吨。", r.Name, t, state, o.OrderNo, total, weight)
	extends, _ := json.Marshal(o)
	userList, err := repositories.GetLeaders(tx, o.RepositoryId)
	for _, user := range userList {
		msgItem := msg.Msg{
			IsRead:  false,
			Content: c,
			Uid:     user.Id,
			Type:    msgType,
			Extends: string(extends),
		}
		if err := msgItem.CreateSelf(tx); err != nil {
			return err
		}
	}

	return nil
}