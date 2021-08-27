package codes

import (
	"github.com/tal-tech/go-zero/core/lang"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// rejectErrors storage rejected errors.
var rejectErrors = make(map[error]lang.PlaceholderType)

// once only allow rejectErrors to be added once.
var once sync.Once

// AddRejectBreakerErrors add customized breaker errors for the client and server.
func AddRejectBreakerErrors(errs ...error) {
	once.Do(func() {
		for _, err := range errs {
			rejectErrors[err] = lang.Placeholder
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
	_, ok := rejectErrors[err]
	return !ok
}
