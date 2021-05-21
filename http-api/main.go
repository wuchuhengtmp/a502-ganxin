package main

import (
	_ "github.com/go-sql-driver/mysql"
	"http-api/cmd"
)

func main() {
	cmd.Run()
}
