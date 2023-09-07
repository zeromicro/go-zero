package errorx

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewInternalError returns status error with Internal error code.
func NewInternalError(msg string) error {
	return status.Error(codes.Internal, msg)
}

// NewInvalidArgumentError returns status error with InvalidArgument error code.
func NewInvalidArgumentError(msg string) error {
	return status.Error(codes.InvalidArgument, msg)
}

// NewNotFoundError returns status error with NotFound error code.
func NewNotFoundError(msg string) error {
	return status.Error(codes.NotFound, msg)
}

// NewAlreadyExistsError returns status error with AlreadyExists error code.
func NewAlreadyExistsError(msg string) error {
	return status.Error(codes.AlreadyExists, msg)
}

// NewUnauthenticatedError returns status error with Unauthenticated error code.
func NewUnauthenticatedError(msg string) error {
	return status.Error(codes.Unauthenticated, msg)
}

// NewResourceExhaustedError returns status error with ResourceExhausted error code.
func NewResourceExhaustedError(msg string) error {
	return status.Error(codes.ResourceExhausted, msg)
}
