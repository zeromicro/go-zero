package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/emicklei/proto"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
	ctxpkg "github.com/zeromicro/go-zero/tools/goctl/util/ctx"
)

func TestServiceNameDetermination(t *testing.T) {
	tests := []struct {
		name             string
		protoName        string
		packageName      string
		hasPackage       bool
		nameFromFilename bool
		expectedName     string
	}{
		{
			name:             "default uses package name when available",
			protoName:        "user.proto",
			packageName:      "userservice",
			hasPackage:       true,
			nameFromFilename: false,
			expectedName:     "userservice",
		},
		{
			name:             "flag enabled uses filename instead of package",
			protoName:        "user.proto",
			packageName:      "userservice",
			hasPackage:       true,
			nameFromFilename: true,
			expectedName:     "user",
		},
		{
			name:             "fallback to filename when package is empty",
			protoName:        "order.proto",
			packageName:      "",
			hasPackage:       true,
			nameFromFilename: false,
			expectedName:     "order",
		},
		{
			name:             "fallback to filename when package is nil",
			protoName:        "product.proto",
			packageName:      "",
			hasPackage:       false,
			nameFromFilename: false,
			expectedName:     "product",
		},
		{
			name:             "flag enabled with nil package uses filename",
			protoName:        "catalog.proto",
			packageName:      "",
			hasPackage:       false,
			nameFromFilename: true,
			expectedName:     "catalog",
		},
		{
			name:             "handles proto file with complex name",
			protoName:        "user-service.proto",
			packageName:      "user",
			hasPackage:       true,
			nameFromFilename: false,
			expectedName:     "user",
		},
		{
			name:             "handles proto file with complex name and flag",
			protoName:        "user-service.proto",
			packageName:      "user",
			hasPackage:       true,
			nameFromFilename: true,
			expectedName:     "user-service",
		},
		{
			name:             "nil context uses package name",
			protoName:        "account.proto",
			packageName:      "accountpb",
			hasPackage:       true,
			nameFromFilename: false,
			expectedName:     "accountpb",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Build the proto struct
			p := parser.Proto{
				Name: tt.protoName,
			}
			if tt.hasPackage {
				p.Package = parser.Package{
					Package: &proto.Package{
						Name: tt.packageName,
					},
				}
			}

			// Build the context
			var ctx *ZRpcContext
			if tt.name != "nil context uses package name" {
				ctx = &ZRpcContext{
					NameFromFilename: tt.nameFromFilename,
				}
			}

			// Call the helper function to determine service name
			serviceName := determineServiceName(p, ctx)
			assert.Equal(t, tt.expectedName, serviceName)
		})
	}
}

func TestServiceNameWithNilContext(t *testing.T) {
	p := parser.Proto{
		Name: "test.proto",
		Package: parser.Package{
			Package: &proto.Package{
				Name: "testpkg",
			},
		},
	}

	// nil context should use package name
	serviceName := determineServiceName(p, nil)
	assert.Equal(t, "testpkg", serviceName)
}

func TestServiceNameFallbackWithEmptyPackage(t *testing.T) {
	p := parser.Proto{
		Name: "myservice.proto",
		Package: parser.Package{
			Package: &proto.Package{
				Name: "", // empty package name
			},
		},
	}

	ctx := &ZRpcContext{
		NameFromFilename: false,
	}

	// Should fall back to filename when package name is empty
	serviceName := determineServiceName(p, ctx)
	assert.Equal(t, "myservice", serviceName)
}

func TestJoinImportPathNormalizesSeparators(t *testing.T) {
	got := joinImportPath("sales-center", filepath.Join("pb", "admin"))
	assert.Equal(t, "sales-center/pb/admin", got)
}

func TestResolveDirPackageWithinCurrentModule(t *testing.T) {
	tmp := t.TempDir()
	mainDir := filepath.Join(tmp, "sales-admin")

	err := os.MkdirAll(filepath.Join(mainDir, "internal", "logic"), 0o755)
	assert.NoError(t, err)
	err = os.WriteFile(filepath.Join(mainDir, "go.mod"), []byte("module sales-admin\n\ngo 1.25.5\n"), 0o644)
	assert.NoError(t, err)

	project := &ctxpkg.ProjectContext{
		Path: "sales-admin",
		Dir:  mainDir,
	}

	got := resolveDirPackage(project, filepath.Join(mainDir, "internal", "logic"), "")
	assert.Equal(t, "sales-admin/internal/logic", got)
}

func TestResolveDirPackageForSiblingModule(t *testing.T) {
	tmp := t.TempDir()
	mainDir := filepath.Join(tmp, "sales-admin")
	siblingDir := filepath.Join(tmp, "sales-center")
	targetDir := filepath.Join(siblingDir, "pb", "admin")

	err := os.MkdirAll(mainDir, 0o755)
	assert.NoError(t, err)
	err = os.MkdirAll(siblingDir, 0o755)
	assert.NoError(t, err)
	err = os.WriteFile(filepath.Join(mainDir, "go.mod"), []byte("module sales-admin\n\ngo 1.25.5\n"), 0o644)
	assert.NoError(t, err)
	err = os.WriteFile(filepath.Join(siblingDir, "go.mod"), []byte("module sales-center\n\ngo 1.25.5\n"), 0o644)
	assert.NoError(t, err)

	project := &ctxpkg.ProjectContext{
		Path: "sales-admin",
		Dir:  mainDir,
	}

	got := resolveDirPackage(project, targetDir, "")
	assert.Equal(t, "sales-center/pb/admin", got)
}

func TestResolveDirPackageFallsBackToProvidedImportPath(t *testing.T) {
	tmp := t.TempDir()
	mainDir := filepath.Join(tmp, "sales-admin")
	targetDir := filepath.Join(tmp, "generated", "admin")

	err := os.MkdirAll(mainDir, 0o755)
	assert.NoError(t, err)
	err = os.WriteFile(filepath.Join(mainDir, "go.mod"), []byte("module sales-admin\n\ngo 1.25.5\n"), 0o644)
	assert.NoError(t, err)

	project := &ctxpkg.ProjectContext{
		Path: "sales-admin",
		Dir:  mainDir,
	}

	got := resolveDirPackage(project, targetDir, "sales-center/pb/admin")
	assert.Equal(t, "sales-center/pb/admin", got)
}

func TestResolveModuleForDirStopsAtRoot(t *testing.T) {
	tmp := t.TempDir()
	targetDir := filepath.Join(tmp, "generated", "admin")

	modulePath, moduleDir, ok := resolveModuleForDir(targetDir)
	assert.False(t, ok)
	assert.Empty(t, modulePath)
	assert.Empty(t, moduleDir)
}

func TestResolveModuleForDirWithNonExistentDirectory(t *testing.T) {
	// Should not panic when the directory doesn't exist on disk
	modulePath, moduleDir, ok := resolveModuleForDir(filepath.Join("Z:", "nonexistent", "deep", "path"))
	assert.False(t, ok)
	assert.Empty(t, modulePath)
	assert.Empty(t, moduleDir)
}

func TestResolveModuleForDirFindsGoModInParent(t *testing.T) {
	tmp := t.TempDir()
	moduleRoot := filepath.Join(tmp, "mymodule")
	deepDir := filepath.Join(moduleRoot, "a", "b", "c", "d")

	err := os.MkdirAll(deepDir, 0o755)
	assert.NoError(t, err)
	err = os.WriteFile(filepath.Join(moduleRoot, "go.mod"), []byte("module example.com/mymodule\n\ngo 1.21\n"), 0o644)
	assert.NoError(t, err)

	modulePath, moduleDir, ok := resolveModuleForDir(deepDir)
	assert.True(t, ok)
	assert.Equal(t, "example.com/mymodule", modulePath)
	assert.Equal(t, moduleRoot, moduleDir)
}

func TestRelativeImportPathRejectsParentTraversal(t *testing.T) {
	base := filepath.Join("a", "b", "c")
	target := filepath.Join("a", "b")

	_, ok := relativeImportPath(base, target)
	assert.False(t, ok)
}

func TestRelativeImportPathSameDir(t *testing.T) {
	dir := filepath.Join("a", "b", "c")
	rel, ok := relativeImportPath(dir, dir)
	assert.True(t, ok)
	assert.Equal(t, "", rel)
}

func TestRelativeImportPathChild(t *testing.T) {
	base := filepath.Join("a", "b")
	target := filepath.Join("a", "b", "c", "d")

	rel, ok := relativeImportPath(base, target)
	assert.True(t, ok)
	assert.Equal(t, "c/d", rel)
}

func TestIsImportPathVariousCases(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", false},
		{".", false},
		{"./relative", false},
		{"../parent", false},
		{"/absolute/unix", false},
		{"C:/windows/path", false},
		{"C:\\windows\\path", false},
		{"D:/some/path", false},
		{"github.com/user/repo", true},
		{"sales-center/pb/admin", true},
		{"mypackage", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, isImportPath(tt.input))
		})
	}
}

func TestReadModulePathWithComments(t *testing.T) {
	tmp := t.TempDir()
	goMod := filepath.Join(tmp, "go.mod")

	content := "// this is a comment\nmodule example.com/commented // inline comment\n\ngo 1.21\n"
	err := os.WriteFile(goMod, []byte(content), 0o644)
	assert.NoError(t, err)

	modulePath, err := readModulePath(goMod)
	assert.NoError(t, err)
	assert.Equal(t, "example.com/commented", modulePath)
}

func TestReadModulePathEmptyFile(t *testing.T) {
	tmp := t.TempDir()
	goMod := filepath.Join(tmp, "go.mod")

	err := os.WriteFile(goMod, []byte(""), 0o644)
	assert.NoError(t, err)

	modulePath, err := readModulePath(goMod)
	assert.NoError(t, err)
	assert.Empty(t, modulePath)
}

func TestReadModulePathNoModuleDirective(t *testing.T) {
	tmp := t.TempDir()
	goMod := filepath.Join(tmp, "go.mod")

	err := os.WriteFile(goMod, []byte("go 1.21\n\nrequire (\n)\n"), 0o644)
	assert.NoError(t, err)

	modulePath, err := readModulePath(goMod)
	assert.NoError(t, err)
	assert.Empty(t, modulePath)
}

func TestResolveDirPackageFallsBackToBasename(t *testing.T) {
	tmp := t.TempDir()
	mainDir := filepath.Join(tmp, "sales-admin")
	targetDir := filepath.Join(tmp, "orphan-dir", "sub")

	err := os.MkdirAll(mainDir, 0o755)
	assert.NoError(t, err)
	err = os.WriteFile(filepath.Join(mainDir, "go.mod"), []byte("module sales-admin\n\ngo 1.21\n"), 0o644)
	assert.NoError(t, err)

	project := &ctxpkg.ProjectContext{
		Path: "sales-admin",
		Dir:  mainDir,
	}

	// No go.mod in orphan-dir, no valid fallback → should use basename
	got := resolveDirPackage(project, targetDir, "")
	assert.Equal(t, "sub", got)
}

func TestResolveDirPackageWhenDirEqualsProjectDir(t *testing.T) {
	tmp := t.TempDir()
	mainDir := filepath.Join(tmp, "myproject")

	err := os.MkdirAll(mainDir, 0o755)
	assert.NoError(t, err)
	err = os.WriteFile(filepath.Join(mainDir, "go.mod"), []byte("module example.com/myproject\n\ngo 1.21\n"), 0o644)
	assert.NoError(t, err)

	project := &ctxpkg.ProjectContext{
		Path: "example.com/myproject",
		Dir:  mainDir,
	}

	got := resolveDirPackage(project, mainDir, "")
	assert.Equal(t, "example.com/myproject", got)
}

func TestJoinImportPathWithEmptyRel(t *testing.T) {
	assert.Equal(t, "example.com/mod", joinImportPath("example.com/mod", ""))
	assert.Equal(t, "example.com/mod", joinImportPath("example.com/mod", "."))
}

func TestSetPbDirKeepsExistingPackageWhenSiblingModuleHasNoGoMod(t *testing.T) {
	tmp := t.TempDir()
	mainDir := filepath.Join(tmp, "sales-admin")
	pbDir := filepath.Join(tmp, "generated", "admin")

	err := os.MkdirAll(mainDir, 0o755)
	assert.NoError(t, err)
	err = os.MkdirAll(pbDir, 0o755)
	assert.NoError(t, err)
	err = os.WriteFile(filepath.Join(mainDir, "go.mod"), []byte("module sales-admin\n\ngo 1.25.5\n"), 0o644)
	assert.NoError(t, err)

	project := &ctxpkg.ProjectContext{
		Path: "sales-admin",
		Dir:  mainDir,
	}
	dirCtx := &defaultDirContext{
		ctx: project,
		inner: map[string]Dir{
			pb: {
				Package: "sales-center/pb/admin",
			},
			protoGo: {
				Package: "sales-center/pb/admin",
			},
		},
	}

	dirCtx.SetPbDir(pbDir, pbDir)

	assert.Equal(t, "sales-center/pb/admin", dirCtx.GetPb().Package)
	assert.Equal(t, "sales-center/pb/admin", dirCtx.GetProtoGo().Package)
	assert.Equal(t, pbDir, dirCtx.GetPb().Filename)
	assert.Equal(t, pbDir, dirCtx.GetProtoGo().Filename)
}
