package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/emicklei/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	conf "github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util/stringx"
)

// mockDirContext is a minimal DirContext for unit-testing genCallGroup.
type mockDirContext struct {
	callDir Dir
	pbDir   Dir
	protoGo Dir
}

func (m *mockDirContext) GetCall() Dir                   { return m.callDir }
func (m *mockDirContext) GetEtc() Dir                    { return Dir{} }
func (m *mockDirContext) GetInternal() Dir               { return Dir{} }
func (m *mockDirContext) GetConfig() Dir                 { return Dir{} }
func (m *mockDirContext) GetLogic() Dir                  { return Dir{} }
func (m *mockDirContext) GetServer() Dir                 { return Dir{} }
func (m *mockDirContext) GetSvc() Dir                    { return Dir{} }
func (m *mockDirContext) GetPb() Dir                     { return m.pbDir }
func (m *mockDirContext) GetProtoGo() Dir                { return m.protoGo }
func (m *mockDirContext) GetMain() Dir                   { return Dir{} }
func (m *mockDirContext) GetServiceName() stringx.String { return stringx.From("test") }
func (m *mockDirContext) SetPbDir(pbDir, grpcDir string) {}

// newTestDirContext builds a mockDirContext that writes generated files under
// callBase, with a pb directory that differs (so alias generation is triggered).
func newTestDirContext(t *testing.T, callBase, pbBase string, services ...string) *mockDirContext {
	t.Helper()
	for _, svc := range services {
		require.NoError(t, os.MkdirAll(filepath.Join(callBase, strings.ToLower(svc)), 0755))
	}
	require.NoError(t, os.MkdirAll(pbBase, 0755))
	return &mockDirContext{
		callDir: Dir{
			Filename: callBase,
			Package:  "example.com/test/call",
			Base:     "call",
			GetChildPackage: func(childPath string) (string, error) {
				return filepath.Join(callBase, strings.ToLower(childPath)), nil
			},
		},
		pbDir: Dir{Filename: pbBase, Package: "example.com/test/pb", Base: "pb"},
		protoGo: Dir{
			// Must differ from service dir names so isCallPkgSameToPbPkg stays
			// false and alias generation is triggered.
			Filename: pbBase,
			Package:  "example.com/test/pb",
			Base:     "pb",
		},
	}
}

// ---- unit tests for collectServiceUsedTypes --------------------------------

// TestCollectServiceUsedTypes_DirectOnly verifies that request and response
// types with no message fields are collected as-is.
func TestCollectServiceUsedTypes_DirectOnly(t *testing.T) {
	messages := []parser.Message{
		{Message: &proto.Message{Name: "AReq"}},
		{Message: &proto.Message{Name: "AResp"}},
		{Message: &proto.Message{Name: "Unrelated"}},
	}
	service := parser.Service{
		Service: &proto.Service{Name: "ServiceA"},
		RPC: []*parser.RPC{
			{RPC: &proto.RPC{Name: "Do", RequestType: "AReq", ReturnsType: "AResp"}},
		},
	}

	got := collectServiceUsedTypes(messages, service)

	assert.True(t, got.Contains("AReq"))
	assert.True(t, got.Contains("AResp"))
	assert.False(t, got.Contains("Unrelated"), "unrelated message must not be collected")
}

// TestCollectServiceUsedTypes_NestedNormalField verifies that a message type
// referenced via a NormalField inside a response is transitively collected
// (regression test for issue #5618).
func TestCollectServiceUsedTypes_NestedNormalField(t *testing.T) {
	messages := []parser.Message{
		{Message: &proto.Message{Name: "AReq"}},
		{Message: &proto.Message{
			Name: "AResp",
			Elements: []proto.Visitee{
				&proto.NormalField{Field: &proto.Field{Name: "items", Type: "AItem"}},
			},
		}},
		{Message: &proto.Message{Name: "AItem"}},
	}
	service := parser.Service{
		Service: &proto.Service{Name: "ServiceA"},
		RPC: []*parser.RPC{
			{RPC: &proto.RPC{Name: "List", RequestType: "AReq", ReturnsType: "AResp"}},
		},
	}

	got := collectServiceUsedTypes(messages, service)

	assert.True(t, got.Contains("AReq"))
	assert.True(t, got.Contains("AResp"))
	assert.True(t, got.Contains("AItem"), "field type AItem must be transitively collected")
}

// TestCollectServiceUsedTypes_MapValueField verifies that the value type of a
// MapField inside a response message is transitively collected.
func TestCollectServiceUsedTypes_MapValueField(t *testing.T) {
	messages := []parser.Message{
		{Message: &proto.Message{Name: "AReq"}},
		{Message: &proto.Message{
			Name: "AResp",
			Elements: []proto.Visitee{
				&proto.MapField{KeyType: "string", Field: &proto.Field{Name: "index", Type: "AItem"}},
			},
		}},
		{Message: &proto.Message{Name: "AItem"}},
	}
	service := parser.Service{
		Service: &proto.Service{Name: "ServiceA"},
		RPC: []*parser.RPC{
			{RPC: &proto.RPC{Name: "GetMap", RequestType: "AReq", ReturnsType: "AResp"}},
		},
	}

	got := collectServiceUsedTypes(messages, service)

	assert.True(t, got.Contains("AResp"))
	assert.True(t, got.Contains("AItem"), "map value type AItem must be transitively collected")
}

// TestCollectServiceUsedTypes_OneofField verifies that message types referenced
// inside a Oneof element are transitively collected.
func TestCollectServiceUsedTypes_OneofField(t *testing.T) {
	oneof := &proto.Oneof{Name: "result"}
	oneof.Elements = []proto.Visitee{
		&proto.OneOfField{Field: &proto.Field{Name: "success", Type: "SuccessMsg"}},
		&proto.OneOfField{Field: &proto.Field{Name: "failure", Type: "FailureMsg"}},
	}
	messages := []parser.Message{
		{Message: &proto.Message{Name: "AReq"}},
		{Message: &proto.Message{
			Name:     "AResp",
			Elements: []proto.Visitee{oneof},
		}},
		{Message: &proto.Message{Name: "SuccessMsg"}},
		{Message: &proto.Message{Name: "FailureMsg"}},
	}
	service := parser.Service{
		Service: &proto.Service{Name: "ServiceA"},
		RPC: []*parser.RPC{
			{RPC: &proto.RPC{Name: "Do", RequestType: "AReq", ReturnsType: "AResp"}},
		},
	}

	got := collectServiceUsedTypes(messages, service)

	assert.True(t, got.Contains("AResp"))
	assert.True(t, got.Contains("SuccessMsg"), "oneof field type SuccessMsg must be collected")
	assert.True(t, got.Contains("FailureMsg"), "oneof field type FailureMsg must be collected")
}

// TestCollectServiceUsedTypes_MultiLevelTransitive verifies that a chain
// AResp → BMsg → CMsg is fully collected (multi-level transitivity).
func TestCollectServiceUsedTypes_MultiLevelTransitive(t *testing.T) {
	messages := []parser.Message{
		{Message: &proto.Message{Name: "AReq"}},
		{Message: &proto.Message{
			Name: "AResp",
			Elements: []proto.Visitee{
				&proto.NormalField{Field: &proto.Field{Name: "b", Type: "BMsg"}},
			},
		}},
		{Message: &proto.Message{
			Name: "BMsg",
			Elements: []proto.Visitee{
				&proto.NormalField{Field: &proto.Field{Name: "c", Type: "CMsg"}},
			},
		}},
		{Message: &proto.Message{Name: "CMsg"}},
	}
	service := parser.Service{
		Service: &proto.Service{Name: "ServiceA"},
		RPC: []*parser.RPC{
			{RPC: &proto.RPC{Name: "Do", RequestType: "AReq", ReturnsType: "AResp"}},
		},
	}

	got := collectServiceUsedTypes(messages, service)

	assert.True(t, got.Contains("AReq"))
	assert.True(t, got.Contains("AResp"))
	assert.True(t, got.Contains("BMsg"), "BMsg must be transitively collected via AResp")
	assert.True(t, got.Contains("CMsg"), "CMsg must be transitively collected via BMsg")
}

// TestCollectServiceUsedTypes_CycleDetection verifies that circular field
// references (AResp ↔ BMsg) do not cause infinite recursion.
func TestCollectServiceUsedTypes_CycleDetection(t *testing.T) {
	messages := []parser.Message{
		{Message: &proto.Message{Name: "AReq"}},
		{Message: &proto.Message{
			Name: "AResp",
			Elements: []proto.Visitee{
				&proto.NormalField{Field: &proto.Field{Name: "b", Type: "BMsg"}},
			},
		}},
		{Message: &proto.Message{
			Name: "BMsg",
			Elements: []proto.Visitee{
				// circular back-reference to AResp
				&proto.NormalField{Field: &proto.Field{Name: "a", Type: "AResp"}},
			},
		}},
	}
	service := parser.Service{
		Service: &proto.Service{Name: "ServiceA"},
		RPC: []*parser.RPC{
			{RPC: &proto.RPC{Name: "Do", RequestType: "AReq", ReturnsType: "AResp"}},
		},
	}

	// Must not panic or loop; both messages are reachable.
	got := collectServiceUsedTypes(messages, service)

	assert.True(t, got.Contains("AResp"))
	assert.True(t, got.Contains("BMsg"))
}

// TestCollectServiceUsedTypes_ExcludesUnrelatedService verifies that messages
// belonging only to another service are not included.
func TestCollectServiceUsedTypes_ExcludesUnrelatedService(t *testing.T) {
	messages := []parser.Message{
		{Message: &proto.Message{Name: "AReq"}},
		{Message: &proto.Message{Name: "AResp"}},
		{Message: &proto.Message{Name: "BReq"}},
		{Message: &proto.Message{Name: "BResp"}},
	}
	service := parser.Service{
		Service: &proto.Service{Name: "ServiceA"},
		RPC: []*parser.RPC{
			{RPC: &proto.RPC{Name: "DoA", RequestType: "AReq", ReturnsType: "AResp"}},
		},
	}

	got := collectServiceUsedTypes(messages, service)

	assert.True(t, got.Contains("AReq"))
	assert.True(t, got.Contains("AResp"))
	assert.False(t, got.Contains("BReq"), "BReq belongs to ServiceB and must be excluded")
	assert.False(t, got.Contains("BResp"), "BResp belongs to ServiceB and must be excluded")
}

// ---- integration tests via genCallGroup ------------------------------------

// TestGenCallGroup_OnlyUsedTypesAliased verifies that in multi-service mode
// each generated client file aliases only its own request/response types and
// their transitive field dependencies (fix for issues #5481 and #5618).
func TestGenCallGroup_OnlyUsedTypesAliased(t *testing.T) {
	tmpDir := t.TempDir()
	callBase := filepath.Join(tmpDir, "call")
	pbBase := filepath.Join(tmpDir, "pb")

	mctx := newTestDirContext(t, callBase, pbBase, "ServiceA", "ServiceB")

	// ServiceA: AResp contains a NormalField of type AItem (issue #5618).
	// ServiceB: BResp has no nested message fields.
	// AItem must appear in ServiceA's file but not ServiceB's.
	protoData := parser.Proto{
		Name:      "multi.proto",
		PbPackage: "pb",
		Message: []parser.Message{
			{Message: &proto.Message{Name: "AReq"}},
			{Message: &proto.Message{
				Name: "AResp",
				Elements: []proto.Visitee{
					&proto.NormalField{Field: &proto.Field{Name: "items", Type: "AItem"}},
				},
			}},
			{Message: &proto.Message{Name: "AItem"}},
			{Message: &proto.Message{Name: "BReq"}},
			{Message: &proto.Message{Name: "BResp"}},
		},
		Service: parser.Services{
			{
				Service: &proto.Service{Name: "ServiceA"},
				RPC: []*parser.RPC{
					{RPC: &proto.RPC{Name: "DoA", RequestType: "AReq", ReturnsType: "AResp"}},
				},
			},
			{
				Service: &proto.Service{Name: "ServiceB"},
				RPC: []*parser.RPC{
					{RPC: &proto.RPC{Name: "DoB", RequestType: "BReq", ReturnsType: "BResp"}},
				},
			},
		},
	}

	cfg, err := conf.NewConfig("")
	require.NoError(t, err)
	require.NoError(t, NewGenerator("gozero", false).genCallGroup(mctx, protoData, cfg))

	aFile := normalizeWS(readGenFile(t, callBase, "servicea", "servicea.go"))
	assert.Contains(t, aFile, "AReq = pb.AReq", "ServiceA must alias AReq")
	assert.Contains(t, aFile, "AResp = pb.AResp", "ServiceA must alias AResp")
	assert.Contains(t, aFile, "AItem = pb.AItem", "ServiceA must alias AItem (transitive NormalField)")
	assert.NotContains(t, aFile, "BReq = pb.BReq", "ServiceA must not alias BReq")
	assert.NotContains(t, aFile, "BResp = pb.BResp", "ServiceA must not alias BResp")

	bFile := normalizeWS(readGenFile(t, callBase, "serviceb", "serviceb.go"))
	assert.Contains(t, bFile, "BReq = pb.BReq", "ServiceB must alias BReq")
	assert.Contains(t, bFile, "BResp = pb.BResp", "ServiceB must alias BResp")
	assert.NotContains(t, bFile, "AReq = pb.AReq", "ServiceB must not alias AReq")
	assert.NotContains(t, bFile, "AResp = pb.AResp", "ServiceB must not alias AResp")
	assert.NotContains(t, bFile, "AItem = pb.AItem", "ServiceB must not alias AItem")
}

// TestGenCallGroup_MapValueAliased verifies that the value type of a MapField
// inside a service response is included in the generated aliases.
func TestGenCallGroup_MapValueAliased(t *testing.T) {
	tmpDir := t.TempDir()
	callBase := filepath.Join(tmpDir, "call")
	pbBase := filepath.Join(tmpDir, "pb")

	mctx := newTestDirContext(t, callBase, pbBase, "ServiceA")

	protoData := parser.Proto{
		Name:      "map.proto",
		PbPackage: "pb",
		Message: []parser.Message{
			{Message: &proto.Message{Name: "AReq"}},
			{Message: &proto.Message{
				Name: "AResp",
				Elements: []proto.Visitee{
					&proto.MapField{KeyType: "string", Field: &proto.Field{Name: "index", Type: "AItem"}},
				},
			}},
			{Message: &proto.Message{Name: "AItem"}},
		},
		Service: parser.Services{
			{
				Service: &proto.Service{Name: "ServiceA"},
				RPC: []*parser.RPC{
					{RPC: &proto.RPC{Name: "GetMap", RequestType: "AReq", ReturnsType: "AResp"}},
				},
			},
		},
	}

	cfg, err := conf.NewConfig("")
	require.NoError(t, err)
	require.NoError(t, NewGenerator("gozero", false).genCallGroup(mctx, protoData, cfg))

	aFile := normalizeWS(readGenFile(t, callBase, "servicea", "servicea.go"))
	assert.Contains(t, aFile, "AResp = pb.AResp")
	assert.Contains(t, aFile, "AItem = pb.AItem", "map value type AItem must be aliased")
}

// TestGenCallGroup_OneofAliased verifies that message types referenced inside a
// Oneof element are included in the generated aliases.
func TestGenCallGroup_OneofAliased(t *testing.T) {
	tmpDir := t.TempDir()
	callBase := filepath.Join(tmpDir, "call")
	pbBase := filepath.Join(tmpDir, "pb")

	mctx := newTestDirContext(t, callBase, pbBase, "ServiceA")

	oneof := &proto.Oneof{Name: "result"}
	oneof.Elements = []proto.Visitee{
		&proto.OneOfField{Field: &proto.Field{Name: "ok", Type: "SuccessMsg"}},
		&proto.OneOfField{Field: &proto.Field{Name: "err", Type: "FailureMsg"}},
	}
	protoData := parser.Proto{
		Name:      "oneof.proto",
		PbPackage: "pb",
		Message: []parser.Message{
			{Message: &proto.Message{Name: "AReq"}},
			{Message: &proto.Message{
				Name:     "AResp",
				Elements: []proto.Visitee{oneof},
			}},
			{Message: &proto.Message{Name: "SuccessMsg"}},
			{Message: &proto.Message{Name: "FailureMsg"}},
		},
		Service: parser.Services{
			{
				Service: &proto.Service{Name: "ServiceA"},
				RPC: []*parser.RPC{
					{RPC: &proto.RPC{Name: "Do", RequestType: "AReq", ReturnsType: "AResp"}},
				},
			},
		},
	}

	cfg, err := conf.NewConfig("")
	require.NoError(t, err)
	require.NoError(t, NewGenerator("gozero", false).genCallGroup(mctx, protoData, cfg))

	aFile := normalizeWS(readGenFile(t, callBase, "servicea", "servicea.go"))
	assert.Contains(t, aFile, "SuccessMsg = pb.SuccessMsg", "oneof type SuccessMsg must be aliased")
	assert.Contains(t, aFile, "FailureMsg = pb.FailureMsg", "oneof type FailureMsg must be aliased")
}

// TestGenCallGroup_MultiLevelTransitiveAliased verifies that a dependency chain
// AResp → BMsg → CMsg causes all three types to be aliased in the client file.
func TestGenCallGroup_MultiLevelTransitiveAliased(t *testing.T) {
	tmpDir := t.TempDir()
	callBase := filepath.Join(tmpDir, "call")
	pbBase := filepath.Join(tmpDir, "pb")

	mctx := newTestDirContext(t, callBase, pbBase, "ServiceA")

	protoData := parser.Proto{
		Name:      "transitive.proto",
		PbPackage: "pb",
		Message: []parser.Message{
			{Message: &proto.Message{Name: "AReq"}},
			{Message: &proto.Message{
				Name: "AResp",
				Elements: []proto.Visitee{
					&proto.NormalField{Field: &proto.Field{Name: "b", Type: "BMsg"}},
				},
			}},
			{Message: &proto.Message{
				Name: "BMsg",
				Elements: []proto.Visitee{
					&proto.NormalField{Field: &proto.Field{Name: "c", Type: "CMsg"}},
				},
			}},
			{Message: &proto.Message{Name: "CMsg"}},
		},
		Service: parser.Services{
			{
				Service: &proto.Service{Name: "ServiceA"},
				RPC: []*parser.RPC{
					{RPC: &proto.RPC{Name: "Do", RequestType: "AReq", ReturnsType: "AResp"}},
				},
			},
		},
	}

	cfg, err := conf.NewConfig("")
	require.NoError(t, err)
	require.NoError(t, NewGenerator("gozero", false).genCallGroup(mctx, protoData, cfg))

	aFile := normalizeWS(readGenFile(t, callBase, "servicea", "servicea.go"))
	assert.Contains(t, aFile, "AResp = pb.AResp")
	assert.Contains(t, aFile, "BMsg = pb.BMsg", "BMsg must be transitively aliased via AResp")
	assert.Contains(t, aFile, "CMsg = pb.CMsg", "CMsg must be transitively aliased via BMsg")
}

// readGenFile reads a generated file relative to callBase and returns its content.
func readGenFile(t *testing.T, callBase string, parts ...string) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(append([]string{callBase}, parts...)...))
	require.NoError(t, err)
	return string(content)
}

// normalizeWS replaces runs of whitespace with a single space.
func normalizeWS(s string) string {
	return strings.Join(strings.Fields(strings.ReplaceAll(s, "\n", " \n ")), " ")
}
