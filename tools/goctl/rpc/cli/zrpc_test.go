package cli

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestZRPCRejectsConflictingClientFlagsBeforeWriting(t *testing.T) {
	oldClient, oldClientOnly := VarBoolClient, VarBoolClientOnly
	oldGoOut, oldGrpcOut, oldZrpcOut := VarStringSliceGoOut, VarStringSliceGoGRPCOut, VarStringZRPCOut
	t.Cleanup(func() {
		VarBoolClient, VarBoolClientOnly = oldClient, oldClientOnly
		VarStringSliceGoOut, VarStringSliceGoGRPCOut, VarStringZRPCOut = oldGoOut, oldGrpcOut, oldZrpcOut
	})
	output := filepath.Join(t.TempDir(), "not-created")
	VarBoolClient, VarBoolClientOnly = false, true
	VarStringSliceGoOut = []string{output}
	VarStringSliceGoGRPCOut = []string{output}
	VarStringZRPCOut = output
	require.EqualError(t, ZRPC(nil, []string{"service.proto"}), "--client-only cannot be combined with --client=false")
	require.NoDirExists(t, output)
}
