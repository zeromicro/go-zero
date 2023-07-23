package env

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/tools/goctl/vars"
)

const (
	bin                = "bin"
	binGo              = "go"
	binProtoc          = "protoc"
	binProtocGenGo     = "protoc-gen-go"
	binProtocGenGrpcGo = "protoc-gen-go-grpc"
	cstOffset          = 60 * 60 * 8 // 8 hours offset for Chinese Standard Time
)

// InChina returns whether the current time is in China Standard Time.
func InChina() bool {
	_, offset := time.Now().Zone()
	return offset == cstOffset
}

// LookUpGo searches an executable go in the directories
// named by the GOROOT/bin or PATH environment variable.
func LookUpGo() (string, error) {
	goRoot := runtime.GOROOT()
	suffix := getExeSuffix()
	xGo := binGo + suffix
	path := filepath.Join(goRoot, bin, xGo)
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}
	return LookPath(xGo)
}

// LookUpProtoc searches an executable protoc in the directories
// named by the PATH environment variable.
func LookUpProtoc() (string, error) {
	suffix := getExeSuffix()
	xProtoc := binProtoc + suffix
	return LookPath(xProtoc)
}

// LookUpProtocGenGo searches an executable protoc-gen-go in the directories
// named by the PATH environment variable.
func LookUpProtocGenGo() (string, error) {
	suffix := getExeSuffix()
	xProtocGenGo := binProtocGenGo + suffix
	return LookPath(xProtocGenGo)
}

// LookUpProtocGenGoGrpc searches an executable protoc-gen-go-grpc in the directories
// named by the PATH environment variable.
func LookUpProtocGenGoGrpc() (string, error) {
	suffix := getExeSuffix()
	xProtocGenGoGrpc := binProtocGenGrpcGo + suffix
	return LookPath(xProtocGenGoGrpc)
}

// LookPath searches for an executable named file in the
// directories named by the PATH environment variable,
// for the os windows, the named file will be spliced with the
// .exe suffix.
func LookPath(xBin string) (string, error) {
	suffix := getExeSuffix()
	if len(suffix) > 0 && !strings.HasSuffix(xBin, suffix) {
		xBin = xBin + suffix
	}

	bin, err := exec.LookPath(xBin)
	if err != nil {
		return "", err
	}
	return bin, nil
}

// CanExec reports whether the current system can start new processes
// using os.StartProcess or (more commonly) exec.Command.
func CanExec() bool {
	switch runtime.GOOS {
	case vars.OsJs, vars.OsIOS:
		return false
	default:
		return true
	}
}

func getExeSuffix() string {
	if runtime.GOOS == vars.OsWindows {
		return ".exe"
	}
	return ""
}
