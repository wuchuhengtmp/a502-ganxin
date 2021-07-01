/**
 * @Desc    The helper is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
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

/**
 * 时间转格式
 */
func Time2Str(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func FormatTime(t time.Time) string {
	year, month, day := t.Date()
	h := t.Hour()
	i := t.Minute()
	s := t.Second()
	return fmt.Sprintf("%d-%d-%d %d:%d:%d", year, month, day, h, i, s)
}

/**
 * 比较2个数据合集
 */
func CompareCollect(newCollect []int64, oldCollect []int64) (addItems []int64, deleteItem []int64) {
	newIdMapId := make(map[int64]int64)
	oldIdMapId := make(map[int64]int64)
	for _, i := range newCollect {
		newIdMapId[i] = i
	}
	for _, i := range oldCollect {
		oldIdMapId[i] = i
	}
	for _, i := range newCollect {
		if _, ok := oldIdMapId[i]; !ok {
			addItems = append(addItems, i)
		}
	}
	for _, i := range oldCollect {
		if _, ok := newIdMapId[i]; !ok {
			 deleteItem = append(deleteItem,  i)
		}
	}

	return
}
