# Example 01: Basic RPC Service

This is the simplest example of generating an RPC service with goctl.

## Proto Definition

A single `greeter.proto` file with one service and one RPC method, no external imports.

The `go_package` uses a full module path:

```protobuf
option go_package = "example.com/demo/greeter";
```

## Generation Commands

### Method 1: Quick Start with `goctl rpc new`

```bash
# Create a complete RPC project with one command
goctl rpc new greeter
```

This generates the proto file and service code together:

```
greeter/
├── etc
│   └── greeter.yaml
├── greeter
│   ├── greeter.pb.go
│   └── greeter_grpc.pb.go
├── greeter.go
├── greeter.proto
├── greeterclient
│   └── greeter.go
└── internal
    ├── config
    │   └── config.go
    ├── logic
    │   └── pinglogic.go
    ├── server
    │   └── greeterserver.go
    └── svc
        └── servicecontext.go
```

### Method 2: Generate from an Existing Proto

First, initialize the output directory with a `go.mod`:

```bash
mkdir -p output && cd output && go mod init example.com/demo && cd ..
```

Then generate the code:

```bash
goctl rpc protoc greeter.proto \
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
│   └── greeter.yaml
├── go.mod
├── greeter
│   ├── greeter.pb.go
│   └── greeter_grpc.pb.go
├── greeter.go
├── greeterclient
│   └── greeter.go
└── internal
    ├── config
    │   └── config.go
    ├── logic
    │   └── sayhellologic.go
    ├── server
    │   └── greeterserver.go
    └── svc
        └── servicecontext.go
```

## Key Points

- This is the simplest scenario: one proto file, one service, one RPC method.
- The `go_package` uses a full module path (`example.com/demo/greeter`), not a relative path.
- The `--module` flag tells goctl the Go module name; `--go_opt=module=...` and `--go-grpc_opt=module=...` tell protoc to strip the module prefix from output paths.
- The `--zrpc_out` flag specifies where the goctl-generated service code goes.
- The `--go_out` and `--go-grpc_out` flags specify where protoc-generated code goes.
- Edit the logic file (`internal/logic/sayhellologic.go`) to implement your business logic.
