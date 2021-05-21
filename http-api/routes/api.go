/**
 * @Desc    The routes is part of http-api
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @DATE    2021/4/27
 * @Listen  MIT
 */
package routes


/*
 * api 路由
 *
 */

import (
	"github.com/gorilla/mux"
	"http-api/app/http/controllers/api"
)

func RegisterApiRoutes(r *mux.Router) {
	rp := r.PathPrefix("/api").Subrouter()
	a := new (api.AuthorizationController)
	// 获取token
	rp.HandleFunc("/authorizations", a.Create).Methods("POST").Name("authorization.create")
}

