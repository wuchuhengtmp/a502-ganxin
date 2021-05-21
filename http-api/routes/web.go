package routes

import (
	"github.com/gorilla/mux"
	"http-api/app/http/controllers"
	"http-api/app/http/middlewares"
	"net/http"
)

func RegisterWebRoutes(r *mux.Router)  {
	pc := new (controllers.PagesController)
	r.NotFoundHandler = http.HandlerFunc(pc.NotFound)
	r.HandleFunc("/", pc.Home).Methods("GET").Name("home")

	r.Use(middlewares.ForceHTML)
}
