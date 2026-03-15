# Example 03: Importing Proto from a Subdirectory

This example demonstrates importing a proto file from a subdirectory.

## Proto Definition

Two proto files with **different** `go_package` values:

- `order.proto` — Defines the `OrderService`, imports `common/types.proto`.

```protobuf
option go_package = "example.com/demo/pb";
```

- `common/types.proto` — Defines reusable pagination and sorting messages.

```protobuf
option go_package = "example.com/demo/pb/common";
```

`order.proto` imports `common/types.proto` from a subdirectory:

```protobuf
import "common/types.proto";
```

Note that the two files have **different** `go_package` values, so they compile into separate Go packages.

## Generation Commands

First, initialize the output directory with a `go.mod`:

```bash
mkdir -p output && cd output && go mod init example.com/demo && cd ..
```

Then generate the code:

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

Generated directory structure:

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

## Key Points

- Two proto files have **different** `go_package` values, so they compile into separate Go packages (`pb/` and `pb/common/`).
- `order.proto` imports `common/types.proto` from a subdirectory.
- When imported protos have a different `go_package`, goctl automatically generates cross-package imports.
- The `-I .` flag tells protoc to search from the current directory, enabling it to find `common/types.proto`.
