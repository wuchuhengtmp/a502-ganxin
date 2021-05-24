package config

import "http-api/pkg/config"

func init()  {
	config.Add("database", config.StrMap {
		"connection": config.Env("DB_CONNECTION"),
		/**
		 * mysql 配置
		 */
		"mysql": map[string]interface{} {
			"host": config.Env("DB_HOST", "127.0.0.1"),
			"port": config.Env("DB_PORT", "3306"),
			"database": config.Env("DB_DATABASE", ""),
			"username": config.Env("DB_USERNAME", ""),
			"password": config.Env("DB_PASSWORD", ""),
			"charset": "utf8mb4",

			"max_idle_connections": config.Env("DB_MAX_IDLE_CONNECTIONS", 100),
			"max_open_connections": config.Env("DB_MAX_OPEN_CONNECTIONS", 25),
			"max_life_seconds": config.Env("DB_MAX_LIFE_SECONDS", 60 * 5),
		},
		/**
		 * sqlServer 配置
		 */
		"sqlserver": map[string]interface{} {
			"host": config.Env("DB_HOST", "127.0.0.1"),
			"port": config.Env("DB_PORT", "1433"),
			"database": config.Env("DB_DATABASE", ""),
			"username": config.Env("DB_USERNAME", ""),
			"password": config.Env("DB_PASSWORD", ""),
		},
	})
}