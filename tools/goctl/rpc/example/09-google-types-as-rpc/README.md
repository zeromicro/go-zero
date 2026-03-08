# Example 09: Google Types as RPC Parameters

This example demonstrates using Google protobuf well-known types **directly** as RPC request or response types (not just as message fields).

## Proto Definition

`service.proto` uses `google.protobuf.Empty` and `google.protobuf.Timestamp` directly as RPC request/response types.

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
goctl rpc protoc service.proto \
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
│   └── healthsvc.yaml
├── go.mod
├── healthservice
│   └── healthservice.go
├── healthsvc.go
├── internal
│   ├── config
│   │   └── config.go
│   ├── logic
│   │   ├── gettimelogic.go
│   │   └── pinglogic.go
│   ├── server
│   │   └── healthserviceserver.go
│   └── svc
│       └── servicecontext.go
└── pb
    ├── service.pb.go
    └── service_grpc.pb.go
```

## Key Points

- Uses Google well-known types (`google.protobuf.Empty`, `google.protobuf.Timestamp`) directly as RPC request/response types (not just message fields).
- goctl correctly maps these to Go types (`emptypb.Empty`, `timestamppb.Timestamp`) and generates proper imports.
- This differs from Example 06 where well-known types are used as message fields.
