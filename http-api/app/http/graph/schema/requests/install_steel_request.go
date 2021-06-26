/**
 * @Desc    安装型钢请求验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/26
 * @Listen  MIT
 */
package requests

import (
	"context"
	"fmt"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/models/steels"
)

func ValidateInstallSteelRequest(ctx context.Context, input graphModel.InstallLocationInput) error {
	// 检验型钢是否存在
	steps := StepsForProject{}
	if err := steps.CheckHasSteel(ctx, input.Identifier); err != nil {
		return err
	}
	// 检验型钢在不在项目中
	if err := steps.CheckIsProjectSteel(ctx, input.Identifier); err != nil {
		return err
	}
	// 检验型钢在不在我管理的项目中
	orderSpecificationSteelItem, err := steps.CheckSteelBelong2MyProject(ctx, input.Identifier);
	if  err != nil {
		return err
	}
	// 检验否是待定使用状态
	if orderSpecificationSteelItem.State != steels.StateProjectWillBeUsed {
		return fmt.Errorf("当前型钢的状态为: %s, 不能安装", steels.StateCodeMapDes[orderSpecificationSteelItem.State])
	}
	// 检验安装码是不是被占用了
	if err := steps.CheckLocationExists(input.Identifier, input.LocationCode); err != nil {
		return err
	}
	// 检验型钢归属的项目的项目管理员是不是我
	if err := steps.CheckSteelBelong2Me(ctx, input.Identifier); err != nil {
		return err
	}

	return nil
}
