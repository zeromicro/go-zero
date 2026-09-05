package breaker

import (
	"context"
	"errors"
	"path"

	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/zrpc/internal/codes"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

const retries = 3

var emptyPickResult balancer.PickResult

type (
	breakerKey struct{}

	// Unwrapper is the interface for unwrapping the inner picker.
	Unwrapper interface {
		Unwrap() balancer.Picker
	}
)

type breakerPicker struct {
	picker    balancer.Picker
	addrMap   map[balancer.SubConn]string
	retryable bool
}

// Unwrap returns the inner picker.
func (p *breakerPicker) Unwrap() balancer.Picker {
	return p.picker
}

// WrapPicker wraps the given picker with circuit breaker.
// retryable indicates whether to retry on another node when circuit breaker is open.
func WrapPicker(info base.PickerBuildInfo, picker balancer.Picker, retryable bool) balancer.Picker {
	addrMap := make(map[balancer.SubConn]string, len(info.ReadySCs))
	for conn, connInfo := range info.ReadySCs {
		addrMap[conn] = connInfo.Address.Addr
	}

	return &breakerPicker{
		picker:    picker,
		addrMap:   addrMap,
		retryable: retryable,
	}
}

func (p *breakerPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if !HasBreaker(info.Ctx) {
		return p.picker.Pick(info)
	}

	if !p.retryable {
		return p.pick(info)
	}

	var (
		err    error
		result balancer.PickResult
	)

	for i := 0; i <= retries; i++ {
		result, err = p.pick(info)
		if err == nil || !errors.Is(err, breaker.ErrServiceUnavailable) {
			break
		}
	}

	return result, err
}

func (p *breakerPicker) pick(info balancer.PickInfo) (balancer.PickResult, error) {
	result, err := p.picker.Pick(info)
	if err != nil {
		return result, err
	}

	addr := p.addrMap[result.SubConn]
	breakerName := path.Join(addr, info.FullMethodName)
	promise, err := breaker.GetBreaker(breakerName).AllowCtx(info.Ctx)
	if err != nil {
		if result.Done != nil {
			result.Done(balancer.DoneInfo{Err: err})
		}
		return emptyPickResult, err
	}

	return balancer.PickResult{
		SubConn: result.SubConn,
		Done:    p.buildDoneFunc(result.Done, promise),
	}, nil
}

func (p *breakerPicker) buildDoneFunc(done func(balancer.DoneInfo), promise breaker.Promise) func(balancer.DoneInfo) {
	return func(info balancer.DoneInfo) {
		if done != nil {
			done(info)
		}

		if info.Err != nil && !codes.Acceptable(info.Err) {
			promise.Reject(info.Err.Error())
		} else {
			promise.Accept()
		}
	}
}

// HasBreaker checks if the circuit breaker is enabled in context.
func HasBreaker(ctx context.Context) bool {
	v, ok := ctx.Value(breakerKey{}).(bool)
	return ok && v
}

// WithBreaker marks the circuit breaker as enabled in context.
func WithBreaker(ctx context.Context) context.Context {
	return context.WithValue(ctx, breakerKey{}, true)
}
