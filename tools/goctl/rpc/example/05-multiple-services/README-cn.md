# 示例 05：多服务模式（`--multiple`）

本示例演示从一个 proto 文件生成多个 RPC 服务。

## Proto 定义

两个 proto 文件共享相同的 `go_package`：

```protobuf
option go_package = "example.com/demo/pb";
```

- `shared.proto` — 定义共享消息类型（`Meta`）。
- `multi.proto` — 定义了**两个**服务：`SearchService` 和 `NotifyService`。

当 proto 文件包含多个 `service` 块时，必须使用 `-m`（或 `--multiple`）标志。

## 生成命令

首先，在输出目录中初始化 `go.mod`：

```bash
mkdir -p output && cd output && go mod init example.com/demo && cd ..
```

然后使用 `-m` 标志生成代码：

```bash
goctl rpc protoc multi.proto \
  --go_out=output \
  --go-grpc_out=output \
  --zrpc_out=output \
  --go_opt=module=example.com/demo \
  --go-grpc_opt=module=example.com/demo \
  --module=example.com/demo \
  -I . \
  -m
```

生成的目录结构：

```
output/
├── client
│   ├── notifyservice
│   │   └── notifyservice.go
│   └── searchservice
│       └── searchservice.go
├── etc
│   └── multisvc.yaml
├── go.mod
├── internal
│   ├── config
│   │   └── config.go
│   ├── logic
│   │   ├── notifyservice
│   │   │   └── notifylogic.go
│   │   └── searchservice
│   │       └── searchlogic.go
│   ├── server
│   │   ├── notifyservice
│   │   │   └── notifyserviceserver.go
│   │   └── searchservice
│   │       └── searchserviceserver.go
│   └── svc
│       └── servicecontext.go
├── multisvc.go
└── pb
    ├── multi.pb.go
    ├── multi_grpc.pb.go
    └── shared.pb.go
```

## 要点说明

- `-m`（或 `--multiple`）标志启用多服务模式。
- 多服务模式下，`client/` 包含按服务名分组的子目录；`logic/` 和 `server/` 也按服务名分组。
- 两个服务共享一个入口文件（`multisvc.go`）和配置。
- 不使用 `--multiple` 时，goctl 只允许每个 proto 文件有一个 `service` 块。
- 所有服务共享同一个 `config.go` 和 `servicecontext.go`。
