package breaker

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/zrpc/internal/balancer/consistenthash"
	"github.com/zeromicro/go-zero/zrpc/internal/balancer/p2c"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/status"
)

func TestBreakerBalancerRegistered(t *testing.T) {
	tests := []struct {
		name         string
		balancerName string
	}{
		{"p2c_service", p2c.Name + BalancerSuffix},
		{"p2c_instance", p2c.Name + BalancerSuffix + "_instance"},
		{"consistenthash_service", consistenthash.Name + BalancerSuffix},
		{"consistenthash_instance", consistenthash.Name + BalancerSuffix + "_instance"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := balancer.Get(tt.balancerName)
			assert.NotNil(t, b)
			assert.Equal(t, tt.balancerName, b.Name())
		})
	}
}

func TestGetName(t *testing.T) {
	tests := []struct {
		baseName string
		strategy string
		expected string
	}{
		{p2c.Name, "service", p2c.Name + BalancerSuffix},
		{p2c.Name, "instance", p2c.Name + BalancerSuffix + "_instance"},
		{consistenthash.Name, "service", consistenthash.Name + BalancerSuffix},
		{consistenthash.Name, "instance", consistenthash.Name + BalancerSuffix + "_instance"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			name := GetName(tt.baseName, tt.strategy)
			assert.Equal(t, tt.expected, name)
		})
	}
}

func TestBreakerPicker_ServiceStrategy(t *testing.T) {
	// Test service-level breaker
	conn1 := &mockSubConn{id: "conn_svc_1"}

	innerPicker := &mockPicker{
		subConn: conn1,
	}

	conns := map[balancer.SubConn]string{
		conn1: "127.0.0.10:8080",
	}

	picker := &breakerPicker{
		picker:     innerPicker,
		conns:      conns,
		target:     "test-service",
		strategy:   strategyService,
		retryTimes: defaultRetryTimes,
	}

	result, err := picker.Pick(balancer.PickInfo{
		FullMethodName: "/test.Service/ServiceMethod",
		Ctx:            context.Background(),
	})

	assert.NoError(t, err)
	assert.NotNil(t, result.SubConn)
	assert.NotNil(t, result.Done)

	// Simulate success
	result.Done(balancer.DoneInfo{Err: nil})
}

func TestBreakerPicker_ServiceStrategyBreaker(t *testing.T) {
	// Test service-level breaker triggers
	conn1 := &mockSubConn{id: "conn_svc_breaker"}

	innerPicker := &mockPicker{
		subConn: conn1,
	}

	conns := map[balancer.SubConn]string{
		conn1: "127.0.0.11:8080",
	}

	picker := &breakerPicker{
		picker:     innerPicker,
		conns:      conns,
		target:     "test-service-breaker",
		strategy:   strategyService,
		retryTimes: defaultRetryTimes,
	}

	// Simulate multiple failures to trigger breaker
	var breakerTriggered bool
	for i := 0; i < 1000; i++ {
		result, err := picker.Pick(balancer.PickInfo{
			FullMethodName: "/test.Service/SvcBreakerMethod",
			Ctx:            context.Background(),
		})
		if err != nil {
			breakerTriggered = true
			break
		}
		if result.Done != nil {
			result.Done(balancer.DoneInfo{
				Err: status.Error(codes.Internal, "internal error"),
			})
		}
	}

	assert.True(t, breakerTriggered, "service breaker should be triggered after many failures")
}

func TestBreakerPicker_InstanceStrategy(t *testing.T) {
	// Test instance-level breaker (normal case)
	conn1 := &mockSubConn{id: "conn_inst_1"}

	innerPicker := &mockPicker{
		subConn: conn1,
	}

	conns := map[balancer.SubConn]string{
		conn1: "127.0.0.12:8080",
	}

	picker := &breakerPicker{
		picker:     innerPicker,
		conns:      conns,
		target:     "test-service",
		strategy:   strategyInstance,
		retryTimes: defaultRetryTimes,
	}

	result, err := picker.Pick(balancer.PickInfo{
		FullMethodName: "/test.Service/InstanceMethod",
		Ctx:            context.Background(),
	})

	assert.NoError(t, err)
	assert.NotNil(t, result.SubConn)
	assert.NotNil(t, result.Done)

	// Simulate success
	result.Done(balancer.DoneInfo{Err: nil})
}

func TestBreakerPicker_InstanceStrategyWithRetry(t *testing.T) {
	// Test instance-level breaker with retry to another instance
	conn1 := &mockSubConn{id: "conn_retry_1"}
	conn2 := &mockSubConn{id: "conn_retry_2"}

	pickCount := 0
	innerPicker := &mockPicker{
		pickFunc: func() (balancer.SubConn, error) {
			pickCount++
			if pickCount == 1 {
				return conn1, nil
			}
			return conn2, nil
		},
	}

	conns := map[balancer.SubConn]string{
		conn1: "127.0.0.13:8080",
		conn2: "127.0.0.13:8081",
	}

	picker := &breakerPicker{
		picker:     innerPicker,
		conns:      conns,
		target:     "test-service",
		strategy:   strategyInstance,
		retryTimes: defaultRetryTimes,
	}

	// Trigger breaker for conn1
	// Google breaker algorithm needs many failures to trigger
	addr1BreakerName := "127.0.0.13:8080/test.Service/RetryMethod"
	for i := 0; i < 1000; i++ {
		p, err := breaker.GetBreaker(addr1BreakerName).Allow()
		if err == nil {
			p.Reject("test error")
		}
	}

	result, err := picker.Pick(balancer.PickInfo{
		FullMethodName: "/test.Service/RetryMethod",
		Ctx:            context.Background(),
	})

	assert.NoError(t, err)
	assert.Equal(t, conn2, result.SubConn)
	assert.True(t, pickCount >= 2)
}

func TestBreakerPicker_InstanceStrategyAllBroken(t *testing.T) {
	// Test when all instances are broken
	conn1 := &mockSubConn{id: "conn_all_broken_inst_v2"}

	innerPicker := &mockPicker{
		subConn: conn1,
	}

	conns := map[balancer.SubConn]string{
		conn1: "127.0.0.24:8080",
	}

	picker := &breakerPicker{
		picker:     innerPicker,
		conns:      conns,
		target:     "test-service",
		strategy:   strategyInstance,
		retryTimes: defaultRetryTimes,
	}

	// Trigger breaker for conn1
	// Google breaker algorithm needs many failures to trigger
	addr1BreakerName := "127.0.0.24:8080/test.Service/AllBrokenInstV2"
	for i := 0; i < 1000; i++ {
		p, err := breaker.GetBreaker(addr1BreakerName).Allow()
		if err == nil {
			p.Reject("test error")
		}
	}

	_, err := picker.Pick(balancer.PickInfo{
		FullMethodName: "/test.Service/AllBrokenInstV2",
		Ctx:            context.Background(),
	})

	assert.ErrorIs(t, err, breaker.ErrServiceUnavailable)
}

func TestBreakerPicker_InnerPickerError(t *testing.T) {
	// Test when inner picker returns error
	innerPicker := &mockPicker{
		err: errors.New("no connection available"),
	}

	picker := &breakerPicker{
		picker:     innerPicker,
		conns:      make(map[balancer.SubConn]string),
		target:     "test-service",
		strategy:   strategyInstance,
		retryTimes: defaultRetryTimes,
	}

	_, err := picker.Pick(balancer.PickInfo{
		FullMethodName: "/test.Service/Method",
		Ctx:            context.Background(),
	})

	assert.Error(t, err)
	assert.Equal(t, "no connection available", err.Error())
}

func TestBreakerPicker_UnknownSubConn(t *testing.T) {
	// Test when SubConn is not in conns map (instance strategy)
	conn1 := &mockSubConn{id: "unknown_conn"}

	innerPicker := &mockPicker{
		subConn: conn1,
	}

	picker := &breakerPicker{
		picker:     innerPicker,
		conns:      make(map[balancer.SubConn]string), // empty map
		target:     "test-service",
		strategy:   strategyInstance,
		retryTimes: defaultRetryTimes,
	}

	result, err := picker.Pick(balancer.PickInfo{
		FullMethodName: "/test.Service/Method",
		Ctx:            context.Background(),
	})

	// Should return result without breaker when addr not found
	assert.NoError(t, err)
	assert.Equal(t, conn1, result.SubConn)
}

type mockSubConn struct {
	id string
}

func (m *mockSubConn) UpdateAddresses(_ []resolver.Address) {}
func (m *mockSubConn) Connect()                             {}
func (m *mockSubConn) Shutdown()                            {}
func (m *mockSubConn) GetOrBuildProducer(_ balancer.ProducerBuilder) (balancer.Producer, func()) {
	return nil, func() {}
}

type mockPicker struct {
	subConn  balancer.SubConn
	err      error
	pickFunc func() (balancer.SubConn, error)
}

func (m *mockPicker) Pick(_ balancer.PickInfo) (balancer.PickResult, error) {
	if m.pickFunc != nil {
		sc, err := m.pickFunc()
		if err != nil {
			return balancer.PickResult{}, err
		}
		return balancer.PickResult{SubConn: sc}, nil
	}

	if m.err != nil {
		return balancer.PickResult{}, m.err
	}
	return balancer.PickResult{SubConn: m.subConn}, nil
}
