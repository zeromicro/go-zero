package rpc

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// Run the real command in a separate process so Cobra flags and CLI globals
// cannot leak between cases. The helper needs no separately built goctl binary.
func TestRPCCommandProcess(t *testing.T) {
	for i, arg := range os.Args {
		if arg != "--" {
			continue
		}
		Cmd.SetArgs(os.Args[i+1:])
		if err := Cmd.Execute(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}
}

func rpcCommand(t *testing.T, dir string, args ...string) *exec.Cmd {
	t.Helper()
	binary, err := os.Executable()
	require.NoError(t, err)
	cmd := exec.Command(binary, append([]string{"-test.run=^TestRPCCommandProcess$", "--"}, args...)...)
	cmd.Dir = dir
	// Let goctl detect modules in the temporary directory instead of inheriting
	// a workspace selected by the developer running the test.
	cmd.Env = append(os.Environ(), "GOWORK=")
	return cmd
}

func TestProtocRejectsConflictingClientFlags(t *testing.T) {
	root := t.TempDir()
	// Distinct outputs ensure none of the three directories is created.
	args := []string{"protoc", "missing.proto", "--client-only", "--client=false"}
	for _, flag := range []string{"go_out", "go-grpc_out", "zrpc_out"} {
		args = append(args, "--"+flag+"="+filepath.Join(root, flag))
	}
	output, err := rpcCommand(t, root, args...).CombinedOutput()
	require.Error(t, err, "%s", output)
	require.Contains(t, string(output), "--client-only cannot be combined with --client=false")
	entries, err := os.ReadDir(root)
	require.NoError(t, err)
	require.Empty(t, entries)
}

func TestProtocClientOnlyCompiles(t *testing.T) {
	for _, name := range []string{"protoc", "protoc-gen-go", "protoc-gen-go-grpc"} {
		if _, err := exec.LookPath(name); err != nil {
			t.Skipf("%s is required for RPC generation: %v", name, err)
		}
	}
	examples, err := filepath.Abs("example")
	require.NoError(t, err)
	// Compile with the dependency versions and checksums already maintained by
	// goctl rather than resolving latest versions in every generated module.
	goMod, err := os.ReadFile("../go.mod")
	require.NoError(t, err)
	goSum, err := os.ReadFile("../go.sum")
	require.NoError(t, err)
	module := strings.Replace(string(goMod), "module github.com/zeromicro/go-zero/tools/goctl",
		"module example.com/demo", 1)
	for _, tc := range []struct {
		name    string
		source  string
		flags   []string
		clients []string
	}{
		{
			name: "basic", source: "01-basic/greeter.proto",
			clients: []string{"greeterclient/greeter.go"},
		},
		{
			name: "multiple", source: "05-multiple-services/multi.proto", flags: []string{"--multiple"},
			clients: []string{"client/searchservice/searchservice.go", "client/notifyservice/notifyservice.go"},
		},
		{
			name: "imports", source: "08-external-proto-diff-pkg/service.proto",
			flags:   []string{"-I", filepath.Join(examples, "08-external-proto-diff-pkg/ext_protos")},
			clients: []string{"dataservice/dataservice.go"},
		},
		{
			name: "streaming", source: "10-streaming/stream.proto",
			clients: []string{"streamservice/streamservice.go"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			output := filepath.Join(root, "output")
			source := filepath.Join(examples, filepath.FromSlash(tc.source))
			args := []string{
				"protoc", source, "--client-only", "--module=example.com/demo",
				"--go_out=" + output, "--go-grpc_out=" + output, "--zrpc_out=" + output,
				"--go_opt=module=example.com/demo", "--go-grpc_opt=module=example.com/demo",
				"-I", filepath.Dir(source),
			}
			result, err := rpcCommand(t, root, append(args, tc.flags...)...).CombinedOutput()
			require.NoError(t, err, "%s", result)
			for _, client := range tc.clients {
				require.FileExists(t, filepath.Join(output, filepath.FromSlash(client)))
			}
			require.NoDirExists(t, filepath.Join(output, "etc"))
			require.NoDirExists(t, filepath.Join(output, "internal"))
			// A server entry point would be a Go file directly under output.
			mains, err := filepath.Glob(filepath.Join(output, "*.go"))
			require.NoError(t, err)
			require.Empty(t, mains)
			require.FileExists(t, filepath.Join(output, "go.mod"))
			require.NoError(t, os.WriteFile(filepath.Join(output, "go.mod"), []byte(module), 0644))
			require.NoError(t, os.WriteFile(filepath.Join(output, "go.sum"), goSum, 0644))
			cmd := exec.Command("go", "test", "-mod=readonly", "-run=^$", "./...")
			cmd.Dir = output
			cmd.Env = append(os.Environ(), "GOWORK=off")
			result, err = cmd.CombinedOutput()
			require.NoError(t, err, "generated client must compile: %s", result)
		})
	}
}
