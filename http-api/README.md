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
#### 原则:  
* 1：依赖只能向下依赖，就是上层可以调用凭单的下层，但下层绝不能调用上层
2: 同层之间也不要相互调用， 如果你不想来点烦人的循环依赖的话
    
    
#### 程序设计弊端:  
* &emsp;模型层的各个文件不在同一个包命下，导致模型中读写业务过于复杂的话，容易导致循环依赖的情况，这才导致不得不在模型层之上再抽象一层服务层来避免这个情况。本来
是没这个必要一，如果模型文件都在一个包名下的话,处理起来会更加简单，但相对的，业务集中在一个文件可能导致文件过大，而又不好拆分。

#### 数据表关系图
https://viewer.diagrams.net/?highlight=0000ff&edit=_blank&layers=1&nav=1&title=a0502-%E5%9E%8B%E9%92%A2.drawio#Uhttps%3A%2F%2Fraw.githubusercontent.com%2Fwuchuheng%2Fgateway_demo%2Fmaster%2Fa0502-%25E5%259E%258B%25E9%2592%25A2.drawio

#### 后台初始化视频的文件
视频文件放在`2021-7-13/1626155017-tutor.mp4`的位置，更换需求的话更换它就可以。