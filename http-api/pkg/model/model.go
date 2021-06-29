package model

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"http-api/pkg/config"
	"http-api/pkg/logger"
)

var DB *gorm.DB

func ConnectDB() *gorm.DB {
	var err error
	dbConnection := config.GetString("database.connection")
	var gormConfig gorm.Dialector
	// sqlServer 数据驱动
	if dbConnection == "sqlserver" {
		gormConfig = getSqlServerDsn()
	} else {
		// 默认msql驱动
		gormConfig = getMysqlServerConfig()
	}
	var level gormlogger.LogLevel
	if config.GetBool("app.debug") {
		// 打印完成日志
		level = gormlogger.Warn
	} else {
		// 只有错误时才显示
		level = gormlogger.Error
	}
	DB, err = gorm.Open(gormConfig, &gorm.Config{
		Logger: gormlogger.Default.LogMode(level),
	})
	logger.LogError(err)

	return DB
}

/**
 * sqlServer 驱动配置
 */
func getSqlServerDsn() gorm.Dialector {
	var (
		host = config.GetString("database.sqlserver.host")
		port = config.GetString("database.sqlserver.port")
		database = config.GetString("database.sqlserver.database")
		username = config.GetString("database.sqlserver.username")
		password = config.GetString("database.sqlserver.password")
	)
	dsn := fmt.Sprintf(
		"sqlserver://%s:%s@%s:%s?database=%s",
		username, password, host, port, database,
	)

	return sqlserver.Open(dsn)
}

/**
 * mysql 配置
 */
func getMysqlServerConfig() gorm.Dialector {
	var (
		host = config.GetString("database.mysql.host")
		port = config.GetString("database.mysql.port")
		database = config.GetString("database.mysql.database")
		username = config.GetString("database.mysql.username")
		password = config.GetString("database.mysql.password")
		charset = config.GetString("database.mysql.charset")
	)
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%t&loc=%s",
		username,
		password,
		host,
		port,
		database,
		charset,
		true,
		"Local",
	)
	gormConfig := mysql.New( mysql.Config{
		DSN: dsn,
	})

	return gormConfig
}
