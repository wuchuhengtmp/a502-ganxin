/**
 * @Desc    删除公司员工解析器
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/4
 * @Listen  MIT
 */
package mutation_resolver

import (
	"context"
	"http-api/app/http/graph/errors"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/companies"
)

/**
 * 删除公司员工解析器
 */
func (m *MutationResolver)DeleteCompanyUser(ctx context.Context, uid int64) (bool, error) {
	validator := requests.DeleteCompanyUserRequest{}
	err := validator.ValidateDeleteCompanyUserRequest(ctx, uid)
	if err != nil {
		return false, errors.ValidateErr(ctx, err)
	}
	if err := companies.DeleteCompanyUserByUid(ctx,uid); err != nil {
		return false, errors.ServerErr(ctx, err)
	}

	return true, nil
}
