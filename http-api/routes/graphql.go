/**
 * @Desc    graphql 路由
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/27
 * @Listen  MIT
 */
package routes

import (
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/schema"
	"http-api/app/http/middlewares"
)

func RegisterGraphRoutes(r *mux.Router) {
	// graphql 沙盒
	r.Handle("/graphql", playground.Handler("GraphQL playground", "/query"))
	// graphql 接口
	r.Handle("/query", schema.Handler())
	r.Use(
		middlewares.AllowCORS,
		auth.GraphMiddleware,
	)
}