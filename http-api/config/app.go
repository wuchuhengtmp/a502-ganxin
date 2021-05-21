package config

import (
	"http-api/pkg/config"
)

func init()  {
	config.Add("app", config.StrMap{
		"name": config.Env("APP_NAME", "GoBlog"),
		"env": config.Env("APP_ENV", "production"),
		"debug": config.Env("APP_DEBUG", false),
		"port": config.Env("APP_PORT", "3000"),
		"key": config.Env("APP_KEY", "w2ir2lliu3q3w3s3fasdfkjfask12kj1jk1311231k"),
	})
}
