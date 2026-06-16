package breaker

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/core/stat"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/status"
)

func init() {
	stat.SetReporter(nil)
}

type mockSubConn struct {
	addr string
}

func (m *mockSubConn) UpdateAddresses([]resolver.Address) {}
func (m *mockSubConn) Connect()                           {}
func (m *mockSubConn) Shutdown()                          {}
func (m *mockSubConn) GetOrBuildProducer(balancer.ProducerBuilder) (balancer.Producer, func()) {
	return nil, nil
}

type mockPicker struct {
	err     error
	index   int
	results []balancer.PickResult
}

func (m *mockPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if m.err != nil {
		return balancer.PickResult{}, m.err
	}
	if len(m.results) == 0 {
		return balancer.PickResult{}, nil
	}
	result := m.results[m.index%len(m.results)]
	m.index++
	return result, nil
}

func TestWithBreaker(t *testing.T) {
	ctx := context.Background()
	assert.False(t, HasBreaker(ctx))

	ctx = WithBreaker(ctx)
	assert.True(t, HasBreaker(ctx))
}

func TestWrapPicker(t *testing.T) {
	sc := &mockSubConn{addr: "127.0.0.1:8080"}
	info := base.PickerBuildInfo{
		ReadySCs: map[balancer.SubConn]base.SubConnInfo{
			sc: {Address: resolver.Address{Addr: "127.0.0.1:8080"}},
		},
	}
	inner := &mockPicker{}

	picker := WrapPicker(info, inner, true)
	assert.NotNil(t, picker)

	bp := picker.(*breakerPicker)
	assert.Equal(t, inner, bp.picker)
	assert.True(t, bp.retryable)
	assert.Equal(t, "127.0.0.1:8080", bp.addrMap[sc])
}

func TestUnwrap(t *testing.T) {
	inner := &mockPicker{}
	picker := &breakerPicker{picker: inner}

	assert.Equal(t, inner, picker.Unwrap())
}

func TestPickWithoutBreaker(t *testing.T) {
	sc := &mockSubConn{addr: "127.0.0.1:8080"}
	inner := &mockPicker{
		results: []balancer.PickResult{{SubConn: sc}},
	}
	picker := &breakerPicker{
		picker:  inner,
		addrMap: map[balancer.SubConn]string{sc: "127.0.0.1:8080"},
	}

	result, err := picker.Pick(balancer.PickInfo{Ctx: context.Background()})
	assert.NoError(t, err)
	assert.Equal(t, sc, result.SubConn)
}

func TestPickNotRetryable(t *testing.T) {
	sc := &mockSubConn{addr: "127.0.0.1:8080"}
	inner := &mockPicker{
		results: []balancer.PickResult{{SubConn: sc}},
	}
	picker := &breakerPicker{
		picker:    inner,
		addrMap:   map[balancer.SubConn]string{sc: "127.0.0.1:8080"},
		retryable: false,
	}

	ctx := WithBreaker(context.Background())
	result, err := picker.Pick(balancer.PickInfo{
		Ctx:            ctx,
		FullMethodName: "/test",
	})
	assert.NoError(t, err)
	assert.NotNil(t, result.Done)
}

func TestPickRetryableSuccess(t *testing.T) {
	sc := &mockSubConn{addr: "127.0.0.1:8080"}
	inner := &mockPicker{
		results: []balancer.PickResult{{SubConn: sc}},
	}
	picker := &breakerPicker{
		picker:    inner,
		addrMap:   map[balancer.SubConn]string{sc: "127.0.0.1:8080"},
		retryable: true,
	}

	ctx := WithBreaker(context.Background())
	result, err := picker.Pick(balancer.PickInfo{
		Ctx:            ctx,
		FullMethodName: "/test",
	})
	assert.NoError(t, err)
	assert.NotNil(t, result.Done)
}

func TestPickRetryablePickerError(t *testing.T) {
	expectedErr := status.Error(codes.Unavailable, "unavailable")
	inner := &mockPicker{err: expectedErr}
	picker := &breakerPicker{
		picker:    inner,
		addrMap:   map[balancer.SubConn]string{},
		retryable: true,
	}

	ctx := WithBreaker(context.Background())
	_, err := picker.Pick(balancer.PickInfo{
		Ctx:            ctx,
		FullMethodName: "/test",
	})
	assert.Equal(t, expectedErr, err)
}

func TestPickNotRetryablePickerError(t *testing.T) {
	expectedErr := status.Error(codes.Unavailable, "unavailable")
	inner := &mockPicker{err: expectedErr}
	picker := &breakerPicker{
		picker:    inner,
		addrMap:   map[balancer.SubConn]string{},
		retryable: false,
	}

	ctx := WithBreaker(context.Background())
	_, err := picker.Pick(balancer.PickInfo{
		Ctx:            ctx,
		FullMethodName: "/test",
	})
	assert.Equal(t, expectedErr, err)
}

func TestPickBreakerOpen(t *testing.T) {
	sc := &mockSubConn{addr: "127.0.0.1:8081"}
	inner := &mockPicker{
		results: []balancer.PickResult{{SubConn: sc}},
	}
	picker := &breakerPicker{
		picker:    inner,
		addrMap:   map[balancer.SubConn]string{sc: "127.0.0.1:8081"},
		retryable: false,
	}

	// trigger breaker
	for i := 0; i < 1000; i++ {
		_ = breaker.DoWithAcceptable("127.0.0.1:8081/test", func() error {
			return status.Error(codes.DeadlineExceeded, "timeout")
		}, func(err error) bool {
			return false
		})
	}

	ctx := WithBreaker(context.Background())
	_, err := picker.Pick(balancer.PickInfo{
		Ctx:            ctx,
		FullMethodName: "/test",
	})
	assert.ErrorIs(t, err, breaker.ErrServiceUnavailable)
}

func TestPickRetryOnBreaker(t *testing.T) {
	sc1 := &mockSubConn{addr: "127.0.0.1:8082"}
	sc2 := &mockSubConn{addr: "127.0.0.1:8083"}
	inner := &mockPicker{
		results: []balancer.PickResult{{SubConn: sc1}, {SubConn: sc2}},
	}
	picker := &breakerPicker{
		picker:  inner,
		addrMap: map[balancer.SubConn]string{sc1: "127.0.0.1:8082", sc2: "127.0.0.1:8083"},
		retryable: true,
	}

	// trigger breaker for sc1
	for i := 0; i < 1000; i++ {
		_ = breaker.DoWithAcceptable("127.0.0.1:8082/retry", func() error {
			return status.Error(codes.DeadlineExceeded, "timeout")
		}, func(err error) bool {
			return false
		})
	}

	ctx := WithBreaker(context.Background())
	result, err := picker.Pick(balancer.PickInfo{
		Ctx:            ctx,
		FullMethodName: "/retry",
	})
	assert.NoError(t, err)
	assert.Equal(t, sc2, result.SubConn)
}

func TestPickAllBreakerOpen(t *testing.T) {
	sc := &mockSubConn{addr: "127.0.0.1:8084"}
	inner := &mockPicker{
		results: []balancer.PickResult{{SubConn: sc}},
	}
	picker := &breakerPicker{
		picker:    inner,
		addrMap:   map[balancer.SubConn]string{sc: "127.0.0.1:8084"},
		retryable: true,
	}

	// trigger breaker
	for i := 0; i < 1000; i++ {
		_ = breaker.DoWithAcceptable("127.0.0.1:8084/all", func() error {
			return status.Error(codes.DeadlineExceeded, "timeout")
		}, func(err error) bool {
			return false
		})
	}

	ctx := WithBreaker(context.Background())
	_, err := picker.Pick(balancer.PickInfo{
		Ctx:            ctx,
		FullMethodName: "/all",
	})
	assert.ErrorIs(t, err, breaker.ErrServiceUnavailable)
}

func TestBuildDoneFuncWithNilDone(t *testing.T) {
	picker := &breakerPicker{}
	brk := breaker.NewBreaker()
	promise, _ := brk.Allow()

	done := picker.buildDoneFunc(nil, promise)
	assert.NotNil(t, done)
	done(balancer.DoneInfo{})
}

func TestBuildDoneFuncWithDone(t *testing.T) {
	picker := &breakerPicker{}
	brk := breaker.NewBreaker()
	promise, _ := brk.Allow()

	called := false
	done := picker.buildDoneFunc(func(info balancer.DoneInfo) {
		called = true
	}, promise)
	done(balancer.DoneInfo{})
	assert.True(t, called)
}

func TestBuildDoneFuncAcceptable(t *testing.T) {
	picker := &breakerPicker{}
	brk := breaker.NewBreaker()
	promise, _ := brk.Allow()

	done := picker.buildDoneFunc(nil, promise)
	done(balancer.DoneInfo{Err: status.Error(codes.NotFound, "not found")})
}

func TestBuildDoneFuncNotAcceptable(t *testing.T) {
	picker := &breakerPicker{}
	brk := breaker.NewBreaker()
	promise, _ := brk.Allow()

	done := picker.buildDoneFunc(nil, promise)
	done(balancer.DoneInfo{Err: status.Error(codes.DeadlineExceeded, "timeout")})
}

func TestPickBreakerOpenWithDone(t *testing.T) {
	sc := &mockSubConn{addr: "127.0.0.1:8085"}
	doneCalled := false
	inner := &mockPicker{
		results: []balancer.PickResult{{
			SubConn: sc,
			Done: func(info balancer.DoneInfo) {
				doneCalled = true
			},
		}},
	}
	picker := &breakerPicker{
		picker:    inner,
		addrMap:   map[balancer.SubConn]string{sc: "127.0.0.1:8085"},
		retryable: false,
	}

	// trigger breaker
	for i := 0; i < 1000; i++ {
		_ = breaker.DoWithAcceptable("127.0.0.1:8085/done", func() error {
			return status.Error(codes.DeadlineExceeded, "timeout")
		}, func(err error) bool {
			return false
		})
	}

	ctx := WithBreaker(context.Background())
	_, err := picker.Pick(balancer.PickInfo{
		Ctx:            ctx,
		FullMethodName: "/done",
	})
	assert.ErrorIs(t, err, breaker.ErrServiceUnavailable)
	assert.True(t, doneCalled)
}
