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
		// gen is the watch generation this key is covered in. reload marks it
		// before the watcher starts, so a monitor racing with the reload does
		// not start a second watcher for the same key.
		gen chan lang.PlaceholderType
	}

	cluster struct {
		endpoints  []string
		key        string
		watchers   map[watchKey]*watchValue
		watchGroup *threading.RoutineGroup
		done       chan lang.PlaceholderType
		// reloadLock serializes generations, so the group and channel are
		// replaced once and a reload does not schedule a key twice.
		reloadLock sync.Mutex
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

	kvs := make([]KV, 0, len(watcher.values))
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
			evKey := string(ev.Kv.Key)
			evVal := string(ev.Kv.Value)
			c.lock.Lock()
			oldVal, exists := watcher.values[evKey]
			watcher.values[evKey] = evVal
			c.lock.Unlock()
			if exists && oldVal == evVal {
				// duplicate PUT with same value, skip to prevent unbounded growth
				continue
			}
			if exists {
				// key moved to a new value, notify delete of old entry first
				for _, l := range listeners {
					l.OnDelete(KV{Key: evKey, Val: oldVal})
				}
			}
			for _, l := range listeners {
				l.OnAdd(KV{Key: evKey, Val: evVal})
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

func (c *cluster) load(done chan lang.PlaceholderType, cli EtcdClient, key watchKey) int64 {
	var resp *clientv3.GetResponse
	for {
		select {
		case <-done:
			return 0
		default:
		}

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
		// The cooldown is interruptible, so reload does not have to wait for a
		// failing load to finish before it can replace the generation.
		select {
		case <-time.After(coolDownUnstable.AroundDuration(coolDownInterval)):
		case <-done:
			return 0
		}
	}

	kvs := make([]KV, 0, len(resp.Kvs))
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

	// load() seeds the listener with the current values, and startWatch() only
	// accepts the generation load() ran against. When a reload wins the race,
	// load() returns at once and the snapshot is taken again in the new
	// generation, so the key ends up with exactly one watcher.
	for {
		done := c.currentDone()
		rev := c.load(done, cli, key)
		if c.startWatch(done, cli, key, rev) {
			return nil
		}
	}
}

func (c *cluster) currentDone() chan lang.PlaceholderType {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.done
}

// startWatch starts a watcher for key in the given generation. It reports
// false only when that generation has been replaced; true means the generation
// is current and the key has or will get a watcher, unless it is no longer
// watched at all.
func (c *cluster) startWatch(done chan lang.PlaceholderType, cli EtcdClient, key watchKey, rev int64) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.done != done {
		return false
	}

	watcher, ok := c.watchers[key]
	if !ok {
		return true
	}
	if watcher.gen == done {
		return true
	}

	watcher.gen = done
	c.watchGroup.Run(func() {
		c.watch(done, cli, key, watcher, rev)
	})

	return true
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
	c.reloadLock.Lock()
	defer c.reloadLock.Unlock()

	// Swap the generation under the lock, but wait for the previous one outside
	// it: watchers take c.lock in handleChanges and handleWatchEvents, so
	// waiting while holding it deadlocks with them.
	type watchStart struct {
		key     watchKey
		watcher *watchValue
	}

	c.lock.Lock()
	oldDone := c.done
	oldGroup := c.watchGroup
	done := make(chan lang.PlaceholderType)
	group := threading.NewRoutineGroup()
	c.done = done
	c.watchGroup = group
	starts := make([]watchStart, 0, len(c.watchers))
	for key, watcher := range c.watchers {
		starts = append(starts, watchStart{key: key, watcher: watcher})
		if watcher.cancel != nil {
			watcher.cancel()
		}
		// Mark the key as covered by this generation before releasing the lock,
		// so a monitor racing with the reload cannot start a duplicate watcher.
		watcher.gen = done
	}
	c.lock.Unlock()

	close(oldDone)
	oldGroup.Wait()

	// start new watches
	for _, start := range starts {
		s := start
		group.Run(func() {
			// The key may have been unmonitored and monitored again while the
			// previous generation was draining. That new watcher owns the key;
			// starting another one here would watch it twice.
			if !c.owns(done, s.key, s.watcher) {
				return
			}

			rev := c.load(done, cli, s.key)
			c.watch(done, cli, s.key, s.watcher, rev)
		})
	}
}

// owns reports whether watcher still serves key in the given generation.
func (c *cluster) owns(done chan lang.PlaceholderType, key watchKey, watcher *watchValue) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.done == done && c.watchers[key] == watcher
}

func (c *cluster) watch(done chan lang.PlaceholderType, cli EtcdClient, key watchKey, watcher *watchValue, rev int64) {
	for {
		err := c.watchStream(done, cli, key, watcher, rev)
		if err == nil {
			return
		}

		if rev != 0 && errors.Is(err, rpctypes.ErrCompacted) {
			logc.Errorf(cli.Ctx(), "etcd watch stream has been compacted, try to reload, rev %d", rev)
			rev = c.load(done, cli, key)
			select {
			case <-done:
				return
			default:
			}
		}

		// log the error and retry with an interruptible cooldown to prevent
		// CPU/disk exhaustion
		logc.Error(cli.Ctx(), err)
		select {
		case <-time.After(coolDownUnstable.AroundDuration(coolDownInterval)):
		case <-done:
			return
		}
	}
}

func (c *cluster) watchStream(done chan lang.PlaceholderType, cli EtcdClient, key watchKey, watcher *watchValue, rev int64) error {
	ctx, cancel, rch, ok := c.setupWatch(done, cli, key, watcher, rev)
	if !ok {
		return nil
	}
	// Cancel the stream once this watcher stops consuming it, otherwise clientv3
	// keeps the substream attached to the shared client.
	defer cancel()

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
		case <-done:
			return nil
		}
	}
}

func (c *cluster) setupWatch(done chan lang.PlaceholderType, cli EtcdClient, key watchKey, watcher *watchValue, rev int64) (context.Context, context.CancelFunc, clientv3.WatchChan, bool) {
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

	// Only the watcher that still owns the key may open a stream. Between the
	// ownership checks in reload and this point the key can be unmonitored and
	// monitored again; that new watcher owns the key now.
	c.lock.Lock()
	if c.done != done || c.watchers[key] != watcher {
		c.lock.Unlock()
		cancel()

		return nil, nil, nil, false
	}
	watcher.cancel = cancel
	c.lock.Unlock()

	rch = cli.Watch(clientv3.WithRequireLeader(ctx), wkey, ops...)

	return ctx, cancel, rch, true
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

// newWatchValue returns a watchValue that make sure values are not nil.
func newWatchValue() *watchValue {
	return &watchValue{
		values: make(map[string]string),
	}
}
