# RPC Examples

English | [中文](README-cn.md) | [한국어](README-ko.md)

This directory contains complete examples for all `goctl rpc` code generation scenarios.

Each example includes:
- `.proto` source files
- `README.md` (English), `README-cn.md` (中文), and `README-ko.md` (한국어) documentation

## Examples

| # | Directory | Scenario | Key Flags |
|---|-----------|----------|-----------|
| 01 | [01-basic](01-basic/) | Basic single service, no imports | — |
| 02 | [02-import-sibling](02-import-sibling/) | Import sibling proto file | `--proto_path=.` |
| 03 | [03-import-subdir](03-import-subdir/) | Import proto from subdirectory | `--proto_path=.` |
| 04 | [04-transitive-import](04-transitive-import/) | Transitive imports (A → B → C) | `--proto_path=.` |
| 05 | [05-multiple-services](05-multiple-services/) | Multiple services in one proto | `--multiple` |
| 06 | [06-wellknown-types](06-wellknown-types/) | Google well-known types in messages | — |
| 07 | [07-external-proto-same-pkg](07-external-proto-same-pkg/) | External proto, same `go_package` | `-I ./ext_protos` |
| 08 | [08-external-proto-diff-pkg](08-external-proto-diff-pkg/) | External proto, different `go_package` | `-I ./ext_protos` |
| 09 | [09-google-types-as-rpc](09-google-types-as-rpc/) | Google well-known types as RPC parameters | — |
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
