/**
 * @Desc    graph中间件
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/27
 * @Listen  MIT
 */
package auth

import (
	"context"
	"http-api/app/models/users"
	"http-api/pkg/jwt"
	"http-api/pkg/model"
	"net/http"
)

const userCtxKey string = "userKey"
const isDeviceKey string = "isDeviceKey"

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
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

/**
 * 获取上下文的用户信息
 */
func GetUser(ctx context.Context) *users.Users {
	raw, _ := ctx.Value(userCtxKey).(users.Users)
	if raw.ID == 0 {
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
