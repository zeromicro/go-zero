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
	c1, _ := GetRegistry().getCluster([]string{"first"})
	c2, _ := GetRegistry().getCluster([]string{"second"})
	c3, _ := GetRegistry().getCluster([]string{"first"})
	assert.Equal(t, c1, c3)
	assert.NotEqual(t, c1, c2)
}

func TestGetClusterKey(t *testing.T) {
	assert.Equal(t, getClusterKey([]string{"localhost:1234", "remotehost:5678"}),
		getClusterKey([]string{"remotehost:5678", "localhost:1234"}))
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
	c.listeners["any"] = []UpdateListener{l}
	c.handleChanges("any", []KV{
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
	}, c.values["any"])
	c.handleChanges("any", []KV{
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
	}, c.values["any"])
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
		values: make(map[string]map[string]string),
	}
	c.load(cli, "any")
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
				values:    make(map[string]map[string]string),
				listeners: make(map[string][]UpdateListener),
				watchCtx:  make(map[string]context.CancelFunc),
				watchFlag: make(map[string]bool),
			}
			listener := NewMockUpdateListener(ctrl)
			c.listeners["any"] = []UpdateListener{listener}
			listener.EXPECT().OnAdd(gomock.Any()).Do(func(kv KV) {
				assert.Equal(t, "hello", kv.Key)
				assert.Equal(t, "world", kv.Val)
				wg.Done()
			}).MaxTimes(1)
			listener.EXPECT().OnDelete(gomock.Any()).Do(func(_ any) {
				wg.Done()
			}).MaxTimes(1)
			go c.watch(cli, "any", 0)
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
			c := newCluster([]string{})
			c.done = make(chan lang.PlaceholderType)
			go func() {
				ch <- resp
				close(c.done)
			}()
			c.watch(cli, "any", 0)
		})
	}
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
	c := newCluster([]string{})
	c.done = make(chan lang.PlaceholderType)
	go func() {
		close(ch)
		close(c.done)
	}()
	c.watch(cli, "any", 0)
}

func TestClusterWatch_CtxCancel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cli := NewMockEtcdClient(ctrl)
	restore := setMockClient(cli)
	defer restore()
	ch := make(chan clientv3.WatchResponse)
	cli.EXPECT().Watch(gomock.Any(), "any/", gomock.Any()).Return(ch).AnyTimes()
	ctx, cancelFunc := context.WithCancel(context.Background())
	cli.EXPECT().Ctx().Return(ctx).AnyTimes()
	c := newCluster([]string{})
	c.done = make(chan lang.PlaceholderType)
	go func() {
		cancelFunc()
		close(ch)
	}()
	c.watch(cli, "any", 0)
}

func TestCluster_ClearWatch(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	c := &cluster{
		watchCtx:  map[string]context.CancelFunc{"foo": cancel},
		watchFlag: map[string]bool{"foo": true},
	}
	c.clearWatch()
	assert.Equal(t, ctx.Err(), context.Canceled)
	assert.Equal(t, 0, len(c.watchCtx))
	assert.Equal(t, 0, len(c.watchFlag))
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
			listeners: map[string][]UpdateListener{},
			values: map[string]map[string]string{
				"foo": {
					"bar": "baz",
				},
			},
			watchCtx:  map[string]context.CancelFunc{},
			watchFlag: map[string]bool{},
		},
	}
	GetRegistry().lock.Unlock()
	assert.Error(t, GetRegistry().Monitor(endpoints, "foo", new(mockListener), false))
}

func TestRegistry_Unmonitor(t *testing.T) {
	l := new(mockListener)
	GetRegistry().lock.Lock()
	GetRegistry().clusters = map[string]*cluster{
		getClusterKey(endpoints): {
			listeners: map[string][]UpdateListener{"foo": {l}},
			values: map[string]map[string]string{
				"foo": {
					"bar": "baz",
				},
			},
		},
	}
	GetRegistry().lock.Unlock()
	l := new(mockListener)
	assert.Error(t, GetRegistry().Monitor(endpoints, "foo", l, false))
	assert.Equal(t, 1, len(GetRegistry().clusters[getClusterKey(endpoints)].listeners["foo"]))
	GetRegistry().Unmonitor(endpoints, "foo", l)
	assert.Equal(t, 0, len(GetRegistry().clusters[getClusterKey(endpoints)].listeners["foo"]))
}

type mockListener struct {
}

func (m *mockListener) OnAdd(_ KV) {
}

func (m *mockListener) OnDelete(_ KV) {
}
