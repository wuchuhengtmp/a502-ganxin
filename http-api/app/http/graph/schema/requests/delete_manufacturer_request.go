/**
 * @Desc    删除制造商请求验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/8
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

func ValidateDeleteManufacturerRequest(ctx context.Context, id int64) error {
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

	me := auth.GetUser(ctx)
	c := codeinfo.CodeInfo{ID: id}
	_ = c.GetSelf()
	if me.CompanyId != c.CompanyId {
		return fmt.Errorf("要删除的制造商与您不是归属于同一家公司，您无权删除")
	}
	if c.Type != codeinfo.Manufacturer {
		return fmt.Errorf("这不是制造商。无法删除")
	}
	s := steels.Steels{}
	ss, _ := s.GetListByManufacturerId(id)
	if len(ss) > 0 {
		return fmt.Errorf("这家商家的信息已经用于型钢中，不可删除")
	}

	return nil
}