package errorx

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/internal/version"
)

var errorFormat = `goctl: generation error: %+v
goctl version: %s
%s`

// GoctlError represents a goctl error.
type GoctlError struct {
	message []string
	err     error
}

func (e *GoctlError) Error() string {
	detail := wrapMessage(e.message...)
	v := fmt.Sprintf("%s %s/%s", version.BuildVersion, runtime.GOOS, runtime.GOARCH)
	return fmt.Sprintf(errorFormat, e.err, v, detail)
}

// Wrap wraps an error with goctl version and message.
func Wrap(err error, message ...string) error {
	e, ok := err.(*GoctlError)
	if ok {
		return e
	}

	return &GoctlError{
		message: message,
		err:     err,
	}
}

func wrapMessage(message ...string) string {
	if len(message) == 0 {
		return ""
	}
	return fmt.Sprintf(`message: %s`, strings.Join(message, "\n"))
}
