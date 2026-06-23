package parser

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveImports_Basic(t *testing.T) {
	// test.proto imports "base.proto", which lives in the same directory.
	absDir, err := filepath.Abs(".")
	assert.NoError(t, err)

	paths, err := ResolveImports("./test.proto", []string{absDir})
	assert.NoError(t, err)
	assert.Len(t, paths, 1)
	assert.True(t, strings.HasSuffix(paths[0], "base.proto"))
	assert.True(t, filepath.IsAbs(paths[0]))
}

func TestResolveImports_SourceExcluded(t *testing.T) {
	// The source file itself must not appear in the result.
	absDir, err := filepath.Abs(".")
	assert.NoError(t, err)

	absSrc, err := filepath.Abs("./test.proto")
	assert.NoError(t, err)

	paths, err := ResolveImports("./test.proto", []string{absDir})
	assert.NoError(t, err)
	for _, p := range paths {
		assert.NotEqual(t, absSrc, p)
	}
}

func TestResolveImports_NotFound(t *testing.T) {
	// Imports that cannot be located in protoPaths are silently skipped.
	paths, err := ResolveImports("./test.proto", []string{"/nonexistent/path"})
	assert.NoError(t, err)
	assert.Empty(t, paths)
}

func TestResolveImports_NoDuplicates(t *testing.T) {
	// Even if the same proto is found via multiple search paths, it should
	// appear only once.
	absDir, err := filepath.Abs(".")
	assert.NoError(t, err)

	paths, err := ResolveImports("./test.proto", []string{absDir, absDir})
	assert.NoError(t, err)

	seen := make(map[string]int)
	for _, p := range paths {
		seen[p]++
	}
	for p, count := range seen {
		assert.Equal(t, 1, count, "duplicate path: %s", p)
	}
}

func TestParseImportedProtos_Basic(t *testing.T) {
	absDir, err := filepath.Abs(".")
	assert.NoError(t, err)

	protos, err := ParseImportedProtos("./test.proto", []string{absDir})
	assert.NoError(t, err)
	assert.Len(t, protos, 1)

	imp := protos[0]
	assert.Equal(t, "github.com/zeromicro/go-zero/tools/goctl/rpc/parser/base", imp.GoPackage)
	assert.Equal(t, "base", imp.PbPackage)
	assert.True(t, filepath.IsAbs(imp.Src))
}

func TestParseImportedProtos_EmptyWhenNoImports(t *testing.T) {
	// test_option.proto has no imports, so the result should be empty.
	absDir, err := filepath.Abs(".")
	assert.NoError(t, err)

	protos, err := ParseImportedProtos("./test_option.proto", []string{absDir})
	assert.NoError(t, err)
	assert.Empty(t, protos)
}

func TestIsWellKnownProto(t *testing.T) {
	assert.True(t, isWellKnownProto("google/protobuf/timestamp.proto"))
	assert.True(t, isWellKnownProto("google/protobuf/empty.proto"))
	assert.False(t, isWellKnownProto("base.proto"))
	assert.False(t, isWellKnownProto("common/types.proto"))
}

func TestLookupProtoFile_Found(t *testing.T) {
	absDir, err := filepath.Abs(".")
	assert.NoError(t, err)

	got, err := lookupProtoFile("base.proto", []string{absDir})
	assert.NoError(t, err)
	assert.True(t, filepath.IsAbs(got))
	assert.True(t, strings.HasSuffix(got, "base.proto"))
}

func TestLookupProtoFile_NotFound(t *testing.T) {
	_, err := lookupProtoFile("nonexistent.proto", []string{"/no/such/dir"})
	assert.Error(t, err)
}
