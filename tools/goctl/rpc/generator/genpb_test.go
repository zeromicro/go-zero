package generator

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

func Test_findPbFile(t *testing.T) {
	dir := t.TempDir()
	protoFile := filepath.Join(dir, "greet.proto")
	err := os.WriteFile(protoFile, []byte(`
syntax = "proto3";

package greet;
option go_package="./greet";

message Req{}
message Resp{}
service Greeter {
  rpc greet(Req) returns (Resp);
}
`), 0o666)
	if err != nil {
		t.Log(err)
		return
	}
	t.Run("", func(t *testing.T) {
		output := t.TempDir()
		grpc := filepath.Join(output, "grpc")
		err := pathx.MkdirIfNotExist(grpc)
		if err != nil {
			t.Log(err)
			return
		}
		cmd := exec.Command("protoc", "-I="+filepath.Dir(protoFile), "--go_out="+output, "--go-grpc_out="+grpc, filepath.Base(protoFile))
		cmd.Dir = output
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			t.Log(err)
			return
		}
		pbDir, err := findPbFile(output, protoFile, false)
		assert.Nil(t, err)
		pbGo := filepath.Join(pbDir, "greet.pb.go")
		assert.True(t, pathx.FileExists(pbGo))

		grpcDir, err := findPbFile(output, protoFile, true)
		assert.Nil(t, err)
		grpcGo := filepath.Join(grpcDir, "greet_grpc.pb.go")
		assert.True(t, pathx.FileExists(grpcGo))
	})

	t.Run("", func(t *testing.T) {
		output := t.TempDir()
		redirect := filepath.Join(output, "pb")
		grpc := filepath.Join(output, "grpc")
		err := pathx.MkdirIfNotExist(grpc)
		if err != nil {
			t.Log(err)
			return
		}
		cmd := exec.Command("protoc", "-I="+filepath.Dir(protoFile), "--go_out="+output,
			"--go-grpc_out="+grpc, filepath.Base(protoFile), "--go_opt=M"+filepath.Base(protoFile)+"="+redirect)
		cmd.Dir = output
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			t.Log(err)
			return
		}
		pbDir, err := findPbFile(output, protoFile, false)
		assert.Nil(t, err)
		pbGo := filepath.Join(pbDir, "greet.pb.go")
		assert.True(t, pathx.FileExists(pbGo))

		grpcDir, err := findPbFile(output, protoFile, true)
		assert.Nil(t, err)
		grpcGo := filepath.Join(grpcDir, "greet_grpc.pb.go")
		assert.True(t, pathx.FileExists(grpcGo))
	})

	t.Run("", func(t *testing.T) {
		output := t.TempDir()
		pbeRedirect := filepath.Join(output, "redirect")
		grpc := filepath.Join(output, "grpc")
		grpcRedirect := filepath.Join(grpc, "redirect")
		err := pathx.MkdirIfNotExist(grpc)
		if err != nil {
			t.Log(err)
			return
		}
		cmd := exec.Command("protoc", "-I="+filepath.Dir(protoFile), "--go_out="+output,
			"--go-grpc_out="+grpc, filepath.Base(protoFile), "--go_opt=M"+filepath.Base(protoFile)+"="+pbeRedirect,
			"--go-grpc_opt=M"+filepath.Base(protoFile)+"="+grpcRedirect)
		cmd.Dir = output
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			t.Log(err)
			return
		}
		pbDir, err := findPbFile(output, protoFile, false)
		assert.Nil(t, err)
		pbGo := filepath.Join(pbDir, "greet.pb.go")
		assert.True(t, pathx.FileExists(pbGo))

		grpcDir, err := findPbFile(output, protoFile, true)
		assert.Nil(t, err)
		grpcGo := filepath.Join(grpcDir, "greet_grpc.pb.go")
		assert.True(t, pathx.FileExists(grpcGo))
	})

	t.Run("", func(t *testing.T) {
		output := t.TempDir()
		pbeRedirect := filepath.Join(output, "redirect")
		grpc := filepath.Join(output, "grpc")
		grpcRedirect := filepath.Join(grpc, "redirect")
		err := pathx.MkdirIfNotExist(grpc)
		if err != nil {
			t.Log(err)
			return
		}
		cmd := exec.Command("protoc", "-I="+filepath.Dir(protoFile), "--go_out="+output,
			"--go-grpc_out="+grpc, filepath.Base(protoFile), "--go_opt=M"+filepath.Base(protoFile)+"="+pbeRedirect,
			"--go-grpc_opt=M"+filepath.Base(protoFile)+"="+grpcRedirect, "--go_opt=paths=import", "--go-grpc_opt=paths=source_relative")
		cmd.Dir = output
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			t.Log(err)
			return
		}
		pbDir, err := findPbFile(output, protoFile, false)
		assert.Nil(t, err)
		pbGo := filepath.Join(pbDir, "greet.pb.go")
		assert.True(t, pathx.FileExists(pbGo))

		grpcDir, err := findPbFile(output, protoFile, true)
		assert.Nil(t, err)
		grpcGo := filepath.Join(grpcDir, "greet_grpc.pb.go")
		assert.True(t, pathx.FileExists(grpcGo))
	})

	t.Run("", func(t *testing.T) {
		output := t.TempDir()
		pbeRedirect := filepath.Join(output, "redirect")
		grpc := filepath.Join(output, "grpc")
		grpcRedirect := filepath.Join(grpc, "redirect")
		err := pathx.MkdirIfNotExist(grpc)
		if err != nil {
			t.Log(err)
			return
		}
		err = pathx.MkdirIfNotExist(pbeRedirect)
		if err != nil {
			t.Log(err)
			return
		}
		err = pathx.MkdirIfNotExist(grpcRedirect)
		if err != nil {
			t.Log(err)
			return
		}
		cmd := exec.Command("protoc", "-I="+filepath.Dir(protoFile), "--go_out="+output,
			"--go-grpc_out="+grpc, filepath.Base(protoFile), "--go_opt=M"+filepath.Base(protoFile)+"="+pbeRedirect,
			"--go-grpc_opt=M"+filepath.Base(protoFile)+"="+grpcRedirect, "--go_opt=paths=import", "--go-grpc_opt=paths=source_relative",
			"--go_out="+pbeRedirect, "--go-grpc_out="+grpcRedirect)
		cmd.Dir = output
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			t.Log(err)
			return
		}
		pbDir, err := findPbFile(output, protoFile, false)
		assert.Nil(t, err)
		pbGo := filepath.Join(pbDir, "greet.pb.go")
		assert.True(t, pathx.FileExists(pbGo))

		grpcDir, err := findPbFile(output, protoFile, true)
		assert.Nil(t, err)
		grpcGo := filepath.Join(grpcDir, "greet_grpc.pb.go")
		assert.True(t, pathx.FileExists(grpcGo))
	})
}
