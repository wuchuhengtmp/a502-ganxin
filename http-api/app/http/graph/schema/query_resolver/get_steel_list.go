/**
 * @Desc    获取型钢列表解析器
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/11
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"http-api/app/http/graph/model"
	"http-api/app/models/codeinfo"
	"http-api/app/models/repositories"
	"http-api/app/models/specificationinfo"
	"http-api/app/models/steels"
)


func (*QueryResolver)GetSteelList(ctx context.Context, input model.PaginationInput) (*steels.GetSteelListRes, error) {
	s := steels.Steels{}
	list, err := s.GetPagination(ctx, input.Page, input.PageSize)
	if  err != nil {
		return nil, err
	}

	res := steels.GetSteelListRes {
		List: list,
		Total: s.GetTotal(ctx),
	}

	return &res, nil
}

type SteelItemResolver struct { }

func (SteelItemResolver) Specifcation(ctx context.Context, obj *steels.Steels) (*specificationinfo.SpecificationInfo, error) {
	return obj.GetSpecification()
}

func (SteelItemResolver)MaterialManufacturer(ctx context.Context, obj *steels.Steels) (*codeinfo.CodeInfo, error) {
	return obj.GetMaterialManufacturer()
}
func (SteelItemResolver)Manufacturer(ctx context.Context, obj *steels.Steels) (*codeinfo.CodeInfo, error) {
	return obj.GetManufacturer()
}
func (SteelItemResolver)Repository(ctx context.Context, obj *steels.Steels) (*repositories.Repositories, error) {
	return obj.GetRepository()
}
