/**
 * @Desc    获取价格解析器
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/9
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"http-api/app/models/configs"
)

func (*QueryResolver) GetPrice(ctx context.Context) (float64, error) {
	cf := configs.Configs{}

	return cf.GetPrice(), nil
}
