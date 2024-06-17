package p2c

import (
	"context"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mathx"
	"github.com/zeromicro/go-zero/core/stringx"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/status"
)

func init() {
	logx.Disable()
}

func TestP2cPicker_PickNil(t *testing.T) {
	builder := new(p2cPickerBuilder)
	picker := builder.Build(base.PickerBuildInfo{})
	_, err := picker.Pick(balancer.PickInfo{
		FullMethodName: "/",
		Ctx:            context.Background(),
	})
	assert.NotNil(t, err)
}

func TestP2cPicker_Pick(t *testing.T) {
	tests := []struct {
		name       string
		candidates int
		err        error
		threshold  float64
	}{
		{
			name:       "empty",
			candidates: 0,
			err:        balancer.ErrNoSubConnAvailable,
		},
		{
			name:       "single",
			candidates: 1,
			threshold:  0.9,
		},
		{
			name:       "two",
			candidates: 2,
			threshold:  0.5,
		},
		{
			name:       "multiple",
			candidates: 100,
			threshold:  0.95,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			const total = 10000
			builder := new(p2cPickerBuilder)
			ready := make(map[balancer.SubConn]base.SubConnInfo)
			for i := 0; i < test.candidates; i++ {
				ready[mockClientConn{
					id: stringx.Rand(),
				}] = base.SubConnInfo{
					Address: resolver.Address{
						Addr: strconv.Itoa(i),
					},
				}
			}

			picker := builder.Build(base.PickerBuildInfo{
				ReadySCs: ready,
			})
			var wg sync.WaitGroup
			wg.Add(total)
			for i := 0; i < total; i++ {
				result, err := picker.Pick(balancer.PickInfo{
					FullMethodName: "/",
					Ctx:            context.Background(),
				})
				assert.Equal(t, test.err, err)

				if test.err != nil {
					return
				}

				if i%100 == 0 {
					err = status.Error(codes.DeadlineExceeded, "deadline")
				}

				go func() {
					runtime.Gosched()
					result.Done(balancer.DoneInfo{
						Err: err,
					})
					wg.Done()
				}()
			}

			wg.Wait()
			dist := make(map[any]int)
			conns := picker.(*p2cPicker).conns
			for _, conn := range conns {
				dist[conn.addr.Addr] = int(conn.requests)
			}

			entropy := mathx.CalcEntropy(dist)
			assert.True(t, entropy > test.threshold, fmt.Sprintf("entropy is %f, less than %f",
				entropy, test.threshold))
		})
	}
}

func TestPickerWithEmptyConns(t *testing.T) {
	var picker p2cPicker
	_, err := picker.Pick(balancer.PickInfo{
		FullMethodName: "/",
		Ctx:            context.Background(),
	})
	assert.ErrorIs(t, err, balancer.ErrNoSubConnAvailable)
}

type mockClientConn struct {
	// add random string member to avoid map key equality.
	id string
}

func (m mockClientConn) GetOrBuildProducer(builder balancer.ProducerBuilder) (
	p balancer.Producer, close func()) {
	return builder.Build(m)
}

func (m mockClientConn) UpdateAddresses(_ []resolver.Address) {
}

func (m mockClientConn) Connect() {
}

func (m mockClientConn) Shutdown() {
}
