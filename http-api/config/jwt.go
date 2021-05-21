package config

import "http-api/pkg/config"

func init()  {
	config.Add("jwt", config.StrMap{
		"secret": config.Env("JWT_SECRET", "c0df2ddcbccb60830fa126951e4dd6a9"),
		"expired": config.Env("JWT_EXPIRED", 60 * 60 * 24 * 7),
	})
}
