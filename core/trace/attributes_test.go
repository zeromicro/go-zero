package trace

import (
	"testing"

	"github.com/stretchr/testify/assert"
	gcodes "google.golang.org/grpc/codes"
)

func TestStatusCodeAttr(t *testing.T) {
	assert.Equal(t, GRPCStatusCodeKey.Int(int(gcodes.DataLoss)), StatusCodeAttr(gcodes.DataLoss))
}
