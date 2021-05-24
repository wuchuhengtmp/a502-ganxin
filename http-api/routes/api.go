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
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"http-api/app/http/graph/generated"
	"http-api/app/http/graph/schema"
	"http-api/app/http/middlewares"
)

func RegisterApiRoutes(r *mux.Router) {
	//rp := r.PathPrefix("/api").Subrouter()
	//a := new (api.AuthorizationController)

	// graphql 沙盒
	r.Handle("/graphql", playground.Handler("GraphQL playground", "/query"))
	// graphql 接口
	r.Handle("/query", handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &schema.Resolver{}})))
	r.Use(middlewares.AllowCORS)
}
