/**
 * @Desc    删除材料商请求验证器
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/7
 * @Listen  MIT
 */
package requests

import (
	"context"
	"fmt"
	"github.com/thedevsaddam/govalidator"
	"http-api/app/http/graph/auth"
	"http-api/app/models/codeinfo"
	"http-api/app/models/steels"
)

func ValidateDeleteMaterialManufacturerRequest(ctx context.Context, id int64) error {
	rules := govalidator.MapData{
		"id": []string{"isCodeInfoId"},
	}
	opts := govalidator.Options{
		Data: &struct { Id int64 `json:"id"` }{id},
		Rules: rules,
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
	// 归属检验
	c := codeinfo.CodeInfo{ID: id}
	_ = c.GetSelf()
	me := auth.GetUser(ctx)
	if me.CompanyId != c.CompanyId {
		return fmt.Errorf("该材料商与您不是归属于同一家公司名下，您无权删除")
	}
	s := steels.Steels{}
	ss, _ := s.GetListByMMID(id)
	if len(ss) > 0  {
		return fmt.Errorf("该材料商已被使用在型钢材，无法删除")
	}

	return nil
}