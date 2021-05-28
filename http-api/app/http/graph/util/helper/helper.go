/**
 * @Desc    The helper is part of http-api
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/28
 * @Listen  MIT
 */
package helper

import (
	"fmt"
	"time"
)

/**
 * 字串转时间类型
 */
func Str2Time(str string) (time.Time, error)  {
	return time.Parse("2006-01-02 15:04:05", str)
}

func FormatTime(t time.Time) string {
	year, month, day := t.Date()
	h := t.Hour()
	i := t.Minute()
	s := t.Second()
	return fmt.Sprintf("%d-%d-%d %d:%d:%d", year, month, day, h, i, s)
}
