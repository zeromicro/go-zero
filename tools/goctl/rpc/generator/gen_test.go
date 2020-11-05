package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/execx"
	"github.com/tal-tech/go-zero/tools/goctl/util"
)

func TestRpcGenerateCaseNilImport(t *testing.T) {
	_ = Clean()
	dispatcher := NewDefaultGenerator()
	if err := dispatcher.Prepare(); err == nil {
		g := NewRpcGenerator(dispatcher, namingLower)
		abs, err := filepath.Abs("./test")
		assert.Nil(t, err)

		err = g.Generate("./test_stream.proto", abs, nil)
		defer func() {
			_ = os.RemoveAll(abs)
		}()
		assert.Nil(t, err)

		_, err = execx.Run("go test "+abs, abs)
		assert.Nil(t, err)
	}
}

func TestRpcGenerateCaseOption(t *testing.T) {
	_ = Clean()
	dispatcher := NewDefaultGenerator()
	if err := dispatcher.Prepare(); err == nil {
		g := NewRpcGenerator(dispatcher, namingLower)
		abs, err := filepath.Abs("./test")
		assert.Nil(t, err)

		err = g.Generate("./test_option.proto", abs, nil)
		defer func() {
			_ = os.RemoveAll(abs)
		}()
		assert.Nil(t, err)

		_, err = execx.Run("go test "+abs, abs)
		assert.Nil(t, err)
	}
}

func TestRpcGenerateCaseWordOption(t *testing.T) {
	_ = Clean()
	dispatcher := NewDefaultGenerator()
	if err := dispatcher.Prepare(); err == nil {
		g := NewRpcGenerator(dispatcher, namingLower)
		abs, err := filepath.Abs("./test")
		assert.Nil(t, err)

		err = g.Generate("./test_word_option.proto", abs, nil)
		defer func() {
			_ = os.RemoveAll(abs)
		}()
		assert.Nil(t, err)

		_, err = execx.Run("go test "+abs, abs)
		assert.Nil(t, err)
	}
}

// test keyword go
func TestRpcGenerateCaseGoOption(t *testing.T) {
	_ = Clean()
	dispatcher := NewDefaultGenerator()
	if err := dispatcher.Prepare(); err == nil {
		g := NewRpcGenerator(dispatcher, namingLower)
		abs, err := filepath.Abs("./test")
		assert.Nil(t, err)

		err = g.Generate("./test_go_option.proto", abs, nil)
		defer func() {
			_ = os.RemoveAll(abs)
		}()
		assert.Nil(t, err)

		_, err = execx.Run("go test "+abs, abs)
		assert.Nil(t, err)
	}
}

func TestRpcGenerateCaseImport(t *testing.T) {
	_ = Clean()
	dispatcher := NewDefaultGenerator()
	if err := dispatcher.Prepare(); err == nil {
		g := NewRpcGenerator(dispatcher, namingLower)
		abs, err := filepath.Abs("./test")
		assert.Nil(t, err)

		err = g.Generate("./test_import.proto", abs, []string{"./base"})
		defer func() {
			_ = os.RemoveAll(abs)
		}()
		assert.Nil(t, err)

		_, err = execx.Run("go test "+abs, abs)
		assert.True(t, func() bool {
			return strings.Contains(err.Error(), "package base is not in GOROOT")
		}())
	}
}

func TestRpcGenerateCaseMultipleIdenticalReturn(t *testing.T) {
	dispatcher := NewDefaultGenerator()
	if err := dispatcher.Prepare(); err == nil {
		g := NewRpcGenerator(dispatcher, namingLower)
		abs, err := filepath.Abs("./test")
		assert.Nil(t, err)

		err = g.Generate("./test_multiple_identical_return.proto", abs, nil)
		defer func() {
			_ = os.RemoveAll(abs)
		}()
		assert.Nil(t, err)

		_, err = execx.Run("go test "+abs, abs)
		assert.Nil(t, err)
	}
}

func TestRpcGenerateCaseNamingEmptyString(t *testing.T) {
	dispatcher := NewDefaultGenerator()
	if err := dispatcher.Prepare(); err == nil {
		// test case empty string
		g := NewRpcGenerator(dispatcher, "")
		abs, err := filepath.Abs("./test")
		assert.Nil(t, err)

		err = g.Generate("./test_stream.proto", abs, nil)
		defer func() {
			_ = os.RemoveAll(abs)
		}()
		assert.Nil(t, err)

		_, err = execx.Run("go test "+abs, abs)
		assert.Nil(t, err)

		assert.True(t, func() bool {
			var ret = true
			etcFilename := filepath.Join(abs, "etc", "test.yaml")
			ret = ret && util.FileExists(etcFilename)
			configFilename := filepath.Join(abs, "internal", "config", "config.go")
			ret = ret && util.FileExists(configFilename)
			logicFilename := filepath.Join(abs, "internal", "logic", "greetlogic.go")
			ret = ret && util.FileExists(logicFilename)
			serverFilename := filepath.Join(abs, "internal", "server", "streamgreeterserver.go")
			ret = ret && util.FileExists(serverFilename)
			pbFilename := filepath.Join(abs, "internal", "stream", "test_stream.pb.go")
			ret = ret && util.FileExists(pbFilename)
			svcFilename := filepath.Join(abs, "internal", "svc", "servicecontext.go")
			ret = ret && util.FileExists(svcFilename)
			callFilename := filepath.Join(abs, "streamgreeter", "streamgreeter.go")
			ret = ret && util.FileExists(callFilename)
			mainFilename := filepath.Join(abs, "test.go")
			ret = ret && util.FileExists(mainFilename)
			return ret
		}())

	}
}

func TestRpcGenerateCaseNamingLower(t *testing.T) {
	dispatcher := NewDefaultGenerator()
	if err := dispatcher.Prepare(); err == nil {
		// test case namingLower
		g := NewRpcGenerator(dispatcher, namingLower)
		abs, err := filepath.Abs("./test")
		assert.Nil(t, err)

		err = g.Generate("./test_stream.proto", abs, nil)
		defer func() {
			_ = os.RemoveAll(abs)
		}()
		assert.Nil(t, err)

		_, err = execx.Run("go test "+abs, abs)
		assert.Nil(t, err)

		assert.True(t, func() bool {
			var ret = true
			etcFilename := filepath.Join(abs, "etc", "Test.yaml")
			ret = ret && util.FileExists(etcFilename)
			configFilename := filepath.Join(abs, "internal", "config", "config.go")
			ret = ret && util.FileExists(configFilename)
			logicFilename := filepath.Join(abs, "internal", "logic", "greetlogic.go")
			ret = ret && util.FileExists(logicFilename)
			serverFilename := filepath.Join(abs, "internal", "server", "streamgreeterserver.go")
			ret = ret && util.FileExists(serverFilename)
			pbFilename := filepath.Join(abs, "internal", "stream", "test_stream.pb.go")
			ret = ret && util.FileExists(pbFilename)
			svcFilename := filepath.Join(abs, "internal", "svc", "servicecontext.go")
			ret = ret && util.FileExists(svcFilename)
			callFilename := filepath.Join(abs, "streamgreeter", "streamgreeter.go")
			ret = ret && util.FileExists(callFilename)
			mainFilename := filepath.Join(abs, "test.go")
			ret = ret && util.FileExists(mainFilename)
			return ret
		}())

	}
}

func TestRpcGenerateCaseNamingCamel(t *testing.T) {
	dispatcher := NewDefaultGenerator()
	if err := dispatcher.Prepare(); err == nil {
		// test case namingCamel
		g := NewRpcGenerator(dispatcher, namingLower)
		abs, err := filepath.Abs("./test")
		assert.Nil(t, err)

		err = g.Generate("./test_stream.proto", abs, nil)
		defer func() {
			_ = os.RemoveAll(abs)
		}()
		assert.Nil(t, err)

		_, err = execx.Run("go test "+abs, abs)
		assert.Nil(t, err)

		assert.True(t, func() bool {
			var ret = true
			etcFilename := filepath.Join(abs, "etc", "Test.yaml")
			ret = ret && util.FileExists(etcFilename)
			configFilename := filepath.Join(abs, "internal", "config", "Config.go")
			ret = ret && util.FileExists(configFilename)
			logicFilename := filepath.Join(abs, "internal", "logic", "GreetLogic.go")
			ret = ret && util.FileExists(logicFilename)
			serverFilename := filepath.Join(abs, "internal", "server", "StreamGreeterServer.go")
			ret = ret && util.FileExists(serverFilename)
			pbFilename := filepath.Join(abs, "internal", "stream", "test_stream.pb.go")
			ret = ret && util.FileExists(pbFilename)
			svcFilename := filepath.Join(abs, "internal", "svc", "ServiceContext.go")
			ret = ret && util.FileExists(svcFilename)
			callFilename := filepath.Join(abs, "streamgreeter", "StreamGreeter.go")
			ret = ret && util.FileExists(callFilename)
			mainFilename := filepath.Join(abs, "Test.go")
			ret = ret && util.FileExists(mainFilename)
			return ret
		}())
	}
}

func TestRpcGenerateCaseNamingSnake(t *testing.T) {
	dispatcher := NewDefaultGenerator()
	if err := dispatcher.Prepare(); err == nil {
		// test case namingSnake
		g := NewRpcGenerator(dispatcher, namingSnake)
		abs, err := filepath.Abs("./test")
		assert.Nil(t, err)

		err = g.Generate("./test_stream.proto", abs, nil)
		defer func() {
			//_ = os.RemoveAll(abs)
		}()
		assert.Nil(t, err)

		_, err = execx.Run("go test "+abs, abs)
		assert.Nil(t, err)

		assert.True(t, func() bool {
			var ret = true
			etcFilename := filepath.Join(abs, "etc", "test.yaml")
			ret = ret && util.FileExists(etcFilename)
			configFilename := filepath.Join(abs, "internal", "config", "config.go")
			ret = ret && util.FileExists(configFilename)
			logicFilename := filepath.Join(abs, "internal", "logic", "greet_logic.go")
			ret = ret && util.FileExists(logicFilename)
			serverFilename := filepath.Join(abs, "internal", "server", "stream_greeter_server.go")
			ret = ret && util.FileExists(serverFilename)
			pbFilename := filepath.Join(abs, "internal", "stream", "test_stream.pb.go")
			ret = ret && util.FileExists(pbFilename)
			svcFilename := filepath.Join(abs, "internal", "svc", "service_context.go")
			ret = ret && util.FileExists(svcFilename)
			callFilename := filepath.Join(abs, "streamgreeter", "stream_greeter.go")
			ret = ret && util.FileExists(callFilename)
			mainFilename := filepath.Join(abs, "test.go")
			ret = ret && util.FileExists(mainFilename)
			return ret
		}())

	}
}

func TestRpcGenerateCaseServiceRpcNamingSnake(t *testing.T) {
	_ = Clean()
	dispatcher := NewDefaultGenerator()
	if err := dispatcher.Prepare(); err == nil {
		g := NewRpcGenerator(dispatcher)
		abs, err := filepath.Abs("./test")
		assert.Nil(t, err)

		err = g.Generate("./test_service_rpc_naming_snake.proto", abs, nil)
		defer func() {
			_ = os.RemoveAll(abs)
		}()
		assert.Nil(t, err)

		_, err = execx.Run("go test "+abs, abs)
		assert.Nil(t, err)
	}
}
