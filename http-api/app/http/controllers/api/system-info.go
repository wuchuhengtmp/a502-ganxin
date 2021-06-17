/**
 * @Desc    The api is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @DATE    2021/4/27
 * @Listen  MIT
 */
package api

import (
	"http-api/pkg/response"
	"net/http"
)

type SystemInfo struct {
	AppName string `json:"appName"`
	AppIcon string `json:"appIcon"`
}

func (*SystemInfo) Show (w http.ResponseWriter, r *http.Request)  {
	sysInfo := SystemInfo{
	}
	response.SuccessResponse(sysInfo, w)
}