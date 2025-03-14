package internal

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
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
	c.load(cli, watchKey{
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
			go c.watch(cli, key, 0)
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
			go func() {
				ch <- resp
				close(c.done)
			}()
			key := watchKey{
				key: "any",
			}
			c.watch(cli, key, 0)
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
	go func() {
		close(ch)
		close(c.done)
	}()
	c.watch(cli, watchKey{
		key: "any",
	}, 0)
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
				watchKey{
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
				watchKey{
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

type mockListener struct {
}

func (m *mockListener) OnAdd(_ KV) {
}

func (m *mockListener) OnDelete(_ KV) {
}
