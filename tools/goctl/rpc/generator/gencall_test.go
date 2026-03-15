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

// TestGenCallGroup_OnlyUsedTypesAliased verifies that in multi-service mode each
// generated client file contains type aliases only for the message types actually
// used by that service's RPCs (fix for issue #5481).
func TestGenCallGroup_OnlyUsedTypesAliased(t *testing.T) {
	tmpDir := t.TempDir()
	callBase := filepath.Join(tmpDir, "call")
	pbBase := filepath.Join(tmpDir, "pb")

	// Pre-create subdirs that genCallGroup will write into.
	require.NoError(t, os.MkdirAll(filepath.Join(callBase, "servicea"), 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(callBase, "serviceb"), 0755))
	require.NoError(t, os.MkdirAll(pbBase, 0755))

	mctx := &mockDirContext{
		callDir: Dir{
			Filename: callBase,
			Package:  "example.com/multitest/call",
			Base:     "call",
			GetChildPackage: func(childPath string) (string, error) {
				// Return a package path whose Base() is the lowercase service name.
				return filepath.Join(callBase, strings.ToLower(childPath)), nil
			},
		},
		pbDir: Dir{
			Filename: pbBase,
			Package:  "example.com/multitest/pb",
			Base:     "pb",
		},
		protoGo: Dir{
			// Must differ from "servicea"/"serviceb" so isCallPkgSameToPbPkg stays false
			// and alias generation is triggered.
			Filename: pbBase,
			Package:  "example.com/multitest/pb",
			Base:     "pb",
		},
	}

	// Proto with two services that use completely disjoint message types.
	protoData := parser.Proto{
		Name:      "multi.proto",
		PbPackage: "pb",
		Message: []parser.Message{
			{Message: &proto.Message{Name: "AReq"}},
			{Message: &proto.Message{Name: "AResp"}},
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

	g := NewGenerator("gozero", false)
	require.NoError(t, g.genCallGroup(mctx, protoData, cfg))

	// servicea/servicea.go — aliases for AReq/AResp only
	aContent, err := os.ReadFile(filepath.Join(callBase, "servicea", "servicea.go"))
	require.NoError(t, err)
	aFile := string(aContent)

	assert.Contains(t, aFile, "AReq", "ServiceA file should alias AReq")
	assert.Contains(t, aFile, "AResp", "ServiceA file should alias AResp")
	assert.NotContains(t, aFile, "BReq", "ServiceA file must not alias BReq")
	assert.NotContains(t, aFile, "BResp", "ServiceA file must not alias BResp")

	// serviceb/serviceb.go — aliases for BReq/BResp only
	bContent, err := os.ReadFile(filepath.Join(callBase, "serviceb", "serviceb.go"))
	require.NoError(t, err)
	bFile := string(bContent)

	assert.Contains(t, bFile, "BReq", "ServiceB file should alias BReq")
	assert.Contains(t, bFile, "BResp", "ServiceB file should alias BResp")
	assert.NotContains(t, bFile, "AReq", "ServiceB file must not alias AReq")
	assert.NotContains(t, bFile, "AResp", "ServiceB file must not alias AResp")
}
