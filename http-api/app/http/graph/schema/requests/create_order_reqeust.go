/**
 * @Desc    创建需求单验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/16
 * @Listen  MIT
 */
package requests

import (
	"context"
	"fmt"
	"github.com/thedevsaddam/govalidator"
	"http-api/app/http/graph/auth"
	grapModel "http-api/app/http/graph/model"
	"http-api/app/models/orders"
	"http-api/app/models/specificationinfo"
	"http-api/app/models/steels"
	"http-api/pkg/model"
	"time"
)

func ValidateCreateOrderValidate(ctx context.Context, input grapModel.CreateOrderInput) error {
	me := auth.GetUser(ctx)
	if input.ExpectedReturnAt.UnixNano() < time.Now().UnixNano() {
		return fmt.Errorf("开始时间不能早于当前时间")
	}
	// 型钢需求验证
	if len(input.SteelList) == 0 {
		return fmt.Errorf("型钢不能为空")
	}
	for _, s := range input.SteelList {
		sp := specificationinfo.SpecificationInfo{ID: s.SpecificationID}
		if err := sp.GetSelf(); err != nil || sp.CompanyId != me.CompanyId {
			return fmt.Errorf("没有规格id为:%d", s.SpecificationID)
		}
		// 库存是否足够
		var total int64
		model.DB.Model(&steels.Steels{}).
			Where("state = ? AND specification_id = ?", steels.StateInStore, s.SpecificationID).
			Count(&total)
		if total < s.Total {
			return fmt.Errorf("型钢规格:%s,库存不足%d", sp.GetSelfSpecification(), s.Total)
		}
		// 减去待发货数量再比较存量
		confirmTotal, err := orders.GetConfirmSteelTotalBySpecificationId(s.SpecificationID)
		if err != nil {
			return err
		}
		if total - confirmTotal < s.Total {
			return fmt.Errorf("型钢规格:%s,库存不足%d", sp.GetSelfSpecification(), s.Total)
		}
	}

	strId := fmt.Sprintf("%d", auth.GetUser(ctx).ID)
	rules := govalidator.MapData{
		"projectId":  []string{"required", "isCompanyProject:" + strId},
		"repositoryId": []string{"required", "isCompanyRepository:" + strId},
	}
	opts := govalidator.Options{
		Data:          &input,
		Rules:         rules,
		TagIdentifier: "json",
	}

	return Validate(opts)
}
