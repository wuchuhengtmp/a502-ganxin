/**
 * @Desc    删除公司
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/2
 * @Listen  MIT
 */
package mutation_resolver

import (
	"context"
	"fmt"
	"http-api/app/http/graph/errors"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/companies"
)

func (m *MutationResolver)DeleteCompany(ctx context.Context, id int64) (bool, error) {
	requestValitation := requests.DeleteCompanyRequest{}
	err := requestValitation.ValidateDeleteCompanyRequest(id)
	if err != nil {
		return false, errors.ValidateErr(ctx, err)
	}
	companiesModel := companies.Companies{}
	err = companiesModel.DeleteById(int64(id))
	if err != nil {
		return false, errors.ServerErr(ctx, fmt.Errorf("删除失败！请联系管理员"))
	}

	return  true, nil
}
