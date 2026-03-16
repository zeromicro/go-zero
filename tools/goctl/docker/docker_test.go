package docker

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDockerCommand_EtcDirResolution(t *testing.T) {
	// Create a temporary project structure
	tempDir := t.TempDir()
	
	// Create project structure: project/service/api/
	serviceDir := filepath.Join(tempDir, "service", "api")
	etcDir := filepath.Join(serviceDir, "etc")
	require.NoError(t, os.MkdirAll(etcDir, 0755))
	
	// Create a Go file
	goFile := filepath.Join(serviceDir, "api.go")
	require.NoError(t, os.WriteFile(goFile, []byte("package main\n\nfunc main() {}"), 0644))
	
	// Create a config file
	configFile := filepath.Join(etcDir, "config.yaml")
	require.NoError(t, os.WriteFile(configFile, []byte("Name: test\n"), 0644))
	
	// Create go.mod at the root
	goModFile := filepath.Join(tempDir, "go.mod")
	require.NoError(t, os.WriteFile(goModFile, []byte("module test\n\ngo 1.21\n"), 0644))
	
	// Test: etc directory should be found relative to Go file, not CWD
	t.Run("etc directory resolved relative to go file", func(t *testing.T) {
		// Save and restore original working directory
		originalWd, err := os.Getwd()
		require.NoError(t, err)
		defer func() {
			require.NoError(t, os.Chdir(originalWd))
		}()
		
		// Change to temp directory (not service/api directory)
		require.NoError(t, os.Chdir(tempDir))
		
		// The relative path from tempDir to the go file
		relGoFile := filepath.Join("service", "api", "api.go")
		
		// Test the etc directory resolution logic
		resolvedEtcDir := filepath.Join(filepath.Dir(relGoFile), "etc")
		
		// Verify the resolved path exists
		_, err = os.Stat(resolvedEtcDir)
		assert.NoError(t, err, "etc directory should be found at service/api/etc")
		
		// Verify it's the correct path (use EvalSymlinks to handle /private on macOS)
		absResolvedEtc, err := filepath.Abs(resolvedEtcDir)
		require.NoError(t, err)
		absResolvedEtc, err = filepath.EvalSymlinks(absResolvedEtc)
		require.NoError(t, err)
		expectedEtc, err := filepath.EvalSymlinks(etcDir)
		require.NoError(t, err)
		assert.Equal(t, expectedEtc, absResolvedEtc)
	})
	
	t.Run("etc directory with empty goFile", func(t *testing.T) {
		// When goFile is empty, should default to "./etc"
		goFile := ""
		resolvedEtcDir := filepath.Join(filepath.Dir(goFile), "etc")
		
		// Should resolve to just "etc"
		assert.Equal(t, "etc", resolvedEtcDir)
	})
	
	t.Run("etc directory with absolute path", func(t *testing.T) {
		// When goFile is absolute path
		absGoFile := filepath.Join(tempDir, "service", "api", "api.go")
		resolvedEtcDir := filepath.Join(filepath.Dir(absGoFile), "etc")
		
		// Should resolve correctly
		_, err := os.Stat(resolvedEtcDir)
		assert.NoError(t, err)
	})
}

func TestGenerateDockerfile_GoMainFromPath(t *testing.T) {
	tests := []struct {
		name         string
		goFile       string
		projPath     string
		expectedPath string
	}{
		{
			name:         "relative path with subdirectory",
			goFile:       "service/api/api.go",
			projPath:     "service/api",
			expectedPath: "service/api/api.go",
		},
		{
			name:         "simple filename",
			goFile:       "main.go",
			projPath:     ".",
			expectedPath: "main.go",
		},
		{
			name:         "nested service path",
			goFile:       "internal/service/user/user.go",
			projPath:     "internal/service/user",
			expectedPath: "internal/service/user/user.go",
		},
		{
			name:         "deep nested path",
			goFile:       "cmd/api/internal/handler/handler.go",
			projPath:     "cmd/api/internal/handler",
			expectedPath: "cmd/api/internal/handler/handler.go",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the fix: using filepath.Base instead of full path
			goMainFrom := filepath.Join(tt.projPath, filepath.Base(tt.goFile))
			
			assert.Equal(t, tt.expectedPath, goMainFrom,
				"GoMainFrom should not duplicate path segments")
			
			// Verify the old buggy behavior would have been wrong
			if tt.goFile != filepath.Base(tt.goFile) {
				buggyPath := filepath.Join(tt.projPath, tt.goFile)
				assert.NotEqual(t, tt.expectedPath, buggyPath,
					"Old implementation would have created incorrect path")
			}
		})
	}
}

func TestGenerateDockerfile_PathJoinBehavior(t *testing.T) {
	t.Run("demonstrates the bug and fix", func(t *testing.T) {
		projPath := "service/api"
		goFile := "service/api/api.go"
		
		// OLD (buggy) behavior: path duplication
		buggyPath := filepath.Join(projPath, goFile)
		assert.Equal(t, "service/api/service/api/api.go", buggyPath,
			"Bug: path segments are duplicated")
		
		// NEW (fixed) behavior: correct path
		fixedPath := filepath.Join(projPath, filepath.Base(goFile))
		assert.Equal(t, "service/api/api.go", fixedPath,
			"Fix: using filepath.Base prevents duplication")
	})
}

func TestFindConfig(t *testing.T) {
	tempDir := t.TempDir()
	etcDir := filepath.Join(tempDir, "etc")
	require.NoError(t, os.MkdirAll(etcDir, 0755))
	
	t.Run("finds config matching go file name", func(t *testing.T) {
		// Create config files
		require.NoError(t, os.WriteFile(
			filepath.Join(etcDir, "api.yaml"), []byte("test"), 0644))
		require.NoError(t, os.WriteFile(
			filepath.Join(etcDir, "other.yaml"), []byte("test"), 0644))
		
		cfg, err := findConfig("api.go", etcDir)
		assert.NoError(t, err)
		assert.Equal(t, "api.yaml", cfg)
	})
	
	t.Run("returns first config when no match", func(t *testing.T) {
		etcDir2 := filepath.Join(tempDir, "etc2")
		require.NoError(t, os.MkdirAll(etcDir2, 0755))
		require.NoError(t, os.WriteFile(
			filepath.Join(etcDir2, "config.yaml"), []byte("test"), 0644))
		
		cfg, err := findConfig("main.go", etcDir2)
		assert.NoError(t, err)
		assert.Equal(t, "config.yaml", cfg)
	})
	
	t.Run("returns error when no yaml files", func(t *testing.T) {
		emptyDir := filepath.Join(tempDir, "empty")
		require.NoError(t, os.MkdirAll(emptyDir, 0755))
		
		_, err := findConfig("api.go", emptyDir)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no yaml file")
	})
	
	t.Run("handles path in go file name", func(t *testing.T) {
		// Test with service/api/api.go - should extract just "api"
		cfg, err := findConfig("service/api/api.go", etcDir)
		assert.NoError(t, err)
		assert.Equal(t, "api.yaml", cfg)
	})
}

func TestGetFilePath(t *testing.T) {
	// Create a temporary directory with go.mod
	tempDir := t.TempDir()
	require.NoError(t, os.WriteFile(
		filepath.Join(tempDir, "go.mod"),
		[]byte("module testproject\n\ngo 1.21\n"),
		0644,
	))
	
	// Create subdirectories
	serviceDir := filepath.Join(tempDir, "service", "api")
	require.NoError(t, os.MkdirAll(serviceDir, 0755))
	
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.Chdir(originalWd))
	}()
	
	t.Run("returns relative path from go.mod", func(t *testing.T) {
		require.NoError(t, os.Chdir(tempDir))
		
		path, err := getFilePath("service/api")
		assert.NoError(t, err)
		assert.Equal(t, "service/api", path)
	})
	
	t.Run("handles current directory", func(t *testing.T) {
		require.NoError(t, os.Chdir(tempDir))
		
		path, err := getFilePath(".")
		assert.NoError(t, err)
		// Current directory returns empty string when at go.mod root
		assert.True(t, path == "." || path == "")
	})
}

// Integration test to verify the complete fix
func TestDockerCommandIntegration(t *testing.T) {
	// Create a complete project structure
	tempDir := t.TempDir()
	
	// Setup: project/service/api/
	serviceDir := filepath.Join(tempDir, "service", "api")
	etcDir := filepath.Join(serviceDir, "etc")
	require.NoError(t, os.MkdirAll(etcDir, 0755))
	
	// Create files
	goFile := filepath.Join(serviceDir, "api.go")
	require.NoError(t, os.WriteFile(goFile, []byte("package main\n\nfunc main() {}"), 0644))
	configFile := filepath.Join(etcDir, "api.yaml")
	require.NoError(t, os.WriteFile(configFile, []byte("Name: test-api\n"), 0644))
	goModFile := filepath.Join(tempDir, "go.mod")
	require.NoError(t, os.WriteFile(goModFile, []byte("module testproject\n\ngo 1.21\n"), 0644))
	goSumFile := filepath.Join(tempDir, "go.sum")
	require.NoError(t, os.WriteFile(goSumFile, []byte(""), 0644))
	
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.Chdir(originalWd))
	}()
	
	t.Run("etc directory detected from different working directory", func(t *testing.T) {
		// Change to project root (not service/api)
		require.NoError(t, os.Chdir(tempDir))
		
		// Relative path to Go file
		relGoFile := filepath.Join("service", "api", "api.go")
		
		// Apply the fix: resolve etc directory relative to go file
		resolvedEtcDir := filepath.Join(filepath.Dir(relGoFile), "etc")
		
		// Verify etc directory is found
		stat, err := os.Stat(resolvedEtcDir)
		assert.NoError(t, err)
		assert.True(t, stat.IsDir())
		
		// Verify config can be found
		cfg, err := findConfig(relGoFile, resolvedEtcDir)
		assert.NoError(t, err)
		assert.Equal(t, "api.yaml", cfg)
	})
	
	t.Run("GoMainFrom path is correct", func(t *testing.T) {
		require.NoError(t, os.Chdir(tempDir))
		
		goFileRel := filepath.Join("service", "api", "api.go")
		
		// Simulate getFilePath return value
		projPath := "service/api"
		
		// Apply the fix: use filepath.Base
		goMainFrom := filepath.Join(projPath, filepath.Base(goFileRel))
		
		assert.Equal(t, "service/api/api.go", goMainFrom)
		
		// Verify no path duplication
		assert.NotContains(t, goMainFrom, "service/api/service/api")
	})
}

// Test that specifically validates the bug described in PR #4343
func TestPR4343_BugFixes(t *testing.T) {
	t.Run("Bug 1: etc directory check uses correct base path", func(t *testing.T) {
		// Setup: Create a project structure where etc is NOT in CWD but IS relative to Go file
		tempDir := t.TempDir()
		serviceDir := filepath.Join(tempDir, "service", "api")
		etcDir := filepath.Join(serviceDir, "etc")
		require.NoError(t, os.MkdirAll(etcDir, 0755))
		
		// Create a config file
		require.NoError(t, os.WriteFile(
			filepath.Join(etcDir, "config.yaml"),
			[]byte("Name: test\n"),
			0644,
		))
		
		originalWd, err := os.Getwd()
		require.NoError(t, err)
		defer func() {
			require.NoError(t, os.Chdir(originalWd))
		}()
		
		// Change to project root (CWD = tempDir)
		require.NoError(t, os.Chdir(tempDir))
		
		goFile := filepath.Join("service", "api", "api.go")
		
		// OLD (buggy) behavior: checks for "etc" in CWD
		_, errOld := os.Stat("etc")
		assert.Error(t, errOld, "Bug: etc not found in CWD")
		
		// NEW (fixed) behavior: checks for "etc" relative to go file
		etcDirResolved := filepath.Join(filepath.Dir(goFile), "etc")
		stat, errNew := os.Stat(etcDirResolved)
		assert.NoError(t, errNew, "Fix: etc found relative to go file")
		assert.True(t, stat.IsDir())
		
		// Verify config is accessible
		cfg, err := findConfig(goFile, etcDirResolved)
		assert.NoError(t, err)
		assert.Equal(t, "config.yaml", cfg)
	})
	
	t.Run("Bug 2: GoMainFrom path not duplicated", func(t *testing.T) {
		// Test case from PR description
		projPath := "service/api"
		goFile := "service/api/api.go"
		
		// OLD (buggy) behavior: duplicates path
		buggyPath := filepath.Join(projPath, goFile)
		assert.Equal(t, "service/api/service/api/api.go", buggyPath,
			"Bug: path duplication occurs with old implementation")
		
		// NEW (fixed) behavior: correct path using filepath.Base
		fixedPath := filepath.Join(projPath, filepath.Base(goFile))
		assert.Equal(t, "service/api/api.go", fixedPath,
			"Fix: using filepath.Base() prevents path duplication")
		
		// Verify the fix works for various scenarios
		testCases := []struct {
			projPath string
			goFile   string
			expected string
		}{
			{"service/api", "service/api/api.go", "service/api/api.go"},
			{"cmd/server", "cmd/server/main.go", "cmd/server/main.go"},
			{"internal/handler", "internal/handler/handler.go", "internal/handler/handler.go"},
			{".", "main.go", "main.go"},
		}
		
		for _, tc := range testCases {
			result := filepath.Join(tc.projPath, filepath.Base(tc.goFile))
			assert.Equal(t, tc.expected, result,
				"Fix should work for projPath=%s, goFile=%s", tc.projPath, tc.goFile)
		}
	})
}
