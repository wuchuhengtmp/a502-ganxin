package database

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"http-api/pkg/logger"
	"time"
)

var DB *sql.DB

func InitDB()  {
	config := mysql.Config{
		User: "wuchuheng_tmp",
		Passwd: "wuchuheng_tmp",
		Addr: "192.168.0.42",
		Net: "tcp",
		DBName: "wuchuheng_tmp",
		AllowNativePasswords: true,
	}
	// 准备数据库连接池
	var err error
	DB, err = sql.Open("mysql", config.FormatDSN())
	logger.LogError(err)

	// 设置最大连接数
	DB.SetMaxOpenConns(25)
	// 设置最大空闲连接数
	DB.SetMaxIdleConns(25)
	// 设置每个链接的过期时间
	DB.SetConnMaxLifetime(5 * time.Minute)

	// 尝试连接，失败会报错
	err = DB.Ping()
	logger.LogError(err)
}
