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

func TestGenLogic_WithFilename(t *testing.T) {
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

	err = genLogic(dir, "demo", "demo", cfg, api)
	assert.NoError(t, err)

	// Verify: Should generate only ONE logic file
	logicFile := filepath.Join(dir, logicDir, "product", "product.go")
	assert.FileExists(t, logicFile)

	content, err := os.ReadFile(logicFile)
	assert.NoError(t, err)
	contentStr := string(content)

	// Verify: All 4 logic functions should be in the same file
	assert.Contains(t, contentStr, "type GetProductLogic struct")
	assert.Contains(t, contentStr, "func (l *GetProductLogic) GetProduct")
	assert.Contains(t, contentStr, "type CreateProductLogic struct")
	assert.Contains(t, contentStr, "func (l *CreateProductLogic) CreateProduct")
	assert.Contains(t, contentStr, "type UpdateProductLogic struct")
	assert.Contains(t, contentStr, "func (l *UpdateProductLogic) UpdateProduct")
	assert.Contains(t, contentStr, "type DeleteProductLogic struct")
	assert.Contains(t, contentStr, "func (l *DeleteProductLogic) DeleteProduct")

	// Verify: Should have package product
	assert.Contains(t, contentStr, "package product")
}

func TestGenLogic_MixedScenarios(t *testing.T) {
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

	err = genLogic(dir, "demo", "demo", cfg, api)
	assert.NoError(t, err)

	// Verify: query operations in one file (with filename)
	queryFile := filepath.Join(dir, logicDir, "user", "userquery.go")
	assert.FileExists(t, queryFile)
	queryContent, _ := os.ReadFile(queryFile)
	assert.Contains(t, string(queryContent), "type GetUserInfoLogic struct")
	assert.Contains(t, string(queryContent), "type GetUserListLogic struct")

	// Verify: command operations in separate files (without filename)
	updateFile := filepath.Join(dir, logicDir, "user", "updateuserlogic.go")
	assert.FileExists(t, updateFile)
	updateContent, _ := os.ReadFile(updateFile)
	assert.Contains(t, string(updateContent), "type UpdateUserLogic struct")

	deleteFile := filepath.Join(dir, logicDir, "user", "deleteuserlogic.go")
	assert.FileExists(t, deleteFile)
	deleteContent, _ := os.ReadFile(deleteFile)
	assert.Contains(t, string(deleteContent), "type DeleteUserLogic struct")
}

func TestGenLogic_ImportDeduplication(t *testing.T) {
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

	err = genLogic(dir, "demo", "demo", cfg, api)
	assert.NoError(t, err)

	// Verify: Check imports are not duplicated
	logicFile := filepath.Join(dir, logicDir, "order", "orderflow.go")
	content, err := os.ReadFile(logicFile)
	assert.NoError(t, err)
	contentStr := string(content)

	// Count occurrences of import statements
	contextImportCount := strings.Count(contentStr, `"context"`)
	svcImportCount := strings.Count(contentStr, `"demo/internal/svc"`)
	typesImportCount := strings.Count(contentStr, `"demo/internal/types"`)
	logxImportCount := strings.Count(contentStr, `"github.com/zeromicro/go-zero/core/logx"`)

	// Each import should appear exactly once in the import block
	assert.Equal(t, 1, contextImportCount, "context import should appear exactly once")
	assert.Equal(t, 1, svcImportCount, "svc import should appear exactly once")
	assert.Equal(t, 1, typesImportCount, "types import should appear exactly once")
	assert.Equal(t, 1, logxImportCount, "logx import should appear exactly once")
}

func TestGenLogic_DifferentGroupsSameFilename(t *testing.T) {
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

	err = genLogic(dir, "demo", "demo", cfg, api)
	assert.NoError(t, err)

	// Verify: Two separate files in different directories
	adminFile := filepath.Join(dir, logicDir, "admin", "common.go")
	userFile := filepath.Join(dir, logicDir, "user", "common.go")

	assert.FileExists(t, adminFile)
	assert.FileExists(t, userFile)

	adminContent, _ := os.ReadFile(adminFile)
	userContent, _ := os.ReadFile(userFile)

	assert.Contains(t, string(adminContent), "type GetAdminInfoLogic struct")
	assert.NotContains(t, string(adminContent), "type GetUserInfoLogic struct")

	assert.Contains(t, string(userContent), "type GetUserInfoLogic struct")
	assert.NotContains(t, string(userContent), "type GetAdminInfoLogic struct")
}

func TestGenLogic_SSEWithFilename(t *testing.T) {
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
	filename: stream_logics
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

	err = genLogic(dir, "demo", "demo", cfg, api)
	assert.NoError(t, err)

	// Verify: SSE logics merged into one file
	logicFile := filepath.Join(dir, logicDir, "stream", "streamlogics.go")
	assert.FileExists(t, logicFile)

	content, err := os.ReadFile(logicFile)
	assert.NoError(t, err)
	contentStr := string(content)

	// Verify: Both SSE logics in the same file
	assert.Contains(t, contentStr, "type StreamEventsLogic struct")
	assert.Contains(t, contentStr, "type StreamLogsLogic struct")

	// Verify: SSE-specific parameters (client channel with pointer)
	assert.Contains(t, contentStr, "client chan<- *types.StreamResp")
}

func TestGenLogic_BackwardCompatibility(t *testing.T) {
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

	err = genLogic(dir, "demo", "demo", cfg, api)
	assert.NoError(t, err)

	// Verify: Traditional naming (logic name as filename)
	// getLogicName returns "greetHandlerLogic", format converts it to "greethandlerlogic.go"
	logicFile := filepath.Join(dir, logicDir, "greet", "greetlogic.go")
	assert.FileExists(t, logicFile)

	content, err := os.ReadFile(logicFile)
	assert.NoError(t, err)

	assert.Contains(t, string(content), "type GreetLogic struct")
}

func TestGenLogic_NoRequestType(t *testing.T) {
	// Test: Routes without request type
	dir := t.TempDir()
	apiFile := filepath.Join(dir, "test.api")

	apiContent := `
syntax = "v1"

type (
	StatusResp {
		Status string ` + "`json:\"status\"`" + `
	}
)

@server (
	group:    health
	prefix:   /api/v1
	filename: health
)
service demo {
	@handler healthCheck
	get /health returns (StatusResp)
	
	@handler readyCheck
	get /ready returns (StatusResp)
}
`

	err := os.WriteFile(apiFile, []byte(apiContent), 0o644)
	assert.NoError(t, err)

	api, err := parser.Parse(apiFile)
	assert.NoError(t, err)

	cfg := &config.Config{
		NamingFormat: "gozero",
	}

	err = genLogic(dir, "demo", "demo", cfg, api)
	assert.NoError(t, err)

	// Verify: Logic file generated with no request parameters
	logicFile := filepath.Join(dir, logicDir, "health", "health.go")
	assert.FileExists(t, logicFile)

	content, err := os.ReadFile(logicFile)
	assert.NoError(t, err)
	contentStr := string(content)

	// Verify: Function signatures without request parameters
	assert.Contains(t, contentStr, "func (l *HealthCheckLogic) HealthCheck()")
	assert.Contains(t, contentStr, "func (l *ReadyCheckLogic) ReadyCheck()")
}

func TestGenLogic_NoResponseType(t *testing.T) {
	// Test: Routes without response type
	dir := t.TempDir()
	apiFile := filepath.Join(dir, "test.api")

	apiContent := `
syntax = "v1"

type (
	ActionReq {
		Id int64 ` + "`path:\"id\"`" + `
	}
)

@server (
	group:    action
	prefix:   /api/v1/action
	filename: actions
)
service demo {
	@handler triggerAction
	post /:id (ActionReq)
	
	@handler cancelAction
	delete /:id (ActionReq)
}
`

	err := os.WriteFile(apiFile, []byte(apiContent), 0o644)
	assert.NoError(t, err)

	api, err := parser.Parse(apiFile)
	assert.NoError(t, err)

	cfg := &config.Config{
		NamingFormat: "gozero",
	}

	err = genLogic(dir, "demo", "demo", cfg, api)
	assert.NoError(t, err)

	// Verify: Logic file generated with error return only
	logicFile := filepath.Join(dir, logicDir, "action", "actions.go")
	assert.FileExists(t, logicFile)

	content, err := os.ReadFile(logicFile)
	assert.NoError(t, err)
	contentStr := string(content)

	// Verify: Function signatures return error only
	assert.Contains(t, contentStr, "func (l *TriggerActionLogic) TriggerAction(req *types.ActionReq) error")
	assert.Contains(t, contentStr, "func (l *CancelActionLogic) CancelAction(req *types.ActionReq) error")
}
