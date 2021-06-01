### 集成测试
* 确保数据库连接正常
* 确保数据库种子数据填充完毕，没有则连接数据库后运行:
``` bash 
$ go run main.go  seeds // 种子数据会自动写入到已连接成功的数据库中
```

最后进行接口的集成测试
``` bash 
$ go test http-api/tests -v -count=1
```