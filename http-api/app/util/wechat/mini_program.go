/**
 * @Desc    The wechat is part of http-api
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @DATE    2021/4/28
 * @Listen  MIT
 */
package wechat

import (
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/miniprogram"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
	"http-api/app/models/configs"
	"sync"
)

var instance *miniprogram.MiniProgram
var once sync.Once

/**
 * 获取小程序接口实例
 */
func GetWCInstance() *miniprogram.MiniProgram {
	once.Do(func() {
		memory := cache.NewMemory()
		wc := wechat.NewWechat()
		cfg := &miniConfig.Config{
			AppID:     configs.GetMiniWechatAppId(),
			AppSecret: configs.GetMiniWechatAppSecret(),
			Cache:     memory,
		}
		instance = wc.GetMiniProgram(cfg)
	})
	return instance
}
