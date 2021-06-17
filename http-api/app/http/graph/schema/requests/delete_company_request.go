/**
 * @Desc    删除公司请求验证
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/31
 * @Listen  MIT
 */
package requests

import (
	"fmt"
	"http-api/app/models/companies"
)

type DeleteCompanyRequest struct { }


func (data *DeleteCompanyRequest) ValidateDeleteCompanyRequest(companyId int64) error {
	companiesModel := companies.Companies{}
	err := companiesModel.GetSelfById(int64(companyId))
	if err != nil {
		return fmt.Errorf("没有这家公司")
	}

	return nil
}
