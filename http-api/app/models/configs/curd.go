/**
 * @Desc    The configs is part of http-api
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @DATE    2021/4/27
 * @Listen  MIT
 */
package configs

func GetAppName() string {
	return getVal(APP_NAME)
}

func GetAppIcon() string {
	return getVal(APP_ICON)
}

func GetMiniWechatAppId() string  {
	return getVal(MINI_WECHAT_APP_ID)
}

func GetMiniWechatAppSecret() string {
	return getVal(MINI_WECHAT_APP_SECRET)
}
