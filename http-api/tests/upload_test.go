/**
 * @Desc    文件上传集成测试
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/2
 * @Listen  MIT
 */
package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"testing"
)

/**
 * 单文件上传集成测试
 */
func TestSingleUpload(t *testing.T)  {
	values := map[string]io.Reader{
		"0":  mustOpen("upload_test_img.png"), // lets assume its this file
		"operations": strings.NewReader(`{"query":"mutation singleUpload ($file: Upload!) {\n  singleUpload(file: $file) {\n    id\n    url\n  }\n}","variables":{"file":null},"operationName":"singleUpload"}`),
		"map": strings.NewReader(`{"0":["variables.file"]}`),
	}
	client := &http.Client{}
	res, err := Upload(client, bashUrl, values)
	if err != nil {
		panic(err)
	}
	if res.StatusCode != 200 {
		hasError(t, fmt.Errorf("单文件上传测试失败"))
	}
	body, _ := ioutil.ReadAll(res.Body)
	// 期望响应的数据格式
	var exprectedResData struct{
		Data struct{
			SingleUpload interface{} `json:"singleUpload"`
		} `json:"data"`
	}
	_ = json.Unmarshal(body, &exprectedResData)
	if exprectedResData.Data.SingleUpload == nil {
		hasError(t, fmt.Errorf("单文件上传测试失败,非期望的响应数据格式"))
	}

}

func Upload(client *http.Client, url string, values map[string]io.Reader) (response *http.Response, err error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add an image file
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				return
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return
		}

	}
	w.Close()
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return
	}
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
	}

	return res, nil
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}