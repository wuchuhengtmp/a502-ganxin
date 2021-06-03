/**
 * @Desc    获取公司人员解析器
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/3
 * @Listen  MIT
 */

package query_resolver

import (
	"context"
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/model"
	"http-api/app/models/companies"
)

func (q *QueryResolver) GetCompanyUser(ctx context.Context) ([]*model.UserItem, error) {
	me := auth.GetUser(ctx)
	res, err := companies.GetCompanyItemsResById(me.CompanyId)

	return res, err
}
