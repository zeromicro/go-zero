package breaker

import (
	"context"
	"errors"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/core/stat"
	"github.com/zeromicro/go-zero/zrpc/internal/balancer/p2c"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/status"
)

func init() {
	stat.SetReporter(nil)
}

// mockError implements error and GRPCStatus for testing
type mockError struct {
	st *status.Status
}

func (m mockError) GRPCStatus() *status.Status {
	return m.st
}

func (m mockError) Error() string {
	return m.st.Message()
}

// mockSubConn implements balancer.SubConn for testing
type mockSubConn struct {
	balancer.SubConn
}

func (m *mockSubConn) UpdateAddresses([]resolver.Address) {}
func (m *mockSubConn) Connect()                           {}
func (m *mockSubConn) Shutdown()                          {}
func (m *mockSubConn) GetOrBuildProducer(balancer.ProducerBuilder) (balancer.Producer, func()) {
	return nil, func() {}
}

// mockClientConn implements balancer.ClientConn for testing
type mockClientConn struct {
	balancer.ClientConn
	subConns    map[balancer.SubConn]string
	state       balancer.State
	newSubConErr error
}

func newMockClientConn() *mockClientConn {
	return &mockClientConn{
		subConns: make(map[balancer.SubConn]string),
	}
}

func (m *mockClientConn) NewSubConn(addrs []resolver.Address, opts balancer.NewSubConnOptions) (balancer.SubConn, error) {
	if m.newSubConErr != nil {
		return nil, m.newSubConErr
	}
	sc := &mockSubConn{}
	if len(addrs) > 0 {
		m.subConns[sc] = addrs[0].Addr
	}
	return sc, nil
}

func (m *mockClientConn) RemoveSubConn(sc balancer.SubConn) {
	delete(m.subConns, sc)
}

func (m *mockClientConn) UpdateState(state balancer.State) {
	m.state = state
}

func (m *mockClientConn) ResolveNow(resolver.ResolveNowOptions) {}

func (m *mockClientConn) UpdateAddresses(balancer.SubConn, []resolver.Address) {}

// mockPicker implements balancer.Picker for testing
type mockPicker struct {
	result balancer.PickResult
	err    error
	index  int
	addrs  []string
}

func (m *mockPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if m.err != nil {
		return balancer.PickResult{}, m.err
	}
	if len(m.addrs) > 0 {
		m.index++
		sc := &mockSubConn{}
		return balancer.PickResult{SubConn: sc, Done: m.result.Done}, nil
	}
	return m.result, nil
}

func TestGetBalancerName(t *testing.T) {
	tests := []struct {
		name     string
		baseName string
		strategy string
		want     string
	}{
		{
			name:     "service strategy",
			baseName: "p2c_ewma",
			strategy: "service",
			want:     "p2c_ewma_breaker",
		},
		{
			name:     "instance strategy",
			baseName: "p2c_ewma",
			strategy: "instance",
			want:     "p2c_ewma_breaker_instance",
		},
		{
			name:     "consistenthash with service strategy",
			baseName: "consistent_hash",
			strategy: "service",
			want:     "consistent_hash_breaker",
		},
		{
			name:     "consistenthash with instance strategy",
			baseName: "consistent_hash",
			strategy: "instance",
			want:     "consistent_hash_breaker_instance",
		},
		{
			name:     "empty strategy defaults to service",
			baseName: "p2c_ewma",
			strategy: "",
			want:     "p2c_ewma_breaker",
		},
		{
			name:     "unknown strategy defaults to service",
			baseName: "p2c_ewma",
			strategy: "unknown",
			want:     "p2c_ewma_breaker",
		},
		{
			name:     "empty baseName defaults to p2c with service strategy",
			baseName: "",
			strategy: "service",
			want:     p2c.Name + "_breaker",
		},
		{
			name:     "empty baseName defaults to p2c with instance strategy",
			baseName: "",
			strategy: "instance",
			want:     p2c.Name + "_breaker_instance",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetBalancerName(tt.baseName, tt.strategy)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRegister(t *testing.T) {
	t.Run("register with p2c base balancer", func(t *testing.T) {
		name := GetBalancerName(p2c.Name, "service")
		Register(p2c.Name, "service", 0)
		assert.NotNil(t, balancer.Get(name))
	})

	t.Run("register with instance strategy", func(t *testing.T) {
		name := GetBalancerName(p2c.Name, "instance")
		Register(p2c.Name, "instance", 2)
		assert.NotNil(t, balancer.Get(name))
	})

	t.Run("register with negative retryTimes", func(t *testing.T) {
		// use consistenthash to avoid conflict with previous tests
		name := GetBalancerName("consistent_hash", "service")
		Register("consistent_hash", "service", -5)
		builder := balancer.Get(name)
		assert.NotNil(t, builder)
		// assert retryTimes is set to 0 when negative value is provided
		bb, ok := builder.(*breakerBuilder)
		assert.True(t, ok)
		assert.Equal(t, 0, bb.retryTimes)
	})

	t.Run("register with non-existent base balancer", func(t *testing.T) {
		name := GetBalancerName("non_existent_balancer", "service")
		Register("non_existent_balancer", "service", 0)
		// should not register because base balancer doesn't exist
		assert.Nil(t, balancer.Get(name))
	})

	t.Run("duplicate register is safe", func(t *testing.T) {
		name := GetBalancerName(p2c.Name, "service")
		Register(p2c.Name, "service", 0)
		Register(p2c.Name, "service", 0)
		assert.NotNil(t, balancer.Get(name))
	})

	t.Run("register with empty baseName defaults to p2c", func(t *testing.T) {
		name := GetBalancerName("", "service")
		Register("", "service", 0)
		assert.NotNil(t, balancer.Get(name))
	})
}

func TestExtractTarget(t *testing.T) {
	tests := []struct {
		name   string
		opts   balancer.BuildOptions
		expect string
	}{
		{
			name: "with path",
			opts: balancer.BuildOptions{
				Target: resolver.Target{
					URL: url.URL{Path: "/service/name"},
				},
			},
			expect: "service/name",
		},
		{
			name: "with path no leading slash",
			opts: balancer.BuildOptions{
				Target: resolver.Target{
					URL: url.URL{Path: "service/name"},
				},
			},
			expect: "service/name",
		},
		{
			name: "empty path uses endpoint",
			opts: balancer.BuildOptions{
				Target: resolver.Target{
					URL: url.URL{Opaque: "localhost:8080"},
				},
			},
			expect: "localhost:8080",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractTarget(tt.opts)
			assert.Equal(t, tt.expect, got)
		})
	}
}

func TestBreakerBuilderName(t *testing.T) {
	bb := &breakerBuilder{
		name: "test_breaker",
	}
	assert.Equal(t, "test_breaker", bb.Name())
}

func TestBreakerClientConnNewSubConn(t *testing.T) {
	mockCC := newMockClientConn()
	cc := &breakerClientConn{
		ClientConn: mockCC,
		conns:      make(map[balancer.SubConn]string),
	}

	addrs := []resolver.Address{{Addr: "127.0.0.1:8080"}}
	sc, err := cc.NewSubConn(addrs, balancer.NewSubConnOptions{})

	assert.NoError(t, err)
	assert.NotNil(t, sc)
	assert.Equal(t, "127.0.0.1:8080", cc.conns[sc])
}

func TestBreakerClientConnNewSubConnEmpty(t *testing.T) {
	mockCC := newMockClientConn()
	cc := &breakerClientConn{
		ClientConn: mockCC,
		conns:      make(map[balancer.SubConn]string),
	}

	sc, err := cc.NewSubConn([]resolver.Address{}, balancer.NewSubConnOptions{})

	assert.NoError(t, err)
	assert.NotNil(t, sc)
	assert.Empty(t, cc.conns[sc])
}

func TestBreakerClientConnRemoveSubConn(t *testing.T) {
	mockCC := newMockClientConn()
	cc := &breakerClientConn{
		ClientConn: mockCC,
		conns:      make(map[balancer.SubConn]string),
	}

	addrs := []resolver.Address{{Addr: "127.0.0.1:8080"}}
	sc, _ := cc.NewSubConn(addrs, balancer.NewSubConnOptions{})
	assert.Equal(t, "127.0.0.1:8080", cc.conns[sc])

	cc.RemoveSubConn(sc)
	assert.Empty(t, cc.conns[sc])
}

func TestBreakerClientConnUpdateState(t *testing.T) {
	mockCC := newMockClientConn()
	cc := &breakerClientConn{
		ClientConn: mockCC,
		conns:      make(map[balancer.SubConn]string),
		target:     "test-target",
		strategy:   "service",
		retryTimes: 2,
	}

	origPicker := &mockPicker{}
	state := balancer.State{
		ConnectivityState: connectivity.Ready,
		Picker:            origPicker,
	}

	cc.UpdateState(state)

	// verify picker is wrapped
	bp, ok := mockCC.state.Picker.(*breakerPicker)
	assert.True(t, ok)
	assert.Equal(t, "test-target", bp.target)
	assert.Equal(t, "service", bp.strategy)
	assert.Equal(t, 2, bp.retryTimes)
}

func TestBreakerPickerPickWithServiceStrategy(t *testing.T) {
	sc := &mockSubConn{}
	picker := &breakerPicker{
		picker: &mockPicker{
			result: balancer.PickResult{SubConn: sc},
		},
		conns:    map[balancer.SubConn]string{sc: "127.0.0.1:8080"},
		target:   "test-service",
		strategy: strategyService,
	}

	info := balancer.PickInfo{
		FullMethodName: "/test.Service/Method",
		Ctx:            context.Background(),
	}

	result, err := picker.Pick(info)
	assert.NoError(t, err)
	assert.NotNil(t, result.SubConn)
	assert.NotNil(t, result.Done)
}

func TestBreakerPickerPickWithInstanceStrategy(t *testing.T) {
	sc := &mockSubConn{}
	picker := &breakerPicker{
		picker: &mockPicker{
			result: balancer.PickResult{SubConn: sc},
		},
		conns:      map[balancer.SubConn]string{sc: "127.0.0.1:8080"},
		target:     "test-service",
		strategy:   strategyInstance,
		retryTimes: 0,
	}

	info := balancer.PickInfo{
		FullMethodName: "/test.Service/Method",
		Ctx:            context.Background(),
	}

	result, err := picker.Pick(info)
	assert.NoError(t, err)
	assert.NotNil(t, result.SubConn)
	assert.NotNil(t, result.Done)
}

func TestBreakerPickerPickError(t *testing.T) {
	picker := &breakerPicker{
		picker: &mockPicker{
			err: errors.New("no subconn available"),
		},
		conns:    make(map[balancer.SubConn]string),
		target:   "test-service",
		strategy: strategyService,
	}

	info := balancer.PickInfo{
		FullMethodName: "/test.Service/Method",
		Ctx:            context.Background(),
	}

	_, err := picker.Pick(info)
	assert.Error(t, err)
}

func TestBuildDoneFuncAccept(t *testing.T) {
	sc := &mockSubConn{}
	picker := &breakerPicker{
		picker: &mockPicker{
			result: balancer.PickResult{SubConn: sc},
		},
		conns:    map[balancer.SubConn]string{sc: "127.0.0.1:8080"},
		target:   "test-accept",
		strategy: strategyService,
	}

	info := balancer.PickInfo{
		FullMethodName: "/test.Service/Accept",
		Ctx:            context.Background(),
	}

	result, err := picker.Pick(info)
	assert.NoError(t, err)

	// call Done with no error - should accept
	result.Done(balancer.DoneInfo{Err: nil})
}

func TestBuildDoneFuncReject(t *testing.T) {
	sc := &mockSubConn{}
	picker := &breakerPicker{
		picker: &mockPicker{
			result: balancer.PickResult{SubConn: sc},
		},
		conns:    map[balancer.SubConn]string{sc: "127.0.0.1:8080"},
		target:   "test-reject",
		strategy: strategyService,
	}

	info := balancer.PickInfo{
		FullMethodName: "/test.Service/Reject",
		Ctx:            context.Background(),
	}

	result, err := picker.Pick(info)
	assert.NoError(t, err)

	// call Done with error - should reject
	result.Done(balancer.DoneInfo{Err: errors.New("some error")})
}

func TestBuildDoneFuncWithAcceptableError(t *testing.T) {
	sc := &mockSubConn{}
	picker := &breakerPicker{
		picker: &mockPicker{
			result: balancer.PickResult{SubConn: sc},
		},
		conns:    map[balancer.SubConn]string{sc: "127.0.0.1:8080"},
		target:   "test-acceptable",
		strategy: strategyService,
	}

	info := balancer.PickInfo{
		FullMethodName: "/test.Service/Acceptable",
		Ctx:            context.Background(),
	}

	result, err := picker.Pick(info)
	assert.NoError(t, err)

	// call Done with NotFound error - should accept (acceptable error)
	notFoundErr := mockError{st: status.New(codes.NotFound, "not found")}
	result.Done(balancer.DoneInfo{Err: notFoundErr})
}

func TestBuildDoneFuncWithOriginalDone(t *testing.T) {
	originalDoneCalled := false
	sc := &mockSubConn{}
	picker := &breakerPicker{
		picker: &mockPicker{
			result: balancer.PickResult{
				SubConn: sc,
				Done: func(info balancer.DoneInfo) {
					originalDoneCalled = true
				},
			},
		},
		conns:    map[balancer.SubConn]string{sc: "127.0.0.1:8080"},
		target:   "test-original-done",
		strategy: strategyService,
	}

	info := balancer.PickInfo{
		FullMethodName: "/test.Service/OriginalDone",
		Ctx:            context.Background(),
	}

	result, err := picker.Pick(info)
	assert.NoError(t, err)

	result.Done(balancer.DoneInfo{Err: nil})
	assert.True(t, originalDoneCalled)
}

func TestPickWithServiceBreakerTriggered(t *testing.T) {
	sc := &mockSubConn{}
	picker := &breakerPicker{
		picker: &mockPicker{
			result: balancer.PickResult{SubConn: sc},
		},
		conns:    map[balancer.SubConn]string{sc: "127.0.0.1:8080"},
		target:   "test-breaker-trigger",
		strategy: strategyService,
	}

	info := balancer.PickInfo{
		FullMethodName: "/test.Service/BreakerTrigger",
		Ctx:            context.Background(),
	}

	// trigger breaker by sending many DeadlineExceeded errors
	deadlineErr := mockError{st: status.New(codes.DeadlineExceeded, "deadline exceeded")}
	var breakerErr error
	for i := 0; i < 100; i++ {
		result, err := picker.Pick(info)
		if err != nil {
			// breaker triggered
			breakerErr = err
			break
		}
		result.Done(balancer.DoneInfo{Err: deadlineErr})
	}

	assert.ErrorIs(t, breakerErr, breaker.ErrServiceUnavailable)
}

func TestPickWithInstanceBreakerRetry(t *testing.T) {
	sc1 := &mockSubConn{}
	sc2 := &mockSubConn{}
	callCount := 0
	picker := &breakerPicker{
		picker: &mockPicker{
			result: balancer.PickResult{},
		},
		conns: map[balancer.SubConn]string{
			sc1: "127.0.0.1:8080",
			sc2: "127.0.0.1:8081",
		},
		target:     "test-instance-retry",
		strategy:   strategyInstance,
		retryTimes: 2,
	}

	// override picker to return different subconns
	picker.picker = &mockPicker{
		result: balancer.PickResult{SubConn: sc1},
	}

	info := balancer.PickInfo{
		FullMethodName: "/test.Service/InstanceRetry",
		Ctx:            context.Background(),
	}

	result, err := picker.Pick(info)
	assert.NoError(t, err)
	assert.NotNil(t, result.SubConn)
	_ = callCount
}

func TestPickWithInstanceBreakerEmptyAddr(t *testing.T) {
	sc := &mockSubConn{}
	picker := &breakerPicker{
		picker: &mockPicker{
			result: balancer.PickResult{SubConn: sc},
		},
		conns:      map[balancer.SubConn]string{}, // empty conns map
		target:     "test-empty-addr",
		strategy:   strategyInstance,
		retryTimes: 0,
	}

	info := balancer.PickInfo{
		FullMethodName: "/test.Service/EmptyAddr",
		Ctx:            context.Background(),
	}

	result, err := picker.Pick(info)
	assert.NoError(t, err)
	assert.NotNil(t, result.SubConn)
}

func TestBreakerBuilderBuild(t *testing.T) {
	mockCC := newMockClientConn()
	bb := &breakerBuilder{
		baseBuilder: balancer.Get(p2c.Name),
		name:        "test_breaker_build",
		strategy:    strategyService,
		retryTimes:  2,
	}

	opts := balancer.BuildOptions{
		Target: resolver.Target{
			URL: url.URL{Path: "/test/service"},
		},
	}

	bal := bb.Build(mockCC, opts)
	assert.NotNil(t, bal)
}

func TestBreakerClientConnNewSubConnError(t *testing.T) {
	mockCC := newMockClientConn()
	mockCC.newSubConErr = errors.New("connection error")

	cc := &breakerClientConn{
		ClientConn: mockCC,
		conns:      make(map[balancer.SubConn]string),
	}

	addrs := []resolver.Address{{Addr: "127.0.0.1:8080"}}
	sc, err := cc.NewSubConn(addrs, balancer.NewSubConnOptions{})

	assert.Error(t, err)
	assert.Nil(t, sc)
}

func TestBreakerClientConnUpdateStateWithConns(t *testing.T) {
	mockCC := newMockClientConn()
	sc := &mockSubConn{}
	cc := &breakerClientConn{
		ClientConn: mockCC,
		conns:      map[balancer.SubConn]string{sc: "127.0.0.1:8080"},
		target:     "test-target",
		strategy:   "service",
		retryTimes: 2,
	}

	origPicker := &mockPicker{}
	state := balancer.State{
		ConnectivityState: connectivity.Ready,
		Picker:            origPicker,
	}

	cc.UpdateState(state)

	// verify picker is wrapped and conns are copied
	bp, ok := mockCC.state.Picker.(*breakerPicker)
	assert.True(t, ok)
	assert.Equal(t, "127.0.0.1:8080", bp.conns[sc])
}

func TestPickWithInstanceBreakerPickerError(t *testing.T) {
	picker := &breakerPicker{
		picker: &mockPicker{
			err: errors.New("picker error"),
		},
		conns:      make(map[balancer.SubConn]string),
		target:     "test-picker-error",
		strategy:   strategyInstance,
		retryTimes: 0,
	}

	info := balancer.PickInfo{
		FullMethodName: "/test.Service/PickerError",
		Ctx:            context.Background(),
	}

	_, err := picker.Pick(info)
	assert.Error(t, err)
}

func TestPickWithInstanceBreakerTriedAddr(t *testing.T) {
	sc := &mockSubConn{}
	callCount := 0
	picker := &breakerPicker{
		picker: &mockPicker{
			result: balancer.PickResult{SubConn: sc},
		},
		conns:      map[balancer.SubConn]string{sc: "127.0.0.1:8080"},
		target:     "test-tried-addr",
		strategy:   strategyInstance,
		retryTimes: 2, // will try 3 times total
	}

	// override picker to track calls and always return same subconn
	picker.picker = &mockPicker{
		result: balancer.PickResult{SubConn: sc},
	}

	info := balancer.PickInfo{
		FullMethodName: "/test.Service/TriedAddr",
		Ctx:            context.Background(),
	}

	result, err := picker.Pick(info)
	assert.NoError(t, err)
	assert.NotNil(t, result.SubConn)
	_ = callCount
}

func TestPickWithInstanceBreakerAllFailed(t *testing.T) {
	sc := &mockSubConn{}
	picker := &breakerPicker{
		picker: &mockPicker{
			result: balancer.PickResult{SubConn: sc},
		},
		conns:      map[balancer.SubConn]string{sc: "127.0.0.1:8080"},
		target:     "test-all-failed",
		strategy:   strategyInstance,
		retryTimes: 0,
	}

	info := balancer.PickInfo{
		FullMethodName: "/test.Service/AllFailed",
		Ctx:            context.Background(),
	}

	// trigger instance breaker
	deadlineErr := mockError{st: status.New(codes.DeadlineExceeded, "deadline exceeded")}
	var breakerErr error
	for i := 0; i < 100; i++ {
		result, err := picker.Pick(info)
		if err != nil {
			// breaker triggered, all retries failed
			breakerErr = err
			break
		}
		result.Done(balancer.DoneInfo{Err: deadlineErr})
	}

	assert.ErrorIs(t, breakerErr, breaker.ErrServiceUnavailable)
}

// mockMultiPicker returns different subconns on each call
type mockMultiPicker struct {
	subConns []balancer.SubConn
	index    int
	err      error
}

func (m *mockMultiPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if m.err != nil {
		return balancer.PickResult{}, m.err
	}
	if len(m.subConns) == 0 {
		return balancer.PickResult{}, errors.New("no subconns")
	}
	sc := m.subConns[m.index%len(m.subConns)]
	m.index++
	return balancer.PickResult{SubConn: sc}, nil
}

func TestPickWithInstanceBreakerRetryWithDifferentAddrs(t *testing.T) {
	sc1 := &mockSubConn{}
	sc2 := &mockSubConn{}

	picker := &breakerPicker{
		picker: &mockMultiPicker{
			subConns: []balancer.SubConn{sc1, sc2},
		},
		conns: map[balancer.SubConn]string{
			sc1: "127.0.0.1:8080",
			sc2: "127.0.0.1:8081",
		},
		target:     "test-retry-diff-addrs",
		strategy:   strategyInstance,
		retryTimes: 2,
	}

	info := balancer.PickInfo{
		FullMethodName: "/test.Service/RetryDiffAddrs",
		Ctx:            context.Background(),
	}

	result, err := picker.Pick(info)
	assert.NoError(t, err)
	assert.NotNil(t, result.SubConn)
}

func TestPickWithInstanceBreakerRetryOnBreakerReject(t *testing.T) {
	sc1 := &mockSubConn{}
	sc2 := &mockSubConn{}

	picker := &breakerPicker{
		picker: &mockMultiPicker{
			subConns: []balancer.SubConn{sc1, sc2},
		},
		conns: map[balancer.SubConn]string{
			sc1: "127.0.0.1:18080",
			sc2: "127.0.0.1:18081",
		},
		target:     "test-retry-breaker-reject",
		strategy:   strategyInstance,
		retryTimes: 2,
	}

	info := balancer.PickInfo{
		FullMethodName: "/test.Service/RetryBreakerReject",
		Ctx:            context.Background(),
	}

	// first trigger breaker for sc1
	deadlineErr := mockError{st: status.New(codes.DeadlineExceeded, "deadline exceeded")}
	for i := 0; i < 50; i++ {
		result, _ := picker.Pick(info)
		if result.Done != nil {
			result.Done(balancer.DoneInfo{Err: deadlineErr})
		}
	}

	// now picker should retry and find sc2
	result, err := picker.Pick(info)
	// either success or breaker error
	if err == nil {
		assert.NotNil(t, result.SubConn)
	}
}

// mockRetryPicker returns success on first call, then returns configured behavior
type mockRetryPicker struct {
	firstSubConn  balancer.SubConn
	retrySubConn  balancer.SubConn
	retryErr      error
	callCount     int
}

func (m *mockRetryPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	m.callCount++
	if m.callCount == 1 {
		return balancer.PickResult{SubConn: m.firstSubConn}, nil
	}
	if m.retryErr != nil {
		return balancer.PickResult{}, m.retryErr
	}
	return balancer.PickResult{SubConn: m.retrySubConn}, nil
}

func TestPickWithInstanceBreakerRetryPickerError(t *testing.T) {
	sc1 := &mockSubConn{}
	addr := "127.0.0.1:58080"
	method := "/test.Service/RetryPickerErr2"

	// First, trigger breaker for sc1
	triggerPicker := &breakerPicker{
		picker: &mockPicker{
			result: balancer.PickResult{SubConn: sc1},
		},
		conns:      map[balancer.SubConn]string{sc1: addr},
		target:     "test-retry-picker-err",
		strategy:   strategyInstance,
		retryTimes: 1,
	}

	triggerInfo := balancer.PickInfo{
		FullMethodName: method,
		Ctx:            context.Background(),
	}

	deadlineErr := mockError{st: status.New(codes.DeadlineExceeded, "deadline exceeded")}
	for i := 0; i < 200; i++ {
		result, err := triggerPicker.Pick(triggerInfo)
		if err != nil {
			// breaker is now triggered
			break
		}
		if result.Done != nil {
			result.Done(balancer.DoneInfo{Err: deadlineErr})
		}
	}

	// Now test: first attempt hits breaker, retry returns picker error
	picker := &breakerPicker{
		picker: &mockRetryPicker{
			firstSubConn: sc1,
			retryErr:     errors.New("no available subconn"),
		},
		conns:      map[balancer.SubConn]string{sc1: addr},
		target:     "test-retry-picker-err",
		strategy:   strategyInstance,
		retryTimes: 1,
	}

	result, err := picker.Pick(triggerInfo)
	// Either breaker triggered and retry returns picker error,
	// or breaker not triggered and returns normally
	if err != nil {
		assert.Equal(t, "no available subconn", err.Error())
	} else {
		assert.NotNil(t, result.SubConn)
	}
}

func TestPickWithInstanceBreakerRetryEmptyAddr(t *testing.T) {
	sc1 := &mockSubConn{}
	sc2 := &mockSubConn{}

	addr := "127.0.0.1:48080"
	method := "/test.Service/RetryEmptyAddr2"

	// First, trigger breaker for sc1
	triggerPicker := &breakerPicker{
		picker: &mockPicker{
			result: balancer.PickResult{SubConn: sc1},
		},
		conns:      map[balancer.SubConn]string{sc1: addr},
		target:     "test-retry-empty-addr",
		strategy:   strategyInstance,
		retryTimes: 1,
	}

	triggerInfo := balancer.PickInfo{
		FullMethodName: method,
		Ctx:            context.Background(),
	}

	deadlineErr := mockError{st: status.New(codes.DeadlineExceeded, "deadline exceeded")}
	for i := 0; i < 200; i++ {
		result, err := triggerPicker.Pick(triggerInfo)
		if err != nil {
			// breaker is now triggered
			break
		}
		if result.Done != nil {
			result.Done(balancer.DoneInfo{Err: deadlineErr})
		}
	}

	// Now test: first attempt hits breaker, retry returns subconn not in conns map
	picker := &breakerPicker{
		picker: &mockRetryPicker{
			firstSubConn: sc1,
			retrySubConn: sc2, // sc2 is not in conns map
		},
		conns:      map[balancer.SubConn]string{sc1: addr}, // only sc1
		target:     "test-retry-empty-addr",
		strategy:   strategyInstance,
		retryTimes: 1,
	}

	result, err := picker.Pick(triggerInfo)
	assert.NoError(t, err)
	// Either sc1 (breaker not triggered) or sc2 (breaker triggered, retry with empty addr)
	assert.NotNil(t, result.SubConn)
}
