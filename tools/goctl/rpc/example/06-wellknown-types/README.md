# Example 06: Well-Known Types

English | [中文](README-cn.md) | [한국어](README-ko.md)

This example demonstrates using a Google protobuf well-known type (`Timestamp`) as a message field.

## Proto Definition

`events.proto` uses `google.protobuf.Timestamp` as a message field type.

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
goctl rpc protoc events.proto \
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
│   └── eventsvc.yaml
├── eventservice
│   └── eventservice.go
├── eventsvc.go
├── go.mod
├── internal
│   ├── config
│   │   └── config.go
│   ├── logic
│   │   ├── createeventlogic.go
│   │   └── listeventslogic.go
│   ├── server
│   │   └── eventserviceserver.go
│   └── svc
│       └── servicecontext.go
└── pb
    ├── events.pb.go
    └── events_grpc.pb.go
```

## Key Points

- Uses the Google well-known type `google.protobuf.Timestamp` as a message field.
- goctl automatically maps well-known types to Go imports; in this example, `Timestamp` maps to `timestamppb`.
- No extra `--proto_path` needed for well-known types if protoc is properly installed.
