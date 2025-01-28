package internal

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/lang"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/mathx"
	"github.com/zeromicro/go-zero/core/syncx"
	"github.com/zeromicro/go-zero/core/threading"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const coolDownDeviation = 0.05

var (
	registry = Registry{
		clusters: make(map[string]*cluster),
	}
	connManager      = syncx.NewResourceManager()
	coolDownUnstable = mathx.NewUnstable(coolDownDeviation)
	errClosed        = errors.New("etcd monitor chan has been closed")
)

// A Registry is a registry that manages the etcd client connections.
type Registry struct {
	clusters map[string]*cluster
	lock     sync.RWMutex
}

// GetRegistry returns a global Registry.
func GetRegistry() *Registry {
	return &registry
}

// GetConn returns an etcd client connection associated with given endpoints.
func (r *Registry) GetConn(endpoints []string) (EtcdClient, error) {
	c, _ := r.getOrCreateCluster(endpoints)
	return c.getClient()
}

// Monitor monitors the key on given etcd endpoints, notify with the given UpdateListener.
func (r *Registry) Monitor(endpoints []string, key string, exactMatch bool, l UpdateListener) error {
	wkey := watchKey{
		key:        key,
		exactMatch: exactMatch,
	}

	c, exists := r.getOrCreateCluster(endpoints)
	// if exists, the existing values should be updated to the listener.
	if exists {
		c.lock.Lock()
		watcher, ok := c.watchers[wkey]
		if ok {
			watcher.listeners = append(watcher.listeners, l)
		}
		c.lock.Unlock()

		if ok {
			kvs := c.getCurrent(wkey)
			for _, kv := range kvs {
				l.OnAdd(kv)
			}

			return nil
		}
	}

	return c.monitor(wkey, l)
}

func (r *Registry) Unmonitor(endpoints []string, key string, exactMatch bool, l UpdateListener) {
	c, exists := r.getCluster(endpoints)
	if !exists {
		return
	}

	wkey := watchKey{
		key:        key,
		exactMatch: exactMatch,
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	watcher, ok := c.watchers[wkey]
	if !ok {
		return
	}

	for i, listener := range watcher.listeners {
		if listener == l {
			watcher.listeners = append(watcher.listeners[:i], watcher.listeners[i+1:]...)
			break
		}
	}

	if len(watcher.listeners) == 0 {
		if watcher.cancel != nil {
			watcher.cancel()
		}
		delete(c.watchers, wkey)
	}
}

func (r *Registry) getCluster(endpoints []string) (*cluster, bool) {
	clusterKey := getClusterKey(endpoints)

	r.lock.RLock()
	c, ok := r.clusters[clusterKey]
	r.lock.RUnlock()

	return c, ok
}

func (r *Registry) getOrCreateCluster(endpoints []string) (c *cluster, exists bool) {
	c, exists = r.getCluster(endpoints)
	if !exists {
		clusterKey := getClusterKey(endpoints)

		r.lock.Lock()
		defer r.lock.Unlock()

		// double-check locking
		c, exists = r.clusters[clusterKey]
		if !exists {
			c = newCluster(endpoints)
			r.clusters[clusterKey] = c
		}
	}

	return
}

type (
	watchKey struct {
		key        string
		exactMatch bool
	}

	watchValue struct {
		listeners []UpdateListener
		values    map[string]string
		cancel    context.CancelFunc
	}

	cluster struct {
		endpoints  []string
		key        string
		watchers   map[watchKey]*watchValue
		watchGroup *threading.RoutineGroup
		done       chan lang.PlaceholderType
		lock       sync.RWMutex
	}
)

func newCluster(endpoints []string) *cluster {
	return &cluster{
		endpoints:  endpoints,
		key:        getClusterKey(endpoints),
		watchers:   make(map[watchKey]*watchValue),
		watchGroup: threading.NewRoutineGroup(),
		done:       make(chan lang.PlaceholderType),
	}
}

func (c *cluster) addListener(key watchKey, l UpdateListener) {
	c.lock.Lock()
	defer c.lock.Unlock()

	watcher, ok := c.watchers[key]
	if ok {
		watcher.listeners = append(watcher.listeners, l)
		return
	}

	val := newWatchValue()
	val.listeners = []UpdateListener{l}
	c.watchers[key] = val
}

func (c *cluster) getClient() (EtcdClient, error) {
	val, err := connManager.GetResource(c.key, func() (io.Closer, error) {
		return c.newClient()
	})
	if err != nil {
		return nil, err
	}

	return val.(EtcdClient), nil
}

func (c *cluster) getCurrent(key watchKey) []KV {
	c.lock.RLock()
	defer c.lock.RUnlock()

	watcher, ok := c.watchers[key]
	if !ok {
		return nil
	}

	var kvs []KV
	for k, v := range watcher.values {
		kvs = append(kvs, KV{
			Key: k,
			Val: v,
		})
	}

	return kvs
}

func (c *cluster) handleChanges(key watchKey, kvs []KV) {
	c.lock.Lock()
	watcher, ok := c.watchers[key]
	if !ok {
		c.lock.Unlock()
		return
	}

	listeners := append([]UpdateListener(nil), watcher.listeners...)
	// watcher.values cannot be nil
	vals := watcher.values
	newVals := make(map[string]string, len(kvs)+len(vals))
	for _, kv := range kvs {
		newVals[kv.Key] = kv.Val
	}
	add, remove := calculateChanges(vals, newVals)
	watcher.values = newVals
	c.lock.Unlock()

	for _, kv := range add {
		for _, l := range listeners {
			l.OnAdd(kv)
		}
	}
	for _, kv := range remove {
		for _, l := range listeners {
			l.OnDelete(kv)
		}
	}
}

func (c *cluster) handleWatchEvents(ctx context.Context, key watchKey, events []*clientv3.Event) {
	c.lock.RLock()
	watcher, ok := c.watchers[key]
	if !ok {
		c.lock.RUnlock()
		return
	}

	listeners := append([]UpdateListener(nil), watcher.listeners...)
	c.lock.RUnlock()

	for _, ev := range events {
		switch ev.Type {
		case clientv3.EventTypePut:
			c.lock.Lock()
			watcher.values[string(ev.Kv.Key)] = string(ev.Kv.Value)
			c.lock.Unlock()
			for _, l := range listeners {
				l.OnAdd(KV{
					Key: string(ev.Kv.Key),
					Val: string(ev.Kv.Value),
				})
			}
		case clientv3.EventTypeDelete:
			c.lock.Lock()
			delete(watcher.values, string(ev.Kv.Key))
			c.lock.Unlock()
			for _, l := range listeners {
				l.OnDelete(KV{
					Key: string(ev.Kv.Key),
					Val: string(ev.Kv.Value),
				})
			}
		default:
			logc.Errorf(ctx, "Unknown event type: %v", ev.Type)
		}
	}
}

func (c *cluster) load(cli EtcdClient, key watchKey) int64 {
	var resp *clientv3.GetResponse
	for {
		var err error
		ctx, cancel := context.WithTimeout(cli.Ctx(), RequestTimeout)
		if key.exactMatch {
			resp, err = cli.Get(ctx, key.key)
		} else {
			resp, err = cli.Get(ctx, makeKeyPrefix(key.key), clientv3.WithPrefix())
		}

		cancel()
		if err == nil {
			break
		}

		logc.Errorf(cli.Ctx(), "%s, key: %s, exactMatch: %t", err.Error(), key.key, key.exactMatch)
		time.Sleep(coolDownUnstable.AroundDuration(coolDownInterval))
	}

	var kvs []KV
	for _, ev := range resp.Kvs {
		kvs = append(kvs, KV{
			Key: string(ev.Key),
			Val: string(ev.Value),
		})
	}

	c.handleChanges(key, kvs)

	return resp.Header.Revision
}

func (c *cluster) monitor(key watchKey, l UpdateListener) error {
	cli, err := c.getClient()
	if err != nil {
		return err
	}

	c.addListener(key, l)
	rev := c.load(cli, key)
	c.watchGroup.Run(func() {
		c.watch(cli, key, rev)
	})

	return nil
}

func (c *cluster) newClient() (EtcdClient, error) {
	cli, err := NewClient(c.endpoints)
	if err != nil {
		return nil, err
	}

	go c.watchConnState(cli)

	return cli, nil
}

func (c *cluster) reload(cli EtcdClient) {
	c.lock.Lock()
	// cancel the previous watches
	close(c.done)
	c.watchGroup.Wait()
	var keys []watchKey
	for wk, wval := range c.watchers {
		keys = append(keys, wk)
		if wval.cancel != nil {
			wval.cancel()
		}
	}

	c.done = make(chan lang.PlaceholderType)
	c.watchGroup = threading.NewRoutineGroup()
	c.lock.Unlock()

	// start new watches
	for _, key := range keys {
		k := key
		c.watchGroup.Run(func() {
			rev := c.load(cli, k)
			c.watch(cli, k, rev)
		})
	}
}

func (c *cluster) watch(cli EtcdClient, key watchKey, rev int64) {
	for {
		err := c.watchStream(cli, key, rev)
		if err == nil {
			return
		}

		if rev != 0 && errors.Is(err, rpctypes.ErrCompacted) {
			logc.Errorf(cli.Ctx(), "etcd watch stream has been compacted, try to reload, rev %d", rev)
			rev = c.load(cli, key)
		}

		// log the error and retry
		logc.Error(cli.Ctx(), err)
	}
}

func (c *cluster) watchStream(cli EtcdClient, key watchKey, rev int64) error {
	ctx, rch := c.setupWatch(cli, key, rev)

	for {
		select {
		case wresp, ok := <-rch:
			if !ok {
				return errClosed
			}
			if wresp.Canceled {
				return fmt.Errorf("etcd monitor chan has been canceled, error: %w", wresp.Err())
			}
			if wresp.Err() != nil {
				return fmt.Errorf("etcd monitor chan error: %w", wresp.Err())
			}

			c.handleWatchEvents(ctx, key, wresp.Events)
		case <-ctx.Done():
			return nil
		case <-c.done:
			return nil
		}
	}
}

func (c *cluster) setupWatch(cli EtcdClient, key watchKey, rev int64) (context.Context, clientv3.WatchChan) {
	var (
		rch  clientv3.WatchChan
		ops  []clientv3.OpOption
		wkey = key.key
	)

	if !key.exactMatch {
		wkey = makeKeyPrefix(key.key)
		ops = append(ops, clientv3.WithPrefix())
	}
	if rev != 0 {
		ops = append(ops, clientv3.WithRev(rev+1))
	}

	ctx, cancel := context.WithCancel(cli.Ctx())
	if watcher, ok := c.watchers[key]; ok {
		watcher.cancel = cancel
	} else {
		val := newWatchValue()
		val.cancel = cancel

		c.lock.Lock()
		c.watchers[key] = val
		c.lock.Unlock()
	}

	rch = cli.Watch(clientv3.WithRequireLeader(ctx), wkey, ops...)

	return ctx, rch
}

func (c *cluster) watchConnState(cli EtcdClient) {
	watcher := newStateWatcher()
	watcher.addListener(func() {
		go c.reload(cli)
	})
	watcher.watch(cli.ActiveConnection())
}

// DialClient dials an etcd cluster with given endpoints.
func DialClient(endpoints []string) (EtcdClient, error) {
	cfg := clientv3.Config{
		Endpoints:           endpoints,
		AutoSyncInterval:    autoSyncInterval,
		DialTimeout:         DialTimeout,
		RejectOldCluster:    true,
		PermitWithoutStream: true,
	}
	if account, ok := GetAccount(endpoints); ok {
		cfg.Username = account.User
		cfg.Password = account.Pass
	}
	if tlsCfg, ok := GetTLS(endpoints); ok {
		cfg.TLS = tlsCfg
	}

	return clientv3.New(cfg)
}

func calculateChanges(oldVals, newVals map[string]string) (add, remove []KV) {
	for k, v := range newVals {
		if val, ok := oldVals[k]; !ok || v != val {
			add = append(add, KV{
				Key: k,
				Val: v,
			})
		}
	}

	for k, v := range oldVals {
		if val, ok := newVals[k]; !ok || v != val {
			remove = append(remove, KV{
				Key: k,
				Val: v,
			})
		}
	}

	return add, remove
}

func getClusterKey(endpoints []string) string {
	sort.Strings(endpoints)
	return strings.Join(endpoints, endpointsSeparator)
}

func makeKeyPrefix(key string) string {
	return fmt.Sprintf("%s%c", key, Delimiter)
}

// NewClient returns a watchValue that make sure values are not nil.
func newWatchValue() *watchValue {
	return &watchValue{
		values: make(map[string]string),
	}
}
