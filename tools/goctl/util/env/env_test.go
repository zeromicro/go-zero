package env

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"github.com/zeromicro/go-zero/tools/goctl/vars"
)

func TestLookUpGo(t *testing.T) {
	xGo, err := LookUpGo()
	if err != nil {
		return
	}

	assert.True(t, pathx.FileExists(xGo))
	output, errOutput, err := execCommand(xGo, "version")
	if err != nil {
		return
	}

	if len(errOutput) > 0 {
		return
	}
	assert.Equal(t, wrapVersion(), output)
}

func TestLookUpProtoc(t *testing.T) {
	xProtoc, err := LookUpProtoc()
	if err != nil {
		return
	}

	assert.True(t, pathx.FileExists(xProtoc))
	output, errOutput, err := execCommand(xProtoc, "--version")
	if err != nil {
		return
	}

	if len(errOutput) > 0 {
		return
	}
	assert.True(t, len(output) > 0)
}

func TestLookUpProtocGenGo(t *testing.T) {
	xProtocGenGo, err := LookUpProtocGenGo()
	if err != nil {
		return
	}
	assert.True(t, pathx.FileExists(xProtocGenGo))
}

func TestLookPath(t *testing.T) {
	xGo, err := LookPath("go")
	if err != nil {
		return
	}
	assert.True(t, pathx.FileExists(xGo))
}

func TestCanExec(t *testing.T) {
	canExec := runtime.GOOS != vars.OsJs && runtime.GOOS != vars.OsIOS
	assert.Equal(t, canExec, CanExec())
}

func execCommand(cmd string, arg ...string) (stdout, stderr string, err error) {
	output := bytes.NewBuffer(nil)
	errOutput := bytes.NewBuffer(nil)
	c := exec.Command(cmd, arg...)
	c.Stdout = output
	c.Stderr = errOutput
	err = c.Run()
	if err != nil {
		return
	}
	if errOutput.Len() > 0 {
		stderr = errOutput.String()
		return
	}
	stdout = strings.TrimSpace(output.String())
	return
}

func wrapVersion() string {
	version := runtime.Version()
	os := runtime.GOOS
	arch := runtime.GOARCH
	return fmt.Sprintf("go version %s %s/%s", version, os, arch)
}
