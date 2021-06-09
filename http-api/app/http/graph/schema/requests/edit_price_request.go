/**
 * @Desc    编辑价格请求验证器
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/9
 * @Listen  MIT
 */
package requests

import (
	"context"
	"fmt"
)

func ValidateEditPriceRequest(ctx context.Context, price float64) error {
	if price < 0.01 {
		return fmt.Errorf("价格不能小于0.01")
	}

	return nil
}
