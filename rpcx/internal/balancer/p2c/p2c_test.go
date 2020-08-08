package p2c

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/mathx"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/status"
)

func init() {
	logx.Disable()
}

func TestP2cPicker_PickNil(t *testing.T) {
	builder := new(p2cPickerBuilder)
	picker := builder.Build(nil)
	_, _, err := picker.Pick(context.Background(), balancer.PickInfo{
		FullMethodName: "/",
		Ctx:            context.Background(),
	})
	assert.NotNil(t, err)
}

func TestP2cPicker_Pick(t *testing.T) {
	tests := []struct {
		name       string
		candidates int
	}{
		{
			name:       "single",
			candidates: 1,
		},
		{
			name:       "multiple",
			candidates: 100,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			builder := new(p2cPickerBuilder)
			ready := make(map[resolver.Address]balancer.SubConn)
			for i := 0; i < test.candidates; i++ {
				ready[resolver.Address{
					Addr: strconv.Itoa(i),
				}] = new(mockClientConn)
			}

			picker := builder.Build(ready)
			for i := 0; i < 10000; i++ {
				_, done, err := picker.Pick(context.Background(), balancer.PickInfo{
					FullMethodName: "/",
					Ctx:            context.Background(),
				})
				assert.Nil(t, err)
				if i%100 == 0 {
					err = status.Error(codes.DeadlineExceeded, "deadline")
				}
				done(balancer.DoneInfo{
					Err: err,
				})
			}

			dist := make(map[interface{}]int)
			conns := picker.(*p2cPicker).conns
			for _, conn := range conns {
				dist[conn.addr.Addr] = int(conn.requests)
			}

			entropy := mathx.CalcEntropy(dist)
			assert.True(t, entropy > .95, fmt.Sprintf("entropy is %f, less than .95", entropy))
		})
	}
}

type mockClientConn struct {
}

func (m mockClientConn) UpdateAddresses(addresses []resolver.Address) {
}

func (m mockClientConn) Connect() {
}
