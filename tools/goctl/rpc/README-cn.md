# goctl rpc — RPC 代码生成

[English](README.md) | 中文

goctl rpc 是 `goctl` 脚手架下的 RPC 服务代码生成模块，基于 `.proto` 文件生成完整的 zRPC 服务代码。你只需编写 proto 定义和业务逻辑，其余代码均由工具自动生成。

## 特性

- **贴近 protoc**：与 protoc 完全兼容，透传所有 protoc 参数
- **外部 Proto 导入**：支持跨目录、跨包的 proto 导入，自动解析传递性依赖
- **多服务模式**：单个 proto 文件中定义多个 service，按服务名自动分组
- **流式支持**：支持服务端流、客户端流和双向流
- **Google 标准类型**：自动识别 `google.protobuf.*` 类型并生成正确的 Go 导入
- **客户端生成**：自动生成封装好的 RPC 客户端代码

## 前置条件

```bash
# 安装 protoc 插件
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## 快速开始

### 方式一：一键创建服务

```bash
goctl rpc new greeter
```

生成完整的项目结构：

```
greeter/
├── etc/
│   └── greeter.yaml
├── greeter/
│   ├── greeter.pb.go
│   └── greeter_grpc.pb.go
├── greeter.go
├── greeter.proto
├── greeterclient/
│   └── greeter.go
└── internal/
    ├── config/
    │   └── config.go
    ├── logic/
    │   └── pinglogic.go
    ├── server/
    │   └── greeterserver.go
    └── svc/
        └── servicecontext.go
```

### 方式二：基于 Proto 文件生成

1. 生成 proto 模板：

```bash
goctl rpc template -o=user.proto
```

2. 初始化输出目录并生成服务代码：

```bash
mkdir -p output && cd output && go mod init example.com/demo && cd ..
goctl rpc protoc user.proto \
  --go_out=output --go-grpc_out=output --zrpc_out=output \
  --go_opt=module=example.com/demo --go-grpc_opt=module=example.com/demo \
  --module=example.com/demo -I .
```

---

## 命令参考

### `goctl rpc protoc`

从 `.proto` 文件生成 zRPC 服务代码。

```bash
goctl rpc protoc <proto_file> [flags]
```

**示例：**

```bash
# 基础用法
goctl rpc protoc user.proto \
  --go_out=output --go-grpc_out=output --zrpc_out=output \
  --go_opt=module=example.com/demo --go-grpc_opt=module=example.com/demo \
  --module=example.com/demo -I .

# 多服务模式
goctl rpc protoc multi.proto \
  --go_out=output --go-grpc_out=output --zrpc_out=output \
  --go_opt=module=example.com/demo --go-grpc_opt=module=example.com/demo \
  --module=example.com/demo -I . -m

# 导入外部 proto
goctl rpc protoc service.proto \
  --go_out=output --go-grpc_out=output --zrpc_out=output \
  --go_opt=module=example.com/demo --go-grpc_opt=module=example.com/demo \
  --module=example.com/demo -I . -I ./shared_protos

# 使用 Google 标准类型
goctl rpc protoc service.proto \
  --go_out=output --go-grpc_out=output --zrpc_out=output \
  --go_opt=module=example.com/demo --go-grpc_opt=module=example.com/demo \
  --module=example.com/demo -I .
```

**参数说明：**

| 参数 | 缩写 | 类型 | 默认值 | 说明 |
|------|------|------|--------|------|
| `--zrpc_out` | | string | **必填** | zRPC 服务代码输出目录 |
| `--go_out` | | string | **必填** | protoc Go 代码输出目录 |
| `--go-grpc_out` | | string | **必填** | protoc gRPC 代码输出目录 |
| `--go_opt` | | string | | protoc-gen-go 选项（如 `module=example.com/demo`） |
| `--go-grpc_opt` | | string | | protoc-gen-go-grpc 选项（如 `module=example.com/demo`） |
| `--proto_path` | `-I` | string[] | | proto 导入搜索目录（可多次指定） |
| `--multiple` | `-m` | bool | `false` | 多服务模式 |
| `--client` | `-c` | bool | `true` | 是否生成 RPC 客户端代码 |
| `--style` | | string | `gozero` | 文件命名风格 |
| `--module` | | string | | 自定义 Go module 名称 |
| `--name-from-filename` | | bool | `false` | 使用文件名而非 package 名命名服务 |
| `--verbose` | `-v` | bool | `false` | 显示详细日志 |
| `--home` | | string | | goctl 模板目录 |
| `--remote` | | string | | 远程模板 Git 仓库地址 |
| `--branch` | | string | | 远程模板分支 |

### `goctl rpc new`

快速创建一个完整的 RPC 服务项目。

```bash
goctl rpc new <service_name> [flags]
```

**参数说明：**

| 参数 | 缩写 | 类型 | 默认值 | 说明 |
|------|------|------|--------|------|
| `--style` | | string | `gozero` | 文件命名风格 |
| `--client` | `-c` | bool | `true` | 是否生成 RPC 客户端代码 |
| `--module` | | string | | 自定义 Go module 名称 |
| `--verbose` | `-v` | bool | `false` | 显示详细日志 |
| `--idea` | | bool | `false` | 生成 IDE 项目标记 |
| `--name-from-filename` | | bool | `false` | 使用文件名而非 package 名命名服务 |
| `--home` | | string | | goctl 模板目录 |
| `--remote` | | string | | 远程模板 Git 仓库地址 |
| `--branch` | | string | | 远程模板分支 |

### `goctl rpc template`

生成 proto 文件模板。

```bash
goctl rpc template -o=<output_file> [flags]
```

**参数说明：**

| 参数 | 类型 | 说明 |
|------|------|------|
| `-o` | string | 输出文件路径（必填） |
| `--home` | string | goctl 模板目录 |
| `--remote` | string | 远程模板 Git 仓库地址 |
| `--branch` | string | 远程模板分支 |

---

## 功能详解

### 多服务模式（`--multiple`）

当 proto 文件包含多个 `service` 定义时，必须使用 `--multiple` 标志。

```protobuf
service SearchService {
  rpc Search(SearchReq) returns (SearchReply);
}

service NotifyService {
  rpc Notify(NotifyReq) returns (NotifyReply);
}
```

**启用 `--multiple` 后的目录变化：**

| 特性 | 默认模式 | `--multiple` 模式 |
|------|---------|-------------------|
| 服务数量 | 仅 1 个 | 1 个或多个 |
| 客户端目录 | 以服务名命名 | 固定为 `client/` |
| 代码组织 | 扁平结构 | 按服务名分组 |

**`--multiple=false`（默认）的目录结构：**

```
output/
├── greeterclient/
│   └── greeter.go
├── internal/
│   ├── logic/
│   │   └── sayhellologic.go
│   └── server/
│       └── greeterserver.go
└── ...
```

**`--multiple=true` 的目录结构：**

```
output/
├── client/
│   ├── searchservice/
│   │   └── searchservice.go
│   └── notifyservice/
│       └── notifyservice.go
├── internal/
│   ├── logic/
│   │   ├── searchservice/
│   │   │   └── searchlogic.go
│   │   └── notifyservice/
│   │       └── notifylogic.go
│   └── server/
│       ├── searchservice/
│       │   └── searchserviceserver.go
│       └── notifyservice/
│           └── notifyserviceserver.go
└── ...
```

### 外部 Proto 导入（`--proto_path`）

通过 `-I` / `--proto_path` 指定额外的 proto 搜索目录，支持以下场景：

- **同目录导入**：`import "types.proto";`
- **子目录导入**：`import "common/types.proto";`
- **外部目录导入**：proto 文件位于项目外部
- **传递性导入**：A 导入 B，B 导入 C，goctl 自动递归解析
- **跨包导入**：不同 `go_package` 的 proto 文件，自动生成正确的 Go 导入

```bash
# 从多个目录搜索 proto 文件
goctl rpc protoc service.proto \
  --go_out=output --go-grpc_out=output --zrpc_out=output \
  --go_opt=module=example.com/demo --go-grpc_opt=module=example.com/demo \
  --module=example.com/demo \
  -I . -I ./shared_protos -I /path/to/external_protos
```

### 服务命名

默认情况下，服务名称来自 proto 的 **package 名称**（例如 `package user;` → 服务名 `user`）。这使得多个 proto 文件可以共享同一个 package：

```
protos/
├── user_base.proto      # package user;
├── user_auth.proto      # package user;
└── user_profile.proto   # package user;
```

三个文件会生成到同一个 `user` 服务中。

如需使用 proto 文件名命名（旧版行为），请添加 `--name-from-filename` 标志。

### 流式 RPC

支持 gRPC 的三种流式模式：

```protobuf
service StreamService {
  rpc ServerStream(Req) returns (stream Reply);       // 服务端流
  rpc ClientStream(stream Req) returns (Reply);       // 客户端流
  rpc BidiStream(stream Req) returns (stream Reply);  // 双向流
}
```

### Google 标准类型

goctl 自动识别并正确处理 Google protobuf 标准类型：

| Proto 类型 | Go 类型 |
|-----------|---------|
| `google.protobuf.Empty` | `emptypb.Empty` |
| `google.protobuf.Timestamp` | `timestamppb.Timestamp` |
| `google.protobuf.Duration` | `durationpb.Duration` |
| `google.protobuf.Any` | `anypb.Any` |
| `google.protobuf.Struct` | `structpb.Struct` |
| `google.protobuf.FieldMask` | `fieldmaskpb.FieldMask` |
| `google.protobuf.*Value` | `wrapperspb.*Value` |

这些类型可直接用作 RPC 参数类型，goctl 会自动生成正确的导入。

---

## 完整示例

详见 [example/](example/) 目录，包含 10 个完整示例，覆盖所有生成场景。

| # | 示例 | 场景 |
|---|------|------|
| 01 | [基础服务](example/01-basic/) | 单服务，无导入 |
| 02 | [同级导入](example/02-import-sibling/) | 导入同目录 proto |
| 03 | [子目录导入](example/03-import-subdir/) | 导入子目录 proto |
| 04 | [传递性导入](example/04-transitive-import/) | A → B → C 依赖链 |
| 05 | [多服务](example/05-multiple-services/) | `--multiple` 模式 |
| 06 | [标准类型](example/06-wellknown-types/) | 消息中使用 Timestamp 等 |
| 07 | [外部 Proto（同包）](example/07-external-proto-same-pkg/) | 外部 proto，相同 go_package |
| 08 | [外部 Proto（跨包）](example/08-external-proto-diff-pkg/) | 外部 proto，不同 go_package |
| 09 | [标准类型作参数](example/09-google-types-as-rpc/) | Empty/Timestamp 作为 RPC 参数 |
| 10 | [流式通信](example/10-streaming/) | 服务端/客户端/双向流 |
