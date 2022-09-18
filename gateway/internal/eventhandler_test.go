package internal

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestEventHandler(t *testing.T) {
	h := NewEventHandler(io.Discard, nil)
	h.OnResolveMethod(nil)
	h.OnSendHeaders(nil)
	h.OnReceiveHeaders(nil)
	h.OnReceiveTrailers(status.New(codes.OK, ""), nil)
	assert.Equal(t, codes.OK, h.Status.Code())
	h.OnReceiveResponse(nil)
}
