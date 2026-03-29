# fix(api): Swagger 生成时验证路径参数与路由占位符的匹配性 (Closes #5428)

## 问题描述

当请求结构体中声明了 `path:"id"` 标记的字段，但路由路径中不包含对应的占位符（如 `/:id` 或 `{id}`）时，goctl 目前仍会生成 `in: path` 的参数，导致 Swagger 文档语义不一致（OpenAPI 规范要求路径参数必须出现在 URL 模板中）。

## 问题示例

```go
type Req {
  ID string `path:"id"`  // 声明了路径参数
}

@handler H
get /foo (Req) returns ()  // 路径 /foo 中没有 {/:id} 占位符
```

生成的 Swagger 结果（错误）：
```json
{
  "type": "string",
  "name": "id",
  "in": "path",    // 参数位置为 path
  "required": true
}
```

但实际路径是 `/foo`，不包含 `{id}`，这是无效的 Swagger 文档。

## 解决方案

在生成路径参数前，解析路由路径中的占位符（支持 `:id` 和 `{id}` 两种格式），仅当路径中存在匹配的占位符时才生成对应的路径参数。对于未匹配的路径参数，静默跳过以避免生成无效的 Swagger 文档。

## 核心改动

### 1. 新增 `extractPathPlaceholders()` 函数 (`tools/goctl/api/swagger/vars.go`)

- 支持解析 `:id` 风格的 go-zero 占位符
- 支持解析 `{id}` 风格的 OpenAPI 占位符
- 返回路由路径中所有有效的路径参数名称集合

### 2. 修改 `parametersFromType()` 函数 (`tools/goctl/api/swagger/parameter.go`)

- 新增 `routePath string` 参数，接收原始路由路径
- 在生成路径参数前调用 `extractPathPlaceholders()` 获取允许的占位符集
- 仅当 `pathParameterTag.Name` 在占位符集合中时才生成参数
- 未匹配的路径参数被静默过滤

### 3. 更新 `spec2Path()` 函数 (`tools/goctl/api/swagger/path.go`)

- 调用 `parametersFromType()` 时传入 `route.Path`

### 4. 新增单元测试

- **`vars_test.go`**: 测试占位符提取逻辑（14 个测试用例）
  - 覆盖空路径、无占位符
  - 单/多占位符
  - 混合格式 (`:id` 和 `{id}`)
  - 边缘情况（空占位符名称、尾随斜杠等）

- **`parameter_path_filter_test.go`**: 测试路径参数验证逻辑（6 个测试用例）
  - 覆盖匹配/不匹配场景
  - 多个路径参数部分匹配的场景
  - 查询参数不受影响的验证

- **修复 `parameter_test.go`**: 兼容新的参数签名

## 测试结果

```text
=== RUN   TestExtractPathPlaceholders
    --- PASS: empty_path
    --- PASS: no_placeholders
    --- PASS: single_:id_style_placeholder
    --- PASS: single_{id}_style_placeholder
    --- PASS: multiple_:id_style_placeholders
    --- PASS: multiple_{id}_style_placeholders
    --- PASS: mixed_style_placeholders
    --- PASS: placeholder_at_root
    --- PASS: placeholder_after_static_segments
    --- PASS: empty_placeholder_name_with_colon
    --- PASS: empty_placeholder_name_with_braces
    --- PASS: trailing_slash_with_placeholder
    --- PASS: complex_path_with_multiple_placeholders
--- PASS

=== RUN   TestExtractPathPlaceholders_Duplicates
--- PASS

=== RUN   TestExtractPathPlaceholders_CaseSensitive
--- PASS

=== RUN   TestParametersFromType_PathParameterValidation
    --- PASS: path parameter matches route placeholder
    --- PASS: path parameter matches {id} style placeholder
    --- PASS: path parameter does NOT match route placeholder - should be filtered out
    --- PASS: multiple path parameters, only some match route placeholders
    --- PASS: no path tag - no path params should be generated
    --- PASS: mixed route path - colon and brace styles
--- PASS

=== RUN   TestParametersFromType_PathParameterFilteringWithQueryString
--- PASS

PASS
ok   github.com/zeromicro/go-zero/tools/goctl/api/swagger 0.003s
```

## 影响范围

- Swagger 生成功能
- 路径参数验证逻辑
- 不影响其他参数类型（query、form、header、body）的生成

## 改动文件

```
tools/goctl/api/swagger/parameter.go               # 修改：增加路径过滤逻辑
tools/goctl/api/swagger/paramter_path_filter_test.go  # 新增：路径参数验证测试
tools/goctl/api/swagger/parameter_test.go          # 修改：兼容新参数签名
tools/goctl/api/swagger/path.go                   # 修改：传递路由路径
tools/goctl/api/swagger/vars.go                    # 新增：extractPathPlaceholders 函数
tools/goctl/api/swagger/vars_test.go              # 新增：占位符提取测试
```

## 分支与提交

**分支**: `fix/swagger-path-only-5428` (clean standalone branch)  
**提交**: `8409aa41`  
**PR 关闭**: #5428