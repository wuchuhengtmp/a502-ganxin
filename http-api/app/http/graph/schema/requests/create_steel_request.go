/**
 * @Desc    型钢入库请求验证器
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/11
 * @Listen  MIT
 */
package requests

import (
	"context"
	"fmt"
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/model"
	"http-api/app/models/codeinfo"
	"http-api/app/models/repositories"
	"http-api/app/models/specificationinfo"
)

func ValidateCreateSteelRequest(ctx context.Context, input model.CreateSteelInput) error {
	me := auth.GetUser(ctx)
	 d := repositories.Repositories{ID:  input.RepositoryID}
	 if err := d.GetSelf(); err != nil {
		return fmt.Errorf("没有这个仓库")
	 }
	 if d.CompanyId != me.CompanyId {
		 return fmt.Errorf("这个仓库与你不是同一家公司归属，你无权操作")
	 }
	 s := specificationinfo.SpecificationInfo{ID: input.SpecificationID}
	 if err := s.GetSelf(); err != nil {
		 return fmt.Errorf("没有这个规格")
	 }
	 if s.CompanyId != me.CompanyId {
		 return fmt.Errorf("这个规格与你不是同一家公司的归属，你无权操作")
	 }
	 c := codeinfo.CodeInfo {ID: input.MaterialManufacturerID}
	 if err := c.GetSelf(); err != nil {
		 return fmt.Errorf("没有这家材料商")
	 }
	 if c.Type != codeinfo.MaterialManufacturer {
		 return fmt.Errorf("没有这家材料商")
	 }
	 if c.CompanyId != me.CompanyId {
		 return fmt.Errorf("材料商与你不是同一家公司的归属，你无权操作")
	 }
	 mf := codeinfo.CodeInfo{ID: input.ManufacturerID}
	 if err := mf.GetSelf(); err != nil {
		return fmt.Errorf("没有这个制作商")
	 }
	 if mf.Type != codeinfo.Manufacturer {
		 return fmt.Errorf("没有这个制作商")
	 }
	 if mf.CompanyId != me.CompanyId {
		 return fmt.Errorf("制作商与你不是同一家公司的归属，你无权操作")
	 }

	return nil
}
