package ctx

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/rpc/execx"
	"github.com/tal-tech/go-zero/tools/goctl/util/console"
)

const (
	constGo          = "go"
	constProtoC      = "protoc"
	constGoMod       = "go env GOMOD"
	constGoPath      = "go env GOPATH"
	constProtoCGenGo = "protoc-gen-go"
)

type (
	Project struct {
		Path  string
		Name  string
		GoMod GoMod
	}

	GoMod struct {
		Module string
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
	_, err = exec.LookPath(constProtoCGenGo)
	if err != nil {
		return nil, err
	}

	var (
		goMod, module string
		goPath        string
		name, path    string
	)

	ret, err := execx.Run(constGoMod)
	if err != nil {
		return nil, err
	}
	goMod = strings.TrimSpace(ret)

	ret, err = execx.Run(constGoPath)
	if err != nil {
		return nil, err
	}

	goPath = strings.TrimSpace(ret)
	src := filepath.Join(goPath, "src")
	if len(goMod) > 0 {
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

	return &Project{
		Name: name,
		Path: path,
		GoMod: GoMod{
			Module: module,
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
