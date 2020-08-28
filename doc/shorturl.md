# 使用go-zero从0到1快速构建高并发的短链服务

## 0. 什么是短链服务？

短链服务就是将长的URL网址，通过程序计算等方式，转换为简短的网址字符串。

写此短链服务是为了从整体上演示go-zero构建完整微服务的过程，算法和实现细节尽可能简化了，所以这不是一个高阶的短链服务。

## 1. 短链微服务架构图

<img src="images/shorturl-arch.png" alt="架构图" width="800" />

## 2. 准备工作

* 准备goctl工具，在任意目录下进行，目的是为了编译goctl工具
  1. `git clone https://github.com/tal-tech/go-zero`
  2. 在`tools/goctl`目录下编译goctl工具`go build goctl.go`
  3. 将生成的goctl放到`$PATH`下，确保goctl命令可运行
* 创建工作目录`shorturl`
* 在`shorturl`目录下执行`go mod init shorturl`初始化`go.mod`

## 3. 编写API Gateway代码

* 通过goctl生成`shorturl.api`并编辑，为了简洁，去除了文件开头的`info`，代码如下：

  ```go
  type (
  	shortenReq struct {
  		url string `form:"url"`
  	}
  
  	shortenResp struct {
  		shortUrl string `json:"shortUrl"`
  	}
  )
  
  type (
  	expandReq struct {
  		key string `form:"key"`
  	}
  
  	expandResp struct {
  		url string `json:"url"`
  	}
  )
  
  service shorturl-api {
  	@server(
  		handler: ShortenHandler
  	)
  	get /shorten(shortenReq) returns(shortenResp)
  
  	@server(
  		handler: ExpandHandler
  	)
  	get /expand(expandReq) returns(expandResp)
  }
  ```

  type用法和go一致，service用来定义get/post/head/delete等api请求，解释如下：

  * `service shorturl-api {`这一行定义了service名字
  * `@server`部分用来定义server端用到的属性
  * `handler`定义了服务端handler名字
  * `get /shorten(shortenReq) returns(shortenResp)`定义了get方法的路由、请求参数、返回参数等

* 使用goctl生成API Gateway代码

  ```shell
  goctl api go -api shorturl.api -dir api
  ```

  生成的文件结构如下：

  ```
  .
  ├── api
  │   ├── etc
  │   │   └── shorturl-api.yaml         // 配置文件
  │   ├── internal
  │   │   ├── config
  │   │   │   └── config.go             // 定义配置
  │   │   ├── handler
  │   │   │   ├── expandhandler.go      // 实现expandHandler
  │   │   │   ├── routes.go             // 定义路由处理
  │   │   │   └── shortenhandler.go     // 实现shortenHandler
  │   │   ├── logic
  │   │   │   ├── expandlogic.go        // 实现ExpandLogic
  │   │   │   └── shortenlogic.go       // 实现ShortenLogic
  │   │   ├── svc
  │   │   │   └── servicecontext.go     // 定义ServiceContext
  │   │   └── types
  │   │       └── types.go              // 定义请求、返回结构体
  │   └── shorturl.go                   // main入口定义
  ├── go.mod
  ├── go.sum
  └── shorturl.api
  ```

* 启动API Gateway服务，默认侦听在8888端口

  ```shell
  go run api/shorturl.go -f api/etc/shorturl-api.yaml
  ```

* 测试API Gateway服务

  ```shell
  curl -i "http://localhost:8888/shorten?url=http://www.xiaoheiban.cn"
  ```

  返回如下：

  ```http
  HTTP/1.1 200 OK
  Content-Type: application/json
  Date: Thu, 27 Aug 2020 14:31:39 GMT
  Content-Length: 15
  
  {"shortUrl":""}
  ```

  可以看到我们API Gateway其实啥也没干，就返回了个空值，接下来我们会在rpc服务里实现业务逻辑

* 可以修改`internal/svc/servicecontext.go`来传递服务依赖（如果需要）

* 实现逻辑可以修改`internal/logic`下的对应文件

* 可以通过`goctl`生成各种客户端语言的api调用代码

## 4. 编写shorten rpc服务

* 在`rpc/shorten`目录下编写`shorten.proto`文件

  可以通过命令生成proto文件模板

  ```shell
  goctl rpc template -o shorten.proto
  ```

  修改后文件内容如下：

  ```protobuf
  syntax = "proto3";
  
  package shorten;
  
  message shortenReq {
      string url = 1;
  }
  
  message shortenResp {
      string key = 1;
  }
  
  service shortener {
      rpc shorten(shortenReq) returns(shortenResp);
  }
  ```

* 用`goctl`生成rpc代码，在`rpc/shorten`目录下执行命令

  ```shell
  goctl rpc proto -src shorten.proto
  ```

  文件结构如下：

  ```
  rpc/shorten
  ├── etc
│   └── shorten.yaml               // 配置文件
  ├── internal
│   ├── config
  │   │   └── config.go              // 配置定义
  │   ├── handler
  │   │   └── shortenerhandler.go    // api handler, 不需要修改
  │   ├── logic
│   │   └── shortenlogic.go        // api业务逻辑在这里实现
  │   └── svc
  │       └── servicecontext.go      // 定义ServiceContext，传递依赖
  ├── pb
  │   └── shorten.pb.go
  ├── shared
  │   ├── shortenermodel.go          // 提供了外部调用方法，无需修改
  │   ├── shortenermodel_mock.go     // mock方法，测试用
  │   └── types.go                   // request/response结构体定义
  ├── shorten.go                     // rpc服务main函数
  └── shorten.proto
  ```
  
  直接可以运行，如下：
  
  ```shell
  $ go run shorten.go -f etc/shorten.yaml
  Starting rpc server at 127.0.0.1:8080...
  ```
  
  `etc/shorten.yaml`文件里可以修改侦听端口等配置

## 5. 编写expand rpc服务

* 在`rpc/expand`目录下编写`expand.proto`文件

  可以通过命令生成proto文件模板

  ```shell
  goctl rpc template -o expand.proto
  ```

  修改后文件内容如下：

  ```protobuf
  syntax = "proto3";
  
  package expand;
  
  message expandReq {
      string key = 1;
  }
  
  message expandResp {
      string url = 1;
  }
  
  service expander {
      rpc expand(expandReq) returns(expandResp);
  }
  ```

* 用`goctl`生成rpc代码，在`rpc/expand`目录下执行命令

  ```shell
  goctl rpc proto -src expand.proto
  ```

  文件结构如下：

  ```
  rpc/expand
  ├── etc
│   └── expand.yaml                // 配置文件
  ├── internal
│   ├── config
  │   │   └── config.go              // 配置定义
  │   ├── handler
  │   │   └── expanderhandler.go     // api handler, 不需要修改
  │   ├── logic
│   │   └── expandlogic.go         // api业务逻辑在这里实现
  │   └── svc
  │       └── servicecontext.go      // 定义ServiceContext，传递依赖
  ├── pb
  │   └── expand.pb.go
  ├── shared
  │   ├── expandermodel.go           // 提供了外部调用方法，无需修改
  │   ├── expandermodel_mock.go      // mock方法，测试用
  │   └── types.go                   // request/response结构体定义
  ├── expand.go                      // rpc服务main函数
  └── expand.proto
  ```
  
  修改`etc/expand.yaml`里面的`ListenOn`的端口为`8081`，因为`8080`已经被`shorten`服务占用了
  
  修改后运行，如下：
  
  ```shell
  $ go run expand.go -f etc/expand.yaml
  Starting rpc server at 127.0.0.1:8081...
  ```
  
  `etc/expand.yaml`文件里可以修改侦听端口等配置

## 6. 修改API Gateway代码调用shorten/expand rpc服务（未完）

## 7. 定义数据库表结构，并生成CRUD+cache代码

* shorturl下创建rpc/model目录：`mkdir -p rpc/model`
* 在roc/model目录下编写创建shorturl表的sql文件`shorturl.sql`，如下：

  ```sql
  CREATE TABLE `shorturl`
  (
    `id` bigint(10) NOT NULL AUTO_INCREMENT,
    `key` varchar(255) NOT NULL DEFAULT '' COMMENT 'shorten key',
    `url` varchar(255) DEFAULT '' COMMENT 'original url',
    PRIMARY KEY(`id`),
    UNIQUE KEY `key_index`(`key`)
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
  ```

* 创建DB和table

  ```sql
  create database gozero;
  ```

  ```sql
  source shorturl.sql;
  ```

* 在`rpc/model`目录下执行如下命令生成CRUD+cache代码，`-c`表示使用`redis cache`

  ```shell
  goctl model mysql ddl -c -src shorturl.sql -dir .
  ```
  
  也可以用`datasource`命令代替`ddl`来指定数据库链接直接从schema生成
  
  生成后的文件结构如下：
  
  ```
  rpc/model
  ├── shorturl.sql
  ├── shorturlmodel.go              // CRUD+cache代码
  └── vars.go                       // 定义常量和变量
  ```

## 8. 修改shorten/expand rpc代码调用crud+cache代码

## 9. 完整调用演示

## 10. Benchmark（未完）

## 11. 总结（未完）

可以看到go-zero不只是一个框架，更是一个建立在框架+工具基础上的，简化和规范了整个微服务构建的技术体系。

我们一直强调**工具大于约定和文档**。

另外，我们在保持简单的同时也尽可能把微服务治理的复杂度封装到了框架内部，极大的降低了开发人员的心智负担，使得业务开发得以快速推进。