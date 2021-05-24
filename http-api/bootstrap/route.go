package bootstrap

import (
	"github.com/gorilla/mux"
	"http-api/app/models/codeinfo"
	"http-api/app/models/specificationinfo"
	"http-api/app/models/configs"
	"http-api/app/models/repositories"
	"http-api/app/models/roles"
	"http-api/app/models/users"
	"http-api/pkg/config"
	"http-api/pkg/model"
	"http-api/pkg/route"
	"http-api/routes"
	"time"
)

func SetupRoute() *mux.Router  {
	router := mux.NewRouter()
	routes.RegisterWebRoutes(router)
	routes.RegisterApiRoutes(router)
	route.SetRoute(router)
	return router
}

func SetupDB() {
	db := model.ConnectDB()
	sqlDB, _ := db.DB()

	// 设置最大连接
	sqlDB.SetMaxOpenConns(config.GetInt("database.mysql.max_open_connections"))
	// 设置最大空闲连接
	sqlDB.SetMaxIdleConns(config.GetInt("database.mysql.max_idle_connections"))
	// 设置每个链接的过期时间
	sqlDB.SetConnMaxLifetime(time.Duration(config.GetInt("database.mysql.max_life_seconds")) * time.Second)
	// 迁移结构
	db.AutoMigrate(
		configs.Configs{},
		users.Users{},
		repositories.Repositories{},
		roles.Roles{},
		specificationinfo.SpecificationInfo{},
		codeinfo.CodeInfo{},
	)
}