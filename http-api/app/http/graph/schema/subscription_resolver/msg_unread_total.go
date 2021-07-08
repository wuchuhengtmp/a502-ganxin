/**
 * @Desc    订阅未读消息总量
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/8
 * @Listen  MIT
 */
package subscription_resolver

import (
	"context"
	"fmt"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/models/users"
	"http-api/pkg/model"
	"time"
)

type SubscriptionResolver struct {
}

func (*SubscriptionResolver) MsgUnreadTotal(ctx context.Context, input graphModel.MsgUnreadTotalInput) (<-chan int64, error) {
	foo := make(chan int64, 1)
	err := model.DB.Model(&users.Users{}).Where("id = ?", input.UID).First(&users.Users{}).Error
	if err != nil {
		if err.Error() == "record not found" {
			return foo, fmt.Errorf("id为: %d 的用户不存在", input.UID)
		}
		return foo, err
	}

	go func() {
		ticker := time.NewTicker(time.Second * 1)
		for {
			<-ticker.C
			foo <- time.Now().Unix()
		}

	}()

	return foo, nil

}
