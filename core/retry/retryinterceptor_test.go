package retry

import (
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestDo(t *testing.T) {
	n := 4
	for i := 0; i < n; i++ {
		count := 0
		err := Do(context.Background(), func(ctx context.Context, opts ...grpc.CallOption) error {
			count++
			return status.Error(codes.ResourceExhausted, "ResourceExhausted")

		}, WithMax(i))
		assert.Error(t, err)
		assert.Equal(t, i+1, count)
	}

}
