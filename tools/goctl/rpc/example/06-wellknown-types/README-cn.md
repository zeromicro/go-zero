# 示例 06：知名类型

本示例演示如何使用 Google protobuf 知名类型（`Timestamp`、`Duration`、`Any`）作为消息字段。

## Proto 定义

`events.proto` 使用 `google.protobuf.Timestamp` 作为消息字段类型。

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
goctl rpc protoc events.proto \
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
│   └── eventsvc.yaml
├── eventservice
│   └── eventservice.go
├── eventsvc.go
├── go.mod
├── internal
│   ├── config
│   │   └── config.go
│   ├── logic
│   │   ├── createeventlogic.go
│   │   └── listeventslogic.go
│   ├── server
│   │   └── eventserviceserver.go
│   └── svc
│       └── servicecontext.go
└── pb
    ├── events.pb.go
    └── events_grpc.pb.go
```

## 要点说明

- 使用 Google 知名类型（`google.protobuf.Timestamp`、`google.protobuf.Duration`、`google.protobuf.Any`）作为消息字段。
- goctl 自动将知名类型映射到 Go 导入包（`timestamppb`、`durationpb`、`anypb` 等）。
- 如果 protoc 已正确安装，知名类型无需额外的 `--proto_path`。
