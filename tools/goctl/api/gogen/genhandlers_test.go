package gogen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/tools/goctl/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/config"
)

func TestGenHandlers_WithFilename(t *testing.T) {
	// Test: Multiple routes with same filename should merge into one file
	dir := t.TempDir()
	apiFile := filepath.Join(dir, "test.api")

	apiContent := `
syntax = "v1"

type (
	ProductReq {
		Id int64 ` + "`path:\"id\"`" + `
	}
	ProductResp {
		Id   int64  ` + "`json:\"id\"`" + `
		Name string ` + "`json:\"name\"`" + `
	}
	CommonResp {
		Code int ` + "`json:\"code\"`" + `
	}
)

@server (
	group:    product
	prefix:   /api/v1/product
	filename: product
)
service demo {
	@handler getProduct
	get /:id (ProductReq) returns (ProductResp)
	
	@handler createProduct
	post / (ProductReq) returns (ProductResp)
	
	@handler updateProduct
	put /:id (ProductReq) returns (ProductResp)
	
	@handler deleteProduct
	delete /:id (ProductReq) returns (CommonResp)
}
`

	err := os.WriteFile(apiFile, []byte(apiContent), 0o644)
	assert.NoError(t, err)

	api, err := parser.Parse(apiFile)
	assert.NoError(t, err)

	cfg := &config.Config{
		NamingFormat: "gozero",
	}

	err = genHandlers(dir, "demo", "demo", cfg, api)
	assert.NoError(t, err)

	// Verify: Should generate only ONE handler file
	handlerFile := filepath.Join(dir, handlerDir, "product", "product.go")
	assert.FileExists(t, handlerFile)

	content, err := os.ReadFile(handlerFile)
	assert.NoError(t, err)
	contentStr := string(content)

	// Verify: All 4 handlers should be in the same file
	assert.Contains(t, contentStr, "func GetProductHandler")
	assert.Contains(t, contentStr, "func CreateProductHandler")
	assert.Contains(t, contentStr, "func UpdateProductHandler")
	assert.Contains(t, contentStr, "func DeleteProductHandler")

	// Verify: Should have package product
	assert.Contains(t, contentStr, "package product")
}

func TestGenHandlers_MixedScenarios(t *testing.T) {
	// Test: Some routes with filename, some without
	dir := t.TempDir()
	apiFile := filepath.Join(dir, "test.api")

	apiContent := `
syntax = "v1"

type (
	UserReq {
		Id int64 ` + "`path:\"id\"`" + `
	}
	UserResp {
		Id   int64  ` + "`json:\"id\"`" + `
		Name string ` + "`json:\"name\"`" + `
	}
)

@server (
	group:    user
	prefix:   /api/v1/user
	filename: user_query
)
service demo {
	@handler getUserInfo
	get /:id (UserReq) returns (UserResp)
	
	@handler getUserList
	get /list returns (UserResp)
}

@server (
	group:  user
	prefix: /api/v1/user
)
service demo {
	@handler updateUser
	put /:id (UserReq) returns (UserResp)
	
	@handler deleteUser
	delete /:id (UserReq) returns (UserResp)
}
`

	err := os.WriteFile(apiFile, []byte(apiContent), 0o644)
	assert.NoError(t, err)

	api, err := parser.Parse(apiFile)
	assert.NoError(t, err)

	cfg := &config.Config{
		NamingFormat: "gozero",
	}

	err = genHandlers(dir, "demo", "demo", cfg, api)
	assert.NoError(t, err)

	// Verify: query operations in one file (with filename)
	queryFile := filepath.Join(dir, handlerDir, "user", "userquery.go")
	assert.FileExists(t, queryFile)
	queryContent, _ := os.ReadFile(queryFile)
	assert.Contains(t, string(queryContent), "func GetUserInfoHandler")
	assert.Contains(t, string(queryContent), "func GetUserListHandler")

	// Verify: command operations in separate files (without filename)
	updateFile := filepath.Join(dir, handlerDir, "user", "updateuserhandler.go")
	assert.FileExists(t, updateFile)
	updateContent, _ := os.ReadFile(updateFile)
	assert.Contains(t, string(updateContent), "func UpdateUserHandler")

	deleteFile := filepath.Join(dir, handlerDir, "user", "deleteuserhandler.go")
	assert.FileExists(t, deleteFile)
	deleteContent, _ := os.ReadFile(deleteFile)
	assert.Contains(t, string(deleteContent), "func DeleteUserHandler")
}

func TestGenHandlers_ImportDeduplication(t *testing.T) {
	// Test: Verify duplicate imports are properly handled
	dir := t.TempDir()
	apiFile := filepath.Join(dir, "test.api")

	apiContent := `
syntax = "v1"

type (
	OrderReq {
		Id int64 ` + "`path:\"id\"`" + `
	}
	OrderResp {
		Id     int64  ` + "`json:\"id\"`" + `
		Status string ` + "`json:\"status\"`" + `
	}
)

@server (
	group:    order
	prefix:   /api/v1/order
	filename: order_flow
)
service demo {
	@handler createOrder
	post / (OrderReq) returns (OrderResp)
	
	@handler payOrder
	post /:id/pay (OrderReq) returns (OrderResp)
	
	@handler confirmOrder
	post /:id/confirm (OrderReq) returns (OrderResp)
}
`

	err := os.WriteFile(apiFile, []byte(apiContent), 0o644)
	assert.NoError(t, err)

	api, err := parser.Parse(apiFile)
	assert.NoError(t, err)

	cfg := &config.Config{
		NamingFormat: "gozero",
	}

	err = genHandlers(dir, "demo", "demo", cfg, api)
	assert.NoError(t, err)

	// Verify: Check imports are not duplicated
	handlerFile := filepath.Join(dir, handlerDir, "order", "orderflow.go")
	content, err := os.ReadFile(handlerFile)
	assert.NoError(t, err)
	contentStr := string(content)

	// Count occurrences of import statements
	svcImportCount := strings.Count(contentStr, `"demo/internal/svc"`)
	logicImportCount := strings.Count(contentStr, `"demo/internal/logic/order"`)
	typesImportCount := strings.Count(contentStr, `"demo/internal/types"`)

	// Each import should appear exactly once in the import block
	assert.Equal(t, 1, svcImportCount, "svc import should appear exactly once")
	assert.Equal(t, 1, logicImportCount, "logic import should appear exactly once")
	assert.Equal(t, 1, typesImportCount, "types import should appear exactly once")
}

func TestGenHandlers_DifferentGroupsSameFilename(t *testing.T) {
	// Test: Routes in different groups with same filename should be in separate files
	dir := t.TempDir()
	apiFile := filepath.Join(dir, "test.api")

	apiContent := `
syntax = "v1"

type (
	AdminReq {
		Id int64 ` + "`path:\"id\"`" + `
	}
	AdminResp {
		Id int64 ` + "`json:\"id\"`" + `
	}
	UserReq {
		Id int64 ` + "`path:\"id\"`" + `
	}
	UserResp {
		Id int64 ` + "`json:\"id\"`" + `
	}
)

@server (
	group:    admin
	prefix:   /api/v1/admin
	filename: common
)
service demo {
	@handler getAdminInfo
	get /:id (AdminReq) returns (AdminResp)
}

@server (
	group:    user
	prefix:   /api/v1/user
	filename: common
)
service demo {
	@handler getUserInfo
	get /:id (UserReq) returns (UserResp)
}
`

	err := os.WriteFile(apiFile, []byte(apiContent), 0o644)
	assert.NoError(t, err)

	api, err := parser.Parse(apiFile)
	assert.NoError(t, err)

	cfg := &config.Config{
		NamingFormat: "gozero",
	}

	err = genHandlers(dir, "demo", "demo", cfg, api)
	assert.NoError(t, err)

	// Verify: Two separate files in different directories
	adminFile := filepath.Join(dir, handlerDir, "admin", "common.go")
	userFile := filepath.Join(dir, handlerDir, "user", "common.go")

	assert.FileExists(t, adminFile)
	assert.FileExists(t, userFile)

	adminContent, _ := os.ReadFile(adminFile)
	userContent, _ := os.ReadFile(userFile)

	assert.Contains(t, string(adminContent), "func GetAdminInfoHandler")
	assert.NotContains(t, string(adminContent), "func GetUserInfoHandler")

	assert.Contains(t, string(userContent), "func GetUserInfoHandler")
	assert.NotContains(t, string(userContent), "func GetAdminInfoHandler")
}

func TestGenHandlers_SSEWithFilename(t *testing.T) {
	// Test: SSE routes with filename annotation
	dir := t.TempDir()
	apiFile := filepath.Join(dir, "test.api")

	apiContent := `
syntax = "v1"

type (
	StreamResp {
		Message string ` + "`json:\"message\"`" + `
	}
)

@server (
	group:    stream
	prefix:   /api/v1/stream
	filename: stream_handlers
	sse:      true
)
service demo {
	@handler streamEvents
	get /events returns (StreamResp)
	
	@handler streamLogs
	get /logs returns (StreamResp)
}
`

	err := os.WriteFile(apiFile, []byte(apiContent), 0o644)
	assert.NoError(t, err)

	api, err := parser.Parse(apiFile)
	assert.NoError(t, err)

	cfg := &config.Config{
		NamingFormat: "gozero",
	}

	err = genHandlers(dir, "demo", "demo", cfg, api)
	assert.NoError(t, err)

	// Verify: SSE handlers merged into one file
	handlerFile := filepath.Join(dir, handlerDir, "stream", "streamhandlers.go")
	assert.FileExists(t, handlerFile)

	content, err := os.ReadFile(handlerFile)
	assert.NoError(t, err)
	contentStr := string(content)

	// Verify: Both SSE handlers in the same file
	assert.Contains(t, contentStr, "func StreamEventsHandler")
	assert.Contains(t, contentStr, "func StreamLogsHandler")

	// Verify: SSE-specific code patterns
	assert.Contains(t, contentStr, "text/event-stream")
}

func TestGenHandlers_BackwardCompatibility(t *testing.T) {
	// Test: Without filename annotation, behavior should remain unchanged
	dir := t.TempDir()
	apiFile := filepath.Join(dir, "test.api")

	apiContent := `
syntax = "v1"

type (
	GreetReq {
		Name string ` + "`path:\"name\"`" + `
	}
	GreetResp {
		Message string ` + "`json:\"message\"`" + `
	}
)

@server (
	group:  greet
	prefix: /api/v1/greet
)
service demo {
	@handler greetHandler
	get /:name (GreetReq) returns (GreetResp)
}
`

	err := os.WriteFile(apiFile, []byte(apiContent), 0o644)
	assert.NoError(t, err)

	api, err := parser.Parse(apiFile)
	assert.NoError(t, err)

	cfg := &config.Config{
		NamingFormat: "gozero",
	}

	err = genHandlers(dir, "demo", "demo", cfg, api)
	assert.NoError(t, err)

	// Verify: Traditional naming (handler name as filename)
	handlerFile := filepath.Join(dir, handlerDir, "greet", "greethandler.go")
	assert.FileExists(t, handlerFile)

	content, err := os.ReadFile(handlerFile)
	assert.NoError(t, err)

	// Handler names in subdirectories are PascalCase
	// greet is a group, so it creates a subdirectory and uses PascalCase
	assert.Contains(t, string(content), "func GreetHandler")
}

func TestGenHandlers_RouteAnnotationPriority(t *testing.T) {
	// Note: Route-level @server annotation is NOT supported by go-zero API parser
	// This test is skipped as the feature requires parser-level changes
	t.Skip("Route-level @server annotation is not supported by API parser yet")
}
