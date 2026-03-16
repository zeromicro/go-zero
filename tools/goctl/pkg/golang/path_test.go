package golang

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetParentPackage(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "goctl-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Test with a directory (should create go.mod with directory name)
	testDir := filepath.Join(tempDir, "testproject")
	err = os.MkdirAll(testDir, 0755)
	require.NoError(t, err)

	parentPkg, rootPkg, err := GetParentPackage(testDir)
	assert.NoError(t, err)
	assert.Equal(t, "testproject", parentPkg)
	assert.Equal(t, "testproject", rootPkg)

	// Verify go.mod was created with directory name
	goModPath := filepath.Join(testDir, "go.mod")
	assert.FileExists(t, goModPath)

	content, err := os.ReadFile(goModPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "module testproject")
}

func TestGetParentPackageWithModule(t *testing.T) {
	tests := []struct {
		name           string
		moduleName     string
		expectedModule string
		expectedPkg    string
	}{
		{
			name:           "custom module name",
			moduleName:     "github.com/example/myproject",
			expectedModule: "github.com/example/myproject",
			expectedPkg:    "github.com/example/myproject",
		},
		{
			name:           "simple module name",
			moduleName:     "myservice",
			expectedModule: "myservice",
			expectedPkg:    "myservice",
		},
		{
			name:           "empty module name falls back to directory",
			moduleName:     "",
			expectedModule: "fallback",
			expectedPkg:    "fallback",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary directory for testing
			tempDir, err := os.MkdirTemp("", "goctl-test-*")
			require.NoError(t, err)
			defer os.RemoveAll(tempDir)

			// Create test directory - use "fallback" name for empty module test
			testDirName := "fallback"
			if tt.name != "empty module name falls back to directory" {
				testDirName = "testdir"
			}

			testDir := filepath.Join(tempDir, testDirName)
			err = os.MkdirAll(testDir, 0755)
			require.NoError(t, err)

			parentPkg, rootPkg, err := GetParentPackageWithModule(testDir, tt.moduleName)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedPkg, parentPkg)
			assert.Equal(t, tt.expectedModule, rootPkg)

			// Verify go.mod was created with correct module name
			goModPath := filepath.Join(testDir, "go.mod")
			assert.FileExists(t, goModPath)

			content, err := os.ReadFile(goModPath)
			require.NoError(t, err)
			assert.Contains(t, string(content), "module "+tt.expectedModule)
		})
	}
}

func TestGetParentPackageWithModule_InvalidDir(t *testing.T) {
	// Test with non-existent directory
	_, _, err := GetParentPackageWithModule("/non/existent/path", "github.com/example/test")
	assert.Error(t, err)
}

func TestGetParentPackage_InvalidDir(t *testing.T) {
	// Test with non-existent directory
	_, _, err := GetParentPackage("/non/existent/path")
	assert.Error(t, err)
}

func TestGetParentPackage_UsesGetParentPackageWithModule(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "goctl-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	testDir := filepath.Join(tempDir, "testproject")
	err = os.MkdirAll(testDir, 0755)
	require.NoError(t, err)

	// Test that GetParentPackage calls GetParentPackageWithModule with empty string
	parentPkg1, rootPkg1, err1 := GetParentPackage(testDir)
	require.NoError(t, err1)

	// Clean up go.mod to test again
	os.Remove(filepath.Join(testDir, "go.mod"))

	parentPkg2, rootPkg2, err2 := GetParentPackageWithModule(testDir, "")
	require.NoError(t, err2)

	// Should produce identical results
	assert.Equal(t, parentPkg1, parentPkg2)
	assert.Equal(t, rootPkg1, rootPkg2)
}

func TestBuildParentPackage(t *testing.T) {
	// This tests the internal buildParentPackage function indirectly
	// through the public API, as it's a private function

	// Create a temporary directory with subdirectory structure
	tempDir, err := os.MkdirTemp("", "goctl-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a nested directory structure
	projectDir := filepath.Join(tempDir, "myproject")
	subDir := filepath.Join(projectDir, "internal", "logic")
	err = os.MkdirAll(subDir, 0755)
	require.NoError(t, err)

	// Test from root directory
	parentPkg, rootPkg, err := GetParentPackageWithModule(projectDir, "github.com/example/myproject")
	assert.NoError(t, err)
	assert.Equal(t, "github.com/example/myproject", parentPkg)
	assert.Equal(t, "github.com/example/myproject", rootPkg)

	// Test from subdirectory
	parentPkg2, rootPkg2, err := GetParentPackageWithModule(subDir, "github.com/example/myproject")
	assert.NoError(t, err)
	assert.Equal(t, "github.com/example/myproject/internal/logic", parentPkg2)
	assert.Equal(t, "github.com/example/myproject", rootPkg2)
}

func TestGetParentPackageWithModule_SpecialCharacters(t *testing.T) {
	tests := []struct {
		name       string
		moduleName string
		valid      bool
	}{
		{
			name:       "domain with path",
			moduleName: "github.com/user/repo",
			valid:      true,
		},
		{
			name:       "domain with version",
			moduleName: "github.com/user/repo/v2",
			valid:      true,
		},
		{
			name:       "private repo",
			moduleName: "private.example.com/team/project",
			valid:      true,
		},
		{
			name:       "simple name with underscore",
			moduleName: "my_project",
			valid:      true,
		},
		{
			name:       "simple name with hyphen",
			moduleName: "my-project",
			valid:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary directory for testing
			tempDir, err := os.MkdirTemp("", "goctl-test-*")
			require.NoError(t, err)
			defer os.RemoveAll(tempDir)

			testDir := filepath.Join(tempDir, "testdir")
			err = os.MkdirAll(testDir, 0755)
			require.NoError(t, err)

			parentPkg, rootPkg, err := GetParentPackageWithModule(testDir, tt.moduleName)

			if tt.valid {
				assert.NoError(t, err)
				assert.Equal(t, tt.moduleName, parentPkg)
				assert.Equal(t, tt.moduleName, rootPkg)

				// Verify go.mod contains the module name
				goModPath := filepath.Join(testDir, "go.mod")
				content, err := os.ReadFile(goModPath)
				require.NoError(t, err)
				assert.Contains(t, string(content), "module "+tt.moduleName)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
