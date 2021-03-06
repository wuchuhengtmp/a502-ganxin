package bootstrap

import (
	"github.com/gorilla/mux"
	"http-api/app/models/codeinfo"
	"http-api/app/models/companies"
	"http-api/app/models/configs"
	"http-api/app/models/devices"
	"http-api/app/models/files"
	"http-api/app/models/logs"
	"http-api/app/models/maintenance"
	"http-api/app/models/maintenance_leader"
	"http-api/app/models/maintenance_record"
	"http-api/app/models/msg"
	"http-api/app/models/order_express"
	order_details "http-api/app/models/order_specification"
	"http-api/app/models/order_specification_steel"
	"http-api/app/models/orders"
	"http-api/app/models/project_leader"
	"http-api/app/models/projects"
	"http-api/app/models/repositories"
	"http-api/app/models/repository_leader"
	"http-api/app/models/roles"
	"http-api/app/models/specificationinfo"
	"http-api/app/models/steel_logs"
	"http-api/app/models/steels"
	"http-api/app/models/users"
	"http-api/pkg/config"
	"http-api/pkg/model"
	"http-api/pkg/route"
	"http-api/routes"
	"time"
)

func SetupRoute() *mux.Router {
	router := mux.NewRouter()
	routes.RegisterWebRoutes(router)
	routes.RegisterApiRoutes(router)
	routes.RegisterGraphRoutes(router)
	route.SetRoute(router)
	return router
}
type MigrationTablesInterface interface {
	TableName() string
}
var MigrationTables = []MigrationTablesInterface {
	configs.Configs{},
	users.Users{},
	repositories.Repositories{},
	roles.Role{},
	specificationinfo.SpecificationInfo{},
	codeinfo.CodeInfo{},
	files.File{},
	devices.Device{},
	companies.Companies{},
	steels.Steels{},
	projects.Projects{},
	order_details.OrderSpecification{},
	orders.Order{},
	logs.Logos{},
	msg.Msg{},
	maintenance_leader.MaintenanceLeader{},
	maintenance_record.MaintenanceRecord{},
	maintenance.Maintenance{},
	repository_leader.RepositoryLeader{},
	project_leader.ProjectLeader{},
	steel_logs.SteelLog{},
	order_specification_steel.OrderSpecificationSteel{},
	order_express.OrderExpress{},
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
	for _, item := range MigrationTables {
		db.AutoMigrate(item)
	}
}
