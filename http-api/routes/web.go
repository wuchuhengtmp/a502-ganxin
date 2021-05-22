package routes

/**
 * @Desc    The main is part of http-api
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/21
 * @Listen  MIT
 */
import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"http-api/app/http/controllers"
	"http-api/app/http/graph/generated"
	"http-api/app/http/graph/schema"
	"http-api/app/http/middlewares"
	"net/http"
)

func RegisterWebRoutes(r *mux.Router)  {
	pc := new (controllers.PagesController)
	r.NotFoundHandler = http.HandlerFunc(pc.NotFound)
	r.HandleFunc("/", pc.Home).Methods("GET").Name("home")

	r.Handle("/graphql", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/query", handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &schema.Resolver{}})))
	r.Use(middlewares.ForceHTML)
}
