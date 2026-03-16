package ctx

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBackground(t *testing.T) {
	workDir := "."
	ctx, err := Prepare(workDir)
	assert.Nil(t, err)
	assert.True(t, true, func() bool {
		return len(ctx.Dir) != 0 && len(ctx.Path) != 0
	}())
}

func TestBackgroundNilWorkDir(t *testing.T) {
	workDir := ""
	_, err := Prepare(workDir)
	assert.NotNil(t, err)
}

func TestPrepareWithModule(t *testing.T) {
	tests := []struct {
		name       string
		moduleName string
		expectMod  string
	}{
		{
			name:       "custom module name",
			moduleName: "github.com/example/testmodule",
			expectMod:  "github.com/example/testmodule",
		},
		{
			name:       "simple module name",
			moduleName: "simplemodule",
			expectMod:  "simplemodule",
		},
		{
			name:       "empty module name uses directory",
			moduleName: "",
			expectMod:  "", // Will be set to directory name
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary directory for testing
			tempDir, err := os.MkdirTemp("", "goctl-ctx-test-*")
			require.NoError(t, err)
			defer os.RemoveAll(tempDir)

			testDir := filepath.Join(tempDir, "testproject")
			err = os.MkdirAll(testDir, 0755)
			require.NoError(t, err)

			ctx, err := PrepareWithModule(testDir, tt.moduleName)
			assert.NoError(t, err)
			assert.NotNil(t, ctx)

			// Check that the context has expected values
			assert.NotEmpty(t, ctx.WorkDir)
			assert.NotEmpty(t, ctx.Name)
			assert.NotEmpty(t, ctx.Path)
			assert.NotEmpty(t, ctx.Dir)

			// Check that go.mod was created
			goModPath := filepath.Join(testDir, "go.mod")
			assert.FileExists(t, goModPath)

			// Verify module name in go.mod
			content, err := os.ReadFile(goModPath)
			require.NoError(t, err)

			expectedModule := tt.expectMod
			if expectedModule == "" {
				expectedModule = "testproject" // directory name fallback
			}

			assert.Contains(t, string(content), "module "+expectedModule)
			assert.Equal(t, expectedModule, ctx.Path)
		})
	}
}

func TestPrepareWithModule_ExistingGoMod(t *testing.T) {
	// Create a temporary directory with existing go.mod
	tempDir, err := os.MkdirTemp("", "goctl-ctx-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	testDir := filepath.Join(tempDir, "existingproject")
	err = os.MkdirAll(testDir, 0755)
	require.NoError(t, err)

	// Create existing go.mod file
	existingGoMod := `module existing.com/project

go 1.21
`
	goModPath := filepath.Join(testDir, "go.mod")
	err = os.WriteFile(goModPath, []byte(existingGoMod), 0644)
	require.NoError(t, err)

	// PrepareWithModule should use existing go.mod, not create new one
	ctx, err := PrepareWithModule(testDir, "github.com/new/module")
	assert.NoError(t, err)
	assert.NotNil(t, ctx)

	// Should use existing module name, not the provided one
	assert.Equal(t, "existing.com/project", ctx.Path)

	// Verify go.mod still contains original content
	content, err := os.ReadFile(goModPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "module existing.com/project")
	assert.NotContains(t, string(content), "module github.com/new/module")
}

func TestPrepareWithModule_InvalidWorkDir(t *testing.T) {
	_, err := PrepareWithModule("/non/existent/path", "github.com/example/test")
	assert.Error(t, err)
}

func TestPrepare_CallsPrepareWithModule(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "goctl-ctx-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	testDir := filepath.Join(tempDir, "testproject")
	err = os.MkdirAll(testDir, 0755)
	require.NoError(t, err)

	// Test that Prepare calls PrepareWithModule with empty string
	ctx1, err1 := Prepare(testDir)
	require.NoError(t, err1)

	// Clean up go.mod to test again
	os.Remove(filepath.Join(testDir, "go.mod"))

	ctx2, err2 := PrepareWithModule(testDir, "")
	require.NoError(t, err2)

	// Should produce identical results
	assert.Equal(t, ctx1.Path, ctx2.Path)
	assert.Equal(t, ctx1.Name, ctx2.Name)
}
