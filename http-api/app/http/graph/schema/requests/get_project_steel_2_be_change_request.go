/**
 * @Desc    获取待修改武钢信息请求验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/27
 * @Listen  MIT
 */
package requests

import (
	"context"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/models/steels"
	"http-api/pkg/model"
)

func GetProjectSteel2BeChangeRequest(ctx context.Context, input graphModel.GetProjectSteel2BeChangeInput) error {
	steps := StepsForProject{}
	// 检验有没有这根型钢
	if err := steps.CheckHasSteel(ctx, input.Identifier); err != nil {
		return err
	}
	// 检验型钢是否归属于我
	if err := steps.CheckIsBelongMeByIdentifier(ctx, input.Identifier); err != nil {
		return err
	}
	// 检验型钢状态是否合法
	steelItem := steels.Steels{}
	if err := model.DB.Model(&steelItem).Where("identifier = ?", input.Identifier).First(&steelItem).Error; err != nil {
		return err
	}
	if err := steps.CheckSteelState(steelItem.State); err != nil {
		return err
	}

	return nil
}
