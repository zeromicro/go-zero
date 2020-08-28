package ctx

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/rpc/execx"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/util/console"
)

var errProtobufNotFound = errors.New("github.com/golang/protobuf is not found,please ensure you has already [go get github.com/golang/protobuf]")

const (
	constGo          = "go"
	constProtoC      = "protoc"
	constGoModOn     = "go env GO111MODULE"
	constGoMod       = "go env GOMOD"
	constGoModCache  = "go env GOMODCACHE"
	constGoPath      = "go env GOPATH"
	constProtoCGenGo = "protoc-gen-go"
)

type (
	Project struct {
		Path     string
		Name     string
		GoPath   string
		Protobuf Protobuf
		GoMod    GoMod
	}

	GoMod struct {
		ModOn      bool
		GoModCache string
		GoMod      string
		Module     string
	}
	Protobuf struct {
		Path string
	}
)

func prepare(log console.Console) (*Project, error) {
	log.Info("checking go env...")
	_, err := exec.LookPath(constGo)
	if err != nil {
		return nil, err
	}

	_, err = exec.LookPath(constProtoC)
	if err != nil {
		return nil, err
	}

	var (
		goModOn                   bool
		goMod, goModCache, module string
		goPath                    string
		name, path                string
		protobufModule            string
	)
	ret, err := execx.Run(constGoModOn)
	if err != nil {
		return nil, err
	}

	goModOn = strings.TrimSpace(ret) == "on"
	ret, err = execx.Run(constGoMod)
	if err != nil {
		return nil, err
	}

	goMod = strings.TrimSpace(ret)
	ret, err = execx.Run(constGoModCache)
	if err != nil {
		return nil, err
	}

	goModCache = strings.TrimSpace(ret)
	ret, err = execx.Run(constGoPath)
	if err != nil {
		return nil, err
	}

	goPath = strings.TrimSpace(ret)
	src := filepath.Join(goPath, "src")
	if len(goMod) > 0 {
		if goModCache == "" {
			goModCache = filepath.Join(goPath, "pkg", "mod")
		}
		path = filepath.Dir(goMod)
		name = filepath.Base(path)
		data, err := ioutil.ReadFile(goMod)
		if err != nil {
			return nil, err
		}

		module, err = matchModule(data)
		if err != nil {
			return nil, err
		}
	} else {
		if goModCache == "" {
			goModCache = src
		}
		pwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		if !strings.HasPrefix(pwd, src) {
			return nil, fmt.Errorf("%s: project is not in go mod and go path", pwd)
		}
		r := strings.TrimPrefix(pwd, src+string(filepath.Separator))
		name = filepath.Dir(r)
		if name == "." {
			name = r
		}
		path = filepath.Join(src, name)
		module = name
	}

	protobuf := filepath.Join(goModCache, protobufModule)
	if !util.FileExists(protobuf) {
		return nil, fmt.Errorf("expected protobuf module in path: %s,please ensure you has already [go get github.com/golang/protobuf]", protobuf)
	}

	var protoCGenGoFilename string
	os := runtime.GOOS
	switch os {
	case "darwin":
		protoCGenGoFilename = filepath.Join(goPath, "bin", "protoc-gen-go")
	case "windows":
		protoCGenGoFilename = filepath.Join(goPath, "bin", "protoc-gen-go.exe")
	default:
		return nil, fmt.Errorf("unexpeted os: %s", os)
	}

	if !util.FileExists(protoCGenGoFilename) {
		sh := "go install " + filepath.Join(protobuf, constProtoCGenGo)
		log.Warning(sh)
		stdout, err := execx.Run(sh)
		if err != nil {
			return nil, err
		}

		log.Info(stdout)
	}
	if !util.FileExists(protoCGenGoFilename) {
		return nil, fmt.Errorf("protoc-gen-go is not found")
	}
	return &Project{
		Name:   name,
		Path:   path,
		GoPath: goPath,
		Protobuf: Protobuf{
			Path: protobuf,
		},
		GoMod: GoMod{
			ModOn:      goModOn,
			GoModCache: goModCache,
			GoMod:      goMod,
			Module:     module,
		},
	}, nil
}

func matchModule(data []byte) (string, error) {
	text := string(data)
	re := regexp.MustCompile(`(?m)^\s*module\s+[a-z0-9/\-.]+$`)
	matches := re.FindAllString(text, -1)
	if len(matches) == 1 {
		target := matches[0]
		index := strings.Index(target, "module")
		return strings.TrimSpace(target[index+6:]), nil
	}

	return "", nil
}
