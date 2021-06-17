/**
 * @Desc    The controllers is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @DATE    2021/4/27
 * @Listen  MIT
 */
package api

import (
	"http-api/pkg/response"
	"net/http"
)

type AuthorizationController struct {}

/**
 * 通知前端重定向到登录页面
 */
func NotifyFrontEndRedirectLogin(token string, w http.ResponseWriter) {
	res := response.Error{
		ErrorCode: response.ErrorCodes.RediRectLogin,
		Errors: map[string][]string {
			"accessToken": {token},
		},
	}
	res.ResponseByHttpWriter(w)
}
