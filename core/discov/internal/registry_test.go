package internal

import (
	"context"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/contextx"
	"github.com/zeromicro/go-zero/core/lang"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stringx"
	"github.com/zeromicro/go-zero/core/threading"
	"go.etcd.io/etcd/api/v3/etcdserverpb"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/mock/mockserver"
	"go.uber.org/mock/gomock"
)

var mockLock sync.Mutex

func init() {
	logx.Disable()
}

func setMockClient(cli EtcdClient) func() {
	mockLock.Lock()
	NewClient = func([]string) (EtcdClient, error) {
		return cli, nil
	}
	return func() {
		NewClient = DialClient
		mockLock.Unlock()
	}
}

func TestGetCluster(t *testing.T) {
	AddAccount([]string{"first"}, "foo", "bar")
	c1, _ := GetRegistry().getOrCreateCluster([]string{"first"})
	c2, _ := GetRegistry().getOrCreateCluster([]string{"second"})
	c3, _ := GetRegistry().getOrCreateCluster([]string{"first"})
	assert.Equal(t, c1, c3)
	assert.NotEqual(t, c1, c2)
}

func TestGetClusterKey(t *testing.T) {
	assert.Equal(t, getClusterKey([]string{"localhost:1234", "remotehost:5678"}),
		getClusterKey([]string{"remotehost:5678", "localhost:1234"}))
}

func TestUnmonitor(t *testing.T) {
	t.Run("no listener", func(t *testing.T) {
		reg := &Registry{
			clusters: map[string]*cluster{},
		}
		assert.NotPanics(t, func() {
			reg.Unmonitor([]string{"any"}, "any", false, nil)
		})
	})

	t.Run("no value", func(t *testing.T) {
		reg := &Registry{
			clusters: map[string]*cluster{
				"any": {
					watchers: map[watchKey]*watchValue{
						{
							key: "any",
						}: {
							values: map[string]string{},
						},
					},
				},
			},
		}
		assert.NotPanics(t, func() {
			reg.Unmonitor([]string{"any"}, "another", false, nil)
		})
	})
}

func TestCluster_HandleChanges(t *testing.T) {
	ctrl := gomock.NewController(t)
	l := NewMockUpdateListener(ctrl)
	l.EXPECT().OnAdd(KV{
		Key: "first",
		Val: "1",
	})
	l.EXPECT().OnAdd(KV{
		Key: "second",
		Val: "2",
	})
	l.EXPECT().OnDelete(KV{
		Key: "first",
		Val: "1",
	})
	l.EXPECT().OnDelete(KV{
		Key: "second",
		Val: "2",
	})
	l.EXPECT().OnAdd(KV{
		Key: "third",
		Val: "3",
	})
	l.EXPECT().OnAdd(KV{
		Key: "fourth",
		Val: "4",
	})
	c := newCluster([]string{"any"})
	key := watchKey{
		key:        "any",
		exactMatch: false,
	}
	c.watchers[key] = &watchValue{
		listeners: []UpdateListener{l},
	}
	c.handleChanges(key, []KV{
		{
			Key: "first",
			Val: "1",
		},
		{
			Key: "second",
			Val: "2",
		},
	})
	assert.EqualValues(t, map[string]string{
		"first":  "1",
		"second": "2",
	}, c.watchers[key].values)
	c.handleChanges(key, []KV{
		{
			Key: "third",
			Val: "3",
		},
		{
			Key: "fourth",
			Val: "4",
		},
	})
	assert.EqualValues(t, map[string]string{
		"third":  "3",
		"fourth": "4",
	}, c.watchers[key].values)
}

func TestCluster_Load(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cli := NewMockEtcdClient(ctrl)
	restore := setMockClient(cli)
	defer restore()
	cli.EXPECT().Get(gomock.Any(), "any/", gomock.Any()).Return(&clientv3.GetResponse{
		Header: &etcdserverpb.ResponseHeader{},
		Kvs: []*mvccpb.KeyValue{
			{
				Key:   []byte("hello"),
				Value: []byte("world"),
			},
		},
	}, nil)
	cli.EXPECT().Ctx().Return(context.Background())
	c := &cluster{
		watchers: make(map[watchKey]*watchValue),
	}
	done := make(chan lang.PlaceholderType)
	defer close(done)
	c.load(done, cli, watchKey{
		key: "any",
	})
}

func TestCluster_Watch(t *testing.T) {
	tests := []struct {
		name      string
		method    int
		eventType mvccpb.Event_EventType
	}{
		{
			name:      "add",
			eventType: clientv3.EventTypePut,
		},
		{
			name:      "delete",
			eventType: clientv3.EventTypeDelete,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			cli := NewMockEtcdClient(ctrl)
			restore := setMockClient(cli)
			defer restore()
			ch := make(chan clientv3.WatchResponse)
			cli.EXPECT().Watch(gomock.Any(), "any/", gomock.Any()).Return(ch)
			cli.EXPECT().Ctx().Return(context.Background())
			var wg sync.WaitGroup
			wg.Add(1)
			c := &cluster{
				watchers: make(map[watchKey]*watchValue),
			}
			key := watchKey{
				key: "any",
			}
			listener := NewMockUpdateListener(ctrl)
			c.watchers[key] = &watchValue{
				listeners: []UpdateListener{listener},
				values:    make(map[string]string),
			}
			listener.EXPECT().OnAdd(gomock.Any()).Do(func(kv KV) {
				assert.Equal(t, "hello", kv.Key)
				assert.Equal(t, "world", kv.Val)
				wg.Done()
			}).MaxTimes(1)
			listener.EXPECT().OnDelete(gomock.Any()).Do(func(_ any) {
				wg.Done()
			}).MaxTimes(1)
			done := make(chan lang.PlaceholderType)
			defer close(done)
			c.done = done
			go c.watch(done, cli, key, c.watchers[key], 0)
			ch <- clientv3.WatchResponse{
				Events: []*clientv3.Event{
					{
						Type: test.eventType,
						Kv: &mvccpb.KeyValue{
							Key:   []byte("hello"),
							Value: []byte("world"),
						},
					},
				},
			}
			wg.Wait()
		})
	}
}

func TestClusterWatch_RespFailures(t *testing.T) {
	resps := []clientv3.WatchResponse{
		{
			Canceled: true,
		},
		{
			// cause resp.Err() != nil
			CompactRevision: 1,
		},
	}

	for _, resp := range resps {
		t.Run(stringx.Rand(), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			cli := NewMockEtcdClient(ctrl)
			restore := setMockClient(cli)
			defer restore()
			ch := make(chan clientv3.WatchResponse)
			cli.EXPECT().Watch(gomock.Any(), "any/", gomock.Any()).Return(ch).AnyTimes()
			cli.EXPECT().Ctx().Return(context.Background()).AnyTimes()
			c := &cluster{
				watchers: make(map[watchKey]*watchValue),
			}
			c.done = make(chan lang.PlaceholderType)
			key := watchKey{
				key: "any",
			}
			c.watchers[key] = newWatchValue()
			go func() {
				ch <- resp
				close(c.done)
			}()
			c.watch(c.done, cli, key, c.watchers[key], 0)
		})
	}
}

func TestCluster_getCurrent(t *testing.T) {
	t.Run("no value", func(t *testing.T) {
		c := &cluster{
			watchers: map[watchKey]*watchValue{
				{
					key: "any",
				}: {
					values: map[string]string{},
				},
			},
		}
		assert.Nil(t, c.getCurrent(watchKey{
			key: "another",
		}))
	})
}

func TestCluster_handleWatchEvents(t *testing.T) {
	t.Run("no value", func(t *testing.T) {
		c := &cluster{
			watchers: map[watchKey]*watchValue{
				{
					key: "any",
				}: {
					values: map[string]string{},
				},
			},
		}
		assert.NotPanics(t, func() {
			c.handleWatchEvents(context.Background(), watchKey{
				key: "another",
			}, nil)
		})
	})
}

func TestCluster_addListener(t *testing.T) {
	t.Run("has listener", func(t *testing.T) {
		c := &cluster{
			watchers: map[watchKey]*watchValue{
				{
					key: "any",
				}: {
					listeners: make([]UpdateListener, 0),
				},
			},
		}
		assert.NotPanics(t, func() {
			c.addListener(watchKey{
				key: "any",
			}, nil)
		})
	})

	t.Run("no listener", func(t *testing.T) {
		c := &cluster{
			watchers: map[watchKey]*watchValue{
				{
					key: "any",
				}: {
					listeners: make([]UpdateListener, 0),
				},
			},
		}
		assert.NotPanics(t, func() {
			c.addListener(watchKey{
				key: "another",
			}, nil)
		})
	})
}

func TestCluster_reload(t *testing.T) {
	c := &cluster{
		watchers:   map[watchKey]*watchValue{},
		watchGroup: threading.NewRoutineGroup(),
		done:       make(chan lang.PlaceholderType),
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cli := NewMockEtcdClient(ctrl)
	restore := setMockClient(cli)
	defer restore()
	assert.NotPanics(t, func() {
		c.reload(cli)
	})
}

func TestCluster_ReloadDoesNotHoldLockWhileWaiting(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cli := NewMockEtcdClient(ctrl)
	cli.EXPECT().Ctx().Return(context.Background()).AnyTimes()
	cli.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(
		&clientv3.GetResponse{Header: &etcdserverpb.ResponseHeader{}}, nil).AnyTimes()
	cli.EXPECT().Watch(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(make(chan clientv3.WatchResponse)).AnyTimes()

	key := watchKey{key: "any"}
	c := newCluster([]string{"reload-lock-test"})
	c.watchers[key] = &watchValue{
		listeners: []UpdateListener{new(mockListener)},
		values:    make(map[string]string),
	}

	// A watcher finishing a load takes c.lock in handleChanges, then keeps
	// running until it is done. reload must wait for it without holding the lock.
	startHandle := make(chan struct{})
	handleDone := make(chan struct{})
	finishWatcher := make(chan struct{})
	watcherDone := make(chan struct{})
	c.watchGroup.Run(func() {
		<-startHandle
		c.handleChanges(key, nil)
		close(handleDone)
		<-finishWatcher
		close(watcherDone)
	})

	done := c.currentDone()
	reloaded := make(chan struct{})
	go func() {
		c.reload(cli)
		close(reloaded)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("reload did not replace the generation")
	}

	// The watcher now reaches handleChanges, which takes c.lock. reload must not
	// be holding it while waiting for the watcher to exit.
	close(startHandle)

	select {
	case <-handleDone:
	case <-time.After(3 * time.Second):
		t.Fatal("reload deadlocked waiting for the watcher")
	}

	// The watcher is still running, so reload must not have returned yet.
	select {
	case <-reloaded:
		t.Fatal("reload returned before the old watcher exited")
	case <-time.After(100 * time.Millisecond):
	}

	close(finishWatcher)
	select {
	case <-reloaded:
	case <-time.After(3 * time.Second):
		t.Fatal("reload did not wait for the watcher")
	}

	select {
	case <-watcherDone:
	case <-time.After(time.Second):
		t.Fatal("reload returned before the old watcher exited")
	}

	c.lock.Lock()
	group := c.watchGroup
	close(c.done)
	c.lock.Unlock()
	group.Wait()
}

func TestCluster_StartWatchOncePerGeneration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cli := NewMockEtcdClient(ctrl)
	cli.EXPECT().Ctx().Return(context.Background()).AnyTimes()
	watchCalled := make(chan struct{}, 2)
	cli.EXPECT().Watch(gomock.Any(), "any/", gomock.Any()).DoAndReturn(
		func(context.Context, string, ...clientv3.OpOption) clientv3.WatchChan {
			watchCalled <- struct{}{}

			return make(chan clientv3.WatchResponse)
		}).AnyTimes()

	key := watchKey{key: "any"}
	c := newCluster([]string{"start-watch-test"})
	c.watchers[key] = newWatchValue()

	done := c.currentDone()
	if !c.startWatch(done, cli, key, 0) {
		t.Fatal("the current generation should accept the watcher")
	}

	select {
	case <-watchCalled:
	case <-time.After(2 * time.Second):
		t.Fatal("watcher did not start")
	}

	if !c.startWatch(done, cli, key, 0) {
		t.Fatal("the generation is still current")
	}

	select {
	case <-watchCalled:
		t.Fatal("a second watcher was started for the same key and generation")
	case <-time.After(50 * time.Millisecond):
	}

	c.lock.Lock()
	c.done = make(chan lang.PlaceholderType)
	c.lock.Unlock()
	assert.False(t, c.startWatch(done, cli, key, 0))

	close(done)
	c.watchGroup.Wait()
}

func TestCluster_ReloadDoesNotDuplicateWatcher(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cli := NewMockEtcdClient(ctrl)
	cli.EXPECT().Ctx().Return(context.Background()).AnyTimes()
	var watches int32
	cli.EXPECT().Watch(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(context.Context, string, ...clientv3.OpOption) clientv3.WatchChan {
			atomic.AddInt32(&watches, 1)

			return make(chan clientv3.WatchResponse)
		}).AnyTimes()
	cli.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(
		&clientv3.GetResponse{Header: &etcdserverpb.ResponseHeader{}}, nil).AnyTimes()

	key := watchKey{key: "any"}
	c := newCluster([]string{"reload-remonitor-test"})
	listener := new(mockListener)
	c.watchers[key] = &watchValue{
		listeners: []UpdateListener{listener},
		values:    make(map[string]string),
	}

	// Keep reload inside oldGroup.Wait while the key is unmonitored and
	// monitored again.
	gate := make(chan struct{})
	c.watchGroup.Run(func() {
		<-gate
	})

	oldDone := c.currentDone()
	reloaded := make(chan struct{})
	go func() {
		c.reload(cli)
		close(reloaded)
	}()

	select {
	case <-oldDone:
	case <-time.After(2 * time.Second):
		t.Fatal("reload did not replace the generation")
	}

	// Unmonitor, then monitor again: this creates a new watchValue and lets its
	// own monitor start the watcher for the new generation.
	c.lock.Lock()
	delete(c.watchers, key)
	c.lock.Unlock()
	c.addListener(key, listener)
	assert.True(t, c.startWatch(c.currentDone(), cli, key, 0))

	assert.Eventually(t, func() bool {
		return atomic.LoadInt32(&watches) == 1
	}, time.Second, time.Millisecond)

	close(gate)
	select {
	case <-reloaded:
	case <-time.After(3 * time.Second):
		t.Fatal("reload did not finish")
	}

	// The reload must leave the key to the watcher that already owns it.
	time.Sleep(100 * time.Millisecond)
	assert.EqualValues(t, 1, atomic.LoadInt32(&watches))

	c.lock.Lock()
	group := c.watchGroup
	close(c.done)
	c.lock.Unlock()
	group.Wait()
}

func TestCluster_ReloadStartsWatcher(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cli := NewMockEtcdClient(ctrl)
	cli.EXPECT().Ctx().Return(context.Background()).AnyTimes()
	cli.EXPECT().Get(gomock.Any(), "any/", gomock.Any()).Return(
		&clientv3.GetResponse{Header: &etcdserverpb.ResponseHeader{}}, nil).AnyTimes()
	var watches int32
	cli.EXPECT().Watch(gomock.Any(), "any/", gomock.Any()).DoAndReturn(
		func(context.Context, string, ...clientv3.OpOption) clientv3.WatchChan {
			atomic.AddInt32(&watches, 1)

			return make(chan clientv3.WatchResponse)
		}).AnyTimes()

	key := watchKey{key: "any"}
	c := newCluster([]string{"reload-start-test"})
	c.watchers[key] = &watchValue{
		listeners: []UpdateListener{new(mockListener)},
		values:    make(map[string]string),
	}

	c.reload(cli)

	assert.Eventually(t, func() bool {
		return atomic.LoadInt32(&watches) == 1
	}, time.Second, time.Millisecond)

	c.lock.Lock()
	group := c.watchGroup
	close(c.done)
	c.lock.Unlock()
	group.Wait()
}

func TestCluster_ReloadSkipsUnmonitoredKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cli := NewMockEtcdClient(ctrl)
	cli.EXPECT().Ctx().Return(context.Background()).AnyTimes()
	cli.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(
		&clientv3.GetResponse{Header: &etcdserverpb.ResponseHeader{}}, nil).AnyTimes()
	var watches int32
	cli.EXPECT().Watch(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(context.Context, string, ...clientv3.OpOption) clientv3.WatchChan {
			atomic.AddInt32(&watches, 1)

			return make(chan clientv3.WatchResponse)
		}).AnyTimes()

	key := watchKey{key: "any"}
	c := newCluster([]string{"reload-unmonitor-test"})
	c.watchers[key] = &watchValue{
		listeners: []UpdateListener{new(mockListener)},
		values:    make(map[string]string),
	}

	gate := make(chan struct{})
	c.watchGroup.Run(func() {
		<-gate
	})

	oldDone := c.currentDone()
	reloaded := make(chan struct{})
	go func() {
		c.reload(cli)
		close(reloaded)
	}()

	select {
	case <-oldDone:
	case <-time.After(2 * time.Second):
		t.Fatal("reload did not replace the generation")
	}

	c.lock.Lock()
	delete(c.watchers, key)
	c.lock.Unlock()
	close(gate)

	select {
	case <-reloaded:
	case <-time.After(3 * time.Second):
		t.Fatal("reload did not finish")
	}

	// The key has no listeners anymore, so reload must not start a watcher for
	// it; that stream would live until the next reload.
	time.Sleep(100 * time.Millisecond)
	assert.EqualValues(t, 0, atomic.LoadInt32(&watches))

	c.lock.Lock()
	group := c.watchGroup
	close(c.done)
	c.lock.Unlock()
	group.Wait()
}

func TestCluster_ReloadDoesNotStartWatcherAfterReplacement(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cli := NewMockEtcdClient(ctrl)
	cli.EXPECT().Ctx().Return(context.Background()).AnyTimes()

	// Hold reload's load until the key has been replaced.
	getEntered := make(chan struct{})
	releaseGet := make(chan struct{})
	var once sync.Once
	cli.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(context.Context, string, ...clientv3.OpOption) (*clientv3.GetResponse, error) {
			once.Do(func() {
				close(getEntered)
			})
			<-releaseGet

			return &clientv3.GetResponse{Header: &etcdserverpb.ResponseHeader{}}, nil
		}).AnyTimes()

	var watches int32
	cli.EXPECT().Watch(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(context.Context, string, ...clientv3.OpOption) clientv3.WatchChan {
			atomic.AddInt32(&watches, 1)

			return make(chan clientv3.WatchResponse)
		}).AnyTimes()

	key := watchKey{key: "any"}
	c := newCluster([]string{"reload-replace-test"})
	c.watchers[key] = &watchValue{
		listeners: []UpdateListener{new(mockListener)},
		values:    make(map[string]string),
	}

	gate := make(chan struct{})
	c.watchGroup.Run(func() {
		<-gate
	})

	oldDone := c.currentDone()
	reloaded := make(chan struct{})
	go func() {
		c.reload(cli)
		close(reloaded)
	}()

	select {
	case <-oldDone:
	case <-time.After(2 * time.Second):
		t.Fatal("reload did not replace the generation")
	}

	close(gate)
	select {
	case <-getEntered:
	case <-time.After(2 * time.Second):
		t.Fatal("reload did not reach load")
	}

	// Replace the watcher while reload's load is in flight.
	c.lock.Lock()
	delete(c.watchers, key)
	c.lock.Unlock()
	c.addListener(key, new(mockListener))
	assert.True(t, c.startWatch(c.currentDone(), cli, key, 0))

	assert.Eventually(t, func() bool {
		return atomic.LoadInt32(&watches) == 1
	}, time.Second, time.Millisecond)

	close(releaseGet)
	select {
	case <-reloaded:
	case <-time.After(3 * time.Second):
		t.Fatal("reload did not finish")
	}

	// The replaced watcher owns the key now; reload must not open a second
	// stream with the stale one.
	time.Sleep(100 * time.Millisecond)
	assert.EqualValues(t, 1, atomic.LoadInt32(&watches))

	c.lock.Lock()
	group := c.watchGroup
	close(c.done)
	c.lock.Unlock()
	group.Wait()
}

func TestCluster_ReloadClaimsKeysForItsGeneration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cli := NewMockEtcdClient(ctrl)
	cli.EXPECT().Ctx().Return(context.Background()).AnyTimes()
	cli.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(
		&clientv3.GetResponse{Header: &etcdserverpb.ResponseHeader{}}, nil).AnyTimes()

	var watches int32
	cli.EXPECT().Watch(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(context.Context, string, ...clientv3.OpOption) clientv3.WatchChan {
			atomic.AddInt32(&watches, 1)

			return make(chan clientv3.WatchResponse)
		}).AnyTimes()

	key := watchKey{key: "any"}
	c := newCluster([]string{"reload-claim-test"})
	c.watchers[key] = &watchValue{
		listeners: []UpdateListener{new(mockListener)},
		values:    make(map[string]string),
	}

	gate := make(chan struct{})
	c.watchGroup.Run(func() {
		<-gate
	})

	oldDone := c.currentDone()
	reloaded := make(chan struct{})
	go func() {
		c.reload(cli)
		close(reloaded)
	}()

	select {
	case <-oldDone:
	case <-time.After(2 * time.Second):
		t.Fatal("reload did not replace the generation")
	}

	// A monitor for the same key reaches startWatch before reload's goroutine
	// runs. reload has already claimed the key for this generation, so
	// startWatch must not start a second watcher.
	assert.True(t, c.startWatch(c.currentDone(), cli, key, 0))

	close(gate)
	select {
	case <-reloaded:
	case <-time.After(3 * time.Second):
		t.Fatal("reload did not finish")
	}

	assert.Eventually(t, func() bool {
		return atomic.LoadInt32(&watches) == 1
	}, time.Second, time.Millisecond)
	time.Sleep(100 * time.Millisecond)
	assert.EqualValues(t, 1, atomic.LoadInt32(&watches))

	c.lock.Lock()
	group := c.watchGroup
	close(c.done)
	c.lock.Unlock()
	group.Wait()
}

func TestClusterWatch_CloseChan(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cli := NewMockEtcdClient(ctrl)
	restore := setMockClient(cli)
	defer restore()
	ch := make(chan clientv3.WatchResponse)
	cli.EXPECT().Watch(gomock.Any(), "any/", gomock.Any()).Return(ch).AnyTimes()
	cli.EXPECT().Ctx().Return(context.Background()).AnyTimes()
	c := &cluster{
		watchers: make(map[watchKey]*watchValue),
	}
	c.done = make(chan lang.PlaceholderType)
	key := watchKey{
		key: "any",
	}
	c.watchers[key] = newWatchValue()
	go func() {
		close(ch)
		close(c.done)
	}()
	c.watch(c.done, cli, key, c.watchers[key], 0)
}

func TestValueOnlyContext(t *testing.T) {
	ctx := contextx.ValueOnlyFrom(context.Background())
	ctx.Done()
	assert.Nil(t, ctx.Err())
}

func TestDialClient(t *testing.T) {
	svr, err := mockserver.StartMockServers(1)
	assert.NoError(t, err)
	svr.StartAt(0)

	certFile := createTempFile(t, []byte(certContent))
	defer os.Remove(certFile)
	keyFile := createTempFile(t, []byte(keyContent))
	defer os.Remove(keyFile)
	caFile := createTempFile(t, []byte(caContent))
	defer os.Remove(caFile)

	endpoints := []string{svr.Servers[0].Address}
	AddAccount(endpoints, "foo", "bar")
	assert.NoError(t, AddTLS(endpoints, certFile, keyFile, caFile, false))

	old := DialTimeout
	DialTimeout = time.Millisecond
	defer func() {
		DialTimeout = old
	}()
	_, err = DialClient(endpoints)
	assert.Error(t, err)
}

func TestRegistry_Monitor(t *testing.T) {
	svr, err := mockserver.StartMockServers(1)
	assert.NoError(t, err)
	svr.StartAt(0)

	endpoints := []string{svr.Servers[0].Address}
	GetRegistry().lock.Lock()
	GetRegistry().clusters = map[string]*cluster{
		getClusterKey(endpoints): {
			watchers: map[watchKey]*watchValue{
				{
					key:        "foo",
					exactMatch: true,
				}: {
					values: map[string]string{
						"bar": "baz",
					},
				},
			},
		},
	}
	GetRegistry().lock.Unlock()
	assert.Error(t, GetRegistry().Monitor(endpoints, "foo", false, new(mockListener)))
}

func TestRegistry_Unmonitor(t *testing.T) {
	svr, err := mockserver.StartMockServers(1)
	assert.NoError(t, err)
	svr.StartAt(0)

	_, cancel := context.WithCancel(context.Background())
	endpoints := []string{svr.Servers[0].Address}
	GetRegistry().lock.Lock()
	GetRegistry().clusters = map[string]*cluster{
		getClusterKey(endpoints): {
			watchers: map[watchKey]*watchValue{
				{
					key:        "foo",
					exactMatch: true,
				}: {
					values: map[string]string{
						"bar": "baz",
					},
					cancel: cancel,
				},
			},
		},
	}
	GetRegistry().lock.Unlock()
	l := new(mockListener)
	assert.NoError(t, GetRegistry().Monitor(endpoints, "foo", true, l))
	watchVals := GetRegistry().clusters[getClusterKey(endpoints)].watchers[watchKey{
		key:        "foo",
		exactMatch: true,
	}]
	assert.Equal(t, 1, len(watchVals.listeners))
	GetRegistry().Unmonitor(endpoints, "foo", true, l)
	watchVals = GetRegistry().clusters[getClusterKey(endpoints)].watchers[watchKey{
		key:        "foo",
		exactMatch: true,
	}]
	assert.Nil(t, watchVals)
}

// TestCluster_ConcurrentMonitor tests the race condition fix in setupWatch
// This test specifically covers the scenario from issue #5394 where:
// - addListener() writes to the watchers map (with lock)
// - setupWatch() reads from the watchers map (now with lock after fix)
// Running with -race flag will detect any race conditions
func TestCluster_ConcurrentMonitor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cli := NewMockEtcdClient(ctrl)
	cli.EXPECT().Ctx().Return(context.Background()).AnyTimes()
	cli.EXPECT().Watch(gomock.Any(), gomock.Any(), gomock.Any()).Return(make(chan clientv3.WatchResponse)).AnyTimes()

	c := &cluster{
		endpoints:  []string{"localhost:2379"},
		key:        "test-cluster",
		watchers:   make(map[watchKey]*watchValue),
		watchGroup: threading.NewRoutineGroup(),
		done:       make(chan lang.PlaceholderType),
		lock:       sync.RWMutex{},
	}

	// Spawn multiple concurrent operations that simulate the race condition:
	// - Some goroutines call addListener (write to map)
	// - Some goroutines call setupWatch (read from map)
	var wg sync.WaitGroup
	numGoroutines := 20
	wg.Add(numGoroutines)

	keys := []watchKey{
		{key: "key-0", exactMatch: false},
		{key: "key-1", exactMatch: false},
		{key: "key-2", exactMatch: false},
	}

	for i := 0; i < numGoroutines; i++ {
		idx := i
		go func() {
			defer wg.Done()
			key := keys[idx%len(keys)]

			if idx%2 == 0 {
				// Half the goroutines add listeners (write operation)
				c.addListener(key, &mockListener{})
			} else {
				// Half the goroutines setup watches (read operation)
				_, _, _, _ = c.setupWatch(c.currentDone(), cli, key, newWatchValue(), 0)
			}
		}()
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Verify that watchers were correctly added
	c.lock.RLock()
	assert.True(t, len(c.watchers) > 0, "watchers should be added")
	for _, watcher := range c.watchers {
		assert.NotNil(t, watcher, "watcher should not be nil")
	}
	c.lock.RUnlock()

	// Clean up
	close(c.done)
}

func TestCluster_handleWatchEvents_DuplicatePut(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	listener := NewMockUpdateListener(ctrl)
	// OnAdd must be called exactly once despite two PUT events with the same key+value.
	listener.EXPECT().OnAdd(KV{Key: "hello", Val: "world"}).Times(1)

	c := newCluster([]string{"any"})
	key := watchKey{key: "any"}
	c.watchers[key] = &watchValue{
		listeners: []UpdateListener{listener},
		values:    make(map[string]string),
	}
	events := []*clientv3.Event{
		{Type: clientv3.EventTypePut, Kv: &mvccpb.KeyValue{Key: []byte("hello"), Value: []byte("world")}},
		{Type: clientv3.EventTypePut, Kv: &mvccpb.KeyValue{Key: []byte("hello"), Value: []byte("world")}},
	}
	c.handleWatchEvents(context.Background(), key, events)
}

func TestCluster_handleWatchEvents_ValueChange(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	listener := NewMockUpdateListener(ctrl)
	gomock.InOrder(
		listener.EXPECT().OnAdd(KV{Key: "hello", Val: "world1"}),
		listener.EXPECT().OnDelete(KV{Key: "hello", Val: "world1"}),
		listener.EXPECT().OnAdd(KV{Key: "hello", Val: "world2"}),
	)

	c := newCluster([]string{"any"})
	key := watchKey{key: "any"}
	c.watchers[key] = &watchValue{
		listeners: []UpdateListener{listener},
		values:    make(map[string]string),
	}
	c.handleWatchEvents(context.Background(), key, []*clientv3.Event{
		{Type: clientv3.EventTypePut, Kv: &mvccpb.KeyValue{Key: []byte("hello"), Value: []byte("world1")}},
	})
	c.handleWatchEvents(context.Background(), key, []*clientv3.Event{
		{Type: clientv3.EventTypePut, Kv: &mvccpb.KeyValue{Key: []byte("hello"), Value: []byte("world2")}},
	})
}

type mockListener struct {
}

func (m *mockListener) OnAdd(_ KV) {
}

func (m *mockListener) OnDelete(_ KV) {
}
