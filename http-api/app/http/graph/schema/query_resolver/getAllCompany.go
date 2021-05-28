/**
 * @Desc    The query_resolver is part of http-api
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/28
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"http-api/app/http/graph/model"
)

type QueryResolver struct { }

/**
 * 获取全部公司解析器
 */
func (q *QueryResolver)GetAllCompany(ctx context.Context) ([]*model.CreateCompanyRes, error) {
	var companies []*model.CreateCompanyRes
	return companies, nil
}
