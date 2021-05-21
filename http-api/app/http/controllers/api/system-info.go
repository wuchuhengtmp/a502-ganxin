/**
 * @Desc    The api is part of http-api
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @DATE    2021/4/27
 * @Listen  MIT
 */
package api

import (
	"http-api/app/models/configs"
	"http-api/pkg/response"
	"net/http"
)

type SystemInfo struct {
	AppName string `json:"appName"`
	AppIcon string `json:"appIcon"`
}

func (*SystemInfo) Show (w http.ResponseWriter, r *http.Request)  {
	sysInfo := SystemInfo{
		AppName: configs.GetAppName(),
		AppIcon: configs.GetAppIcon(),
	}
	response.SuccessResponse(sysInfo, w)
}