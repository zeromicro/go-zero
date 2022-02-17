package discov

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/discov/internal"
	"github.com/zeromicro/go-zero/core/lang"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stringx"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func init() {
	logx.Disable()
}

func TestPublisher_register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	const id = 1
	cli := internal.NewMockEtcdClient(ctrl)
	restore := setMockClient(cli)
	defer restore()
	cli.EXPECT().Ctx().AnyTimes()
	cli.EXPECT().Grant(gomock.Any(), timeToLive).Return(&clientv3.LeaseGrantResponse{
		ID: id,
	}, nil)
	cli.EXPECT().Put(gomock.Any(), makeEtcdKey("thekey", id), "thevalue", gomock.Any())
	pub := NewPublisher(nil, "thekey", "thevalue",
		WithPubEtcdAccount(stringx.Rand(), "bar"))
	_, err := pub.register(cli)
	assert.Nil(t, err)
}

func TestPublisher_registerWithId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	const id = 2
	cli := internal.NewMockEtcdClient(ctrl)
	restore := setMockClient(cli)
	defer restore()
	cli.EXPECT().Ctx().AnyTimes()
	cli.EXPECT().Grant(gomock.Any(), timeToLive).Return(&clientv3.LeaseGrantResponse{
		ID: 1,
	}, nil)
	cli.EXPECT().Put(gomock.Any(), makeEtcdKey("thekey", id), "thevalue", gomock.Any())
	pub := NewPublisher(nil, "thekey", "thevalue", WithId(id))
	_, err := pub.register(cli)
	assert.Nil(t, err)
}

func TestPublisher_registerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cli := internal.NewMockEtcdClient(ctrl)
	restore := setMockClient(cli)
	defer restore()
	cli.EXPECT().Ctx().AnyTimes()
	cli.EXPECT().Grant(gomock.Any(), timeToLive).Return(nil, errors.New("error"))
	pub := NewPublisher(nil, "thekey", "thevalue")
	val, err := pub.register(cli)
	assert.NotNil(t, err)
	assert.Equal(t, clientv3.NoLease, val)
}

func TestPublisher_revoke(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	const id clientv3.LeaseID = 1
	cli := internal.NewMockEtcdClient(ctrl)
	restore := setMockClient(cli)
	defer restore()
	cli.EXPECT().Ctx().AnyTimes()
	cli.EXPECT().Revoke(gomock.Any(), id)
	pub := NewPublisher(nil, "thekey", "thevalue")
	pub.lease = id
	pub.revoke(cli)
}

func TestPublisher_revokeError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	const id clientv3.LeaseID = 1
	cli := internal.NewMockEtcdClient(ctrl)
	restore := setMockClient(cli)
	defer restore()
	cli.EXPECT().Ctx().AnyTimes()
	cli.EXPECT().Revoke(gomock.Any(), id).Return(nil, errors.New("error"))
	pub := NewPublisher(nil, "thekey", "thevalue")
	pub.lease = id
	pub.revoke(cli)
}

func TestPublisher_keepAliveAsyncError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	const id clientv3.LeaseID = 1
	cli := internal.NewMockEtcdClient(ctrl)
	restore := setMockClient(cli)
	defer restore()
	cli.EXPECT().Ctx().AnyTimes()
	cli.EXPECT().KeepAlive(gomock.Any(), id).Return(nil, errors.New("error"))
	pub := NewPublisher(nil, "thekey", "thevalue")
	pub.lease = id
	assert.NotNil(t, pub.keepAliveAsync(cli))
}

func TestPublisher_keepAliveAsyncQuit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	const id clientv3.LeaseID = 1
	cli := internal.NewMockEtcdClient(ctrl)
	cli.EXPECT().ActiveConnection()
	cli.EXPECT().Close()
	defer cli.Close()
	cli.ActiveConnection()
	restore := setMockClient(cli)
	defer restore()
	cli.EXPECT().Ctx().AnyTimes()
	cli.EXPECT().KeepAlive(gomock.Any(), id)
	var wg sync.WaitGroup
	wg.Add(1)
	cli.EXPECT().Revoke(gomock.Any(), id).Do(func(_, _ interface{}) {
		wg.Done()
	})
	pub := NewPublisher(nil, "thekey", "thevalue")
	pub.lease = id
	pub.Stop()
	assert.Nil(t, pub.keepAliveAsync(cli))
	wg.Wait()
}

func TestPublisher_keepAliveAsyncPause(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	const id clientv3.LeaseID = 1
	cli := internal.NewMockEtcdClient(ctrl)
	restore := setMockClient(cli)
	defer restore()
	cli.EXPECT().Ctx().AnyTimes()
	cli.EXPECT().KeepAlive(gomock.Any(), id)
	pub := NewPublisher(nil, "thekey", "thevalue")
	var wg sync.WaitGroup
	wg.Add(1)
	cli.EXPECT().Revoke(gomock.Any(), id).Do(func(_, _ interface{}) {
		pub.Stop()
		wg.Done()
	})
	pub.lease = id
	assert.Nil(t, pub.keepAliveAsync(cli))
	pub.Pause()
	wg.Wait()
}

func TestPublisher_Resume(t *testing.T) {
	publisher := new(Publisher)
	publisher.resumeChan = make(chan lang.PlaceholderType)
	go func() {
		publisher.Resume()
	}()
	go func() {
		time.Sleep(time.Minute)
		t.Fail()
	}()
	<-publisher.resumeChan
}
