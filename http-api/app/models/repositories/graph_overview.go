/**
 * @Desc    graphql仓库概览字段解析
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/16
 * @Listen  MIT
 */
package repositories

type GetRepositoryOverviewRes struct {
	Total  int64   `json:"total"`
	Weight float64 `json:"weight"`
}
