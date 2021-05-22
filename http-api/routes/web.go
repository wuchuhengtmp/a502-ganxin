package routes

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"http-api/app/http/controllers"
	"http-api/app/http/graphql"
	"http-api/app/http/graphql/generated"
	"http-api/app/http/middlewares"
	"net/http"
)

func RegisterWebRoutes(r *mux.Router)  {
	pc := new (controllers.PagesController)
	r.NotFoundHandler = http.HandlerFunc(pc.NotFound)
	r.HandleFunc("/", pc.Home).Methods("GET").Name("home")

	r.Handle("/graphql", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/query", handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graphql.Resolver{}})))
	r.Use(middlewares.ForceHTML)
}
