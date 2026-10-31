package breaker

import (
	"path"
	"sync"

	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/core/logx"
	_ "github.com/zeromicro/go-zero/zrpc/internal/balancer/consistenthash" // ensure consistenthash is registered
	"github.com/zeromicro/go-zero/zrpc/internal/balancer/p2c"
	"github.com/zeromicro/go-zero/zrpc/internal/codes"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/resolver"
)

const (
	// BalancerSuffix is the suffix for breaker-enabled balancer.
	BalancerSuffix = "_breaker"

	// strategyService uses target/method as breaker name (all instances share one breaker)
	strategyService = "service"
	// strategyInstance uses addr/method as breaker name (each instance has its own breaker)
	strategyInstance = "instance"
)

var (
	emptyPickResult balancer.PickResult
	registerLock    sync.Mutex
)

// GetBalancerName returns the balancer name for the given base name and strategy.
func GetBalancerName(baseName, strategy string) string {
	if baseName == "" {
		baseName = p2c.Name
	}
	if strategy == strategyInstance {
		return baseName + BalancerSuffix + "_" + strategyInstance
	}
	return baseName + BalancerSuffix
}

// Register registers a breaker-enabled balancer that wraps the base balancer.
// It's safe to call multiple times - if already registered, it returns immediately.
func Register(baseName, strategy string, retryTimes int) {
	registerLock.Lock()
	defer registerLock.Unlock()

	if len(baseName) == 0 {
		baseName = p2c.Name
	}
	if retryTimes < 0 {
		retryTimes = 0
	}

	name := GetBalancerName(baseName, strategy)
	if balancer.Get(name) != nil {
		return
	}

	baseBuilder := balancer.Get(baseName)
	if baseBuilder == nil {
		logx.Errorf("balancer %q not found, make sure it's registered before creating client", baseName)
		return
	}

	balancer.Register(&breakerBuilder{
		baseBuilder: baseBuilder,
		name:        name,
		strategy:    strategy,
		retryTimes:  retryTimes,
	})
}

// breakerBuilder wraps a base balancer builder with circuit breaker capability.
type breakerBuilder struct {
	baseBuilder balancer.Builder
	name        string
	strategy    string
	retryTimes  int
}

func (b *breakerBuilder) Name() string {
	return b.name
}

func (b *breakerBuilder) Build(cc balancer.ClientConn, opts balancer.BuildOptions) balancer.Balancer {
	wrappedCC := &breakerClientConn{
		ClientConn: cc,
		conns:      make(map[balancer.SubConn]string),
		target:     extractTarget(opts),
		strategy:   b.strategy,
		retryTimes: b.retryTimes,
	}
	return b.baseBuilder.Build(wrappedCC, opts)
}

// extractTarget extracts the target name from BuildOptions.
// It first tries to use the URL path, falling back to the endpoint if empty.
func extractTarget(opts balancer.BuildOptions) string {
	target := opts.Target.URL.Path
	if len(target) > 0 && target[0] == '/' {
		target = target[1:]
	}
	if len(target) == 0 {
		target = opts.Target.Endpoint()
	}
	return target
}

// breakerClientConn wraps ClientConn to track SubConn addresses and wrap pickers.
type breakerClientConn struct {
	balancer.ClientConn
	lock       sync.RWMutex
	conns      map[balancer.SubConn]string
	target     string
	strategy   string
	retryTimes int
}

// NewSubConn creates a new SubConn and tracks its address for breaker naming.
func (cc *breakerClientConn) NewSubConn(addrs []resolver.Address, opts balancer.NewSubConnOptions) (balancer.SubConn, error) {
	sc, err := cc.ClientConn.NewSubConn(addrs, opts)
	if err != nil {
		return nil, err
	}

	cc.lock.Lock()
	if len(addrs) > 0 {
		cc.conns[sc] = addrs[0].Addr
	}
	cc.lock.Unlock()

	return sc, nil
}

// RemoveSubConn removes the SubConn from tracking and delegates to the underlying ClientConn.
func (cc *breakerClientConn) RemoveSubConn(sc balancer.SubConn) {
	cc.lock.Lock()
	delete(cc.conns, sc)
	cc.lock.Unlock()

	cc.ClientConn.RemoveSubConn(sc)
}

// UpdateState wraps the picker with breaker logic and updates the state.
func (cc *breakerClientConn) UpdateState(state balancer.State) {
	cc.lock.RLock()
	conns := make(map[balancer.SubConn]string, len(cc.conns))
	for k, v := range cc.conns {
		conns[k] = v
	}
	cc.lock.RUnlock()

	// Wrap the picker with breaker logic
	state.Picker = &breakerPicker{
		picker:     state.Picker,
		conns:      conns,
		target:     cc.target,
		strategy:   cc.strategy,
		retryTimes: cc.retryTimes,
	}
	cc.ClientConn.UpdateState(state)
}

// breakerPicker wraps a picker with circuit breaker logic.
type breakerPicker struct {
	picker     balancer.Picker
	conns      map[balancer.SubConn]string
	target     string
	strategy   string
	retryTimes int
}

// Pick selects a SubConn with circuit breaker protection.
// It delegates to pickWithServiceBreaker or pickWithInstanceBreaker based on strategy.
func (p *breakerPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if p.strategy == strategyService {
		return p.pickWithServiceBreaker(info)
	}
	return p.pickWithInstanceBreaker(info)
}

// pickWithServiceBreaker uses target/method as breaker name.
// All instances share one breaker for the same target/method combination.
func (p *breakerPicker) pickWithServiceBreaker(info balancer.PickInfo) (balancer.PickResult, error) {
	breakerName := path.Join(p.target, info.FullMethodName)
	promise, err := breaker.GetBreaker(breakerName).AllowCtx(info.Ctx)
	if err != nil {
		return emptyPickResult, err
	}

	result, err := p.picker.Pick(info)
	if err != nil {
		promise.Accept()
		return result, err
	}

	result.Done = p.buildDoneFunc(result.Done, promise)
	return result, nil
}

// pickWithInstanceBreaker uses addr/method as breaker name.
// Each instance has its own breaker, allowing retry on different instances when one is broken.
func (p *breakerPicker) pickWithInstanceBreaker(info balancer.PickInfo) (balancer.PickResult, error) {
	// First attempt
	result, err := p.picker.Pick(info)
	if err != nil {
		return result, err
	}

	addr := p.conns[result.SubConn]
	if addr == "" {
		return result, nil
	}

	promise, err := breaker.GetBreaker(path.Join(addr, info.FullMethodName)).AllowCtx(info.Ctx)
	if err == nil {
		// First instance is available, return directly
		result.Done = p.buildDoneFunc(result.Done, promise)
		return result, nil
	}

	// First instance is broken, try other instances
	metricInstanceBreakerTriggered.Inc(addr, info.FullMethodName)
	triedAddrs := map[string]struct{}{addr: {}}
	lastBreakerErr := err

	for i := 0; i < p.retryTimes; i++ {
		result, err := p.picker.Pick(info)
		if err != nil {
			return result, err
		}

		addr := p.conns[result.SubConn]
		if addr == "" {
			return result, nil
		}
		if _, tried := triedAddrs[addr]; tried {
			continue
		}
		triedAddrs[addr] = struct{}{}

		promise, err := breaker.GetBreaker(path.Join(addr, info.FullMethodName)).AllowCtx(info.Ctx)
		if err != nil {
			lastBreakerErr = err
			metricInstanceBreakerTriggered.Inc(addr, info.FullMethodName)
			continue
		}

		result.Done = p.buildDoneFunc(result.Done, promise)
		return result, nil
	}

	return emptyPickResult, lastBreakerErr
}

// buildDoneFunc wraps the original Done callback with breaker result reporting.
// It calls Accept for successful requests or acceptable errors, Reject otherwise.
func (p *breakerPicker) buildDoneFunc(
	originalDone func(balancer.DoneInfo),
	promise breaker.Promise,
) func(balancer.DoneInfo) {
	return func(di balancer.DoneInfo) {
		if originalDone != nil {
			originalDone(di)
		}

		if di.Err != nil && !codes.Acceptable(di.Err) {
			promise.Reject(di.Err.Error())
		} else {
			promise.Accept()
		}
	}
}
