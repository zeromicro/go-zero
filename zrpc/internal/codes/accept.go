package codes

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Acceptable checks if given error is acceptable.
func Acceptable(err error) bool {
	switch status.Code(err) {
	case codes.DeadlineExceeded, codes.Internal, codes.Unavailable, codes.DataLoss,
		codes.Unimplemented, codes.ResourceExhausted:
		return false
	default:
		return true
	}
}
