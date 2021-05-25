/**
 * @Desc    文件系统配置
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/24
 * @Listen  MIT
 */
package config

import "http-api/pkg/config"

func init()  {
	config.Add("fileSystem", config.StrMap {
		// 本地存储配置
		"local": struct {
			Domain interface{}// 访问域名
			PrefixPath string // 目录前缀
		}{
			Domain: config.Env("FILE_SYSTEM_DOMAIN", "http://127.0.0.1"),
			PrefixPath: "uploads/local", // 相当于 public/local 目录
		},
	})
}

