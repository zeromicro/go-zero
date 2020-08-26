package ctx

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/tal-tech/go-zero/tools/goctl/rpc/execx"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/project"
	"github.com/tal-tech/go-zero/tools/goctl/util"
)

var (
	protoGenGo        = "protoc-gen-go"
	errNotFound       = errors.New("protoc-gen-go is not found")
	errPluginNotFound = errors.New("protobuf is not found in go mod,please make sure you have import the protobuf")
)

type GoMod struct {
	// protobuf path: such as github.com/golang/protobuf@1.2.3
	protobuf string
	// module: such as github.com/tal-tech/go-zero
	module string
}

func prepare() (*GoMod, error) {
	_, err := exec.LookPath("go")
	if err != nil {
		return nil, err
	}
	_, err = exec.LookPath("protoc")
	if err != nil {
		return nil, err
	}
	info, err := project.Info()
	if err != nil {
		return nil, err
	}

	// find in go mod,if go version >= 1.5
	if info.IsGoModProject {
		data, err := ioutil.ReadFile(info.GoModFilePath)
		if err != nil {
			return nil, err
		}
		protobufRelative, err := matchProtocGenGo(data)
		switch err {
		case nil:
			protobufDir := filepath.Join(info.LibraryPath, protobufRelative)
			if util.FileExists(protobufDir) {
				plugin := filepath.Join(protobufDir, protoGenGo)
				if util.FileExists(plugin) {
					module, err := matchModule(data)
					if err != nil {
						return nil, err
					}
					return &GoMod{protobuf: plugin, module: module}, nil
				}
			}
			return nil, errPluginNotFound
		case errNotFound:
			return nil, errPluginNotFound
		default:
			return nil, err
		}
	}
	// else: find in go path
	protobufDir := filepath.Join(info.LibraryPath, "github.com", "golang", "protobuf")
	if util.FileExists(protobufDir) {
		plugin := filepath.Join(protobufDir, protoGenGo)
		if util.FileExists(plugin) {
			return &GoMod{protobuf: plugin, module: info.Name}, nil
		}
	}
	// else: go get latest protobuf
	sh := "go get -u github.com/golang/protobuf/protoc-gen-go"
	doWithTimeOut(30*time.Second, func() {
		err = execx.RunShOrBat(sh)
	}, func() {
		err = fmt.Errorf("timeout: %s", sh)
	})
	if err == nil {
		if util.FileExists(protobufDir) {
			plugin := filepath.Join(protobufDir, protoGenGo)
			if util.FileExists(plugin) {
				return &GoMod{protobuf: plugin, module: info.Name}, nil
			}
		}
		return nil, errPluginNotFound
	}
	return nil, err
}

// github.com/golang/protobuf@{version}
func matchProtocGenGo(data []byte) (string, error) {
	text := string(data)
	re := regexp.MustCompile(`(?m)(github.com/golang/protobuf)\s+(v[0-9.]+)`)
	matches := re.FindAllStringSubmatch(text, -1)
	if len(matches) == 0 {
		return "", errNotFound
	}
	groups := matches[0]
	if len(groups) < 3 {
		return "", errNotFound
	}
	return fmt.Sprintf("%s@%s", groups[1], groups[2]), nil
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
func (p *GoMod) Protobuf() string {
	return p.protobuf
}

func (p *GoMod) Install() error {
	_, err := execx.Run("go", fmt.Sprintf("install %v", p))
	return err
}

func doWithTimeOut(duration time.Duration, handleFunc func(), timeOutFunc func()) {
	doneChan := make(chan int)
	once := &sync.Once{}
	go func(ch chan int, o *sync.Once) {
		handleFunc()
		o.Do(func() {
			close(doneChan)
		})
	}(doneChan, once)
	timer := time.NewTimer(duration)
	select {
	case <-timer.C:
		timeOutFunc()
		timer.Stop()
		return
	case <-doneChan:
		return
	default:
	}
}
