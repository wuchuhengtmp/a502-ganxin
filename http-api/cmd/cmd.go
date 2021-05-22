/**
 * @Desc    The cmd is part of http-api
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @DATE    2021/4/27
 * @Listen  MIT
 */
package cmd

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/urfave/cli"
	"http-api/app/http/middlewares"
	"http-api/bootstrap"
	"http-api/config"
	pkgC "http-api/pkg/config"
	"http-api/pkg/logger"
	"http-api/pkg/model"
	"http-api/seeders"
	"log"
	"net/http"
	"os"
)

func Run ()  {
	app := cli.NewApp()
	app.Name = "the back end for A502-钢型平台后端服务"
	app.Version = "0.0.1"
	app.Usage = "A502-钢型平台后端服务"
	app.Commands = []cli.Command{
		cli.Command{
			Name: "http_api",
			Usage: "up the http server for api",
			Action: RunWeb,
		},
		cli.Command{
			Name: "seeds",
			Usage: "To generate seeds within the database",
			Action: RunMigrateSeed,
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Println(err)
	}
}

func init()  {
	config.Initialize()
}

var router = mux.NewRouter().StrictSlash(true)

// 启动web服务
func RunWeb (c *cli.Context)  {
	bootstrap.SetupDB()
	router = bootstrap.SetupRoute()
	done := make(chan bool)
	go http.ListenAndServe(":" + pkgC.GetString("APP_PORT"), middlewares.RemoveTrailingSlash(router))
	fmt.Printf(`
		🚀 🚀 🚀 Server is running!
		Listening on port %s
		Explore at http://localhost:%s
		Explore graphql at http://localhost:%s/graphql

		`,
	pkgC.GetString("APP_PORT"), pkgC.GetString("APP_PORT"), pkgC.GetString("APP_PORT"))
	<-done
}

func RunMigrateSeed(c *cli.Context) {
	bootstrap.SetupDB()
	for _, seed := range seeders.All() {
		log.Println(seed.Name)
		err := seed.Run(model.DB)
		logger.LogError(err)
	}
}
