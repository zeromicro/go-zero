# Rpc Generation

Goctl Rpc是`goctl`脚手架下的一个rpc服务代码生成模块，支持proto模板生成和rpc服务代码生成，通过此工具生成代码你只需要关注业务逻辑编写而不用去编写一些重复性的代码。这使得我们把精力重心放在业务上，从而加快了开发效率且降低了代码出错率。

## 特性

* 简单易用
* 快速提升开发效率
* 出错率低
* 贴近 protoc


## 快速开始

### 方式一：快速生成greet服务

  通过命令 `goctl rpc new ${servieName}`生成

  如生成greet rpc服务：

  ```Bash
  goctl rpc new greet
  ```

  执行后代码结构如下:

```text
.
└── greet
    ├── etc
    │   └── greet.yaml
    ├── greet
    │   ├── greet.go
    │   ├── greet.pb.go
    │   └── greet_grpc.pb.go
    ├── greet.go
    ├── greet.proto
    └── internal
        ├── config
        │   └── config.go
        ├── logic
        │   └── pinglogic.go
        ├── server
        │   └── greetserver.go
        └── svc
            └── servicecontext.go
```

### 方式二：通过指定proto生成rpc服务

* 生成proto模板

```Bash
$ goctl rpc template -o=user.proto
```
  
```proto
syntax = "proto3";

package user;
option go_package="./user";

message Request {
  string ping = 1;
}

message Response {
  string pong = 1;
}

service User {
  rpc Ping(Request) returns(Response);
}
```
  

* 生成rpc服务代码

```bash
$ goctl rpc protoc  user.proto --go_out=. --go-grpc_out=. --zrpc_out=.
```


## 用法

### rpc 服务生成用法

```Bash
$ goctl rpc protoc -h
Generate grpc code

Usage:
  goctl rpc protoc [flags]

Examples:
goctl rpc protoc xx.proto --go_out=./pb --go-grpc_out=./pb --zrpc_out=.

Flags:
      --branch string     The branch of the remote repo, it does work with --remote
  -h, --help              help for protoc
      --home string       The goctl home path of the template, --home and --remote cannot be set at the same time, if they are, --remote has higher priority
  -m, --multiple          Generated in multiple rpc service mode
      --remote string     The remote git repo of the template, --home and --remote cannot be set at the same time, if they are, --remote has higher priority
                          	The git repo directory must be consistent with the https://github.com/zeromicro/go-zero-template directory structure
      --style string      The file naming format, see [https://github.com/zeromicro/go-zero/tree/master/tools/goctl/config/readme.md] (default "gozero")
  -v, --verbose           Enable log output
      --zrpc_out string   The zrpc output directory
```

### 参数说明

* --branch 指定远程仓库模板分支
* --home 指定goctl模板根目录
* -m, --multiple 指定生成多个rpc服务模式, 默认为 false, 如果为  false, 则只支持生成一个rpc service, 如果为 true, 则支持生成多个 rpc service，且多个 rpc service 会分组。
* --style 指定文件输出格式
* -v, --verbose 显示日志
* --zrpc_out 指定zrpc输出目录

> ## --multiple
> 是否开启多个 rpc service 生成，如果开启，则满足一下新特性
> 1. 支持 1 到多个 rpc service 
> 2. 生成 rpc 服务会按照服务名称分组（尽管只有一个 rpc service）
> 3. rpc client 的文件目录变更为固定名称 `client`
> 
> 如果不开启，则和旧版本 rpc 生成逻辑一样（兼容）
> 1. 有且只能有一个 rpc service


## rpc 服务生成 example
详情见 [example/rpc](https://github.com/zeromicro/go-zero/tree/master/tools/goctl/example)

## --multiple 为 true 和 false 的目录区别
源 proto 文件

```protobuf
syntax = "proto3";

package hello;

option go_package = "./hello";

message HelloReq {
  string in = 1;
}

message HelloResp {
  string msg = 1;
}

service Greet {
  rpc SayHello(HelloReq) returns (HelloResp);
}
```

### --multiple=true

```text
hello
├── client // 区别1：rpc client 目录固定为 client 名称
│   └── greet // 区别2：会按照 rpc service 名称分组
│       └── greet.go
├── etc
│   └── hello.yaml
├── hello.go
├── internal
│   ├── config
│   │   └── config.go
│   ├── logic
│   │   └── greet // 区别2：会按照 rpc service 名称分组
│   │       └── sayhellologic.go
│   ├── server
│   │   └── greet // 区别2：会按照 rpc service 名称分组
│   │       └── greetserver.go
│   └── svc
│       └── servicecontext.go
└── pb
    └── hello
        ├── hello.pb.go
        └── hello_grpc.pb.go
```

### --multiple=false (旧版本目录，向后兼容)
```text
hello
├── etc
│   └── hello.yaml
├── greet
│   └── greet.go
├── hello.go
├── internal
│   ├── config
│   │   └── config.go
│   ├── logic
│   │   └── sayhellologic.go
│   ├── server
│   │   └── greetserver.go
│   └── svc
│       └── servicecontext.go
└── pb
    └── hello
        ├── hello.pb.go
        └── hello_grpc.pb.go
```