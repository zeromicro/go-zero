# Example 04: Transitive Imports

English | [中文](README-cn.md) | [한국어](README-ko.md)

This example demonstrates transitive proto imports, where A imports B and B imports C.

## Proto Definition

Three proto files form a transitive import chain, all sharing the same `go_package`:

```protobuf
option go_package = "example.com/demo/pb";
```

- `base.proto` — Layer C: defines base types (`BaseResp`).
- `middleware.proto` — Layer B: imports `base.proto`, defines `RequestMeta`.
- `main.proto` — Layer A: imports `middleware.proto`, defines the `PingService` (entry point).

Import chain: `main.proto` → `middleware.proto` → `base.proto`

## Generation Commands

First, initialize the output directory with a `go.mod`:

```bash
mkdir -p output && cd output && go mod init example.com/demo && cd ..
```

Then generate the code:

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

Generated directory structure:

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

## Key Points

- Three proto files (`base.proto` → `middleware.proto` → `main.proto`) form a transitive import chain.
- goctl recursively resolves all transitive imports automatically.
- All three files share the same `go_package = "example.com/demo/pb"`.
- You only need to specify the entry proto file — goctl and protoc handle the rest.
- Circular imports are detected and will cause an error (same as protoc behavior).
