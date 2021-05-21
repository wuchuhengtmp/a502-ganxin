package model

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"http-api/pkg/config"
	"http-api/pkg/logger"
)

var DB *gorm.DB

func ConnectDB() *gorm.DB {
	var err error
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
		username, password, host, port, database, charset, true, "Local",
	)
	gormConfig := mysql.New( mysql.Config{
		DSN: dsn,
	})
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
