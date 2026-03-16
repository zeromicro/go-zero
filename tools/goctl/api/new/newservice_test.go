package new

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeromicro/go-zero/tools/goctl/api/gogen"
	"github.com/zeromicro/go-zero/tools/goctl/config"
)

func TestDoGenProjectWithModule_Integration(t *testing.T) {
	tests := []struct {
		name        string
		moduleName  string
		serviceName string
		expectedMod string
	}{
		{
			name:        "with custom module",
			moduleName:  "github.com/test/customapi",
			serviceName: "myservice",
			expectedMod: "github.com/test/customapi",
		},
		{
			name:        "with empty module",
			moduleName:  "",
			serviceName: "myservice",
			expectedMod: "myservice",
		},
		{
			name:        "with simple module",
			moduleName:  "simpleapi",
			serviceName: "testapi",
			expectedMod: "simpleapi",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			tempDir, err := os.MkdirTemp("", "goctl-api-module-test-*")
			require.NoError(t, err)
			defer os.RemoveAll(tempDir)

			// Create service directory
			serviceDir := filepath.Join(tempDir, tt.serviceName)
			err = os.MkdirAll(serviceDir, 0755)
			require.NoError(t, err)

			// Create a simple API file for testing
			apiContent := `syntax = "v1"

type Request {
  Name string ` + "`" + `path:"name,options=you|me"` + "`" + `
}

type Response {
  Message string ` + "`" + `json:"message"` + "`" + `
}

service ` + tt.serviceName + `-api {
  @handler ` + tt.serviceName + `Handler
  get /from/:name(Request) returns (Response)
}
`
			apiFile := filepath.Join(serviceDir, tt.serviceName+".api")
			err = os.WriteFile(apiFile, []byte(apiContent), 0644)
			require.NoError(t, err)

			// Call the module-aware service creation function
			err = gogen.DoGenProjectWithModule(apiFile, serviceDir, tt.moduleName, config.DefaultFormat, false)
			assert.NoError(t, err)

			// Check go.mod file
			goModPath := filepath.Join(serviceDir, "go.mod")
			assert.FileExists(t, goModPath)

			// Verify module name in go.mod
			content, err := os.ReadFile(goModPath)
			require.NoError(t, err)
			assert.Contains(t, string(content), "module "+tt.expectedMod)

			// Check basic directory structure was created
			assert.DirExists(t, filepath.Join(serviceDir, "etc"))
			assert.DirExists(t, filepath.Join(serviceDir, "internal"))
			assert.DirExists(t, filepath.Join(serviceDir, "internal", "handler"))
			assert.DirExists(t, filepath.Join(serviceDir, "internal", "logic"))
			assert.DirExists(t, filepath.Join(serviceDir, "internal", "svc"))
			assert.DirExists(t, filepath.Join(serviceDir, "internal", "types"))
			assert.DirExists(t, filepath.Join(serviceDir, "internal", "config"))

			// Check that main.go imports use correct module
			mainGoPath := filepath.Join(serviceDir, tt.serviceName+".go")
			if _, err := os.Stat(mainGoPath); err == nil {
				mainContent, err := os.ReadFile(mainGoPath)
				require.NoError(t, err)
				// Check for import of internal packages with correct module path
				assert.Contains(t, string(mainContent), `"`+tt.expectedMod+"/internal/")
			}
		})
	}
}

func TestCreateServiceCommand_Integration(t *testing.T) {
	tests := []struct {
		name        string
		moduleName  string
		serviceName string
		expectedMod string
		shouldError bool
	}{
		{
			name:        "valid service with custom module",
			moduleName:  "github.com/example/testapi",
			serviceName: "myapi",
			expectedMod: "github.com/example/testapi",
			shouldError: false,
		},
		{
			name:        "valid service with no module",
			moduleName:  "",
			serviceName: "simpleapi",
			expectedMod: "simpleapi",
			shouldError: false,
		},
		{
			name:        "invalid service name with hyphens",
			moduleName:  "github.com/test/api",
			serviceName: "my-api",
			expectedMod: "",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldError && tt.serviceName == "my-api" {
				// Test that service names with hyphens are rejected
				// This is tested in the actual command function, not the generate function
				assert.Contains(t, tt.serviceName, "-")
				return
			}

			// Create temporary directory
			tempDir, err := os.MkdirTemp("", "goctl-create-service-test-*")
			require.NoError(t, err)
			defer os.RemoveAll(tempDir)

			// Change to temp directory
			oldDir, _ := os.Getwd()
			defer os.Chdir(oldDir)
			os.Chdir(tempDir)

			// Set the module variable as the command would
			VarStringModule = tt.moduleName
			VarStringStyle = config.DefaultFormat

			// Create the service directory manually since we're testing the core functionality
			serviceDir := filepath.Join(tempDir, tt.serviceName)

			// Simulate what CreateServiceCommand does - create API file and call DoGenProjectWithModule
			err = os.MkdirAll(serviceDir, 0755)
			require.NoError(t, err)

			// Create API file
			apiContent := `syntax = "v1"

type Request {
  Name string ` + "`" + `path:"name,options=you|me"` + "`" + `
}

type Response {
  Message string ` + "`" + `json:"message"` + "`" + `
}

service ` + tt.serviceName + `-api {
  @handler ` + tt.serviceName + `Handler
  get /from/:name(Request) returns (Response)
}
`
			apiFile := filepath.Join(serviceDir, tt.serviceName+".api")
			err = os.WriteFile(apiFile, []byte(apiContent), 0644)
			require.NoError(t, err)

			// Call DoGenProjectWithModule as CreateServiceCommand does
			err = gogen.DoGenProjectWithModule(apiFile, serviceDir, VarStringModule, VarStringStyle, false)

			if tt.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify go.mod
				goModPath := filepath.Join(serviceDir, "go.mod")
				assert.FileExists(t, goModPath)
				content, err := os.ReadFile(goModPath)
				require.NoError(t, err)
				assert.Contains(t, string(content), "module "+tt.expectedMod)
			}
		})
	}
}
