package errorx

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

var errorFormat = `goctl: generation error: %+v
goctl version: %s
%s`

type GoctlError struct {
	message []string
	err     error
}

func (e *GoctlError) Error() string {
	buildVersion := os.Getenv("GOCTL_VERSION")
	detail := wrapMessage(e.message...)
	version := fmt.Sprintf("%s %s/%s", buildVersion, runtime.GOOS, runtime.GOARCH)
	return fmt.Sprintf(errorFormat, e.err, version, detail)
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
