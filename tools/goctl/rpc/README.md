# goctl rpc — RPC Code Generation

English | [中文](README-cn.md)

goctl rpc is the RPC service code generation module of the `goctl` scaffold. It generates a complete zRPC service from `.proto` files. You only need to write the proto definition and business logic — all boilerplate code is generated automatically.

## Features

- **protoc compatible**: Fully compatible with protoc, all protoc arguments are passed through
- **External proto imports**: Cross-directory and cross-package proto imports with automatic transitive dependency resolution
- **Multiple services**: Define multiple services in a single proto file, auto-grouped by service name
- **Streaming support**: Server streaming, client streaming, and bidirectional streaming
- **Google well-known types**: Automatic recognition of `google.protobuf.*` types with correct Go imports
- **Client generation**: Auto-generated RPC client wrapper code

## Prerequisites

```bash
# Install protoc plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## Quick Start

### Method 1: Create a Service Instantly

```bash
goctl rpc new greeter
```

Generates a complete project structure:

```
greeter/
├── etc/
│   └── greeter.yaml
├── greeter/
│   ├── greeter.pb.go
│   └── greeter_grpc.pb.go
├── greeter.go
├── greeter.proto
├── greeterclient/
│   └── greeter.go
└── internal/
    ├── config/
    │   └── config.go
    ├── logic/
    │   └── pinglogic.go
    ├── server/
    │   └── greeterserver.go
    └── svc/
        └── servicecontext.go
```

### Method 2: Generate from a Proto File

1. Generate a proto template:

```bash
goctl rpc template -o=user.proto
```

2. Initialize the output directory and generate service code:

```bash
mkdir -p output && cd output && go mod init example.com/demo && cd ..
goctl rpc protoc user.proto \
  --go_out=output --go-grpc_out=output --zrpc_out=output \
  --go_opt=module=example.com/demo --go-grpc_opt=module=example.com/demo \
  --module=example.com/demo -I .
```

---

## Command Reference

### `goctl rpc protoc`

Generate zRPC service code from a `.proto` file.

```bash
goctl rpc protoc <proto_file> [flags]
```

**Examples:**

```bash
# Basic usage
goctl rpc protoc user.proto \
  --go_out=output --go-grpc_out=output --zrpc_out=output \
  --go_opt=module=example.com/demo --go-grpc_opt=module=example.com/demo \
  --module=example.com/demo -I .

# Multiple services mode
goctl rpc protoc multi.proto \
  --go_out=output --go-grpc_out=output --zrpc_out=output \
  --go_opt=module=example.com/demo --go-grpc_opt=module=example.com/demo \
  --module=example.com/demo -I . -m

# Import external protos
goctl rpc protoc service.proto \
  --go_out=output --go-grpc_out=output --zrpc_out=output \
  --go_opt=module=example.com/demo --go-grpc_opt=module=example.com/demo \
  --module=example.com/demo -I . -I ./shared_protos

# Use Google well-known types
goctl rpc protoc service.proto \
  --go_out=output --go-grpc_out=output --zrpc_out=output \
  --go_opt=module=example.com/demo --go-grpc_opt=module=example.com/demo \
  --module=example.com/demo -I .
```

**Flags:**

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--zrpc_out` | | string | **required** | Output directory for zRPC service code |
| `--go_out` | | string | **required** | Output directory for protoc Go code |
| `--go-grpc_out` | | string | **required** | Output directory for protoc gRPC code |
| `--go_opt` | | string | | Options for protoc-gen-go (e.g., `module=example.com/demo`) |
| `--go-grpc_opt` | | string | | Options for protoc-gen-go-grpc (e.g., `module=example.com/demo`) |
| `--proto_path` | `-I` | string[] | | Proto import search directories (repeatable) |
| `--multiple` | `-m` | bool | `false` | Multiple services mode |
| `--client` | `-c` | bool | `true` | Generate RPC client code |
| `--style` | | string | `gozero` | File naming style |
| `--module` | | string | | Custom Go module name |
| `--name-from-filename` | | bool | `false` | Use filename instead of package name for service naming |
| `--verbose` | `-v` | bool | `false` | Enable verbose logging |
| `--home` | | string | | goctl template directory |
| `--remote` | | string | | Remote template Git repository URL |
| `--branch` | | string | | Remote template branch |

### `goctl rpc new`

Quickly create a complete RPC service project.

```bash
goctl rpc new <service_name> [flags]
```

**Flags:**

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--style` | | string | `gozero` | File naming style |
| `--client` | `-c` | bool | `true` | Generate RPC client code |
| `--module` | | string | | Custom Go module name |
| `--verbose` | `-v` | bool | `false` | Enable verbose logging |
| `--idea` | | bool | `false` | Generate IDE project marker |
| `--name-from-filename` | | bool | `false` | Use filename instead of package name for service naming |
| `--home` | | string | | goctl template directory |
| `--remote` | | string | | Remote template Git repository URL |
| `--branch` | | string | | Remote template branch |

### `goctl rpc template`

Generate a proto file template.

```bash
goctl rpc template -o=<output_file> [flags]
```

**Flags:**

| Flag | Type | Description |
|------|------|-------------|
| `-o` | string | Output file path (required) |
| `--home` | string | goctl template directory |
| `--remote` | string | Remote template Git repository URL |
| `--branch` | string | Remote template branch |

---

## Feature Details

### Multiple Services Mode (`--multiple`)

When a proto file contains multiple `service` definitions, the `--multiple` flag is required.

```protobuf
service SearchService {
  rpc Search(SearchReq) returns (SearchReply);
}

service NotifyService {
  rpc Notify(NotifyReq) returns (NotifyReply);
}
```

**Directory differences with `--multiple`:**

| Feature | Default mode | `--multiple` mode |
|---------|-------------|-------------------|
| Services per proto | Exactly 1 | 1 or more |
| Client directory | Named after service | Fixed `client/` directory |
| Code organization | Flat structure | Grouped by service name |

**`--multiple=false` (default) directory structure:**

```
output/
├── greeterclient/
│   └── greeter.go
├── internal/
│   ├── logic/
│   │   └── sayhellologic.go
│   └── server/
│       └── greeterserver.go
└── ...
```

**`--multiple=true` directory structure:**

```
output/
├── client/
│   ├── searchservice/
│   │   └── searchservice.go
│   └── notifyservice/
│       └── notifyservice.go
├── internal/
│   ├── logic/
│   │   ├── searchservice/
│   │   │   └── searchlogic.go
│   │   └── notifyservice/
│   │       └── notifylogic.go
│   └── server/
│       ├── searchservice/
│       │   └── searchserviceserver.go
│       └── notifyservice/
│           └── notifyserviceserver.go
└── ...
```

### External Proto Imports (`--proto_path`)

Use `-I` / `--proto_path` to specify additional proto search directories. Supported scenarios:

- **Same-directory import**: `import "types.proto";`
- **Subdirectory import**: `import "common/types.proto";`
- **External directory import**: Proto files outside the project
- **Transitive imports**: A imports B, B imports C — goctl resolves recursively
- **Cross-package imports**: Different `go_package` values generate correct Go imports automatically

```bash
# Search multiple directories for proto files
goctl rpc protoc service.proto \
  --go_out=output --go-grpc_out=output --zrpc_out=output \
  --go_opt=module=example.com/demo --go-grpc_opt=module=example.com/demo \
  --module=example.com/demo \
  -I . -I ./shared_protos -I /path/to/external_protos
```

### Service Naming

By default, the service name is derived from the proto **package name** (e.g., `package user;` → service name `user`). This allows multiple proto files to share the same package:

```
protos/
├── user_base.proto      # package user;
├── user_auth.proto      # package user;
└── user_profile.proto   # package user;
```

All three files generate into a single `user` service.

To use the proto filename for naming (legacy behavior), add the `--name-from-filename` flag.

### Streaming RPC

All three gRPC streaming patterns are supported:

```protobuf
service StreamService {
  rpc ServerStream(Req) returns (stream Reply);       // Server streaming
  rpc ClientStream(stream Req) returns (Reply);       // Client streaming
  rpc BidiStream(stream Req) returns (stream Reply);  // Bidirectional streaming
}
```

### Google Well-Known Types

goctl automatically recognizes and handles Google protobuf well-known types:

| Proto Type | Go Type |
|-----------|---------|
| `google.protobuf.Empty` | `emptypb.Empty` |
| `google.protobuf.Timestamp` | `timestamppb.Timestamp` |
| `google.protobuf.Duration` | `durationpb.Duration` |
| `google.protobuf.Any` | `anypb.Any` |
| `google.protobuf.Struct` | `structpb.Struct` |
| `google.protobuf.FieldMask` | `fieldmaskpb.FieldMask` |
| `google.protobuf.*Value` | `wrapperspb.*Value` |

These types can be used directly as RPC parameter types — goctl generates the correct imports automatically.

---

## Examples

See the [example/](example/) directory for 10 complete examples covering all generation scenarios.

| # | Example | Scenario |
|---|---------|----------|
| 01 | [Basic service](example/01-basic/) | Single service, no imports |
| 02 | [Sibling import](example/02-import-sibling/) | Import from same directory |
| 03 | [Subdirectory import](example/03-import-subdir/) | Import from subdirectory |
| 04 | [Transitive import](example/04-transitive-import/) | A → B → C dependency chain |
| 05 | [Multiple services](example/05-multiple-services/) | `--multiple` mode |
| 06 | [Well-known types](example/06-wellknown-types/) | Timestamp etc. in messages |
| 07 | [External proto (same pkg)](example/07-external-proto-same-pkg/) | External proto, same go_package |
| 08 | [External proto (diff pkg)](example/08-external-proto-diff-pkg/) | External proto, different go_package |
| 09 | [Google types as params](example/09-google-types-as-rpc/) | Empty/Timestamp as RPC parameters |
| 10 | [Streaming](example/10-streaming/) | Server/client/bidirectional streaming |
