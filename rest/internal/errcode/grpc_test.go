package errcode

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCodeFromGrpcError(t *testing.T) {
	tests := []struct {
		name string
		code codes.Code
		want int
	}{
		{
			name: "OK",
			code: codes.OK,
			want: http.StatusOK,
		},
		{
			name: "Invalid argument",
			code: codes.InvalidArgument,
			want: http.StatusBadRequest,
		},
		{
			name: "Failed precondition",
			code: codes.FailedPrecondition,
			want: http.StatusBadRequest,
		},
		{
			name: "Out of range",
			code: codes.OutOfRange,
			want: http.StatusBadRequest,
		},
		{
			name: "Unauthorized",
			code: codes.Unauthenticated,
			want: http.StatusUnauthorized,
		},
		{
			name: "Permission denied",
			code: codes.PermissionDenied,
			want: http.StatusForbidden,
		},
		{
			name: "Not found",
			code: codes.NotFound,
			want: http.StatusNotFound,
		},
		{
			name: "Canceled",
			code: codes.Canceled,
			want: http.StatusRequestTimeout,
		},
		{
			name: "Already exists",
			code: codes.AlreadyExists,
			want: http.StatusConflict,
		},
		{
			name: "Aborted",
			code: codes.Aborted,
			want: http.StatusConflict,
		},
		{
			name: "Resource exhausted",
			code: codes.ResourceExhausted,
			want: http.StatusTooManyRequests,
		},
		{
			name: "Internal",
			code: codes.Internal,
			want: http.StatusInternalServerError,
		},
		{
			name: "Data loss",
			code: codes.DataLoss,
			want: http.StatusInternalServerError,
		},
		{
			name: "Unknown",
			code: codes.Unknown,
			want: http.StatusInternalServerError,
		},
		{
			name: "Unimplemented",
			code: codes.Unimplemented,
			want: http.StatusNotImplemented,
		},
		{
			name: "Unavailable",
			code: codes.Unavailable,
			want: http.StatusServiceUnavailable,
		},
		{
			name: "Deadline exceeded",
			code: codes.DeadlineExceeded,
			want: http.StatusGatewayTimeout,
		},
		{
			name: "Beyond defined error",
			code: codes.Code(^uint32(0)),
			want: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, CodeFromGrpcError(status.Error(test.code, "foo")))
		})
	}
}

func TestIsGrpcError(t *testing.T) {
	assert.True(t, IsGrpcError(status.Error(codes.Unknown, "foo")))
	assert.False(t, IsGrpcError(errors.New("foo")))
	assert.False(t, IsGrpcError(nil))
}
