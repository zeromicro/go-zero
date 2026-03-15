# Example 05: Multiple Services (`--multiple`)

This example demonstrates generating multiple RPC services from a single proto file.

## Proto Definition

Two proto files share the same `go_package`:

```protobuf
option go_package = "example.com/demo/pb";
```

- `shared.proto` — Defines shared message types (`Meta`).
- `multi.proto` — Defines **two** services: `SearchService` and `NotifyService`.

The `-m` (or `--multiple`) flag is required when a proto file contains more than one `service` block.

## Generation Commands

First, initialize the output directory with a `go.mod`:

```bash
mkdir -p output && cd output && go mod init example.com/demo && cd ..
```

Then generate the code with the `-m` flag:

```bash
goctl rpc protoc multi.proto \
  --go_out=output \
  --go-grpc_out=output \
  --zrpc_out=output \
  --go_opt=module=example.com/demo \
  --go-grpc_opt=module=example.com/demo \
  --module=example.com/demo \
  -I . \
  -m
```

Generated directory structure:

```
output/
├── client
│   ├── notifyservice
│   │   └── notifyservice.go
│   └── searchservice
│       └── searchservice.go
├── etc
│   └── multisvc.yaml
├── go.mod
├── internal
│   ├── config
│   │   └── config.go
│   ├── logic
│   │   ├── notifyservice
│   │   │   └── notifylogic.go
│   │   └── searchservice
│   │       └── searchlogic.go
│   ├── server
│   │   ├── notifyservice
│   │   │   └── notifyserviceserver.go
│   │   └── searchservice
│   │       └── searchserviceserver.go
│   └── svc
│       └── servicecontext.go
├── multisvc.go
└── pb
    ├── multi.pb.go
    ├── multi_grpc.pb.go
    └── shared.pb.go
```

## Key Points

- The `-m` (or `--multiple`) flag enables multiple-service mode.
- In multiple mode, `client/` contains per-service subdirectories; `logic/` and `server/` are also grouped by service name.
- Both services share a single entry point (`multisvc.go`) and config.
- Without `--multiple`, goctl only allows one `service` block per proto file.
- All services share the same `config.go` and `servicecontext.go`.
