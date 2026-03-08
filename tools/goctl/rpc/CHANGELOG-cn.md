# 变更日志

## 未发布

### 新功能

#### 外部 Proto 导入支持（`--proto_path` / `-I`）

新增通过 `-I` / `--proto_path` 标志导入外部目录中的 proto 文件，支持完整的传递性依赖解析。

**涉及文件：**
- `generator/gen.go` — `ZRpcContext` 新增 `ProtoPaths` 字段；新增 `resolveImportedProtos()` 在代码生成前填充 `ImportedProtos`。
- `generator/genpb.go` — 新增 `buildProtocCmd()` 自动发现并追加传递性导入的 proto 文件到 `protoc` 命令；新增 `relativeToProtoPath()` 计算正确的相对路径。
- `parser/import.go` — 新文件（主要新增）。实现 `ResolveImports()` 递归解析传递性导入，`ParseImportedProtos()` 提取导入 proto 的 `go_package` / `package` 元数据，`BuildProtoPackageMap()` 构建按 proto 包名的 O(1) 查找表。
- `parser/proto.go` — `Proto` 结构体新增 `ImportedProtos []ImportedProto` 字段。
- `cli/cli.go` — `RPCNew` 传递 `ProtoPaths` 到 `ZRpcContext`。
- `cli/zrpc.go` — 将 `VarStringSliceProtoPath` 传递到 `ZRpcContext.ProtoPaths`。

**前后对比：**

| | 变更前 | 变更后 |
|---|---|---|
| 从外部目录导入 Proto | ❌ 不支持，所有类型必须在同一文件中定义 | ✅ 使用 `-I ./ext_protos` 添加搜索路径 |
| 传递性导入（A → B → C） | ❌ 仅识别直接导入 | ✅ 递归解析所有传递性依赖 |
| 导入 proto 的 `.pb.go` 生成 | ❌ 需手动为每个文件单独运行 protoc | ✅ 自动将导入的 proto 追加到 protoc 命令 |
| Proto 搜索路径 | ❌ 仅源文件所在目录 | ✅ 支持多个 `-I` 路径，与 protoc 一致 |

**行为说明：**
- 递归遍历 proto 文件中的所有 `import` 声明，跳过 `google/*` 知名类型。
- 在每个 `-I` 目录中搜索被导入的文件，未找到的系统级 proto 静默跳过。
- 将发现的 proto 文件追加到 `protoc` 命令，使其 `.pb.go` 文件与主 proto 一同生成。

---

#### 跨包类型解析

当导入的 proto 与主 proto 具有**不同**的 `go_package` 时，goctl 现在能够自动在 server、logic 和 client 代码中生成正确的 Go 导入路径和限定类型引用。

**涉及文件：**
- `generator/typeref.go` — 新文件，核心类型解析引擎：
  - `resolveRPCTypeRef()` — 将 proto RPC 类型（简单类型、同包点号类型、跨包点号类型、Google WKT）解析为带正确导入路径的 Go 类型引用。
  - `resolveCallTypeRef()` — 客户端代码生成变体，支持类型别名。
  - `googleWKTTable` — 全部 16 种 Google 知名类型到 Go 等价类型的映射表。
- `generator/genserver.go` — `genFunctions()` 调用 `resolveRPCTypeRef()` 解析请求/响应类型并收集额外导入路径。
- `generator/genlogic.go` — `genLogicFunction()` 使用 `resolveRPCTypeRef()`；新增 `addLogicImports()` 按需添加主 pb 导入和跨包导入。
- `generator/gencall.go` — `genFunction()` 和 `getInterfaceFuncs()` 使用 `resolveCallTypeRef()` 处理类型别名和额外导入；新增 `buildExtraImportLines()` 辅助函数。
- `generator/call.tpl` — 新增 `{{.extraImports}}` 占位符用于跨包导入行。

**前后对比：**

| Proto 类型 | 变更前 | 变更后 |
|---|---|---|
| `GetReq`（同文件） | `pb.GetReq` | `pb.GetReq`（无变化） |
| `ext.ExtReq`（相同 `go_package`） | ❌ 报错："request type must defined in" | ✅ `pb.ExtReq` — 合并到主包 |
| `common.TypesReq`（不同 `go_package`） | ❌ 报错："request type must defined in" | ✅ `common.TypesReq` + 自动生成 `import "example.com/demo/pb/common"` |
| `google.protobuf.Empty` | ❌ 报错："request type must defined in" | ✅ `emptypb.Empty` + 自动生成导入 |

**行为说明：**
- 简单类型（如 `GetReq`）解析为 `pb.GetReq`，无额外导入。
- 同包点号类型（如 `ext.ExtReq`，其中 `ext` 与主 proto 有相同的 `go_package`）解析为 `pb.ExtReq`。
- 跨包点号类型（如 `common.TypesReq`，其中 `common` 有不同的 `go_package`）解析为 `common.TypesReq`，并自动添加正确的 Go 导入路径。

---

#### Google 知名类型作为 RPC 参数

Google protobuf 知名类型现在可以直接用作 RPC 的请求/响应类型（而不仅仅是消息字段）。

**涉及文件：**
- `generator/typeref.go` — `resolveGoogleWKT()` + `googleWKTTable` 处理所有标准类型。

**前后对比：**

| Proto 类型 | 变更前（作为 RPC 参数） | 变更后（作为 RPC 参数） |
|---|---|---|
| `google.protobuf.Empty` | ❌ 报错 | ✅ `emptypb.Empty` |
| `google.protobuf.Timestamp` | ❌ 报错 | ✅ `timestamppb.Timestamp` |
| `google.protobuf.Duration` | ❌ 报错 | ✅ `durationpb.Duration` |
| `google.protobuf.Any` | ❌ 报错 | ✅ `anypb.Any` |
| `google.protobuf.Struct` | ❌ 报错 | ✅ `structpb.Struct` |
| `google.protobuf.FieldMask` | ❌ 报错 | ✅ `fieldmaskpb.FieldMask` |
| `google.protobuf.*Value` | ❌ 报错 | ✅ `wrapperspb.*Value` |

> 注：这些类型此前已可用作**消息字段**。本次变更使其可直接用作 **RPC 请求/响应类型**。

**完整类型映射表：**

| Proto 类型 | Go 类型 |
|---|---|
| `google.protobuf.Empty` | `emptypb.Empty` |
| `google.protobuf.Timestamp` | `timestamppb.Timestamp` |
| `google.protobuf.Duration` | `durationpb.Duration` |
| `google.protobuf.Any` | `anypb.Any` |
| `google.protobuf.Struct` | `structpb.Struct` |
| `google.protobuf.Value` | `structpb.Value` |
| `google.protobuf.ListValue` | `structpb.ListValue` |
| `google.protobuf.FieldMask` | `fieldmaskpb.FieldMask` |
| `google.protobuf.*Value`（包装类型） | `wrapperspb.*Value` |

---

### 不兼容变更

#### RPC 定义中允许使用点号类型名

此前 goctl 会拒绝 RPC 请求/响应类型中包含点号的情况（如 `base.Req`），要求所有类型必须定义在同一个 proto 文件中。此限制已移除。

**前后对比：**

| Proto 定义 | 变更前 | 变更后 |
|---|---|---|
| `rpc Fetch(base.Req) returns (base.Reply)` | ❌ 解析错误："request type must defined in xxx.proto" | ✅ 解析成功，`base.Req` 通过导入的 proto 解析 |
| `rpc Ping(google.protobuf.Empty) returns (Reply)` | ❌ 解析错误："request type must defined in xxx.proto" | ✅ 解析成功，解析为 `emptypb.Empty` |

**涉及文件：**
- `parser/service.go` — 移除了拒绝点号类型名的验证循环（原错误信息为 `"request type must defined in"` / `"returns type must defined in"`）。
- `parser/parser_test.go` — `TestDefaultProtoParseCaseInvalidRequestType` 和 `TestDefaultProtoParseCaseInvalidResponseType` 重命名并更新，验证点号类型现在可以正常解析。
