# go-zero

[![Go](https://github.com/tal-tech/go-zero/workflows/Go/badge.svg?branch=master)](https://github.com/tal-tech/go-zero/actions)
[![codecov](https://codecov.io/gh/tal-tech/go-zero/branch/master/graph/badge.svg)](https://codecov.io/gh/tal-tech/go-zero)
[![Go Report Card](https://goreportcard.com/badge/github.com/tal-tech/go-zero)](https://goreportcard.com/report/github.com/tal-tech/go-zero)
[![Release](https://img.shields.io/github/v/release/tal-tech/go-zero.svg?style=flat-square)](https://github.com/tal-tech/go-zero)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## 1. go-zero框架背景

18年初，晓黑板后端在经过频繁的宕机后，决定从`Java+MongoDB`的单体架构迁移到微服务架构，经过仔细思考和对比，我们决定：

* 基于Go语言
  * 高效的性能
  * 简洁的语法
  * 广泛验证的工程效率
  * 极致的部署体验
  * 极低的服务端资源成本
* 自研微服务框架
  * 个人有过很多微服务框架自研经验
  * 需要有更快速的问题定位能力
  * 更便捷的增加新特性

## 2. go-zero框架设计思考

对于微服务框架的设计，我们期望保障微服务稳定性的同时，也要特别注重研发效率。所以设计之初，我们就有如下一些准则：

* 保持简单
* 高可用
* 高并发
* 易扩展
* 弹性设计，面向故障编程
* 尽可能对业务开发友好，封装复杂度
* 尽可能约束做一件事只有一种方式

我们经历不到半年时间，彻底完成了从`Java+MongoDB`到`Golang+MySQL`为主的微服务体系迁移，并于18年8月底完全上线，稳定保障了晓黑板后续增长，确保了整个服务的高可用。

## 3. go-zero项目实现和特点

go-zero是一个集成了各种工程实践的包含web和rpc框架，有如下主要特点：

* 强大的工具支持，尽可能少的代码编写
* 极简的接口
* 完全兼容net/http
* 支持中间件，方便扩展
* 高性能
* 面向故障编程，弹性设计
* 内建服务发现、负载均衡
* 内建限流、熔断、降载，且自动触发，自动恢复
* API参数自动校验
* 超时级联控制
* 自动缓存控制
* 链路跟踪、统计报警等
* 高并发支撑，稳定保障了晓黑板疫情期间每天的流量洪峰

如下图，我们从多个层面保障了整体服务的高可用：

![弹性设计](doc/images/resilience.jpg)

## 4. go-zero框架收益

* 保障大并发服务端的稳定性，经受了充分的实战检验
* 极简的API定义
* 一键生成Go, iOS, Android, Dart, TypeScript, JavaScript代码，并可直接运行
* 服务端自动校验参数合法性

## 5. go-zero近期开发计划

* 自动生成API mock server，便于客户端开发
* 自动生成服务端功能测试

## 6. Installation

1. 在项目目录下通过如下命令安装：

   ```shell
   go get -u github.com/tal-tech/go-zero
   ```

2. 代码里导入go-zero

   ```go
   import "github.com/tal-tech/go-zero"
   ```

## 7. Quick Start

1. 编译goctl工具

   ```shell
   go build tools/goctl/goctl.go
   ```

   把goctl放到$PATH的目录下，确保goctl可执行

2. 定义API文件，比如greet.api，可以在vs code里安装`goctl`插件，支持api语法

   ```go
   type Request struct {
     Name string `path:"name"`
   }
   
   type Response struct {
     Message string `json:"message"`
   }
   
   service greet-api {
     @server(
       handler: GreetHandler
     )
     get /greet/from/:name(Request) returns (Response);
   }
   ```

   也可以通过goctl生成api模本文件，命令如下：

   ```shell
   goctl api -o greet.api
   ```

3. 生成go服务端代码

   ```shell
   goctl api go -api greet.api -dir greet
   ```

   生成的文件结构如下：

   ```
   ├── greet
   │   ├── etc
   │   │   └── greet-api.json        // 配置文件
   │   ├── greet.go                  // main文件
   │   └── internal
   │       ├── config
   │       │   └── config.go         // 配置定义
   │       ├── handler
   │       │   ├── greethandler.go   // get/put/post/delete等路由定义文件
   │       │   └── routes.go         // 路由列表
   │       ├── logic
   │       │   └── greetlogic.go     // 请求逻辑处理文件
   │       ├── svc
   │       │   └── servicecontext.go // 请求上下文，可以传入mysql, redis等依赖
   │       └── types
   │           └── types.go          // 请求、返回等类型定义
   └── greet.api                     // api描述文件
   ```
   生成的代码可以直接运行：
   
```shell
   cd greet
   go run greet.go -f etc/greet-api.json
```

默认侦听在8888端口（可以在配置文件里修改），可以通过curl请求：

```shell
   ➜  go-zero git:(master) curl -w "\ncode: %{http_code}\n" http://localhost:8888/greet/from/kevin
   {"code":0}
   code: 200
```

编写业务代码：

* 可以在servicecontext.go里面传递依赖给logic，比如mysql, redis等
   * 在api定义的get/post/put/delete等请求对应的logic里增加业务处理逻辑
   
4. 可以根据api文件生成前端需要的Java, TypeScript, Dart, JavaScript代码

   ```shell
   goctl api java -api greet.api -dir greet
   goctl api dart -api greet.api -dir greet
   ...
   ```

## 8. Benchmark

![benchmark](doc/images/benchmark.png)

[测试代码见这里](https://github.com/smallnest/go-web-framework-benchmark)

## 9. 文档

* [goctl使用帮助](doc/goctl.md)
* [关键字替换和敏感词过滤工具](doc/keywords.md)

## 10. 微信交流群

添加我的微信：kevwan，请注明go-zero，我拉进go-zero社区群🤝
