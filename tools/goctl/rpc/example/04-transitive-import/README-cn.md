# 示例 04：传递性导入

本示例演示 proto 的传递性导入，即 A 导入 B，B 导入 C。

## Proto 定义

三个 proto 文件形成传递导入链，共享相同的 `go_package`：

```protobuf
option go_package = "example.com/demo/pb";
```

- `base.proto` — C 层：定义基础类型（`BaseResp`）。
- `middleware.proto` — B 层：导入 `base.proto`，定义 `RequestMeta`。
- `main.proto` — A 层：导入 `middleware.proto`，定义 `PingService`（入口文件）。

导入链：`main.proto` → `middleware.proto` → `base.proto`

## 生成命令

首先，在输出目录中初始化 `go.mod`：

```bash
mkdir -p output && cd output && go mod init example.com/demo && cd ..
```

然后生成代码：

```bash
goctl rpc protoc main.proto \
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
│   └── pingsvc.yaml
├── go.mod
├── internal
│   ├── config
│   │   └── config.go
│   ├── logic
│   │   └── pinglogic.go
│   ├── server
│   │   └── pingserviceserver.go
│   └── svc
│       └── servicecontext.go
├── pb
│   ├── base.pb.go
│   ├── main.pb.go
│   ├── main_grpc.pb.go
│   └── middleware.pb.go
├── pingservice
│   └── pingservice.go
└── pingsvc.go
```

## 要点说明

- 三个 proto 文件（`base.proto` → `middleware.proto` → `main.proto`）形成传递导入链。
- goctl 自动递归解析所有传递导入。
- 三个文件共享相同的 `go_package = "example.com/demo/pb"`。
- 只需指定入口 proto 文件，goctl 和 protoc 会自动处理其余部分。
- 循环导入会被检测并报错（与 protoc 行为一致）。
