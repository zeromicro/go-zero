# Example 07: External Proto — Same `go_package`

English | [中文](README-cn.md) | [한국어](README-ko.md)

This example demonstrates importing proto files from an external directory where both files share the **same** `go_package`.

## Proto Definition

Both `service.proto` and `ext.proto` use the same `go_package`:

```protobuf
option go_package = "example.com/demo/pb";
```

Source layout:

```
07-external-proto-same-pkg/
├── ext_protos
│   └── ext.proto        # External proto (go_package = "example.com/demo/pb")
├── service.proto        # Service definition (go_package = "example.com/demo/pb")
├── README.md
├── README-cn.md
└── README-ko.md
```

- `ext.proto` lives in a separate directory (`ext_protos/`), but has the same `go_package` as `service.proto`.
- `service.proto` imports `ext.proto` and uses `ext.ExtReq` / `ext.ExtReply` as RPC types.

## Generation Commands

First, initialize the output directory with a `go.mod`:

```bash
mkdir -p output && cd output && go mod init example.com/demo && cd ..
```

Then generate the code (note `-I ./ext_protos`):

```bash
goctl rpc protoc service.proto \
  --go_out=output \
  --go-grpc_out=output \
  --zrpc_out=output \
  --go_opt=module=example.com/demo \
  --go-grpc_opt=module=example.com/demo \
  --module=example.com/demo \
  -I . -I ./ext_protos
```

Generated directory structure:

```
output/
├── etc
│   └── svc.yaml
├── go.mod
├── internal
│   ├── config
│   │   └── config.go
│   ├── logic
│   │   └── querylogic.go
│   ├── server
│   │   └── queryserviceserver.go
│   └── svc
│       └── servicecontext.go
├── pb
│   ├── ext.pb.go
│   ├── service.pb.go
│   └── service_grpc.pb.go
├── queryservice
│   └── queryservice.go
└── svc.go
```

## Key Points

- `ext.proto` lives in a separate directory (`ext_protos/`), but has the same `go_package` as `service.proto`.
- Use `-I ./ext_protos` to add the external directory to the proto search path.
- When the external proto has the **same** `go_package`, all types merge into one Go package — no cross-package imports are needed.
