package consistenthash

import (
	"context"

	"github.com/zeromicro/go-zero/core/hash"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	Name = "consistent_hash"

	defaultReplicaCount = 100
)

var emptyPickResult balancer.PickResult

func init() {
	balancer.Register(newBuilder())
}

type (
	// hashKey is the key type for consistent hash in context.
	hashKey struct{}
	// pickerBuilder is a builder for picker.
	pickerBuilder struct{}
	// picker is a picker that uses consistent hash to pick a sub connection.
	picker struct {
		hashRing *hash.ConsistentHash
		conns    map[string]balancer.SubConn
	}
)

func (b *pickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	readySCs := info.ReadySCs
	if len(readySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}

	conns := make(map[string]balancer.SubConn, len(readySCs))
	hashRing := hash.NewCustomConsistentHash(defaultReplicaCount, hash.Hash)
	for conn, connInfo := range readySCs {
		addr := connInfo.Address.Addr
		conns[addr] = conn
		hashRing.Add(addr)
	}

	return &picker{
		hashRing: hashRing,
		conns:    conns,
	}
}

func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder(Name, &pickerBuilder{}, base.Config{HealthCheck: true})
}

func (p *picker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	hashKey := GetHashKey(info.Ctx)
	if len(hashKey) == 0 {
		return emptyPickResult, status.Error(codes.InvalidArgument,
			"[consistent_hash] missing hash key in context")
	}

	if addrAny, ok := p.hashRing.Get(hashKey); ok {
		addr, ok := addrAny.(string)
		if !ok {
			return emptyPickResult, status.Error(codes.Internal,
				"[consistent_hash] invalid addr type in consistent hash")
		}

		subConn, ok := p.conns[addr]
		if !ok {
			return emptyPickResult, status.Errorf(codes.Internal,
				"[consistent_hash] no subConn for addr: %s", addr)
		}

		return balancer.PickResult{SubConn: subConn}, nil
	}

	return emptyPickResult, status.Errorf(codes.Unavailable,
		"[consistent_hash] no matching conn for hashKey: %s", hashKey)
}

// SetHashKey sets the hash key into context.
func SetHashKey(ctx context.Context, key string) context.Context {
	return context.WithValue(ctx, hashKey{}, key)
}

// GetHashKey gets the hash key from context.
func GetHashKey(ctx context.Context) string {
	v, _ := ctx.Value(hashKey{}).(string)
	return v
}
