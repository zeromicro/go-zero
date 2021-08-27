package codes

import (
	"github.com/tal-tech/go-zero/core/lang"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// breakerErrors storage rejected errors.
var breakerErrors = make(map[error]lang.PlaceholderType)

// once only allow breakerErrors to be added once.
var once sync.Once

// AddBreakerErrors add customized breaker errors for the client and server.
func AddBreakerErrors(errs ...error) {
	once.Do(func() {
		for _, err := range errs {
			breakerErrors[err] = lang.Placeholder
		}
	})
}

// Acceptable checks if given error is acceptable.
func Acceptable(err error) bool {
	switch status.Code(err) {
	case codes.DeadlineExceeded, codes.Internal, codes.Unavailable, codes.DataLoss:
		return false
	case codes.Unknown:
		return acceptableUnknown(err)
	default:
		return true
	}
}

func acceptableUnknown(err error) bool {
	_, ok := breakerErrors[err]
	return !ok
}
