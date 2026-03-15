# 示例 01：基础 RPC 服务

这是使用 goctl 生成 RPC 服务的最简单示例。

## Proto 定义

一个 `greeter.proto` 文件，包含一个服务和一个 RPC 方法，无外部导入。

`go_package` 使用完整的模块路径：

```protobuf
option go_package = "example.com/demo/greeter";
```

## 生成命令

### 方式一：使用 `goctl rpc new` 快速创建

```bash
# 一条命令创建完整的 RPC 项目
goctl rpc new greeter
```

该命令会同时生成 proto 文件和服务代码：

```
greeter/
├── etc
│   └── greeter.yaml
├── greeter
│   ├── greeter.pb.go
│   └── greeter_grpc.pb.go
├── greeter.go
├── greeter.proto
├── greeterclient
│   └── greeter.go
└── internal
    ├── config
    │   └── config.go
    ├── logic
    │   └── pinglogic.go
    ├── server
    │   └── greeterserver.go
    └── svc
        └── servicecontext.go
```

### 方式二：基于已有 Proto 文件生成

首先，在输出目录中初始化 `go.mod`：

```bash
mkdir -p output && cd output && go mod init example.com/demo && cd ..
```

然后生成代码：

```bash
goctl rpc protoc greeter.proto \
  --go_out=output \
  --go-grpc_out=output \
  --zrpc_out=output \
  --go_opt=module=example.com/demo \
  --go-grpc_opt=module=example.com/demo \
  --module=example.com/demo \
  -I .
```

生成的目录结构：

```
output/
├── etc
│   └── greeter.yaml
├── go.mod
├── greeter
│   ├── greeter.pb.go
│   └── greeter_grpc.pb.go
├── greeter.go
├── greeterclient
│   └── greeter.go
└── internal
    ├── config
    │   └── config.go
    ├── logic
    │   └── sayhellologic.go
    ├── server
    │   └── greeterserver.go
    └── svc
        └── servicecontext.go
```

## 要点说明

- 这是最简单的场景：一个 proto 文件、一个服务、一个 RPC 方法。
- `go_package` 使用完整的模块路径（`example.com/demo/greeter`），而非相对路径。
- `--module` 告诉 goctl Go 模块名；`--go_opt=module=...` 和 `--go-grpc_opt=module=...` 告诉 protoc 从输出路径中去除模块前缀。
- `--zrpc_out` 指定 goctl 生成的服务代码输出目录。
- `--go_out` 和 `--go-grpc_out` 指定 protoc 生成代码的输出目录。
- 编辑逻辑文件（`internal/logic/sayhellologic.go`）来实现业务逻辑。
