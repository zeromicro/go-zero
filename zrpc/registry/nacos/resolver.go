package nacos

import (
	"context"
	"fmt"
	"sort"

	"github.com/nacos-group/nacos-sdk-go/common/logger"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/resolver"
)

type resolvr struct {
	cancelFunc context.CancelFunc
}

func (r *resolvr) ResolveNow(resolver.ResolveNowOptions) {}

// Close closes the resolver.
func (r *resolvr) Close() {
	r.cancelFunc()
}

type watcher struct {
	ctx    context.Context
	cancel context.CancelFunc
	out    chan<- []string
}

func newWatcher(ctx context.Context, cancel context.CancelFunc, out chan<- []string) *watcher {
	return &watcher{
		ctx:    ctx,
		cancel: cancel,
		out:    out,
	}
}

func (nw *watcher) CallBackHandle(services []model.SubscribeService, err error) {
	if err != nil {
		logger.Error("[Nacos resolver] watcher call back handle error:%v", err)
		return
	}
	ee := make([]string, 0, len(services))
	for _, s := range services {
		ee = append(ee, fmt.Sprintf("%s:%d", s.Ip, s.Port))
	}
	nw.out <- ee
}

func populateEndpoints(ctx context.Context, clientConn resolver.ClientConn, input <-chan []string) {
	for {
		select {
		case cc := <-input:
			connsSet := make(map[string]struct{}, len(cc))
			for _, c := range cc {
				connsSet[c] = struct{}{}
			}
			conns := make([]resolver.Address, 0, len(connsSet))
			for c := range connsSet {
				conns = append(conns, resolver.Address{Addr: c})
			}
			sort.Sort(byAddressString(conns)) // Don't replace the same address list in the balancer
			_ = clientConn.UpdateState(resolver.State{Addresses: conns})
		case <-ctx.Done():
			logx.Info("[Nacos resolver] Watch has been finished")
			return
		}
	}
}

// byAddressString sorts resolver.Address by Address Field  sorting in increasing order.
type byAddressString []resolver.Address

func (p byAddressString) Len() int           { return len(p) }
func (p byAddressString) Less(i, j int) bool { return p[i].Addr < p[j].Addr }
func (p byAddressString) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
