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

	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/zeromicro/go-zero/core/contextx"
	"github.com/zeromicro/go-zero/core/lang"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/syncx"
	"github.com/zeromicro/go-zero/core/threading"
)

var (
	registry = Registry{
		clusters: make(map[string]*cluster),
	}
	connManager = syncx.NewResourceManager()
	errClosed   = errors.New("etcd monitor chan has been closed")
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
	c, _ := r.getCluster(endpoints)
	return c.getClient()
}

// Monitor monitors the key on given etcd endpoints, notify with the given UpdateListener.
func (r *Registry) Monitor(endpoints []string, key string, l UpdateListener, exactMatch bool) error {
	c, exists := r.getCluster(endpoints)
	// if exists, the existing values should be updated to the listener.
	if exists {
		kvs := c.getCurrent(key)
		for _, kv := range kvs {
			l.OnAdd(kv)
		}
	}

	return c.monitor(key, l, exactMatch)
}

func (r *Registry) getCluster(endpoints []string) (c *cluster, exists bool) {
	clusterKey := getClusterKey(endpoints)
	r.lock.RLock()
	c, exists = r.clusters[clusterKey]
	r.lock.RUnlock()

	if !exists {
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

type cluster struct {
	endpoints  []string
	key        string
	values     map[string]map[string]string
	listeners  map[string][]UpdateListener
	watchGroup *threading.RoutineGroup
	done       chan lang.PlaceholderType
	lock       sync.RWMutex
	exactMatch bool
}

func newCluster(endpoints []string) *cluster {
	return &cluster{
		endpoints:  endpoints,
		key:        getClusterKey(endpoints),
		values:     make(map[string]map[string]string),
		listeners:  make(map[string][]UpdateListener),
		watchGroup: threading.NewRoutineGroup(),
		done:       make(chan lang.PlaceholderType),
	}
}

func (c *cluster) context(cli EtcdClient) context.Context {
	return contextx.ValueOnlyFrom(cli.Ctx())
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

func (c *cluster) getCurrent(key string) []KV {
	c.lock.RLock()
	defer c.lock.RUnlock()

	var kvs []KV
	for k, v := range c.values[key] {
		kvs = append(kvs, KV{
			Key: k,
			Val: v,
		})
	}

	return kvs
}

func (c *cluster) handleChanges(key string, kvs []KV) {
	var add []KV
	var remove []KV

	c.lock.Lock()
	listeners := append([]UpdateListener(nil), c.listeners[key]...)
	vals, ok := c.values[key]
	if !ok {
		add = kvs
		vals = make(map[string]string)
		for _, kv := range kvs {
			vals[kv.Key] = kv.Val
		}
		c.values[key] = vals
	} else {
		m := make(map[string]string)
		for _, kv := range kvs {
			m[kv.Key] = kv.Val
		}
		for k, v := range vals {
			if val, ok := m[k]; !ok || v != val {
				remove = append(remove, KV{
					Key: k,
					Val: v,
				})
			}
		}
		for k, v := range m {
			if val, ok := vals[k]; !ok || v != val {
				add = append(add, KV{
					Key: k,
					Val: v,
				})
			}
		}
		c.values[key] = m
	}
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

func (c *cluster) handleWatchEvents(key string, events []*clientv3.Event) {
	c.lock.RLock()
	listeners := append([]UpdateListener(nil), c.listeners[key]...)
	c.lock.RUnlock()

	for _, ev := range events {
		switch ev.Type {
		case clientv3.EventTypePut:
			c.lock.Lock()
			if vals, ok := c.values[key]; ok {
				vals[string(ev.Kv.Key)] = string(ev.Kv.Value)
			} else {
				c.values[key] = map[string]string{string(ev.Kv.Key): string(ev.Kv.Value)}
			}
			c.lock.Unlock()
			for _, l := range listeners {
				l.OnAdd(KV{
					Key: string(ev.Kv.Key),
					Val: string(ev.Kv.Value),
				})
			}
		case clientv3.EventTypeDelete:
			c.lock.Lock()
			if vals, ok := c.values[key]; ok {
				delete(vals, string(ev.Kv.Key))
			}
			c.lock.Unlock()
			for _, l := range listeners {
				l.OnDelete(KV{
					Key: string(ev.Kv.Key),
					Val: string(ev.Kv.Value),
				})
			}
		default:
			logx.Errorf("Unknown event type: %v", ev.Type)
		}
	}
}

func (c *cluster) load(cli EtcdClient, key string) int64 {
	var resp *clientv3.GetResponse
	for {
		var err error
		ctx, cancel := context.WithTimeout(c.context(cli), RequestTimeout)
		if c.exactMatch {
			resp, err = cli.Get(ctx, key)
		} else {
			resp, err = cli.Get(ctx, makeKeyPrefix(key), clientv3.WithPrefix())
		}

		cancel()
		if err == nil {
			break
		}

		logx.Errorf("%s, key is %s", err.Error(), key)
		time.Sleep(coolDownInterval)
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

func (c *cluster) monitor(key string, l UpdateListener, exactMatch bool) error {
	c.lock.Lock()
	c.listeners[key] = append(c.listeners[key], l)
	c.exactMatch = exactMatch
	c.lock.Unlock()

	cli, err := c.getClient()
	if err != nil {
		return err
	}

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
	close(c.done)
	c.watchGroup.Wait()
	c.done = make(chan lang.PlaceholderType)
	c.watchGroup = threading.NewRoutineGroup()
	var keys []string
	for k := range c.listeners {
		keys = append(keys, k)
	}
	c.lock.Unlock()

	for _, key := range keys {
		k := key
		c.watchGroup.Run(func() {
			rev := c.load(cli, k)
			c.watch(cli, k, rev)
		})
	}
}

func (c *cluster) watch(cli EtcdClient, key string, rev int64) {
	for {
		err := c.watchStream(cli, key, rev)
		if err == nil {
			return
		}

		if rev != 0 && errors.Is(err, rpctypes.ErrCompacted) {
			logx.Errorf("etcd watch stream has been compacted, try to reload, rev %d", rev)
			rev = c.load(cli, key)
		}

		// log the error and retry
		logx.Error(err)
	}
}

func (c *cluster) watchStream(cli EtcdClient, key string, rev int64) error {
	var (
		rch      clientv3.WatchChan
		ops      []clientv3.OpOption
		watchKey = key
	)
	if !c.exactMatch {
		watchKey = makeKeyPrefix(key)
		ops = append(ops, clientv3.WithPrefix())
	}
	if rev != 0 {
		ops = append(ops, clientv3.WithRev(rev+1))
	}

	rch = cli.Watch(clientv3.WithRequireLeader(c.context(cli)), watchKey, ops...)

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

			c.handleWatchEvents(key, wresp.Events)
		case <-c.done:
			return nil
		}
	}
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

func getClusterKey(endpoints []string) string {
	sort.Strings(endpoints)
	return strings.Join(endpoints, endpointsSeparator)
}

func makeKeyPrefix(key string) string {
	return fmt.Sprintf("%s%c", key, Delimiter)
}
