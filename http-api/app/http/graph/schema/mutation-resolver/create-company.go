/**
 * @Desc    创建公司解析器
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/28
 * @Listen  MIT
 */
package mutation_resolver

import (
	"context"
	"http-api/app/http/graph/model"
)

func (m *MutationResolver)CreateCompany(ctx context.Context, input model.CreateCompanyInput) (*model.CreateCompanyRes, error) {
	res := model.CreateCompanyRes{ }
	res.LogoFile = &model.SingleUploadRes{}
	res.BackgroundFile = &model.SingleUploadRes{}

	return &res, nil
}
