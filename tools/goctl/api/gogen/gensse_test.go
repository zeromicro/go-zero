package gogen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSSEGeneration(t *testing.T) {
	// Create a temporary directory for test
	dir := t.TempDir()

	// Create a test API file with SSE annotation
	apiContent := `syntax = "v1"

type SseReq {
    Message string ` + "`json:\"message\"`" + `
}

type SseResp {
    Data string ` + "`json:\"data\"`" + `
}

@server (
    sse: true
)
service Test {
    @handler Sse
    get /sse (SseReq) returns (SseResp)
}
`
	apiFile := filepath.Join(dir, "test.api")
	err := os.WriteFile(apiFile, []byte(apiContent), 0644)
	assert.NoError(t, err)

	// Generate code
	err = DoGenProject(apiFile, dir, "gozero", false)
	assert.NoError(t, err)

	// Read generated handler file
	handlerPath := filepath.Join(dir, "internal/handler/ssehandler.go")
	handlerContent, err := os.ReadFile(handlerPath)
	assert.NoError(t, err)

	// Read generated logic file
	logicPath := filepath.Join(dir, "internal/logic/sselogic.go")
	logicContent, err := os.ReadFile(logicPath)
	assert.NoError(t, err)

	handlerStr := string(handlerContent)
	logicStr := string(logicContent)

	// Verify SSE-specific patterns in handler
	// Handler should call: err := l.Sse(&req, client)
	assert.Contains(t, handlerStr, "err := l.Sse(&req, client)",
		"Handler should call logic with client channel parameter")

	// Handler should NOT have the regular pattern: resp, err := l.Sse(&req)
	assert.NotContains(t, handlerStr, "resp, err := l.Sse(&req)",
		"Handler should not use regular pattern with resp return")

	// Handler should use threading.GoSafeCtx
	assert.Contains(t, handlerStr, "threading.GoSafeCtx",
		"Handler should use threading.GoSafeCtx for SSE")

	// Handler should create client channel
	assert.Contains(t, handlerStr, "client := make(chan",
		"Handler should create client channel")

	// Verify SSE-specific patterns in logic
	// Logic should have signature: Sse(req *types.SseReq, client chan<- *types.SseResp) error
	assert.Contains(t, logicStr, "func (l *SseLogic) Sse(req *types.SseReq, client chan<- *types.SseResp) error",
		"Logic should have SSE signature with client channel parameter")

	// Logic should NOT have regular signature: Sse(req *types.SseReq) (resp *types.SseResp, err error)
	assert.NotContains(t, logicStr, "(resp *types.SseResp, err error)",
		"Logic should not have regular signature with resp return")
}

func TestNonSSEGeneration(t *testing.T) {
	// Create a temporary directory for test
	dir := t.TempDir()

	// Create a test API file WITHOUT SSE annotation
	apiContent := `syntax = "v1"

type SseReq {
    Message string ` + "`json:\"message\"`" + `
}

type SseResp {
    Data string ` + "`json:\"data\"`" + `
}

service Test {
    @handler Sse
    get /sse (SseReq) returns (SseResp)
}
`
	apiFile := filepath.Join(dir, "test.api")
	err := os.WriteFile(apiFile, []byte(apiContent), 0644)
	assert.NoError(t, err)

	// Generate code
	err = DoGenProject(apiFile, dir, "gozero", false)
	assert.NoError(t, err)

	// Read generated handler file
	handlerPath := filepath.Join(dir, "internal/handler/ssehandler.go")
	handlerContent, err := os.ReadFile(handlerPath)
	assert.NoError(t, err)

	// Read generated logic file
	logicPath := filepath.Join(dir, "internal/logic/sselogic.go")
	logicContent, err := os.ReadFile(logicPath)
	assert.NoError(t, err)

	handlerStr := string(handlerContent)
	logicStr := string(logicContent)

	// Verify regular (non-SSE) patterns in handler
	// Handler should call: resp, err := l.Sse(&req)
	assert.Contains(t, handlerStr, "resp, err := l.Sse(&req)",
		"Handler should use regular pattern with resp return")

	// Handler should NOT have SSE pattern: err := l.Sse(&req, client)
	assert.NotContains(t, handlerStr, "err := l.Sse(&req, client)",
		"Handler should not use SSE pattern")

	// Handler should NOT use threading.GoSafeCtx
	assert.NotContains(t, handlerStr, "threading.GoSafeCtx",
		"Handler should not use threading.GoSafeCtx for regular routes")

	// Verify regular (non-SSE) patterns in logic
	// Logic should have signature: Sse(req *types.SseReq) (resp *types.SseResp, err error)
	assert.Contains(t, logicStr, "(resp *types.SseResp, err error)",
		"Logic should have regular signature with resp return")

	// Logic should NOT have SSE signature with client parameter
	linesToCheck := strings.Split(logicStr, "\n")
	hasSSESignature := false
	for _, line := range linesToCheck {
		if strings.Contains(line, "func (l *SseLogic) Sse") && strings.Contains(line, "client chan<-") {
			hasSSESignature = true
			break
		}
	}
	assert.False(t, hasSSESignature,
		"Logic should not have SSE signature with client channel parameter")
}
