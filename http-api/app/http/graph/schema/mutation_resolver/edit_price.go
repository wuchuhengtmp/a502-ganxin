/**
 * @Desc    The mutation_resolver is part of http-api
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/9
 * @Listen  MIT
 */
package mutation_resolver

import (
	"context"
	"fmt"
	"http-api/app/http/graph/errors"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/configs"
	"strconv"
)

func (*MutationResolver)EditPrice(ctx context.Context, price float64) (float64, error) {
	if err := requests.ValidateEditPriceRequest(ctx, price); err != nil {
		return 0, errors.ValidateErr(ctx, err)
	}
	cf := configs.Configs{Value: fmt.Sprintf("%.4f", price)}
	if err := cf.EditPrice(ctx); err != nil {
		return 0, errors.ServerErr(ctx, err)
	}
	price, _ = strconv.ParseFloat(cf.Value, 64)

	return price, nil
}