package codes

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Acceptable(err error) bool {
	switch status.Code(err) {
	case codes.DeadlineExceeded, codes.Internal, codes.Unavailable, codes.DataLoss:
		return false
	default:
		return true
	}
}
