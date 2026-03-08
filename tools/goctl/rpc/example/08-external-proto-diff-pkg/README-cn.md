# 示例 08：外部 Proto — 不同 `go_package`

本示例演示从外部目录导入 proto 文件，且文件具有**不同**的 `go_package` 值，需要在生成的 Go 代码中进行跨包导入。

## Proto 定义

proto 文件使用不同的 `go_package` 值：

- `service.proto`：`go_package = "example.com/demo/pb"`
- `ext_protos/common/types.proto`：`go_package = "example.com/demo/pb/common"`

源码布局：

```
08-external-proto-diff-pkg/
├── ext_protos
│   └── common
│       └── types.proto    # 外部 proto（go_package = "example.com/demo/pb/common"）
├── service.proto          # 服务定义（go_package = "example.com/demo/pb"）
├── README.md
└── README-cn.md
```

- `types.proto` 的 `go_package = "example.com/demo/pb/common"` — **不同**的 Go 包。
- `service.proto` 直接使用 `common.ExtReq` / `common.ExtReply` 作为 RPC 参数类型。

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
├── dataservice
│   └── dataservice.go
├── etc
│   └── svc.yaml
├── go.mod
├── internal
│   ├── config
│   │   └── config.go
│   ├── logic
│   │   └── fetchlogic.go
│   ├── server
│   │   └── dataserviceserver.go
│   └── svc
│       └── servicecontext.go
├── pb
│   ├── common
│   │   └── types.pb.go
│   ├── service.pb.go
│   └── service_grpc.pb.go
└── svc.go
```

## 要点说明

- 当外部 proto 有**不同**的 `go_package` 时，goctl 会自动生成跨包 Go 导入。
- goctl 通过解析导入 proto 的 `go_package` 选项，将 proto 包名（如 `common`）映射到正确的 Go 导入路径。
- `service.proto` 直接使用 `common.ExtReq` / `common.ExtReply` 作为 RPC 参数类型。
