/**
 * @Desc    错误处理
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/28
 * @Listen  MIT
 */
package errors

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

/**
 * 没有token或token无效
 */
func InvalidToken(ctx context.Context) (bool, error)  {
	 err := &gqlerror.Error {
		Path: graphql.GetPath(ctx),
		Message: "没有token或token无效",
		Extensions: map[string]interface{}{
			"code": 4000,
		},
	}

	return false, err
}

/**
 * 无权调用
 */
func AccessDenied(ctx context.Context) (bool, error)  {
	err := &gqlerror.Error {
		Path: graphql.GetPath(ctx),
		Message: "拒绝访问",
		Extensions: map[string]interface{}{
			"code": 4100,
		},
	}

	return false, err
}

/**
 * 验证错误
 */
func ValidateErr(ctx context.Context, err error)  error  {
	gqlErr := &gqlerror.Error {
		Path: graphql.GetPath(ctx),
		Message: err.Error(),
		Extensions: map[string]interface{}{
			"code": 4200,
		},
	}

	return gqlErr
}
