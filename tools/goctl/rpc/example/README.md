# RPC Examples

This directory contains complete examples for all `goctl rpc` code generation scenarios.

Each example includes:
- `.proto` source files
- `README.md` (English) and `README-cn.md` (中文) documentation

## Examples

| # | Directory | Scenario | Key Flags |
|---|-----------|----------|-----------|
| 01 | [01-basic](01-basic/) | Basic single service, no imports | — |
| 02 | [02-import-sibling](02-import-sibling/) | Import sibling proto file | `--proto_path=.` |
| 03 | [03-import-subdir](03-import-subdir/) | Import proto from subdirectory | `--proto_path=.` |
| 04 | [04-transitive-import](04-transitive-import/) | Transitive imports (A → B → C) | `--proto_path=.` |
| 05 | [05-multiple-services](05-multiple-services/) | Multiple services in one proto | `--multiple` |
| 06 | [06-wellknown-types](06-wellknown-types/) | Google well-known types in messages | `--proto_path=$PROTOC_INCLUDE` |
| 07 | [07-external-proto-same-pkg](07-external-proto-same-pkg/) | External proto, same `go_package` | `-I ./ext_protos` |
| 08 | [08-external-proto-diff-pkg](08-external-proto-diff-pkg/) | External proto, different `go_package` | `-I ./ext_protos` |
| 09 | [09-google-types-as-rpc](09-google-types-as-rpc/) | Google types as RPC parameters | `--proto_path=$PROTOC_INCLUDE` |
| 10 | [10-streaming](10-streaming/) | Server/client/bidirectional streaming | — |

## Prerequisites

- [Go](https://go.dev/) 1.22+
- [protoc](https://github.com/protocolbuffers/protobuf/releases) (Protocol Buffers compiler)
- [protoc-gen-go](https://pkg.go.dev/google.golang.org/protobuf/cmd/protoc-gen-go) and [protoc-gen-go-grpc](https://pkg.go.dev/google.golang.org/grpc/cmd/protoc-gen-go-grpc)
- [goctl](https://github.com/zeromicro/go-zero/tree/master/tools/goctl)

## Quick Start

```bash
# Install protoc plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Try the basic example
cd 01-basic
mkdir -p output && cd output && go mod init example.com/demo && cd ..
goctl rpc protoc greeter.proto \
  --go_out=output --go-grpc_out=output --zrpc_out=output \
  --go_opt=module=example.com/demo --go-grpc_opt=module=example.com/demo \
  --module=example.com/demo -I .
```

---

# RPC 示例

本目录包含所有 `goctl rpc` 代码生成场景的完整示例。

每个示例包含：
- `.proto` 源文件
- `README.md`（英文）和 `README-cn.md`（中文）文档

## 示例列表

| # | 目录 | 场景 | 关键标志 |
|---|------|------|---------|
| 01 | [01-basic](01-basic/) | 基础单服务，无导入 | — |
| 02 | [02-import-sibling](02-import-sibling/) | 导入同级 proto 文件 | `--proto_path=.` |
| 03 | [03-import-subdir](03-import-subdir/) | 导入子目录中的 proto | `--proto_path=.` |
| 04 | [04-transitive-import](04-transitive-import/) | 传递性导入（A → B → C） | `--proto_path=.` |
| 05 | [05-multiple-services](05-multiple-services/) | 单 proto 多服务 | `--multiple` |
| 06 | [06-wellknown-types](06-wellknown-types/) | 消息中使用 Google 标准类型 | `--proto_path=$PROTOC_INCLUDE` |
| 07 | [07-external-proto-same-pkg](07-external-proto-same-pkg/) | 外部 proto，相同 `go_package` | `-I ./ext_protos` |
| 08 | [08-external-proto-diff-pkg](08-external-proto-diff-pkg/) | 外部 proto，不同 `go_package` | `-I ./ext_protos` |
| 09 | [09-google-types-as-rpc](09-google-types-as-rpc/) | Google 类型作为 RPC 参数 | `--proto_path=$PROTOC_INCLUDE` |
| 10 | [10-streaming](10-streaming/) | 服务端/客户端/双向流 | — |

## 前置条件

- [Go](https://go.dev/) 1.22+
- [protoc](https://github.com/protocolbuffers/protobuf/releases)（Protocol Buffers 编译器）
- [protoc-gen-go](https://pkg.go.dev/google.golang.org/protobuf/cmd/protoc-gen-go) 和 [protoc-gen-go-grpc](https://pkg.go.dev/google.golang.org/grpc/cmd/protoc-gen-go-grpc)
- [goctl](https://github.com/zeromicro/go-zero/tree/master/tools/goctl)

## 快速开始

```bash
# 安装 protoc 插件
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 试试基础示例
cd 01-basic
mkdir -p output && cd output && go mod init example.com/demo && cd ..
goctl rpc protoc greeter.proto \
  --go_out=output --go-grpc_out=output --zrpc_out=output \
  --go_opt=module=example.com/demo --go-grpc_opt=module=example.com/demo \
  --module=example.com/demo -I .
```
