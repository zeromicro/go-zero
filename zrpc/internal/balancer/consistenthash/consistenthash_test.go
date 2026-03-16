package consistenthash

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/hash"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
)

type fakeSubConn struct{ id int }

func (f *fakeSubConn) Connect()                             {}
func (f *fakeSubConn) UpdateAddresses(_ []resolver.Address) {}
func (f *fakeSubConn) Shutdown()                            {}
func (f *fakeSubConn) GetOrBuildProducer(b balancer.ProducerBuilder) (balancer.Producer, func()) {
	return nil, func() {}
}

func TestPickerBuilder_EmptyReadySCs(t *testing.T) {
	b := &pickerBuilder{}
	p := b.Build(base.PickerBuildInfo{ReadySCs: map[balancer.SubConn]base.SubConnInfo{}})

	_, err := p.Pick(balancer.PickInfo{})
	assert.Equal(t, balancer.ErrNoSubConnAvailable, err)
}

func TestPickerBuilder_BuildAndRing(t *testing.T) {
	subConn1 := &fakeSubConn{id: 1}
	subConn2 := &fakeSubConn{id: 2}
	addr1 := "127.0.0.1:8080"
	addr2 := "127.0.0.1:8081"

	b := &pickerBuilder{}
	info := base.PickerBuildInfo{
		ReadySCs: map[balancer.SubConn]base.SubConnInfo{
			subConn1: {Address: resolver.Address{Addr: addr1}},
			subConn2: {Address: resolver.Address{Addr: addr2}},
		},
	}

	p := b.Build(info).(*picker)
	assert.NotNil(t, p.hashRing)
	assert.Len(t, p.conns, 2)
}

func TestPicker_HashConsistency(t *testing.T) {
	subConn1 := &fakeSubConn{id: 1}
	subConn2 := &fakeSubConn{id: 2}

	pb := &pickerBuilder{}
	info := base.PickerBuildInfo{
		ReadySCs: map[balancer.SubConn]base.SubConnInfo{
			subConn1: {Address: resolver.Address{Addr: "127.0.0.1:8080"}},
			subConn2: {Address: resolver.Address{Addr: "127.0.0.1:8081"}},
		},
	}
	p := pb.Build(info).(*picker)
	ctx := SetHashKey(context.Background(), "user_123")
	res1, err := p.Pick(balancer.PickInfo{Ctx: ctx})
	assert.NoError(t, err)
	assert.NotNil(t, res1.SubConn)

	// Multiple requests with the same key remain consistent
	for i := 0; i < 5; i++ {
		resN, err := p.Pick(balancer.PickInfo{Ctx: ctx})
		assert.NoError(t, err)
		assert.Equal(t, res1.SubConn, resN.SubConn)
	}
}

func TestPicker_MissingKey(t *testing.T) {
	subConn := &fakeSubConn{id: 1}

	pb := &pickerBuilder{}
	info := base.PickerBuildInfo{
		ReadySCs: map[balancer.SubConn]base.SubConnInfo{
			subConn: {Address: resolver.Address{Addr: "127.0.0.1:8080"}},
		},
	}
	p := pb.Build(info).(*picker)

	// No hash key in context
	_, err := p.Pick(balancer.PickInfo{Ctx: context.Background()})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "[consistent_hash] missing hash key in context")
}

func TestPicker_NoMatchingConn(t *testing.T) {
	emptyRing := newCustomRingForTest()
	p := &picker{
		hashRing: emptyRing,
		conns:    map[string]balancer.SubConn{},
	}

	_, err := p.Pick(balancer.PickInfo{Ctx: SetHashKey(context.Background(), "someone")})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "[consistent_hash] no matching conn for hashKey: someone")
}

func TestPicker_InvalidAddrType(t *testing.T) {
	ring := newCustomRingForTest()
	ring.Add(12345)

	subConn := &fakeSubConn{id: 1}
	p := &picker{
		hashRing: ring,
		conns: map[string]balancer.SubConn{
			"12345": subConn,
		},
	}

	_, err := p.Pick(balancer.PickInfo{Ctx: SetHashKey(context.Background(), "anykey")})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "[consistent_hash] invalid addr type in consistent hash")
}

func TestPicker_NoSubConnForAddr(t *testing.T) {
	ring := newCustomRingForTest()
	ring.Add("ghost:9999")

	exist := &fakeSubConn{id: 1}
	p := &picker{
		hashRing: ring,
		conns: map[string]balancer.SubConn{
			"real:8080": exist,
		},
	}

	_, err := p.Pick(balancer.PickInfo{Ctx: SetHashKey(context.Background(), "anykey")})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "[consistent_hash] no subConn for addr: ghost:9999")
}

func TestSetAndGetHashKey(t *testing.T) {
	ctx := context.Background()
	key := "abc123"

	ctx = SetHashKey(ctx, key)
	got := GetHashKey(ctx)
	assert.Equal(t, key, got)

	assert.Empty(t, GetHashKey(context.Background()))
}

func BenchmarkPicker_HashConsistency(b *testing.B) {
	subConn1 := &fakeSubConn{id: 1}
	subConn2 := &fakeSubConn{id: 2}

	pb := &pickerBuilder{}
	info := base.PickerBuildInfo{
		ReadySCs: map[balancer.SubConn]base.SubConnInfo{
			subConn1: {Address: resolver.Address{Addr: "127.0.0.1:8080"}},
			subConn2: {Address: resolver.Address{Addr: "127.0.0.1:8081"}},
		},
	}
	p := pb.Build(info).(*picker)

	ctx := SetHashKey(context.Background(), "hot_user_123")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, err := p.Pick(balancer.PickInfo{Ctx: ctx})
		if err != nil || res.SubConn == nil {
			b.Fatalf("unexpected result: res=%v err=%v", res.SubConn, err)
		}
	}
}

func newCustomRingForTest() *hash.ConsistentHash {
	return hash.NewCustomConsistentHash(defaultReplicaCount, hash.Hash)
}
