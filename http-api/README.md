### 集成测试
* 确保数据库连接正常
* 确保数据库种子数据填充完毕，没有则连接数据库后运行:
``` bash 
$ go run main.go  seeds // 种子数据会自动写入到已连接成功的数据库中
$ go run main.go http_api // 启动服务
```

最后进行接口的集成测试
``` bash 
$ go test http-api/tests -v -count=1
```

## 开发迭代注意事项
### `Graphql`接口的程序分层原则
`grapqh request` --> `controller`层 --> `request validator`层 --> `service`层 --> `model`层；
原则:
    1：依赖只能向下依赖，就是上层可以调用凭单的下层，但下层绝不能调用上层
    2: 同层之间也不要相互调用， 如果你不想来点烦人的循环依赖的话