package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/emicklei/proto"
)

// Import embeds proto.Import
type Import struct {
	*proto.Import
}

// ImportedProto holds the package information of a transitively imported proto file.
type ImportedProto struct {
	// Src is the absolute path to the proto file.
	Src string
	// ProtoPackage is the value of the proto "package" declaration.
	// It is the qualifier used in dotted type references, e.g. "ext" in "ext.ExtReq".
	ProtoPackage string
	// GoPackage is the value of the option go_package field, or the proto
	// package name when go_package is absent.
	GoPackage string
	// PbPackage is the sanitized Go package name derived from GoPackage.
	PbPackage string
}

// BuildProtoPackageMap returns a map from proto package name to ImportedProto,
// enabling O(1) lookup of Go package info given a proto type qualifier like "ext".
func BuildProtoPackageMap(importedProtos []ImportedProto) map[string]ImportedProto {
	m := make(map[string]ImportedProto, len(importedProtos))
	for _, imp := range importedProtos {
		if imp.ProtoPackage != "" {
			m[imp.ProtoPackage] = imp
		}
	}
	return m
}

// ResolveImports returns the absolute paths of all transitively imported proto
// files reachable from src, excluding well-known types (google/*).
// It searches for imported files in protoPaths (equivalent to protoc -I flags).
// Files that cannot be found in protoPaths are silently skipped so that
// system-level or well-known protos do not cause errors.
func ResolveImports(src string, protoPaths []string) ([]string, error) {
	absSrc, err := filepath.Abs(src)
	if err != nil {
		return nil, err
	}

	visited := make(map[string]bool)
	visited[absSrc] = true // exclude the source itself from the result
	var result []string
	if err := collectImports(absSrc, protoPaths, visited, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// ParseImportedProtos resolves and parses all transitively imported proto
// files, returning their package information for use in code generation.
func ParseImportedProtos(src string, protoPaths []string) ([]ImportedProto, error) {
	paths, err := ResolveImports(src, protoPaths)
	if err != nil {
		return nil, err
	}

	result := make([]ImportedProto, 0, len(paths))
	for _, p := range paths {
		goPackage, pbPackage, protoPackage, err := parseGoPackage(p)
		if err != nil {
			return nil, err
		}
		result = append(result, ImportedProto{
			Src:          p,
			ProtoPackage: protoPackage,
			GoPackage:    goPackage,
			PbPackage:    pbPackage,
		})
	}
	return result, nil
}

// collectImports recursively walks import declarations of src, appending newly
// discovered absolute proto file paths to result.
func collectImports(src string, protoPaths []string, visited map[string]bool, result *[]string) error {
	importFilenames, err := parseImportFilenames(src)
	if err != nil {
		return err
	}

	for _, filename := range importFilenames {
		if isWellKnownProto(filename) {
			continue
		}

		abs, err := lookupProtoFile(filename, protoPaths)
		if err != nil {
			// Not found in the provided proto paths — may be a system-level proto.
			// Skip rather than fail, mirroring protoc's own behaviour.
			continue
		}

		if visited[abs] {
			continue
		}
		visited[abs] = true
		*result = append(*result, abs)

		if err := collectImports(abs, protoPaths, visited, result); err != nil {
			return err
		}
	}
	return nil
}

// parseImportFilenames opens src and returns the Filename field of every
// import statement without performing any file-system lookups.
func parseImportFilenames(src string) ([]string, error) {
	r, err := os.Open(src)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	p := proto.NewParser(r)
	set, err := p.Parse()
	if err != nil {
		return nil, err
	}

	var imports []string
	proto.Walk(set, proto.WithImport(func(i *proto.Import) {
		imports = append(imports, i.Filename)
	}))
	return imports, nil
}

// parseGoPackage reads only the go_package option and package declaration from
// src, returning the derived GoPackage, PbPackage, and ProtoPackage without
// requiring a service definition (imported protos often have no service block).
func parseGoPackage(src string) (goPackage, pbPackage, protoPackage string, err error) {
	r, err := os.Open(src)
	if err != nil {
		return "", "", "", err
	}
	defer r.Close()

	p := proto.NewParser(r)
	set, err := p.Parse()
	if err != nil {
		return "", "", "", err
	}

	var packageName string
	proto.Walk(set,
		proto.WithOption(func(opt *proto.Option) {
			if opt.Name == "go_package" {
				goPackage = opt.Constant.Source
			}
		}),
		proto.WithPackage(func(pkg *proto.Package) {
			packageName = pkg.Name
		}),
	)

	if len(goPackage) == 0 {
		goPackage = packageName
	}
	pbPackage = GoSanitized(filepath.Base(goPackage))
	protoPackage = packageName
	return goPackage, pbPackage, protoPackage, nil
}

// lookupProtoFile searches for filename inside each directory of protoPaths,
// returning its absolute path on the first match.
func lookupProtoFile(filename string, protoPaths []string) (string, error) {
	for _, dir := range protoPaths {
		candidate := filepath.Join(dir, filename)
		if _, err := os.Stat(candidate); err == nil {
			return filepath.Abs(candidate)
		}
	}
	return "", fmt.Errorf("proto file %q not found in proto paths %v", filename, protoPaths)
}

// isWellKnownProto reports whether filename refers to a well-known type
// bundled with protoc (e.g. google/protobuf/timestamp.proto).
func isWellKnownProto(filename string) bool {
	return strings.HasPrefix(filename, "google/")
}
