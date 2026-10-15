package generator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	protoparser "github.com/emicklei/proto"
	"github.com/stretchr/testify/require"
	conf "github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util/ctx"
)

func TestMkdirClientOnly(t *testing.T) {
	for _, clientOnly := range []bool{false, true} {
		for _, multiple := range []bool{false, true} {
			t.Run(fmt.Sprintf("clientOnly=%t/multiple=%t", clientOnly, multiple), func(t *testing.T) {
				root := t.TempDir()
				output := filepath.Join(root, "rpc")
				project := &ctx.ProjectContext{WorkDir: output, Dir: root, Path: "example.com/gateway"}
				cfg, err := conf.NewConfig("gozero")
				require.NoError(t, err)
				proto := parser.Proto{
					Name: "service.proto", GoPackage: "pb",
					Service: parser.Services{{Service: &protoparser.Service{Name: "Greeter"}}},
				}
				context := &ZRpcContext{
					IsGenClient: true, ClientOnly: clientOnly, Multiple: multiple,
					ProtoGenGrpcDir: filepath.Join(output, "pb"), ProtoGenGoDir: filepath.Join(output, "pb"),
				}
				dirs, err := mkdir(project, proto, cfg, context)
				require.NoError(t, err)
				require.DirExists(t, dirs.GetCall().Filename)
				require.DirExists(t, dirs.GetPb().Filename)
				for _, name := range []string{"etc", "internal", "internal/config", "internal/logic", "internal/server", "internal/svc"} {
					path := filepath.Join(output, filepath.FromSlash(name))
					if clientOnly {
						require.NoDirExists(t, path)
					} else {
						require.DirExists(t, path)
					}
				}
			})
		}
	}
}

func TestGenerateRejectsClientOnlyWithoutClient(t *testing.T) {
	output := filepath.Join(t.TempDir(), "output")
	err := NewGenerator("gozero", false).Generate(&ZRpcContext{Output: output, ClientOnly: true})
	require.EqualError(t, err, "--client-only cannot be combined with --client=false")
	require.NoDirExists(t, output)
}

func TestGenerateClientOnly(t *testing.T) {
	for _, name := range []string{"protoc", "protoc-gen-go", "protoc-gen-go-grpc"} {
		if _, err := exec.LookPath(name); err != nil {
			t.Skipf("%s is required for RPC generation: %v", name, err)
		}
	}
	for _, existingModule := range []bool{false, true} {
		for _, multiple := range []bool{false, true} {
			t.Run(fmt.Sprintf("existingModule=%t/multiple=%t", existingModule, multiple), func(t *testing.T) {
				root := t.TempDir()
				output := filepath.Join(root, "rpc")
				require.NoError(t, os.MkdirAll(output, 0755))
				module := []byte("module example.com/gateway\n\ngo 1.25.0\n")
				if existingModule {
					require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), module, 0644))
				}
				source := filepath.Join(root, "service.proto")
				proto := `syntax = "proto3";
package sample;
option go_package = "example.com/gateway/rpc/pb";
message Request { string name = 1; }
message Response { string message = 1; }
service Greeter { rpc Greet(Request) returns (Response); }
`
				if multiple {
					proto += "service Other { rpc Echo(Request) returns (Response); }\n"
				}
				require.NoError(t, os.WriteFile(source, []byte(proto), 0644))
				context := &ZRpcContext{
					Src: source, Output: output, GoOutput: output, GrpcOutput: output,
					IsGooglePlugin: true, IsGenClient: true, ClientOnly: true, Multiple: multiple,
					Module: "example.com/gateway/rpc", ProtoPaths: []string{root},
					ProtocCmd: fmt.Sprintf("protoc -I=%q %q --go_out=%q --go-grpc_out=%q --go_opt=module=example.com/gateway/rpc --go-grpc_opt=module=example.com/gateway/rpc", filepath.ToSlash(root), filepath.ToSlash(source), filepath.ToSlash(output), filepath.ToSlash(output)),
				}
				require.NoError(t, NewGenerator("gozero", false).Generate(context))
				require.FileExists(t, filepath.Join(output, "pb", "service.pb.go"))
				require.FileExists(t, filepath.Join(output, "pb", "service_grpc.pb.go"))
				if multiple {
					require.FileExists(t, filepath.Join(output, "client", "greeter", "greeter.go"))
					require.FileExists(t, filepath.Join(output, "client", "other", "other.go"))
				} else {
					require.FileExists(t, filepath.Join(output, "greeter", "greeter.go"))
				}
				require.NoDirExists(t, filepath.Join(output, "etc"))
				require.NoDirExists(t, filepath.Join(output, "internal"))
				require.NoFileExists(t, filepath.Join(output, "sample.go"))
				if existingModule {
					require.NoFileExists(t, filepath.Join(output, "go.mod"))
					content, err := os.ReadFile(filepath.Join(root, "go.mod"))
					require.NoError(t, err)
					require.Equal(t, module, content)
				} else {
					require.FileExists(t, filepath.Join(output, "go.mod"))
				}

				// Regenerating a client must also leave an existing server intact.
				serverFiles := []string{"etc/sample.yaml", "internal/svc/servicecontext.go", "sample.go"}
				for _, name := range serverFiles {
					path := filepath.Join(output, filepath.FromSlash(name))
					require.NoError(t, os.MkdirAll(filepath.Dir(path), 0755))
					require.NoError(t, os.WriteFile(path, []byte("existing server content\n"), 0644))
				}
				require.NoError(t, NewGenerator("gozero", false).Generate(context))
				for _, name := range serverFiles {
					content, err := os.ReadFile(filepath.Join(output, filepath.FromSlash(name)))
					require.NoError(t, err)
					require.Equal(t, "existing server content\n", string(content))
				}
			})
		}
	}
}
