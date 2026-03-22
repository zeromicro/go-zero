# Example 10: Streaming RPC

This example demonstrates all three gRPC streaming patterns: server streaming, client streaming, and bidirectional streaming.

## Proto Definition

`stream.proto` defines three RPC methods demonstrating each streaming pattern.

The `go_package` uses a full module path:

```protobuf
option go_package = "example.com/demo/pb";
```

## Generation Commands

First, initialize the output directory with a `go.mod`:

```bash
mkdir -p output && cd output && go mod init example.com/demo && cd ..
```

Then generate the code:

```bash
goctl rpc protoc stream.proto \
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
│   └── streamsvc.yaml
├── go.mod
├── internal
│   ├── config
│   │   └── config.go
│   ├── logic
│   │   ├── bidistreamlogic.go
│   │   ├── clientstreamlogic.go
│   │   └── serverstreamlogic.go
│   ├── server
│   │   └── streamserviceserver.go
│   └── svc
│       └── servicecontext.go
├── pb
│   ├── stream.pb.go
│   └── stream_grpc.pb.go
├── streamservice
│   └── streamservice.go
└── streamsvc.go
```

## Key Points

- Supports three streaming patterns: server streaming (`stream` on response), client streaming (`stream` on request), and bidirectional streaming (`stream` on both).
- goctl generates separate logic files for each streaming RPC method.
- Streaming client code is not auto-generated; use the gRPC client directly.
