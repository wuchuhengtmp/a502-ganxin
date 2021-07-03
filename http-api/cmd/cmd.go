/**
 * @Desc    The cmd is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @DATE    2021/4/27
 * @Listen  MIT
 */
package cmd

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/urfave/cli/v2"
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
	"time"
)

func Run ()  {
	app := cli.NewApp()
	app.Name = "the back end for A502-钢型平台后端服务"
	app.Version = "0.0.1"
	app.Usage = "A502-钢型平台后端服务"
	app.Commands = []*cli.Command{
		&cli.Command{
			Name: "http_api",
			Usage: "up the http server for api",
			Action: RunWeb,
		},
		&cli.Command{
			Name: "seeds",
			Usage: "To generate seeds within the database",
			Action: RunMigrateSeed,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name: "force",
					Usage: "强制清除数据并重新生成，慎用!!!",
					Value: false,
					Aliases: []string{"f"},
				},
			},
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
func RunWeb (c *cli.Context) error {
	bootstrap.SetupDB()
	router = bootstrap.SetupRoute()
	go func() {
		err := http.ListenAndServe(":" + pkgC.GetString("APP_PORT"), middlewares.RemoveTrailingSlash(router))
		log.Fatalln(err)
	}()
	// :xxx 这里的端口占用采用延时，判断，还是不好， 最好能拿到服务启动后的后续hook才能更友好，更准确
	time.Sleep(time.Second * 1)
	fmt.Printf(`
		🚀 🚀 🚀 Server is running!
		Listening on port %s
		Explore at http://localhost:%s
		Explore graphql at http://localhost:%s/graphql

		`,
		pkgC.GetString("APP_PORT"), pkgC.GetString("APP_PORT"), pkgC.GetString("APP_PORT"))
	select{}
}

func RunMigrateSeed(c *cli.Context) error {
	bootstrap.SetupDB()
	// 强制清空数据表并重新生成
	if c.Bool("force") {
		var tables string
		for i, item := range bootstrap.MigrationTables {
			var flag  string
			if i + 1 != len(bootstrap.MigrationTables ) {
				flag = " , "
			} else {
				flag = " "
			}
			tables +=  item.TableName() + flag
		}
		sql := fmt.Sprintf("DROP TABLE %s", tables)
		model.DB.Exec(sql)
		for _, item := range bootstrap.MigrationTables {
			_ = model.DB.AutoMigrate(item)
		}
	}
	for _, seed := range seeders.All() {
		log.Println(seed.Name)
		err := seed.Run(model.DB)
		logger.LogError(err)
	}
	return nil
}
