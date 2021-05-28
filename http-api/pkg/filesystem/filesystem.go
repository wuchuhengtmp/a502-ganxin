/**
 * @Desc    文件操作
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/24
 * @Listen  MIT
 */
package filesystem

import (
	"encoding/json"
	"fmt"
	"http-api/pkg/config"
	"io/ioutil"
	"os"
	"path/filepath"
)

type FileInstance struct {
	Path   string // 文件的相对目录
	Disk   string  // 硬盘名
	Domain string // 访问的域名,如: https://wuchuheng.com
}

/**
 * 获取文件的url链接
 */
func (f FileInstance) Url(path string) string {
	s := fmt.Sprintf("%s/%s/%s", f.Domain, f.Path, path)
	return s
}
const AccessDir = "public"
/**
 * 文件写入
 */
func (f FileInstance)Put(path string, content []byte) error {
	rootDir, _:= os.Getwd()
	fileName := fmt.Sprintf("%s/%s/%s/%s", rootDir, AccessDir, f.Path, path)
	dir := filepath.Dir(fileName)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}
	err := ioutil.WriteFile(fileName, content, 0644)
	return err
}

/**
 * 默认硬盘文件直接写入
 */
func Put(path string, content []byte) error {
	rootDir, _:= os.Getwd()
	f := getDefaultConfig()
	fileName := fmt.Sprintf("%s/%s/%s/%s", rootDir, AccessDir, f.Path, path)
	dir := filepath.Dir(fileName)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}
	err := ioutil.WriteFile(fileName, content, 0644)
	return err
}

func Disk(disk string) (diskInstance FileInstance) {
	diskInstance.Disk = disk
	switch disk {
	// 本地配置
	case "local":
		localConfig := getLocalConfig()
		diskInstance.Path = localConfig.PrefixPath
		diskInstance.Domain = localConfig.Domain
	}
	return diskInstance
}

type LocalConfig struct {
	Domain string
	PrefixPath string
}
/**
 * 获取本地的硬盘配置
 */
func getLocalConfig() LocalConfig {
	fileConfig := config.Get("fileSystem")
	localConfStr, _ := json.Marshal(fileConfig.(map[string]interface{})["local"])
	localConf := LocalConfig{}
	_ :json.Unmarshal( localConfStr, &localConf)
	return localConf
}

/**
 * 获取默认硬盘配置
 */
func getDefaultConfig() (diskInstance FileInstance) {
	disk := GetDefaultDisk()
	diskInstance.Disk = disk
	switch disk {
		case "local":
			localConfig := getLocalConfig()
			diskInstance.Path = localConfig.PrefixPath
			diskInstance.Domain = localConfig.Domain
			break
			// ... match another some disk
	}
	return diskInstance
}

func GetDefaultDisk() string {
	fileConfig := config.Get("fileSystem")
	defaultDisk := fileConfig.(map[string]interface{})["default"]
	return defaultDisk.(string)
}

