/**
 * @Desc    获取日志列表请求验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/8
 * @Listen  MIT
 */
package requests

import (
	"context"
	"fmt"
	graphModel "http-api/app/http/graph/model"
)

func ValidateGetLogListRequest(ctx context.Context, input graphModel.GetLogListInput) error {
	// 检验分页数据
	if !input.IsShowAll {
		if input.Page == nil {
			return fmt.Errorf("页码不能为空")
		}
	}
	return nil
}
