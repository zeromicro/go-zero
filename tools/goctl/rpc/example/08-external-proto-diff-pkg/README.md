# Example 08: External Proto — Different `go_package`

This example demonstrates importing proto files from an external directory where the files have **different** `go_package` values, requiring cross-package imports in the generated Go code.

## Proto Definition

The proto files use different `go_package` values:

- `service.proto`: `go_package = "example.com/demo/pb"`
- `ext_protos/common/types.proto`: `go_package = "example.com/demo/pb/common"`

Source layout:

```
08-external-proto-diff-pkg/
├── ext_protos
│   └── common
│       └── types.proto    # External proto (go_package = "example.com/demo/pb/common")
├── service.proto          # Service definition (go_package = "example.com/demo/pb")
├── README.md
└── README-cn.md
```

- `types.proto` has `go_package = "example.com/demo/pb/common"` — a **different** Go package.
- `service.proto` uses `common.ExtReq` / `common.ExtReply` directly as RPC parameter types.

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
├── dataservice
│   └── dataservice.go
├── etc
│   └── svc.yaml
├── go.mod
├── internal
│   ├── config
│   │   └── config.go
│   ├── logic
│   │   └── fetchlogic.go
│   ├── server
│   │   └── dataserviceserver.go
│   └── svc
│       └── servicecontext.go
├── pb
│   ├── common
│   │   └── types.pb.go
│   ├── service.pb.go
│   └── service_grpc.pb.go
└── svc.go
```

## Key Points

- When the external proto has a **different** `go_package`, goctl generates cross-package Go imports automatically.
- goctl resolves the proto package name (e.g., `common`) to the correct Go import path by parsing the imported proto's `go_package` option.
- `service.proto` uses `common.ExtReq` / `common.ExtReply` directly as RPC parameter types.
