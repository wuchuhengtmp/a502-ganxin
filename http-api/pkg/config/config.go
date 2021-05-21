package config

import (
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"http-api/pkg/logger"
)

var Viper *viper.Viper

type StrMap map[string]interface{}

func init() {
	Viper = viper.New()
	// 2. 设置文件名称
	Viper.SetConfigName(".env")
	// 3. 配置类型，支持 "json", "toml", "yaml", "yml", "properties",
	Viper.SetConfigType("env")
	// 4. 环境变量配置文件查找的路径，相对于 main.go
	Viper.AddConfigPath(".")
	// 5. 开始读根目录下的 .env 文件，读不到会报错
	err := Viper.ReadInConfig()
	logger.LogError(err)
	// 6. 设置环境变量前缀，用以区分 Go 的系统环境变量
	Viper.SetEnvPrefix("appenv")
	// 7. Viper.Get() 时，优先读取环境变量
	Viper.AutomaticEnv()
}

// 读取环境变量,支持默认值
func Env(envName string, defaultValue...interface{}) interface{} {
	if len(defaultValue)  > 0{
		return Get(envName, defaultValue[0])
	}
	return Get(envName)
}

// 获取配置， 允许使用方法获取如 app.name
func Get(path string, defaultValue ...interface{}) interface{} {
	if !Viper.IsSet(path) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}
	return Viper.Get(path)
}

// GetString 获取 string 类型配置信息
func GetString(path string, defaultValue ...interface{}) string {
	 return cast.ToString(Get(path, defaultValue...))
}

// GetInt64 获取 Int64 类型配置信息
func GetInt64(path string, defaultValue ...interface{}) int64 {
	return cast.ToInt64(Get(path, defaultValue...))
}

// GetInt 获取 Int64 类型配置信息
func GetInt(path string, defaultValue ...interface{}) int {
	return cast.ToInt(Get(path, defaultValue...))
}

// GetUint 获取 Uint 类型配置信息
func GetUint(path string, defaultValue ...interface{}) uint {
	return cast.ToUint(Get(path, defaultValue...))
}

// GetUint 获取 Uint 类型配置信息
func GetBool(path string, defaultValue ...interface{}) bool {
	return cast.ToBool(Get(path, defaultValue...))
}

// 添加配置项
func Add(name string, configuration map[string]interface{})  {
	Viper.Set(name, configuration)
}
