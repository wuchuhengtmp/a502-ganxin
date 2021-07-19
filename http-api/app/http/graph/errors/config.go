/**
 * @Desc    错误处理
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/28
 * @Listen  MIT
 */
package errors

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"http-api/pkg/logger"
)

const (
	InvalidErrCode    = 4000 // 验证错误码
	AccessDenyErrCode = 4100 // 权限限制错误码
	ServerErrCode     = 5000 // 服务器错误码
	DeviceDeniedCode  = 4200 // 设备禁用错误码
)

var CodeMapDes map[int64]string

func init() {
	CodeMapDes = make(map[int64]string)
	CodeMapDes[InvalidErrCode] = "验证错误码"
	CodeMapDes[AccessDenyErrCode] = "权限限制错误码"
	CodeMapDes[ServerErrCode] = "服务器错误码"
	CodeMapDes[DeviceDeniedCode] = "设备禁用错误码"
}

/**
 * 没有token或token无效
 */
func InvalidToken(ctx context.Context) (bool, error) {
	err := &gqlerror.Error{
		Path:    graphql.GetPath(ctx),
		Message: "没有token或token无效",
		Extensions: map[string]interface{}{
			"code": InvalidErrCode,
		},
	}

	return false, err
}

/**
 * 无权调用
 */
func AccessDenied(ctx context.Context, msg string) (bool, error) {
	err := &gqlerror.Error{
		Path:    graphql.GetPath(ctx),
		Message: fmt.Sprintf("%s", msg),
		Extensions: map[string]interface{}{
			"code": 4100,
		},
	}

	return false, err
}

/**
 * 设备禁用
 */
func DeviceDenied(ctx context.Context, msg string) (bool, error) {
	err := &gqlerror.Error{
		Path:    graphql.GetPath(ctx),
		Message: fmt.Sprintf("%s", msg),
		Extensions: map[string]interface{}{
			"code": DeviceDeniedCode,
		},
	}

	return false, err
}

/**
 * 验证错误
 */
func ValidateErr(ctx context.Context, err error) error {
	gqlErr := &gqlerror.Error{
		Path:    graphql.GetPath(ctx),
		Message: err.Error(),
		Extensions: map[string]interface{}{
			"code": AccessDenyErrCode,
		},
	}

	return gqlErr
}

const (
	ServerErrorMsg = "操作出错，请联系管理员"
)

/**
 * 后台出现错误
 */
func ServerErr(ctx context.Context, err error) error {
	logger.LogError(err)
	gqlErr := &gqlerror.Error{
		Path:    graphql.GetPath(ctx),
		Message: err.Error(),
		Extensions: map[string]interface{}{
			"code": ServerErrCode,
		},
	}

	return gqlErr
}
