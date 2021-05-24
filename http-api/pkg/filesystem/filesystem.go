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
)

type FileInstance struct {
	Path string
	Disk string
}

/**
 * 获取文件的url链接
 */
func (f FileInstance) Url(path string) string {
	fileConfig := config.Get("fileSystem")
	localConfStr, _ := json.Marshal(fileConfig.(map[string]interface{})["local"])
	localConf := struct {
		Domain string
		PrefixPath string
	}{}
	_ :json.Unmarshal( localConfStr, &localConf)

	return fmt.Sprintf("%s/%s/%s", localConf.Domain, localConf.PrefixPath, path)
}

func Disk(disk string) FileInstance {
	return FileInstance{ Disk: disk }
}
