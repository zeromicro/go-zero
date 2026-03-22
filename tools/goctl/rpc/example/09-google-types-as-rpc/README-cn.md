# 示例 09：Google 类型作为 RPC 参数

本示例演示将 Google protobuf 知名类型**直接**用作 RPC 请求或响应类型（而不仅仅是消息字段）。

## Proto 定义

`service.proto` 使用 `google.protobuf.Empty` 和 `google.protobuf.Timestamp` 直接作为 RPC 请求/响应类型。

`go_package` 使用完整的模块路径：

```protobuf
option go_package = "example.com/demo/pb";
```

## 生成命令

首先，在输出目录中初始化 `go.mod`：

```bash
mkdir -p output && cd output && go mod init example.com/demo && cd ..
```

然后生成代码：

```bash
goctl rpc protoc service.proto \
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
│   └── healthsvc.yaml
├── go.mod
├── healthservice
│   └── healthservice.go
├── healthsvc.go
├── internal
│   ├── config
│   │   └── config.go
│   ├── logic
│   │   ├── gettimelogic.go
│   │   └── pinglogic.go
│   ├── server
│   │   └── healthserviceserver.go
│   └── svc
│       └── servicecontext.go
└── pb
    ├── service.pb.go
    └── service_grpc.pb.go
```

## 要点说明

- 使用 Google 知名类型（`google.protobuf.Empty`、`google.protobuf.Timestamp`）直接作为 RPC 请求/响应类型（不仅仅是消息字段）。
- goctl 正确将其映射到 Go 类型（`emptypb.Empty`、`timestamppb.Timestamp`）并生成正确的导入。
- 这与示例 06 不同，示例 06 中知名类型用作消息字段。
