package clientinterceptors

import (
	"testing"

	"zero/core/breaker"
	"zero/core/stat"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func init() {
	stat.SetReporter(nil)
}

type mockError struct {
	st *status.Status
}

func (m mockError) GRPCStatus() *status.Status {
	return m.st
}

func (m mockError) Error() string {
	return "mocked error"
}

func TestBreakerInterceptorNotFound(t *testing.T) {
	err := mockError{st: status.New(codes.NotFound, "any")}
	for i := 0; i < 1000; i++ {
		assert.Equal(t, err, breaker.DoWithAcceptable("call", func() error {
			return err
		}, acceptable))
	}
}

func TestBreakerInterceptorDeadlineExceeded(t *testing.T) {
	err := mockError{st: status.New(codes.DeadlineExceeded, "any")}
	errs := make(map[error]int)
	for i := 0; i < 1000; i++ {
		e := breaker.DoWithAcceptable("call", func() error {
			return err
		}, acceptable)
		errs[e]++
	}
	assert.Equal(t, 2, len(errs))
	assert.True(t, errs[err] > 0)
	assert.True(t, errs[breaker.ErrServiceUnavailable] > 0)
}
