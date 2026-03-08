# 示例 07：外部 Proto — 相同 `go_package`

本示例演示从外部目录导入 proto 文件，且两个文件共享**相同**的 `go_package`。

## Proto 定义

`service.proto` 和 `ext.proto` 使用相同的 `go_package`：

```protobuf
option go_package = "example.com/demo/pb";
```

源码布局：

```
07-external-proto-same-pkg/
├── ext_protos
│   └── ext.proto        # 外部 proto（go_package = "example.com/demo/pb"）
├── service.proto        # 服务定义（go_package = "example.com/demo/pb"）
├── README.md
└── README-cn.md
```

- `ext.proto` 位于独立目录（`ext_protos/`），但与 `service.proto` 有相同的 `go_package`。
- `service.proto` 导入 `ext.proto`，使用 `ext.ExtReq` / `ext.ExtReply` 作为 RPC 类型。

## 生成命令

首先，在输出目录中初始化 `go.mod`：

```bash
mkdir -p output && cd output && go mod init example.com/demo && cd ..
```

然后生成代码（注意 `-I ./ext_protos`）：

```bash
goctl rpc protoc service.proto \
  --go_out=output \
  --go-grpc_out=output \
  --zrpc_out=output \
  --go_opt=module=example.com/demo \
  --go-grpc_opt=module=example.com/demo \
  --module=example.com/demo \
  -I . -I ./ext_protos
```

生成的目录结构：

```
output/
├── etc
│   └── svc.yaml
├── go.mod
├── internal
│   ├── config
│   │   └── config.go
│   ├── logic
│   │   └── querylogic.go
│   ├── server
│   │   └── queryserviceserver.go
│   └── svc
│       └── servicecontext.go
├── pb
│   ├── ext.pb.go
│   ├── service.pb.go
│   └── service_grpc.pb.go
├── queryservice
│   └── queryservice.go
└── svc.go
```

## 要点说明

- `ext.proto` 位于独立目录（`ext_protos/`），但与 `service.proto` 有相同的 `go_package`。
- 使用 `-I ./ext_protos` 将外部目录添加到 proto 搜索路径。
- 当外部 proto 有**相同**的 `go_package` 时，所有类型合并到一个 Go 包中——无需跨包导入。
