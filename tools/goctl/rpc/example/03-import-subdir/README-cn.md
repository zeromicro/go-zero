# 示例 03：导入子目录中的 Proto 文件

本示例演示如何导入子目录中的 proto 文件。

## Proto 定义

两个 proto 文件有**不同**的 `go_package` 值：

- `order.proto` — 定义 `OrderService`，导入 `common/types.proto`。

```protobuf
option go_package = "example.com/demo/pb";
```

- `common/types.proto` — 定义可复用的分页和排序消息。

```protobuf
option go_package = "example.com/demo/pb/common";
```

`order.proto` 从子目录导入 `common/types.proto`：

```protobuf
import "common/types.proto";
```

注意两个文件的 `go_package` **不同**，因此会编译到不同的 Go 包中。

## 生成命令

首先，在输出目录中初始化 `go.mod`：

```bash
mkdir -p output && cd output && go mod init example.com/demo && cd ..
```

然后生成代码：

```bash
goctl rpc protoc order.proto \
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
│   └── ordersvc.yaml
├── go.mod
├── internal
│   ├── config
│   │   └── config.go
│   ├── logic
│   │   ├── getorderlogic.go
│   │   └── listorderslogic.go
│   ├── server
│   │   └── orderserviceserver.go
│   └── svc
│       └── servicecontext.go
├── orderservice
│   └── orderservice.go
├── ordersvc.go
└── pb
    ├── common
    │   └── types.pb.go
    ├── order.pb.go
    └── order_grpc.pb.go
```

## 要点说明

- 两个 proto 文件有**不同**的 `go_package` 值，编译到不同的 Go 包中（`pb/` 和 `pb/common/`）。
- `order.proto` 从子目录导入 `common/types.proto`。
- 当导入的 proto 文件有不同的 `go_package` 时，goctl 会自动生成跨包导入。
- `-I .` 告诉 protoc 从当前目录开始搜索，使其能够找到 `common/types.proto`。
