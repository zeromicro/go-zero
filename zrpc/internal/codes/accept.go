package codes

import (
	"context"
	"github.com/tal-tech/go-zero/core/lang"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// rejectErr A collection of errors that should be rejected.
	rejectErr = make(map[error]lang.PlaceholderType)
	// one Make sure to be executed only once.
	one sync.Once
)

// RejectErr Setting should be rejected error.
func RejectErr(errs ...error) {
	one.Do(func() {
		for _, err := range errs {
			rejectErr[err] = lang.PlaceholderType{}
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
	switch err {
	case context.DeadlineExceeded:
		return false
	default:
		_, ok := rejectErr[err]
		return !ok
	}
}
