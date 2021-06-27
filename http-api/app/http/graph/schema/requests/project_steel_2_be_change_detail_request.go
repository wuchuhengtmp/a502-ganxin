/**
 * @Desc    待修改武钢详细信息验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/27
 * @Listen  MIT
 */
package requests

import (
	"context"
	"fmt"
	"http-api/app/http/graph/auth"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/models/specificationinfo"
	"http-api/app/models/steels"
	"http-api/pkg/model"
)

func GetProjectSteel2BeChangeDetailRequest(ctx context.Context, input graphModel.ProjectSteel2BeChangeInput) error {
	steps := StepsForProject{}
	for _, identifier := range input.IdentifierList {
		// 检验有没有这根型钢
		if err := steps.CheckHasSteel(ctx, identifier); err != nil {
			return err
		}
		// 检验型钢是否属于我管理的项目下
		steelItem, err := steps.CheckSteelBelong2MyProject(ctx, identifier)
		if  err != nil {
			return err
		}
		// 检验状态是否是项目的合法状态之一
		if err := steps.CheckSteelState(steelItem.State); err != nil {
			return fmt.Errorf("标识码为:%s 的型钢状态为:%s 不能修改", identifier, steels.StateCodeMapDes[steelItem.State])
		}
	}
	// 检验过滤的是否是项目的合法状态之一
	if input.State != nil && steps.CheckSteelState(*input.State) != nil {
		return fmt.Errorf("状态为:%d 不是项目的合法状态", input.State)
	}
	// 检验有没这个规格id
	me := auth.GetUser(ctx)
	if input.SpecificationID != nil {
		err := model.DB.Model(&specificationinfo.SpecificationInfo{}).
			Where("id = ?", *input.SpecificationID).
			Where("company = ?", me.CompanyId).
			First(&specificationinfo.SpecificationInfo{}).
			Error
		if err != nil {
			return fmt.Errorf("规格id为：%d 的规格不存在", *input.SpecificationID)
		}
	}

	return nil
}