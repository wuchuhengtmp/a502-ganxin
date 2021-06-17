/**
 * @Desc    The main is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/21
 * @Listen  MIT
 */
package main

import (
	_ "github.com/go-sql-driver/mysql"
	"http-api/cmd"
)

func main() {
	cmd.Run()
}
