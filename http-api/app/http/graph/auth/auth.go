/**
 * @Desc    graph中间件
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/27
 * @Listen  MIT
 */
package auth

import (
	"context"
	"fmt"
	"http-api/app/models/users"
	"http-api/pkg/jwt"
	"http-api/pkg/model"
	"net/http"
)

const userCtxKey string = "userKey"
const isDeviceKey string = "isDeviceKey"
const tokenKey string = "hasTokenKey"

func GraphMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		user := users.Users{ }
		isDevice := false
		// 把用户信息写注入到上下文中
		if len(token) > 7 {
			db := model.DB;
			payload, err := jwt.ParseByTokenStr(token[7:])
			if payload!= nil {
				isDevice = payload.IsDevice
			}
			if err == nil {
				userModel := users.Users{}
				err = db.Model(&userModel).Where("id = ?", payload.Uid).First(&userModel).Error
				if err == nil {
					user = userModel
				}
			}
		}
		ctx := context.WithValue(r.Context(), userCtxKey, user)
		ctx = context.WithValue(ctx, isDeviceKey, isDevice)
		ctx = context.WithValue(ctx, tokenKey, token)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

/**
 * 获取上下文的用户信息
 */
func GetUser(ctx context.Context) *users.Users {
	raw, _ := ctx.Value(userCtxKey).(users.Users)
	if raw.Id == 0 {
		return nil
	}

	return &raw
}
/**
 * 是否是手持设备
 */
func IsDevice(ctx context.Context) bool {
	raw, _ := ctx.Value(isDeviceKey).(bool)

	return raw
}

func Token(ctx context.Context) string  {
	raw, _ := ctx.Value(tokenKey).(string)

	return raw
}

/**
 * 验证token
 */
func ValidateToken(ctx context.Context) error {
	token := Token(ctx)
	if len(token) == 0 {
		return fmt.Errorf("token 不能为空")
	}
	if len(token) > 7 {
		db := model.DB;
		payload, err := jwt.ParseByTokenStr(token[7:])
		if err == nil {
			userModel := users.Users{}
			err = db.Model(&userModel).Where("id = ?", payload.Uid).First(&userModel).Error
			if err == nil {
				return nil
			}
		}
	} else {
		return fmt.Errorf("不是Bearer 格式的token")
	}

return fmt.Errorf("无效token")
}
