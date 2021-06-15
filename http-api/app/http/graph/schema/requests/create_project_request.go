/**
 * @Desc    创建项目请求验证器
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/15
 * @Listen  MIT
 */
package requests

import (
	"context"
	"fmt"
	"github.com/thedevsaddam/govalidator"
	"http-api/app/http/graph/model"
	"http-api/app/models/roles"
	"http-api/app/models/users"
	"time"
)

func ValidateCreateProjectRequest(ctx context.Context, input model.CreateProjectInput) error {
	rules := govalidator.MapData{
		"name":    []string{"required", "min:1"},
		"address": []string{"required"},
		"city":    []string{"required"},
	}
	message := govalidator.MapData{
		"name":    []string{"required:项目名不能为空"},
		"address": []string{"required:地址"},
		"city":    []string{"required:城市不能为空"},
	}
	opts := govalidator.Options{
		Data:          &input,
		Rules:         rules,
		Messages:      message,
		TagIdentifier: "json",
	}
	res := govalidator.New(opts).ValidateStruct()
	if len(res) > 0 {
		for _, fieldErrors := range res {
			for _, err := range fieldErrors {
				return fmt.Errorf("%s", err)
			}
		}
	}
	if input.StartAt.Unix() < time.Now().Unix() {
		return fmt.Errorf("开始时间不能小于当前时间")
	}
	if len(input.LeaderIDS) < 1 {
		return fmt.Errorf("负责人至少要一位")
	}
	for _, uid := range input.LeaderIDS {
		u := users.Users{ID: uid}
		if err := u.GetSelfById(uid); err != nil {
			return fmt.Errorf("没有用户id为:%d", uid)
		}
		// 用户角色鉴定和是否待业
		role, _ := u.GetRole()
		switch role.Tag {
		case roles.RoleAdmin:
			return fmt.Errorf("用户id:%d, 超级管理员不能作为项目管理员加入", uid)
		case roles.RoleCompanyAdmin:
			return fmt.Errorf("用户id:%d, 公司管理员不能作为项目管理员加入", uid)
		case roles.RoleRepositoryAdmin:
			return fmt.Errorf("用户id:%d, 仓库管理员不能作为项目管理员加入", uid)
		case roles.RoleMaintenanceAdmin:
			return fmt.Errorf("用户id:%d, 维修管理员不能作为项目管理员加入", uid)
		}
	}

	return nil
}
