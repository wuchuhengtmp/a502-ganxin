/**
 * @Desc    The controllers is part of http-api
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @DATE    2021/4/27
 * @Listen  MIT
 */
package api

import (
	"http-api/app/models/users"
	"http-api/app/requests/api"
	"http-api/app/util/wechat"
	"http-api/pkg/jwt"
	"http-api/pkg/response"
	"net/http"
)

type AuthorizationController struct {}

/**
 * 登录授权
 */
func (*AuthorizationController) Create (w http.ResponseWriter, r *http.Request)  {
	defer r.Body.Close()
	validateData := api.AuthorizationLoginRequest{}
	errs := validateData.ValidateAuthorizationLoginRequest(r.Body)
	if len(errs) > 0 {
		errRes := response.Error{
			Errors: errs,
			ErrorCode: response.ErrorCodes.LoginFail,
		}
		errRes.ResponseByHttpWriter(w)

	} else {
		wc := wechat.GetWCInstance()
		res, err := wc.GetAuth().Code2Session(validateData.Code)
		var wechatModel = users.WechatModel{}
		if err != nil {
			errRes := response.Error{
				Errors: map[string][]string{
					"code": { err.Error() },
				} ,
				ErrorCode: response.ErrorCodes.LoginFail,
			}
			errRes.ResponseByHttpWriter(w)
		} else if !wechatModel.IsUserByOpenId(res.OpenID) {
			// 添加新用户
			wechatModel.WcOpenId = res.OpenID
			wechatModel.WcSessionKey = res.SessionKey
			wechatModel.WcUnionId = res.UnionID
			user, _ := wechatModel.AddUser()
			token, _ := jwt.GenerateTokenByUID(user.ID)
			NotifyFrontEndRedirectLogin(token, w)

		} else {
			user, _ := wechatModel.GetUserByOpenId(res.OpenID)
			token, _ := jwt.GenerateTokenByUID(user.ID)
			if len(user.AvatarUrl) == 0 {
				NotifyFrontEndRedirectLogin(token, w)
			} else {
				resData := struct{ AccessToken string `json:"accessToken"` }{
					AccessToken: token,
				}
				response.SuccessResponse(resData, w)
			}
		}
	}
}

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
