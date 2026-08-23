# feat(api/swagger): 在 `.api` 文件中支持自定义 `x-` 扩展字段并生成 Swagger

关联 issue：go-zero/go-zero#5628

## 摘要
本改动让 API 作者可以直接在 `.api` 文件里声明自定义 Swagger 扩展字段（`x-*`），语法是在路由的 `@handler` 之前添加 `@x-<key>` 注解：

```api
service order {
    @doc "Create Order"
    @x-log-enabled true
    @x-rate-limit "100/min"
    @x-owners `["alice","bob"]`
    @handler createOrder
    post /order/create (CreateOrderRequest) returns (CreateOrderResponse)
}
```

生成的 Swagger JSON 会在 operation 上带上这些扩展字段：

```json
"post": {
    "x-log-enabled": true,
    "x-rate-limit": "100/min",
    "x-owners": ["alice", "bob"]
}
```

## 改动内容
- **词法扫描**（`tools/goctl/pkg/parser/api/scanner/scanner.go`）：识别 `@x-<key>` 词法单元。
- **Token / AST**（`tools/goctl/pkg/parser/api/token/token.go`、`tools/goctl/pkg/parser/api/ast/servicestatement.go`）：新增 `AT_X` token 类型，以及 `AtXStmt` / `AtXAnnotations` AST 节点。
- **解析器**（`tools/goctl/pkg/parser/api/parser/parser.go`）：在 `@handler` 前解析零个或多个 `@x-*` 注解，并在 `ParseForUintTest` 中暴露该能力。
- **分析器 / spec**（`tools/goctl/pkg/parser/api/parser/analyzer.go`、`tools/goctl/api/spec/spec.go`）：将注解写入 `spec.Route.Extensions`，并移除一处重复的 `@doc` 转换。
- **Swagger 生成**（`tools/goctl/api/swagger/swagger.go`、`tools/goctl/api/swagger/path.go`）：将扩展值解析为布尔值、数字、字符串或 JSON 数组/对象，并写入 `spec.Operation.Extensions`。
- **测试**（`tools/goctl/pkg/parser/api/parser/parser_test.go`、`tools/goctl/api/swagger/swagger_test.go`）：覆盖独立 `@x-*` 解析、service 级解析以及端到端 Swagger 生成，包括字符串、布尔、数字、JSON 数组/对象等取值。

## 扩展字段值解析规则
按以下顺序解析值：
1. JSON 数组或对象（支持双引号字符串和反引号 raw string）。
2. 带引号或反引号的字符串字面量。
3. 布尔字面量（`true` / `false`）。
4. 数字字面量（如 `30`、`3.14`）。
5. 其它情况按普通字符串返回。

## 验证
```bash
cd tools/goctl
go test ./pkg/parser/api/... ./api/swagger/...
go build ./...
```

所有测试通过，构建成功。
