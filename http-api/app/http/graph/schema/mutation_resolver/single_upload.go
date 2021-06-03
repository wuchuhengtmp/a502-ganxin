/**
 * @Desc    单文件上传解析器
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/27
 * @Listen  MIT
 */
package mutation_resolver

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"http-api/app/http/graph/model"
	"http-api/app/models/files"
	"http-api/pkg/filesystem"
	"io/ioutil"
	"time"
)

/**
 * 单文件上传
 */
func (m *MutationResolver)SingleUpload(ctx context.Context, file graphql.Upload) (*model.FileItem, error) {
	year, month, day := time.Now().Date()
	path := fmt.Sprintf("%d-%d-%d/%d-%s", year, month, day, time.Now().Unix(), file.Filename)
	content, _ := ioutil.ReadAll(file.File)
	filesystem.Put(path, content)
	newFile := files.File{
		Path: path,
	}
	newFile.CreateFile()
	res := model.FileItem{}
	res.ID =  newFile.ID
	res.URL = newFile.GetUrl()

	return &res, nil
}
