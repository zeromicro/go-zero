package generator

import (
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
)

// rpcTypeRef holds the resolved Go type reference for an RPC request/response type.
type rpcTypeRef struct {
	// GoRef is the qualified Go type name, e.g. "pb.GetReq", "common.TypesReq", "emptypb.Empty".
	GoRef string
	// ImportPath is an extra Go import path required for cross-package types.
	// Empty when the type lives in a package that is already imported.
	ImportPath string
}

// resolveRPCTypeRef resolves a proto RPC type (possibly dotted) to its Go type
// reference and the optional extra import it requires.
//
//   - Simple types ("GetReq") → mainPbPackage.GetReq, no extra import.
//   - Same-package dotted types ("ext.ExtReq", same go_package) → mainPbPackage.ExtReq.
//   - Cross-package dotted types ("common.TypesReq") → common.TypesReq + import path.
//   - google.protobuf.X types → well-known Go type + import path.
func resolveRPCTypeRef(protoType, mainPbPackage, mainGoPackage string, pkgMap map[string]parser.ImportedProto) rpcTypeRef {
	if !strings.Contains(protoType, ".") {
		return rpcTypeRef{GoRef: fmt.Sprintf("%s.%s", mainPbPackage, parser.CamelCase(protoType))}
	}

	if strings.HasPrefix(protoType, "google.protobuf.") {
		typeName := strings.TrimPrefix(protoType, "google.protobuf.")
		return resolveGoogleWKT(typeName)
	}

	dot := strings.Index(protoType, ".")
	protoPkg, typeName := protoType[:dot], protoType[dot+1:]

	if imp, ok := pkgMap[protoPkg]; ok {
		camelType := parser.CamelCase(typeName)
		if imp.GoPackage == mainGoPackage {
			// Same Go package as main proto — no extra import needed.
			return rpcTypeRef{GoRef: fmt.Sprintf("%s.%s", mainPbPackage, camelType)}
		}
		return rpcTypeRef{
			GoRef:      fmt.Sprintf("%s.%s", imp.PbPackage, camelType),
			ImportPath: imp.GoPackage,
		}
	}

	// Fallback: treat as same package with CamelCase applied to full dotted name.
	return rpcTypeRef{GoRef: fmt.Sprintf("%s.%s", mainPbPackage, parser.CamelCase(protoType))}
}

// resolveCallTypeRef is tailored for gencall.go's type-alias system.
// Returns:
//   - typeName: the identifier to place in function signatures (alias name or full "pkg.Type" ref).
//   - aliasEntry: if non-empty, a "TypeName = pkg.TypeName" alias declaration to add to the type block.
//   - importPath: if non-empty, the extra import path needed.
func resolveCallTypeRef(protoType, mainPbPackage, mainGoPackage string, pkgMap map[string]parser.ImportedProto) (typeName, aliasEntry, importPath string) {
	if !strings.Contains(protoType, ".") {
		// Simple type — alias is produced by the existing proto.Message loop.
		return parser.CamelCase(protoType), "", ""
	}

	if strings.HasPrefix(protoType, "google.protobuf.") {
		tn := strings.TrimPrefix(protoType, "google.protobuf.")
		ref := resolveGoogleWKT(tn)
		return ref.GoRef, "", ref.ImportPath
	}

	dot := strings.Index(protoType, ".")
	protoPkg, tn := protoType[:dot], protoType[dot+1:]
	camelType := parser.CamelCase(tn)

	if imp, ok := pkgMap[protoPkg]; ok {
		if imp.GoPackage == mainGoPackage {
			// Same Go package: add an alias so the function signature uses a simple name.
			entry := fmt.Sprintf("%s = %s.%s", camelType, mainPbPackage, camelType)
			return camelType, entry, ""
		}
		// Different Go package: use fully-qualified ref directly; no alias needed.
		return fmt.Sprintf("%s.%s", imp.PbPackage, camelType), "", imp.GoPackage
	}

	return parser.CamelCase(protoType), "", ""
}

// googleWKTTable maps google.protobuf type names to their generated Go equivalents.
var googleWKTTable = map[string]rpcTypeRef{
	"Empty":       {GoRef: "emptypb.Empty", ImportPath: "google.golang.org/protobuf/types/known/emptypb"},
	"Timestamp":   {GoRef: "timestamppb.Timestamp", ImportPath: "google.golang.org/protobuf/types/known/timestamppb"},
	"Duration":    {GoRef: "durationpb.Duration", ImportPath: "google.golang.org/protobuf/types/known/durationpb"},
	"Any":         {GoRef: "anypb.Any", ImportPath: "google.golang.org/protobuf/types/known/anypb"},
	"StringValue": {GoRef: "wrapperspb.StringValue", ImportPath: "google.golang.org/protobuf/types/known/wrapperspb"},
	"Int32Value":  {GoRef: "wrapperspb.Int32Value", ImportPath: "google.golang.org/protobuf/types/known/wrapperspb"},
	"Int64Value":  {GoRef: "wrapperspb.Int64Value", ImportPath: "google.golang.org/protobuf/types/known/wrapperspb"},
	"BoolValue":   {GoRef: "wrapperspb.BoolValue", ImportPath: "google.golang.org/protobuf/types/known/wrapperspb"},
	"BytesValue":  {GoRef: "wrapperspb.BytesValue", ImportPath: "google.golang.org/protobuf/types/known/wrapperspb"},
	"FloatValue":  {GoRef: "wrapperspb.FloatValue", ImportPath: "google.golang.org/protobuf/types/known/wrapperspb"},
	"DoubleValue": {GoRef: "wrapperspb.DoubleValue", ImportPath: "google.golang.org/protobuf/types/known/wrapperspb"},
	"UInt32Value": {GoRef: "wrapperspb.UInt32Value", ImportPath: "google.golang.org/protobuf/types/known/wrapperspb"},
	"UInt64Value": {GoRef: "wrapperspb.UInt64Value", ImportPath: "google.golang.org/protobuf/types/known/wrapperspb"},
	"Struct":      {GoRef: "structpb.Struct", ImportPath: "google.golang.org/protobuf/types/known/structpb"},
	"Value":       {GoRef: "structpb.Value", ImportPath: "google.golang.org/protobuf/types/known/structpb"},
	"ListValue":   {GoRef: "structpb.ListValue", ImportPath: "google.golang.org/protobuf/types/known/structpb"},
	"FieldMask":   {GoRef: "fieldmaskpb.FieldMask", ImportPath: "google.golang.org/protobuf/types/known/fieldmaskpb"},
}

func resolveGoogleWKT(typeName string) rpcTypeRef {
	if r, ok := googleWKTTable[typeName]; ok {
		return r
	}
	return rpcTypeRef{GoRef: "interface{}"}
}
