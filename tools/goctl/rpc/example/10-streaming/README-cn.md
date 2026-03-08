# 示例 10：流式 RPC

本示例演示 gRPC 的三种流式通信模式：服务端流、客户端流和双向流。

## Proto 定义

`stream.proto` 定义了三个 RPC 方法，演示每种流式模式。

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
goctl rpc protoc stream.proto \
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
│   └── streamsvc.yaml
├── go.mod
├── internal
│   ├── config
│   │   └── config.go
│   ├── logic
│   │   ├── bidistreamlogic.go
│   │   ├── clientstreamlogic.go
│   │   └── serverstreamlogic.go
│   ├── server
│   │   └── streamserviceserver.go
│   └── svc
│       └── servicecontext.go
├── pb
│   ├── stream.pb.go
│   └── stream_grpc.pb.go
├── streamservice
│   └── streamservice.go
└── streamsvc.go
```

## 要点说明

- 支持三种流式模式：服务端流（响应带 `stream`）、客户端流（请求带 `stream`）和双向流（两端都带 `stream`）。
- goctl 为每个流式 RPC 方法生成独立的逻辑文件。
- 流式客户端代码不会自动生成，需直接使用 gRPC 客户端。
