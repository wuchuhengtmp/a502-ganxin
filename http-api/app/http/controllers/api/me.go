/**
 * @Desc    The controllers is part of http-api
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @DATE    2021/4/27
 * @Listen  MIT
 */
package api

import (
	"http-api/app/requests/api"
	"http-api/pkg/response"
	"net/http"
)

type MeController struct {}


/**
 * 登录授权
 */
func (*MeController) Update (w http.ResponseWriter, r *http.Request)  {
	defer r.Body.Close()
	validateData := api.MeUpdateRequest{}
	errs := validateData.ValidateAuthorizationLoginRequest(r.Body)
	if len(errs) > 0 {
		errRes := response.Error{
			Errors: errs,
			ErrorCode: response.ErrorCodes.LoginFail,
		}
		errRes.ResponseByHttpWriter(w)
	} else {

	}
}