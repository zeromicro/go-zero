# Changelog

## Unreleased

### New Features

#### External Proto Import Support (`--proto_path` / `-I`)

Added support for importing proto files from external directories via `-I` / `--proto_path` flags, with full transitive dependency resolution.

**Affected files:**
- `generator/gen.go` — Added `ProtoPaths` field to `ZRpcContext`; added `resolveImportedProtos()` to populate `ImportedProtos` before code generation.
- `generator/genpb.go` — Added `buildProtocCmd()` to automatically discover and append transitively imported proto files to the `protoc` command; added `relativeToProtoPath()` to compute correct relative paths for protoc output.
- `parser/import.go` — New file (major addition). Implements `ResolveImports()` for recursive transitive import resolution, `ParseImportedProtos()` for extracting `go_package` / `package` metadata from imported protos, and `BuildProtoPackageMap()` for O(1) lookup by proto package name.
- `parser/proto.go` — Added `ImportedProtos []ImportedProto` field to the `Proto` struct.
- `cli/cli.go` — Passes `ProtoPaths` from `RPCNew` to `ZRpcContext`.
- `cli/zrpc.go` — Passes `VarStringSliceProtoPath` to `ZRpcContext.ProtoPaths`.

**Before vs After:**

| | Before | After |
|---|---|---|
| Proto imports from external dirs | ❌ Not supported, all types must be in the same file | ✅ Use `-I ./ext_protos` to add search paths |
| Transitive imports (A → B → C) | ❌ Only direct imports recognized | ✅ Recursively resolves all transitive dependencies |
| Imported proto `.pb.go` generation | ❌ Manual, must run protoc separately for each file | ✅ Automatic, imported protos appended to protoc command |
| Proto search paths | ❌ Only source file directory | ✅ Multiple `-I` paths, same as protoc |

**Behavior:**
- Transitively walks all `import` declarations in proto files, skipping `google/*` well-known types.
- Searches each `-I` directory for imported files, silently skipping system-level protos not found in user paths.
- Appends discovered proto files to the `protoc` command so their `.pb.go` files are generated alongside the main proto.

---

#### Cross-Package Type Resolution

When an imported proto has a **different** `go_package` from the main proto, goctl now automatically generates the correct Go import paths and qualified type references in server, logic, and client code.

**Affected files:**
- `generator/typeref.go` — New file. Core type resolution engine:
  - `resolveRPCTypeRef()` — Resolves proto RPC types (simple, same-package dotted, cross-package dotted, Google WKT) to Go type references with correct import paths.
  - `resolveCallTypeRef()` — Variant for client code generation with type alias support.
  - `googleWKTTable` — Mapping table for all 16 Google well-known types to their Go equivalents.
- `generator/genserver.go` — `genFunctions()` now calls `resolveRPCTypeRef()` for request/response types and collects extra import paths.
- `generator/genlogic.go` — `genLogicFunction()` uses `resolveRPCTypeRef()`; added `addLogicImports()` to conditionally include main pb import and cross-package imports.
- `generator/gencall.go` — `genFunction()` and `getInterfaceFuncs()` use `resolveCallTypeRef()` for type aliases and extra imports; added `buildExtraImportLines()` helper.
- `generator/call.tpl` — Added `{{.extraImports}}` placeholder for cross-package import lines.

**Before vs After:**

| Proto type | Before | After |
|---|---|---|
| `GetReq` (same file) | `pb.GetReq` | `pb.GetReq` (unchanged) |
| `ext.ExtReq` (same `go_package`) | ❌ Error: "request type must defined in" | ✅ `pb.ExtReq` — merged into main package |
| `common.TypesReq` (different `go_package`) | ❌ Error: "request type must defined in" | ✅ `common.TypesReq` + auto-generated `import "example.com/demo/pb/common"` |
| `google.protobuf.Empty` | ❌ Error: "request type must defined in" | ✅ `emptypb.Empty` + auto-generated import |

**Behavior:**
- Simple types (e.g., `GetReq`) resolve to `pb.GetReq` with no extra import.
- Same-package dotted types (e.g., `ext.ExtReq` where `ext` has the same `go_package`) resolve to `pb.ExtReq`.
- Cross-package dotted types (e.g., `common.TypesReq` where `common` has a different `go_package`) resolve to `common.TypesReq` with the correct Go import path added automatically.

---

#### Google Well-Known Types as RPC Parameters

Google protobuf well-known types can now be used directly as RPC request/response types (not just as message fields).

**Affected files:**
- `generator/typeref.go` — `resolveGoogleWKT()` + `googleWKTTable` handles all standard types.

**Before vs After:**

| Proto Type | Before (as RPC param) | After (as RPC param) |
|---|---|---|
| `google.protobuf.Empty` | ❌ Error | ✅ `emptypb.Empty` |
| `google.protobuf.Timestamp` | ❌ Error | ✅ `timestamppb.Timestamp` |
| `google.protobuf.Duration` | ❌ Error | ✅ `durationpb.Duration` |
| `google.protobuf.Any` | ❌ Error | ✅ `anypb.Any` |
| `google.protobuf.Struct` | ❌ Error | ✅ `structpb.Struct` |
| `google.protobuf.FieldMask` | ❌ Error | ✅ `fieldmaskpb.FieldMask` |
| `google.protobuf.*Value` | ❌ Error | ✅ `wrapperspb.*Value` |

> Note: These types were already usable as **message fields** before. This change allows them as **RPC request/response types** directly.

**Supported types:**

| Proto Type | Go Type |
|---|---|
| `google.protobuf.Empty` | `emptypb.Empty` |
| `google.protobuf.Timestamp` | `timestamppb.Timestamp` |
| `google.protobuf.Duration` | `durationpb.Duration` |
| `google.protobuf.Any` | `anypb.Any` |
| `google.protobuf.Struct` | `structpb.Struct` |
| `google.protobuf.Value` | `structpb.Value` |
| `google.protobuf.ListValue` | `structpb.ListValue` |
| `google.protobuf.FieldMask` | `fieldmaskpb.FieldMask` |
| `google.protobuf.*Value` (wrappers) | `wrapperspb.*Value` |

---

### Breaking Changes

#### Dotted Type Names Now Allowed in RPC Definitions

Previously, goctl rejected any RPC request/response type containing a dot (e.g., `base.Req`), requiring all types to be defined in the same proto file. This restriction has been removed.

**Before vs After:**

| Proto Definition | Before | After |
|---|---|---|
| `rpc Fetch(base.Req) returns (base.Reply)` | ❌ Parse error: "request type must defined in xxx.proto" | ✅ Parsed successfully, `base.Req` resolved via imported proto |
| `rpc Ping(google.protobuf.Empty) returns (Reply)` | ❌ Parse error: "request type must defined in xxx.proto" | ✅ Parsed successfully, resolved to `emptypb.Empty` |

**Affected files:**
- `parser/service.go` — Removed the validation loop that rejected dotted type names with `"request type must defined in"` / `"returns type must defined in"` errors.
- `parser/parser_test.go` — `TestDefaultProtoParseCaseInvalidRequestType` and `TestDefaultProtoParseCaseInvalidResponseType` renamed and updated to verify that dotted types now parse successfully.
