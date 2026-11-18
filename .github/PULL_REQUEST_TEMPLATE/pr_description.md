# 概述

本 PR 包含两个主要修复:

1. **Base64 换行符兼容性修复 (Closes #4710)**
   - 修复服务端无法正确解析带有换行符的 Base64 字符串的问题
   - 兼容 Android 客户端使用 `Base64.encodeToString()` 默认添加换行符的行为

2. **Swagger 路径参数验证修复 (Closes #5428)**
   - 修复 Swagger 生成时错误添加路径参数的问题
   - 只在路由路径包含匹配占位符时才添加路径参数

---

## 修复详情

### Base64 换行符兼容性修复

**问题描述:**
Android 的 `Base64.encodeToString()` 方法默认会添加换行符（MIME 风格），导致服务端在解码 `X-Content-Security` 请求头时失败。HTTP 请求头不能包含换行符或其他特殊字符。

**解决方案:**
1. 修改 `core/codec/rsa.go` 中的 `DecryptBase64()` 方法，在解码前移除 Base64 字符串中的换行符 (`\n`)、回车符 (`\r`) 和空白字符
2. 修改 `rest/internal/security/contentsecurity.go` 中的 `ParseContentSecurity()` 函数，对 `base64Key` 和 `signature` 字段进行同样的换行符清理

**兼容性:**
- 此修改向后兼容，不影响已发送干净 Base64 字符串的客户端
- 同时支持标准 Base64 和带换行符的 Base64 字符串

### Swagger 路径参数验证修复

**问题描述:**
Swagger 生成时，即使路由路径中没有对应的占位符，也会错误地添加路径参数，导致生成的 Swagger 文档无效。

**解决方案:**
1. 新增 `extractPathPlaceholders()` 函数，解析路由路径中的占位符
   - 支持 `:id` 风格（go-zero）
   - 支持 `{id}` 风格（OpenAPI）
2. 修改 `parametersFromType()` 函数，接收路由路径参数，仅在路径包含匹配占位符时添加路径参数

---

## 文件变更

| 文件 | 变更说明 |
|------|---------|
| `core/codec/rsa.go` | 添加 Base64 字符串换行符清理逻辑 |
| `core/codec/rsa_test.go` | 新增带换行符 Base64 解码测试 |
| `rest/internal/security/contentsecurity.go` | 对 key 和 signature 字段清理换行符 |
| `rest/internal/security/contentsecurity_test.go` | 新增多组换行符兼容性测试 |
| `tools/goctl/api/swagger/parameter.go` | 添加路径参数验证逻辑 |
| `tools/goctl/api/swagger/path.go` | 传递路由路径到参数生成函数 |
| `tools/goctl/api/swagger/vars.go` | 新增 `extractPathPlaceholders()` 函数 |
| `tools/goctl/api/swagger/vars_test.go` | 新增占位符提取测试 |
| `tools/goctl/api/swagger/parameter_path_filter_test.go` | 新增路径参数过滤测试 |
| `tools/goctl/api/swagger/parameter_test.go` | 更新测试用例 |

---

## 测试

所有相关测试已通过:
- `core/codec/...` ✅
- `rest/internal/security/...` ✅  
- `tools/goctl/api/swagger/...` ✅

---

Co-Authored-By: Oz <oz-agent@warp.dev>
