# 示例 02：导入同级 Proto 文件

本示例演示如何导入同一目录下的 proto 文件。

## Proto 定义

同一目录下的两个 proto 文件共享相同的 `go_package`：

- `types.proto` — 定义共享消息类型（`User`）。
- `user.proto` — 定义 RPC 服务，导入 `types.proto`。

两个文件使用相同的 `go_package`，采用完整模块路径：

```protobuf
option go_package = "example.com/demo/pb";
```

`user.proto` 通过以下方式导入 `types.proto`：

```protobuf
import "types.proto";
```

## 生成命令

首先，在输出目录中初始化 `go.mod`：

```bash
mkdir -p output && cd output && go mod init example.com/demo && cd ..
```

然后生成代码：

```bash
goctl rpc protoc user.proto \
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
│   └── usersvc.yaml
├── go.mod
├── internal
│   ├── config
│   │   └── config.go
│   ├── logic
│   │   ├── createuserlogic.go
│   │   └── getuserlogic.go
│   ├── server
│   │   └── userserviceserver.go
│   └── svc
│       └── servicecontext.go
├── pb
│   ├── types.pb.go
│   ├── user.pb.go
│   └── user_grpc.pb.go
├── userservice
│   └── userservice.go
└── usersvc.go
```

## 要点说明

- 两个 proto 文件（`user.proto` 和 `types.proto`）共享相同的 `go_package = "example.com/demo/pb"`，编译到同一个 Go 包中。
- `user.proto` 通过 `import "types.proto"` 导入 `types.proto`。
- 多个 proto 文件共享相同的 `go_package` 时，它们会编译到同一个 Go 包中。
- 只需将包含 `service` 定义的 proto 文件传递给 `goctl rpc protoc`。
- 导入的 proto 文件会被 protoc 自动编译，并由 goctl 自动解析。
