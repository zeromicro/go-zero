# Example 02: Importing a Sibling Proto File

This example demonstrates importing a proto file from the same directory.

## Proto Definition

Two proto files in the same directory share the same `go_package`:

- `types.proto` — Defines shared message types (`User`).
- `user.proto` — Defines the RPC service, importing `types.proto`.

Both files use the same `go_package` with a full module path:

```protobuf
option go_package = "example.com/demo/pb";
```

`user.proto` imports `types.proto` via:

```protobuf
import "types.proto";
```

## Generation Commands

First, initialize the output directory with a `go.mod`:

```bash
mkdir -p output && cd output && go mod init example.com/demo && cd ..
```

Then generate the code:

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

Generated directory structure:

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

## Key Points

- Two proto files (`user.proto` and `types.proto`) share the same `go_package = "example.com/demo/pb"`, compiled into a single Go package.
- `user.proto` imports `types.proto` via `import "types.proto"`.
- When multiple proto files share the same `go_package`, they compile into a single Go package.
- Only the proto file containing `service` definitions needs to be passed to `goctl rpc protoc`.
- The imported proto is automatically compiled by protoc and resolved by goctl.
