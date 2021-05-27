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
	"github.com/gorilla/mux"
	"http-api/app/http/controllers"
	"net/http"
)

func RegisterWebRoutes(r *mux.Router)  {
	pc := new (controllers.PagesController)
	r.NotFoundHandler = http.HandlerFunc(pc.NotFound)
	r.PathPrefix("/uploads/").Handler(http.FileServer(http.Dir("./public")))

	r.HandleFunc("/", pc.Home).Methods("GET").Name("home")
	// todo 强制text/html 会造成静态图片不能正常显示的情况
	//r.Use(middlewares.ForceHTML)
}
