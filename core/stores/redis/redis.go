package redis

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	red "github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mapping"
	"github.com/zeromicro/go-zero/core/syncx"
)

const (
	// ClusterType means redis cluster.
	ClusterType = "cluster"
	// NodeType means redis node.
	NodeType = "node"
	// Nil is an alias of redis.Nil.
	Nil = red.Nil

	blockingQueryTimeout = 5 * time.Second
	readWriteTimeout     = 2 * time.Second
	defaultSlowThreshold = time.Millisecond * 100
	defaultPingTimeout   = time.Second
)

var (
	// ErrNilNode is an error that indicates a nil redis node.
	ErrNilNode    = errors.New("nil redis node")
	slowThreshold = syncx.ForAtomicDuration(defaultSlowThreshold)
)

type (
	// Option defines the method to customize a Redis.
	Option func(r *Redis)

	// A Pair is a key/pair set used in redis zset.
	Pair struct {
		Key   string
		Score int64
	}

	// A FloatPair is a key/pair for float set used in redis zet.
	FloatPair struct {
		Key   string
		Score float64
	}

	// Redis defines a redis node/cluster. It is thread-safe.
	Redis struct {
		Addr  string
		Type  string
		User  string
		Pass  string
		tls   bool
		brk   breaker.Breaker
		hooks []red.Hook
	}

	// RedisNode interface represents a redis node.
	RedisNode interface {
		red.Cmdable
		red.BitMapCmdable
	}

	// GeoLocation is used with GeoAdd to add geospatial location.
	GeoLocation = red.GeoLocation
	// GeoRadiusQuery is used with GeoRadius to query geospatial index.
	GeoRadiusQuery = red.GeoRadiusQuery
	// GeoPos is used to represent a geo position.
	GeoPos = red.GeoPos

	// Pipeliner is an alias of redis.Pipeliner.
	Pipeliner = red.Pipeliner

	// Z represents sorted set member.
	Z = red.Z
	// ZStore is an alias of redis.ZStore.
	ZStore = red.ZStore

	// IntCmd is an alias of redis.IntCmd.
	IntCmd = red.IntCmd
	// FloatCmd is an alias of redis.FloatCmd.
	FloatCmd = red.FloatCmd
	// StringCmd is an alias of redis.StringCmd.
	StringCmd = red.StringCmd
	// Script is an alias of redis.Script.
	Script = red.Script

	// Hook is an alias of redis.Hook.
	Hook = red.Hook
	// DialHook is an alias of redis.DialHook.
	DialHook = red.DialHook
	// ProcessHook is an alias of redis.ProcessHook.
	ProcessHook = red.ProcessHook
	// ProcessPipelineHook is an alias of redis.ProcessPipelineHook.
	ProcessPipelineHook = red.ProcessPipelineHook

	// Cmder is an alias of redis.Cmder.
	Cmder = red.Cmder
)

// MustNewRedis returns a Redis with given options.
func MustNewRedis(conf RedisConf, opts ...Option) *Redis {
	rds, err := NewRedis(conf, opts...)
	logx.Must(err)
	return rds
}

// New returns a Redis with given options.
// Deprecated: use MustNewRedis or NewRedis instead.
func New(addr string, opts ...Option) *Redis {
	return newRedis(addr, opts...)
}

// NewRedis returns a Redis with given options.
func NewRedis(conf RedisConf, opts ...Option) (*Redis, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	if conf.Type == ClusterType {
		opts = append([]Option{Cluster()}, opts...)
	}
	if len(conf.User) > 0 {
		opts = append([]Option{WithUser(conf.User)}, opts...)
	}
	if len(conf.Pass) > 0 {
		opts = append([]Option{WithPass(conf.Pass)}, opts...)
	}
	if conf.Tls {
		opts = append([]Option{WithTLS()}, opts...)
	}

	rds := newRedis(conf.Host, opts...)
	if !conf.NonBlock {
		if err := rds.checkConnection(conf.PingTimeout); err != nil {
			return nil, errorx.Wrap(err, fmt.Sprintf("redis connect error, addr: %s", conf.Host))
		}
	}

	return rds, nil
}

func newRedis(addr string, opts ...Option) *Redis {
	r := &Redis{
		Addr: addr,
		Type: NodeType,
		brk:  breaker.NewBreaker(),
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

// NewScript returns a new Script instance.
func NewScript(script string) *Script {
	return red.NewScript(script)
}

// BitCount is redis bitcount command implementation.
func (s *Redis) BitCount(key string, start, end int64) (int64, error) {
	return s.BitCountCtx(context.Background(), key, start, end)
}

// BitCountCtx is redis bitcount command implementation.
func (s *Redis) BitCountCtx(ctx context.Context, key string, start, end int64) (int64, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	return conn.BitCount(ctx, key, &red.BitCount{
		Start: start,
		End:   end,
	}).Result()
}

// BitOpAnd is redis bit operation (and) command implementation.
func (s *Redis) BitOpAnd(destKey string, keys ...string) (int64, error) {
	return s.BitOpAndCtx(context.Background(), destKey, keys...)
}

// BitOpAndCtx is redis bit operation (and) command implementation.
func (s *Redis) BitOpAndCtx(ctx context.Context, destKey string, keys ...string) (int64, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	return conn.BitOpAnd(ctx, destKey, keys...).Result()
}

// BitOpNot is redis bit operation (not) command implementation.
func (s *Redis) BitOpNot(destKey, key string) (int64, error) {
	return s.BitOpNotCtx(context.Background(), destKey, key)
}

// BitOpNotCtx is redis bit operation (not) command implementation.
func (s *Redis) BitOpNotCtx(ctx context.Context, destKey, key string) (int64, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	return conn.BitOpNot(ctx, destKey, key).Result()
}

// BitOpOr is redis bit operation (or) command implementation.
func (s *Redis) BitOpOr(destKey string, keys ...string) (int64, error) {
	return s.BitOpOrCtx(context.Background(), destKey, keys...)
}

// BitOpOrCtx is redis bit operation (or) command implementation.
func (s *Redis) BitOpOrCtx(ctx context.Context, destKey string, keys ...string) (int64, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	return conn.BitOpOr(ctx, destKey, keys...).Result()
}

// BitOpXor is redis bit operation (xor) command implementation.
func (s *Redis) BitOpXor(destKey string, keys ...string) (int64, error) {
	return s.BitOpXorCtx(context.Background(), destKey, keys...)
}

// BitOpXorCtx is redis bit operation (xor) command implementation.
func (s *Redis) BitOpXorCtx(ctx context.Context, destKey string, keys ...string) (int64, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	return conn.BitOpXor(ctx, destKey, keys...).Result()
}

// BitPos is redis bitpos command implementation.
func (s *Redis) BitPos(key string, bit, start, end int64) (int64, error) {
	return s.BitPosCtx(context.Background(), key, bit, start, end)
}

// BitPosCtx is redis bitpos command implementation.
func (s *Redis) BitPosCtx(ctx context.Context, key string, bit, start, end int64) (int64, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	return conn.BitPos(ctx, key, bit, start, end).Result()
}

// Blpop uses passed in redis connection to execute blocking queries.
// Doesn't benefit from pooling redis connections of blocking queries
func (s *Redis) Blpop(node RedisNode, key string) (string, error) {
	return s.BlpopCtx(context.Background(), node, key)
}

// BlpopCtx uses passed in redis connection to execute blocking queries.
// Doesn't benefit from pooling redis connections of blocking queries
func (s *Redis) BlpopCtx(ctx context.Context, node RedisNode, key string) (string, error) {
	return s.BlpopWithTimeoutCtx(ctx, node, blockingQueryTimeout, key)
}

// BlpopEx uses passed in redis connection to execute blpop command.
// The difference against Blpop is that this method returns a bool to indicate success.
func (s *Redis) BlpopEx(node RedisNode, key string) (string, bool, error) {
	return s.BlpopExCtx(context.Background(), node, key)
}

// BlpopExCtx uses passed in redis connection to execute blpop command.
// The difference against Blpop is that this method returns a bool to indicate success.
func (s *Redis) BlpopExCtx(ctx context.Context, node RedisNode, key string) (string, bool, error) {
	if node == nil {
		return "", false, ErrNilNode
	}

	vals, err := node.BLPop(ctx, blockingQueryTimeout, key).Result()
	if err != nil {
		return "", false, err
	}

	if len(vals) < 2 {
		return "", false, fmt.Errorf("no value on key: %s", key)
	}

	return vals[1], true, nil
}

// BlpopWithTimeout uses passed in redis connection to execute blpop command.
// Control blocking query timeout
func (s *Redis) BlpopWithTimeout(node RedisNode, timeout time.Duration, key string) (string, error) {
	return s.BlpopWithTimeoutCtx(context.Background(), node, timeout, key)
}

// BlpopWithTimeoutCtx uses passed in redis connection to execute blpop command.
// Control blocking query timeout
func (s *Redis) BlpopWithTimeoutCtx(ctx context.Context, node RedisNode, timeout time.Duration,
	key string) (string, error) {
	if node == nil {
		return "", ErrNilNode
	}

	vals, err := node.BLPop(ctx, timeout, key).Result()
	if err != nil {
		return "", err
	}

	if len(vals) < 2 {
		return "", fmt.Errorf("no value on key: %s", key)
	}

	return vals[1], nil
}

// Decr is the implementation of redis decr command.
func (s *Redis) Decr(key string) (int64, error) {
	return s.DecrCtx(context.Background(), key)
}

// DecrCtx is the implementation of redis decr command.
func (s *Redis) DecrCtx(ctx context.Context, key string) (int64, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	return conn.Decr(ctx, key).Result()
}

// Decrby is the implementation of redis decrby command.
func (s *Redis) Decrby(key string, decrement int64) (int64, error) {
	return s.DecrbyCtx(context.Background(), key, decrement)
}

// DecrbyCtx is the implementation of redis decrby command.
func (s *Redis) DecrbyCtx(ctx context.Context, key string, decrement int64) (int64, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	return conn.DecrBy(ctx, key, decrement).Result()
}

// Del deletes keys.
func (s *Redis) Del(keys ...string) (int, error) {
	return s.DelCtx(context.Background(), keys...)
}

// DelCtx deletes keys.
func (s *Redis) DelCtx(ctx context.Context, keys ...string) (int, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	v, err := conn.Del(ctx, keys...).Result()
	if err != nil {
		return 0, err
	}

	return int(v), nil
}

// Eval is the implementation of redis eval command.
func (s *Redis) Eval(script string, keys []string, args ...any) (any, error) {
	return s.EvalCtx(context.Background(), script, keys, args...)
}

// EvalCtx is the implementation of redis eval command.
func (s *Redis) EvalCtx(ctx context.Context, script string, keys []string,
	args ...any) (any, error) {
	conn, err := getRedis(s)
	if err != nil {
		return nil, err
	}

	return conn.Eval(ctx, script, keys, args...).Result()
}

// EvalSha is the implementation of redis evalsha command.
func (s *Redis) EvalSha(sha string, keys []string, args ...any) (any, error) {
	return s.EvalShaCtx(context.Background(), sha, keys, args...)
}

// EvalShaCtx is the implementation of redis evalsha command.
func (s *Redis) EvalShaCtx(ctx context.Context, sha string, keys []string,
	args ...any) (any, error) {
	conn, err := getRedis(s)
	if err != nil {
		return nil, err
	}

	return conn.EvalSha(ctx, sha, keys, args...).Result()
}

// Exists is the implementation of redis exists command.
func (s *Redis) Exists(key string) (bool, error) {
	return s.ExistsCtx(context.Background(), key)
}

// ExistsCtx is the implementation of redis exists command.
func (s *Redis) ExistsCtx(ctx context.Context, key string) (bool, error) {
	conn, err := getRedis(s)
	if err != nil {
		return false, err
	}

	v, err := conn.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return v == 1, nil
}

// ExistsMany is the implementation of redis exists command.
// checks the existence of multiple keys in Redis using the EXISTS command.
func (s *Redis) ExistsMany(keys ...string) (int64, error) {
	return s.ExistsManyCtx(context.Background(), keys...)
}

// ExistsManyCtx is the implementation of redis exists command.
// checks the existence of multiple keys in Redis using the EXISTS command.
func (s *Redis) ExistsManyCtx(ctx context.Context, keys ...string) (int64, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	return conn.Exists(ctx, keys...).Result()
}

// Expire is the implementation of redis expire command.
func (s *Redis) Expire(key string, seconds int) error {
	return s.ExpireCtx(context.Background(), key, seconds)
}

// ExpireCtx is the implementation of redis expire command.
func (s *Redis) ExpireCtx(ctx context.Context, key string, seconds int) error {
	conn, err := getRedis(s)
	if err != nil {
		return err
	}

	return conn.Expire(ctx, key, time.Duration(seconds)*time.Second).Err()
}

// Expireat is the implementation of redis expireat command.
func (s *Redis) Expireat(key string, expireTime int64) error {
	return s.ExpireatCtx(context.Background(), key, expireTime)
}

// ExpireatCtx is the implementation of redis expireat command.
func (s *Redis) ExpireatCtx(ctx context.Context, key string, expireTime int64) error {
	conn, err := getRedis(s)
	if err != nil {
		return err
	}

	return conn.ExpireAt(ctx, key, time.Unix(expireTime, 0)).Err()
}

// GeoAdd is the implementation of redis geoadd command.
func (s *Redis) GeoAdd(key string, geoLocation ...*GeoLocation) (int64, error) {
	return s.GeoAddCtx(context.Background(), key, geoLocation...)
}

// GeoAddCtx is the implementation of redis geoadd command.
func (s *Redis) GeoAddCtx(ctx context.Context, key string, geoLocation ...*GeoLocation) (
	int64, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	return conn.GeoAdd(ctx, key, geoLocation...).Result()
}

// GeoDist is the implementation of redis geodist command.
func (s *Redis) GeoDist(key, member1, member2, unit string) (float64, error) {
	return s.GeoDistCtx(context.Background(), key, member1, member2, unit)
}

// GeoDistCtx is the implementation of redis geodist command.
func (s *Redis) GeoDistCtx(ctx context.Context, key, member1, member2, unit string) (
	float64, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	return conn.GeoDist(ctx, key, member1, member2, unit).Result()
}

// GeoHash is the implementation of redis geohash command.
func (s *Redis) GeoHash(key string, members ...string) ([]string, error) {
	return s.GeoHashCtx(context.Background(), key, members...)
}

// GeoHashCtx is the implementation of redis geohash command.
func (s *Redis) GeoHashCtx(ctx context.Context, key string, members ...string) (
	[]string, error) {
	conn, err := getRedis(s)
	if err != nil {
		return nil, err
	}

	return conn.GeoHash(ctx, key, members...).Result()
}

// GeoRadius is the implementation of redis georadius command.
func (s *Redis) GeoRadius(key string, longitude, latitude float64, query *GeoRadiusQuery) (
	[]GeoLocation, error) {
	return s.GeoRadiusCtx(context.Background(), key, longitude, latitude, query)
}

// GeoRadiusCtx is the implementation of redis georadius command.
func (s *Redis) GeoRadiusCtx(ctx context.Context, key string, longitude, latitude float64,
	query *GeoRadiusQuery) ([]GeoLocation, error) {
	conn, err := getRedis(s)
	if err != nil {
		return nil, err
	}

	return conn.GeoRadius(ctx, key, longitude, latitude, query).Result()
}

// GeoRadiusByMember is the implementation of redis georadiusbymember command.
func (s *Redis) GeoRadiusByMember(key, member string, query *GeoRadiusQuery) ([]GeoLocation, error) {
	return s.GeoRadiusByMemberCtx(context.Background(), key, member, query)
}

// GeoRadiusByMemberCtx is the implementation of redis georadiusbymember command.
func (s *Redis) GeoRadiusByMemberCtx(ctx context.Context, key, member string,
	query *GeoRadiusQuery) ([]GeoLocation, error) {
	conn, err := getRedis(s)
	if err != nil {
		return nil, err
	}

	return conn.GeoRadiusByMember(ctx, key, member, query).Result()
}

// GeoPos is the implementation of redis geopos command.
func (s *Redis) GeoPos(key string, members ...string) ([]*GeoPos, error) {
	return s.GeoPosCtx(context.Background(), key, members...)
}

// GeoPosCtx is the implementation of redis geopos command.
func (s *Redis) GeoPosCtx(ctx context.Context, key string, members ...string) (
	[]*GeoPos, error) {
	conn, err := getRedis(s)
	if err != nil {
		return nil, err
	}

	return conn.GeoPos(ctx, key, members...).Result()
}

// Get is the implementation of redis get command.
func (s *Redis) Get(key string) (string, error) {
	return s.GetCtx(context.Background(), key)
}

// GetCtx is the implementation of redis get command.
func (s *Redis) GetCtx(ctx context.Context, key string) (string, error) {
	conn, err := getRedis(s)
	if err != nil {
		return "", err
	}

	if val, err := conn.Get(ctx, key).Result(); errors.Is(err, red.Nil) {
		return "", nil
	} else if err != nil {
		return "", err
	} else {
		return val, nil
	}
}

// GetBit is the implementation of redis getbit command.
func (s *Redis) GetBit(key string, offset int64) (int, error) {
	return s.GetBitCtx(context.Background(), key, offset)
}

// GetBitCtx is the implementation of redis getbit command.
func (s *Redis) GetBitCtx(ctx context.Context, key string, offset int64) (int, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	v, err := conn.GetBit(ctx, key, offset).Result()
	if err != nil {
		return 0, err
	}

	return int(v), nil
}

// GetDel is the implementation of redis getdel command.
// Available since: redis version 6.2.0
func (s *Redis) GetDel(key string) (string, error) {
	return s.GetDelCtx(context.Background(), key)
}

// GetDelCtx is the implementation of redis getdel command.
// Available since: redis version 6.2.0
func (s *Redis) GetDelCtx(ctx context.Context, key string) (string, error) {
	conn, err := getRedis(s)
	if err != nil {
		return "", err
	}

	val, err := conn.GetDel(ctx, key).Result()
	if errors.Is(err, red.Nil) {
		return "", nil
	}

	return val, err
}

// GetSet is the implementation of redis getset command.
func (s *Redis) GetSet(key, value string) (string, error) {
	return s.GetSetCtx(context.Background(), key, value)
}

// GetSetCtx is the implementation of redis getset command.
func (s *Redis) GetSetCtx(ctx context.Context, key, value string) (string, error) {
	conn, err := getRedis(s)
	if err != nil {
		return "", err
	}

	val, err := conn.GetSet(ctx, key, value).Result()
	if errors.Is(err, red.Nil) {
		return "", nil
	}

	return val, err
}

// Hdel is the implementation of redis hdel command.
func (s *Redis) Hdel(key string, fields ...string) (bool, error) {
	return s.HdelCtx(context.Background(), key, fields...)
}

// HdelCtx is the implementation of redis hdel command.
func (s *Redis) HdelCtx(ctx context.Context, key string, fields ...string) (bool, error) {
	conn, err := getRedis(s)
	if err != nil {
		return false, err
	}

	v, err := conn.HDel(ctx, key, fields...).Result()
	if err != nil {
		return false, err
	}

	return v >= 1, nil
}

// Hexists is the implementation of redis hexists command.
func (s *Redis) Hexists(key, field string) (bool, error) {
	return s.HexistsCtx(context.Background(), key, field)
}

// HexistsCtx is the implementation of redis hexists command.
func (s *Redis) HexistsCtx(ctx context.Context, key, field string) (bool, error) {
	conn, err := getRedis(s)
	if err != nil {
		return false, err
	}

	return conn.HExists(ctx, key, field).Result()
}

// Hget is the implementation of redis hget command.
func (s *Redis) Hget(key, field string) (string, error) {
	return s.HgetCtx(context.Background(), key, field)
}

// HgetCtx is the implementation of redis hget command.
func (s *Redis) HgetCtx(ctx context.Context, key, field string) (string, error) {
	conn, err := getRedis(s)
	if err != nil {
		return "", err
	}

	return conn.HGet(ctx, key, field).Result()
}

// Hgetall is the implementation of redis hgetall command.
func (s *Redis) Hgetall(key string) (map[string]string, error) {
	return s.HgetallCtx(context.Background(), key)
}

// HgetallCtx is the implementation of redis hgetall command.
func (s *Redis) HgetallCtx(ctx context.Context, key string) (map[string]string, error) {
	conn, err := getRedis(s)
	if err != nil {
		return nil, err
	}

	return conn.HGetAll(ctx, key).Result()
}

// Hincrby is the implementation of redis hincrby command.
func (s *Redis) Hincrby(key, field string, increment int) (int, error) {
	return s.HincrbyCtx(context.Background(), key, field, increment)
}

// HincrbyCtx is the implementation of redis hincrby command.
func (s *Redis) HincrbyCtx(ctx context.Context, key, field string, increment int) (int, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	v, err := conn.HIncrBy(ctx, key, field, int64(increment)).Result()
	if err != nil {
		return 0, err
	}

	return int(v), nil
}

// HincrbyFloat is the implementation of redis hincrbyfloat command.
func (s *Redis) HincrbyFloat(key, field string, increment float64) (float64, error) {
	return s.HincrbyFloatCtx(context.Background(), key, field, increment)
}

// HincrbyFloatCtx is the implementation of redis hincrbyfloat command.
func (s *Redis) HincrbyFloatCtx(ctx context.Context, key, field string, increment float64) (
	float64, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	return conn.HIncrByFloat(ctx, key, field, increment).Result()
}

// Hkeys is the implementation of redis hkeys command.
func (s *Redis) Hkeys(key string) ([]string, error) {
	return s.HkeysCtx(context.Background(), key)
}

// HkeysCtx is the implementation of redis hkeys command.
func (s *Redis) HkeysCtx(ctx context.Context, key string) ([]string, error) {
	conn, err := getRedis(s)
	if err != nil {
		return nil, err
	}

	return conn.HKeys(ctx, key).Result()
}

// Hlen is the implementation of redis hlen command.
func (s *Redis) Hlen(key string) (int, error) {
	return s.HlenCtx(context.Background(), key)
}

// HlenCtx is the implementation of redis hlen command.
func (s *Redis) HlenCtx(ctx context.Context, key string) (int, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	v, err := conn.HLen(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	return int(v), nil
}

// Hmget is the implementation of redis hmget command.
func (s *Redis) Hmget(key string, fields ...string) ([]string, error) {
	return s.HmgetCtx(context.Background(), key, fields...)
}

// HmgetCtx is the implementation of redis hmget command.
func (s *Redis) HmgetCtx(ctx context.Context, key string, fields ...string) ([]string, error) {
	conn, err := getRedis(s)
	if err != nil {
		return nil, err
	}

	v, err := conn.HMGet(ctx, key, fields...).Result()
	if err != nil {
		return nil, err
	}

	return toStrings(v), nil
}

// Hset is the implementation of redis hset command.
func (s *Redis) Hset(key, field, value string) error {
	return s.HsetCtx(context.Background(), key, field, value)
}

// HsetCtx is the implementation of redis hset command.
func (s *Redis) HsetCtx(ctx context.Context, key, field, value string) error {
	conn, err := getRedis(s)
	if err != nil {
		return err
	}

	return conn.HSet(ctx, key, field, value).Err()
}

// Hsetnx is the implementation of redis hsetnx command.
func (s *Redis) Hsetnx(key, field, value string) (bool, error) {
	return s.HsetnxCtx(context.Background(), key, field, value)
}

// HsetnxCtx is the implementation of redis hsetnx command.
func (s *Redis) HsetnxCtx(ctx context.Context, key, field, value string) (bool, error) {
	conn, err := getRedis(s)
	if err != nil {
		return false, err
	}

	return conn.HSetNX(ctx, key, field, value).Result()
}

// Hmset is the implementation of redis hmset command.
func (s *Redis) Hmset(key string, fieldsAndValues map[string]string) error {
	return s.HmsetCtx(context.Background(), key, fieldsAndValues)
}

// HmsetCtx is the implementation of redis hmset command.
func (s *Redis) HmsetCtx(ctx context.Context, key string, fieldsAndValues map[string]string) error {
	conn, err := getRedis(s)
	if err != nil {
		return err
	}

	vals := make(map[string]any, len(fieldsAndValues))
	for k, v := range fieldsAndValues {
		vals[k] = v
	}

	return conn.HMSet(ctx, key, vals).Err()
}

// Hscan is the implementation of redis hscan command.
func (s *Redis) Hscan(key string, cursor uint64, match string, count int64) (
	[]string, uint64, error) {
	return s.HscanCtx(context.Background(), key, cursor, match, count)
}

// HscanCtx is the implementation of redis hscan command.
func (s *Redis) HscanCtx(ctx context.Context, key string, cursor uint64, match string, count int64) (
	[]string, uint64, error) {
	conn, err := getRedis(s)
	if err != nil {
		return nil, 0, err
	}

	return conn.HScan(ctx, key, cursor, match, count).Result()
}

// Hvals is the implementation of redis hvals command.
func (s *Redis) Hvals(key string) ([]string, error) {
	return s.HvalsCtx(context.Background(), key)
}

// HvalsCtx is the implementation of redis hvals command.
func (s *Redis) HvalsCtx(ctx context.Context, key string) ([]string, error) {
	conn, err := getRedis(s)
	if err != nil {
		return nil, err
	}

	return conn.HVals(ctx, key).Result()
}

// Incr is the implementation of redis incr command.
func (s *Redis) Incr(key string) (int64, error) {
	return s.IncrCtx(context.Background(), key)
}

// IncrCtx is the implementation of redis incr command.
func (s *Redis) IncrCtx(ctx context.Context, key string) (int64, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	return conn.Incr(ctx, key).Result()
}

// Incrby is the implementation of redis incrby command.
func (s *Redis) Incrby(key string, increment int64) (int64, error) {
	return s.IncrbyCtx(context.Background(), key, increment)
}

// IncrbyCtx is the implementation of redis incrby command.
func (s *Redis) IncrbyCtx(ctx context.Context, key string, increment int64) (int64, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	return conn.IncrBy(ctx, key, increment).Result()
}

// IncrbyFloat is the implementation of redis hincrbyfloat command.
func (s *Redis) IncrbyFloat(key string, increment float64) (float64, error) {
	return s.IncrbyFloatCtx(context.Background(), key, increment)
}

// IncrbyFloatCtx is the implementation of redis hincrbyfloat command.
func (s *Redis) IncrbyFloatCtx(ctx context.Context, key string, increment float64) (float64, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	return conn.IncrByFloat(ctx, key, increment).Result()
}

// Keys is the implementation of redis keys command.
func (s *Redis) Keys(pattern string) ([]string, error) {
	return s.KeysCtx(context.Background(), pattern)
}

// KeysCtx is the implementation of redis keys command.
func (s *Redis) KeysCtx(ctx context.Context, pattern string) ([]string, error) {
	conn, err := getRedis(s)
	if err != nil {
		return nil, err
	}

	return conn.Keys(ctx, pattern).Result()
}

// Llen is the implementation of redis llen command.
func (s *Redis) Llen(key string) (int, error) {
	return s.LlenCtx(context.Background(), key)
}

// LlenCtx is the implementation of redis llen command.
func (s *Redis) LlenCtx(ctx context.Context, key string) (int, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	v, err := conn.LLen(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	return int(v), nil
}

// Lindex is the implementation of redis lindex command.
func (s *Redis) Lindex(key string, index int64) (string, error) {
	return s.LindexCtx(context.Background(), key, index)
}

// LindexCtx is the implementation of redis lindex command.
func (s *Redis) LindexCtx(ctx context.Context, key string, index int64) (string, error) {
	conn, err := getRedis(s)
	if err != nil {
		return "", err
	}

	return conn.LIndex(ctx, key, index).Result()
}

// Lpop is the implementation of redis lpop command.
func (s *Redis) Lpop(key string) (string, error) {
	return s.LpopCtx(context.Background(), key)
}

// LpopCtx is the implementation of redis lpop command.
func (s *Redis) LpopCtx(ctx context.Context, key string) (string, error) {
	conn, err := getRedis(s)
	if err != nil {
		return "", err
	}

	return conn.LPop(ctx, key).Result()
}

// LpopCount is the implementation of redis lpopCount command.
func (s *Redis) LpopCount(key string, count int) ([]string, error) {
	return s.LpopCountCtx(context.Background(), key, count)
}

// LpopCountCtx is the implementation of redis lpopCount command.
func (s *Redis) LpopCountCtx(ctx context.Context, key string, count int) ([]string, error) {
	conn, err := getRedis(s)
	if err != nil {
		return nil, err
	}

	return conn.LPopCount(ctx, key, count).Result()
}

// Lpush is the implementation of redis lpush command.
func (s *Redis) Lpush(key string, values ...any) (int, error) {
	return s.LpushCtx(context.Background(), key, values...)
}

// LpushCtx is the implementation of redis lpush command.
func (s *Redis) LpushCtx(ctx context.Context, key string, values ...any) (int, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	v, err := conn.LPush(ctx, key, values...).Result()
	if err != nil {
		return 0, err
	}

	return int(v), nil
}

// Lrange is the implementation of redis lrange command.
func (s *Redis) Lrange(key string, start, stop int) ([]string, error) {
	return s.LrangeCtx(context.Background(), key, start, stop)
}

// LrangeCtx is the implementation of redis lrange command.
func (s *Redis) LrangeCtx(ctx context.Context, key string, start, stop int) ([]string, error) {
	conn, err := getRedis(s)
	if err != nil {
		return nil, err
	}

	return conn.LRange(ctx, key, int64(start), int64(stop)).Result()
}

// Lrem is the implementation of redis lrem command.
func (s *Redis) Lrem(key string, count int, value string) (int, error) {
	return s.LremCtx(context.Background(), key, count, value)
}

// LremCtx is the implementation of redis lrem command.
func (s *Redis) LremCtx(ctx context.Context, key string, count int, value string) (int, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	v, err := conn.LRem(ctx, key, int64(count), value).Result()
	if err != nil {
		return 0, err
	}

	return int(v), nil
}

// Ltrim is the implementation of redis ltrim command.
func (s *Redis) Ltrim(key string, start, stop int64) error {
	return s.LtrimCtx(context.Background(), key, start, stop)
}

// LtrimCtx is the implementation of redis ltrim command.
func (s *Redis) LtrimCtx(ctx context.Context, key string, start, stop int64) error {
	conn, err := getRedis(s)
	if err != nil {
		return err
	}

	return conn.LTrim(ctx, key, start, stop).Err()
}

// Mget is the implementation of redis mget command.
func (s *Redis) Mget(keys ...string) ([]string, error) {
	return s.MgetCtx(context.Background(), keys...)
}

// MgetCtx is the implementation of redis mget command.
func (s *Redis) MgetCtx(ctx context.Context, keys ...string) ([]string, error) {
	conn, err := getRedis(s)
	if err != nil {
		return nil, err
	}

	v, err := conn.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	return toStrings(v), nil
}

// Mset is the implementation of redis mset command.
func (s *Redis) Mset(fieldsAndValues ...any) (string, error) {
	return s.MsetCtx(context.Background(), fieldsAndValues...)
}

// MsetCtx is the implementation of redis mset command.
func (s *Redis) MsetCtx(ctx context.Context, fieldsAndValues ...any) (string, error) {
	conn, err := getRedis(s)
	if err != nil {
		return "", err
	}

	return conn.MSet(ctx, fieldsAndValues...).Result()
}

// Persist is the implementation of redis persist command.
func (s *Redis) Persist(key string) (bool, error) {
	return s.PersistCtx(context.Background(), key)
}

// PersistCtx is the implementation of redis persist command.
func (s *Redis) PersistCtx(ctx context.Context, key string) (bool, error) {
	conn, err := getRedis(s)
	if err != nil {
		return false, err
	}

	return conn.Persist(ctx, key).Result()
}

// Pfadd is the implementation of redis pfadd command.
func (s *Redis) Pfadd(key string, values ...any) (bool, error) {
	return s.PfaddCtx(context.Background(), key, values...)
}

// PfaddCtx is the implementation of redis pfadd command.
func (s *Redis) PfaddCtx(ctx context.Context, key string, values ...any) (bool, error) {
	conn, err := getRedis(s)
	if err != nil {
		return false, err
	}

	v, err := conn.PFAdd(ctx, key, values...).Result()
	if err != nil {
		return false, err
	}

	return v >= 1, nil
}

// Pfcount is the implementation of redis pfcount command.
func (s *Redis) Pfcount(key string) (int64, error) {
	return s.PfcountCtx(context.Background(), key)
}

// PfcountCtx is the implementation of redis pfcount command.
func (s *Redis) PfcountCtx(ctx context.Context, key string) (int64, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	return conn.PFCount(ctx, key).Result()
}

// Pfmerge is the implementation of redis pfmerge command.
func (s *Redis) Pfmerge(dest string, keys ...string) error {
	return s.PfmergeCtx(context.Background(), dest, keys...)
}

// PfmergeCtx is the implementation of redis pfmerge command.
func (s *Redis) PfmergeCtx(ctx context.Context, dest string, keys ...string) error {
	conn, err := getRedis(s)
	if err != nil {
		return err
	}

	_, err = conn.PFMerge(ctx, dest, keys...).Result()
	return err
}

// Ping is the implementation of redis ping command.
func (s *Redis) Ping() bool {
	return s.PingCtx(context.Background())
}

// PingCtx is the implementation of redis ping command.
func (s *Redis) PingCtx(ctx context.Context) bool {
	// ignore error, error means false
	conn, err := getRedis(s)
	if err != nil {
		return false
	}

	v, err := conn.Ping(ctx).Result()
	if err != nil {
		return false
	}

	return v == "PONG"
}

// Pipelined lets fn execute pipelined commands.
func (s *Redis) Pipelined(fn func(Pipeliner) error) error {
	return s.PipelinedCtx(context.Background(), fn)
}

// PipelinedCtx lets fn execute pipelined commands.
// Results need to be retrieved by calling Pipeline.Exec()
func (s *Redis) PipelinedCtx(ctx context.Context, fn func(Pipeliner) error) error {
	conn, err := getRedis(s)
	if err != nil {
		return err
	}

	_, err = conn.Pipelined(ctx, fn)
	return err
}

func (s *Redis) Publish(channel string, message interface{}) (int64, error) {
	return s.PublishCtx(context.Background(), channel, message)
}

func (s *Redis) PublishCtx(ctx context.Context, channel string, message interface{}) (int64, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}
	return conn.Publish(ctx, channel, message).Result()
}

// Rpop is the implementation of redis rpop command.
func (s *Redis) Rpop(key string) (string, error) {
	return s.RpopCtx(context.Background(), key)
}

// RpopCtx is the implementation of redis rpop command.
func (s *Redis) RpopCtx(ctx context.Context, key string) (string, error) {
	conn, err := getRedis(s)
	if err != nil {
		return "", err
	}

	return conn.RPop(ctx, key).Result()
}

// RpopCount is the implementation of redis rpopCount command.
func (s *Redis) RpopCount(key string, count int) ([]string, error) {
	return s.RpopCountCtx(context.Background(), key, count)
}

// RpopCountCtx is the implementation of redis rpopCount command.
func (s *Redis) RpopCountCtx(ctx context.Context, key string, count int) ([]string, error) {
	conn, err := getRedis(s)
	if err != nil {
		return nil, err
	}

	return conn.RPopCount(ctx, key, count).Result()
}

// Rpush is the implementation of redis rpush command.
func (s *Redis) Rpush(key string, values ...any) (int, error) {
	return s.RpushCtx(context.Background(), key, values...)
}

// RpushCtx is the implementation of redis rpush command.
func (s *Redis) RpushCtx(ctx context.Context, key string, values ...any) (int, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	v, err := conn.RPush(ctx, key, values...).Result()
	if err != nil {
		return 0, err
	}

	return int(v), nil
}

func (s *Redis) RPopLPush(source string, destination string) (string, error) {
	return s.RPopLPushCtx(context.Background(), source, destination)
}

func (s *Redis) RPopLPushCtx(ctx context.Context, source string, destination string) (string, error) {
	conn, err := getRedis(s)
	if err != nil {
		return "", err
	}
	return conn.RPopLPush(ctx, source, destination).Result()
}

// Sadd is the implementation of redis sadd command.
func (s *Redis) Sadd(key string, values ...any) (int, error) {
	return s.SaddCtx(context.Background(), key, values...)
}

// SaddCtx is the implementation of redis sadd command.
func (s *Redis) SaddCtx(ctx context.Context, key string, values ...any) (int, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	v, err := conn.SAdd(ctx, key, values...).Result()
	if err != nil {
		return 0, err
	}

	return int(v), nil
}

// Scan is the implementation of redis scan command.
func (s *Redis) Scan(cursor uint64, match string, count int64) ([]string, uint64, error) {
	return s.ScanCtx(context.Background(), cursor, match, count)
}

// ScanCtx is the implementation of redis scan command.
func (s *Redis) ScanCtx(ctx context.Context, cursor uint64, match string, count int64) (
	[]string, uint64, error) {
	conn, err := getRedis(s)
	if err != nil {
		return nil, 0, err
	}

	return conn.Scan(ctx, cursor, match, count).Result()
}

// SetBit is the implementation of redis setbit command.
func (s *Redis) SetBit(key string, offset int64, value int) (int, error) {
	return s.SetBitCtx(context.Background(), key, offset, value)
}

// SetBitCtx is the implementation of redis setbit command.
func (s *Redis) SetBitCtx(ctx context.Context, key string, offset int64, value int) (int, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	v, err := conn.SetBit(ctx, key, offset, value).Result()
	if err != nil {
		return 0, err
	}

	return int(v), nil
}

// Sscan is the implementation of redis sscan command.
func (s *Redis) Sscan(key string, cursor uint64, match string, count int64) (
	[]string, uint64, error) {
	return s.SscanCtx(context.Background(), key, cursor, match, count)
}

// SscanCtx is the implementation of redis sscan command.
func (s *Redis) SscanCtx(ctx context.Context, key string, cursor uint64, match string, count int64) (
	[]string, uint64, error) {
	conn, err := getRedis(s)
	if err != nil {
		return nil, 0, err
	}

	return conn.SScan(ctx, key, cursor, match, count).Result()
}

// Scard is the implementation of redis scard command.
func (s *Redis) Scard(key string) (int64, error) {
	return s.ScardCtx(context.Background(), key)
}

// ScardCtx is the implementation of redis scard command.
func (s *Redis) ScardCtx(ctx context.Context, key string) (int64, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	return conn.SCard(ctx, key).Result()
}

// ScriptLoad is the implementation of redis script load command.
func (s *Redis) ScriptLoad(script string) (string, error) {
	return s.ScriptLoadCtx(context.Background(), script)
}

// ScriptLoadCtx is the implementation of redis script load command.
func (s *Redis) ScriptLoadCtx(ctx context.Context, script string) (string, error) {
	conn, err := getRedis(s)
	if err != nil {
		return "", err
	}

	return conn.ScriptLoad(ctx, script).Result()
}

// ScriptRun is the implementation of *redis.Script run command.
func (s *Redis) ScriptRun(script *Script, keys []string, args ...any) (any, error) {
	return s.ScriptRunCtx(context.Background(), script, keys, args...)
}

// ScriptRunCtx is the implementation of *redis.Script run command.
func (s *Redis) ScriptRunCtx(ctx context.Context, script *Script, keys []string,
	args ...any) (any, error) {
	conn, err := getRedis(s)
	if err != nil {
		return nil, err
	}

	return script.Run(ctx, conn, keys, args...).Result()
}

// Set is the implementation of redis set command.
func (s *Redis) Set(key, value string) error {
	return s.SetCtx(context.Background(), key, value)
}

// SetCtx is the implementation of redis set command.
func (s *Redis) SetCtx(ctx context.Context, key, value string) error {
	conn, err := getRedis(s)
	if err != nil {
		return err
	}

	return conn.Set(ctx, key, value, 0).Err()
}

// Setex is the implementation of redis setex command.
func (s *Redis) Setex(key, value string, seconds int) error {
	return s.SetexCtx(context.Background(), key, value, seconds)
}

// SetexCtx is the implementation of redis setex command.
func (s *Redis) SetexCtx(ctx context.Context, key, value string, seconds int) error {
	conn, err := getRedis(s)
	if err != nil {
		return err
	}

	return conn.Set(ctx, key, value, time.Duration(seconds)*time.Second).Err()
}

// Setnx is the implementation of redis setnx command.
func (s *Redis) Setnx(key, value string) (bool, error) {
	return s.SetnxCtx(context.Background(), key, value)
}

// SetnxCtx is the implementation of redis setnx command.
func (s *Redis) SetnxCtx(ctx context.Context, key, value string) (bool, error) {
	conn, err := getRedis(s)
	if err != nil {
		return false, err
	}

	return conn.SetNX(ctx, key, value, 0).Result()
}

// SetnxEx is the implementation of redis setnx command with expire.
func (s *Redis) SetnxEx(key, value string, seconds int) (bool, error) {
	return s.SetnxExCtx(context.Background(), key, value, seconds)
}

// SetnxExCtx is the implementation of redis setnx command with expire.
func (s *Redis) SetnxExCtx(ctx context.Context, key, value string, seconds int) (bool, error) {
	conn, err := getRedis(s)
	if err != nil {
		return false, err
	}

	return conn.SetNX(ctx, key, value, time.Duration(seconds)*time.Second).Result()
}

// Sismember is the implementation of redis sismember command.
func (s *Redis) Sismember(key string, value any) (bool, error) {
	return s.SismemberCtx(context.Background(), key, value)
}

// SismemberCtx is the implementation of redis sismember command.
func (s *Redis) SismemberCtx(ctx context.Context, key string, value any) (bool, error) {
	conn, err := getRedis(s)
	if err != nil {
		return false, err
	}

	return conn.SIsMember(ctx, key, value).Result()
}

// Smembers is the implementation of redis smembers command.
func (s *Redis) Smembers(key string) ([]string, error) {
	return s.SmembersCtx(context.Background(), key)
}

// SmembersCtx is the implementation of redis smembers command.
func (s *Redis) SmembersCtx(ctx context.Context, key string) ([]string, error) {
	conn, err := getRedis(s)
	if err != nil {
		return nil, err
	}

	return conn.SMembers(ctx, key).Result()
}

// Spop is the implementation of redis spop command.
func (s *Redis) Spop(key string) (string, error) {
	return s.SpopCtx(context.Background(), key)
}

// SpopCtx is the implementation of redis spop command.
func (s *Redis) SpopCtx(ctx context.Context, key string) (string, error) {
	conn, err := getRedis(s)
	if err != nil {
		return "", err
	}

	return conn.SPop(ctx, key).Result()
}

// Srandmember is the implementation of redis srandmember command.
func (s *Redis) Srandmember(key string, count int) ([]string, error) {
	return s.SrandmemberCtx(context.Background(), key, count)
}

// SrandmemberCtx is the implementation of redis srandmember command.
func (s *Redis) SrandmemberCtx(ctx context.Context, key string, count int) ([]string, error) {
	conn, err := getRedis(s)
	if err != nil {
		return nil, err
	}

	return conn.SRandMemberN(ctx, key, int64(count)).Result()
}

// Srem is the implementation of redis srem command.
func (s *Redis) Srem(key string, values ...any) (int, error) {
	return s.SremCtx(context.Background(), key, values...)
}

// SremCtx is the implementation of redis srem command.
func (s *Redis) SremCtx(ctx context.Context, key string, values ...any) (int, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	v, err := conn.SRem(ctx, key, values...).Result()
	if err != nil {
		return 0, err
	}

	return int(v), nil
}

// String returns the string representation of s.
func (s *Redis) String() string {
	return s.Addr
}

// Sunion is the implementation of redis sunion command.
func (s *Redis) Sunion(keys ...string) ([]string, error) {
	return s.SunionCtx(context.Background(), keys...)
}

// SunionCtx is the implementation of redis sunion command.
func (s *Redis) SunionCtx(ctx context.Context, keys ...string) ([]string, error) {
	conn, err := getRedis(s)
	if err != nil {
		return nil, err
	}

	return conn.SUnion(ctx, keys...).Result()
}

// Sunionstore is the implementation of redis sunionstore command.
func (s *Redis) Sunionstore(destination string, keys ...string) (int, error) {
	return s.SunionstoreCtx(context.Background(), destination, keys...)
}

// SunionstoreCtx is the implementation of redis sunionstore command.
func (s *Redis) SunionstoreCtx(ctx context.Context, destination string, keys ...string) (
	int, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	v, err := conn.SUnionStore(ctx, destination, keys...).Result()
	if err != nil {
		return 0, err
	}

	return int(v), nil
}

// Sdiff is the implementation of redis sdiff command.
func (s *Redis) Sdiff(keys ...string) ([]string, error) {
	return s.SdiffCtx(context.Background(), keys...)
}

// SdiffCtx is the implementation of redis sdiff command.
func (s *Redis) SdiffCtx(ctx context.Context, keys ...string) ([]string, error) {
	conn, err := getRedis(s)
	if err != nil {
		return nil, err
	}

	return conn.SDiff(ctx, keys...).Result()
}

// Sdiffstore is the implementation of redis sdiffstore command.
func (s *Redis) Sdiffstore(destination string, keys ...string) (int, error) {
	return s.SdiffstoreCtx(context.Background(), destination, keys...)
}

// SdiffstoreCtx is the implementation of redis sdiffstore command.
func (s *Redis) SdiffstoreCtx(ctx context.Context, destination string, keys ...string) (
	int, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	v, err := conn.SDiffStore(ctx, destination, keys...).Result()
	if err != nil {
		return 0, err
	}

	return int(v), nil
}

// Sinter is the implementation of redis sinter command.
func (s *Redis) Sinter(keys ...string) ([]string, error) {
	return s.SinterCtx(context.Background(), keys...)
}

// SinterCtx is the implementation of redis sinter command.
func (s *Redis) SinterCtx(ctx context.Context, keys ...string) ([]string, error) {
	conn, err := getRedis(s)
	if err != nil {
		return nil, err
	}

	return conn.SInter(ctx, keys...).Result()
}

// Sinterstore is the implementation of redis sinterstore command.
func (s *Redis) Sinterstore(destination string, keys ...string) (int, error) {
	return s.SinterstoreCtx(context.Background(), destination, keys...)
}

// SinterstoreCtx is the implementation of redis sinterstore command.
func (s *Redis) SinterstoreCtx(ctx context.Context, destination string, keys ...string) (
	int, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	v, err := conn.SInterStore(ctx, destination, keys...).Result()
	if err != nil {
		return 0, err
	}

	return int(v), nil
}

// Ttl is the implementation of redis ttl command.
func (s *Redis) Ttl(key string) (int, error) {
	return s.TtlCtx(context.Background(), key)
}

// TtlCtx is the implementation of redis ttl command.
func (s *Redis) TtlCtx(ctx context.Context, key string) (int, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	duration, err := conn.TTL(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	if duration >= 0 {
		return int(duration / time.Second), nil
	}

	// -2 means key does not exist
	// -1 means key exists but has no expire
	return int(duration), nil
}

func (s *Redis) TxPipeline() (pipe Pipeliner, err error) {
	conn, err := getRedis(s)
	if err != nil {
		return nil, err
	}
	return conn.TxPipeline(), nil
}

func (s *Redis) Unlink(keys ...string) (int64, error) {
	return s.UnlinkCtx(context.Background(), keys...)
}

func (s *Redis) UnlinkCtx(ctx context.Context, keys ...string) (int64, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}
	return conn.Unlink(ctx, keys...).Result()
}

// Zadd is the implementation of redis zadd command.
func (s *Redis) Zadd(key string, score int64, value string) (bool, error) {
	return s.ZaddCtx(context.Background(), key, score, value)
}

// ZaddCtx is the implementation of redis zadd command.
func (s *Redis) ZaddCtx(ctx context.Context, key string, score int64, value string) (bool, error) {
	return s.ZaddFloatCtx(ctx, key, float64(score), value)
}

// ZaddFloat is the implementation of redis zadd command.
func (s *Redis) ZaddFloat(key string, score float64, value string) (bool, error) {
	return s.ZaddFloatCtx(context.Background(), key, score, value)
}

// ZaddFloatCtx is the implementation of redis zadd command.
func (s *Redis) ZaddFloatCtx(ctx context.Context, key string, score float64, value string) (
	bool, error) {
	conn, err := getRedis(s)
	if err != nil {
		return false, err
	}

	v, err := conn.ZAdd(ctx, key, red.Z{
		Score:  score,
		Member: value,
	}).Result()
	if err != nil {
		return false, err
	}

	return v == 1, nil
}

// Zaddnx is the implementation of redis zadd nx command.
func (s *Redis) Zaddnx(key string, score int64, value string) (bool, error) {
	return s.ZaddnxCtx(context.Background(), key, score, value)
}

// ZaddnxCtx is the implementation of redis zadd nx command.
func (s *Redis) ZaddnxCtx(ctx context.Context, key string, score int64, value string) (bool, error) {
	return s.ZaddnxFloatCtx(ctx, key, float64(score), value)
}

// ZaddnxFloat is the implementation of redis zaddnx command.
func (s *Redis) ZaddnxFloat(key string, score float64, value string) (bool, error) {
	return s.ZaddFloatCtx(context.Background(), key, score, value)
}

// ZaddnxFloatCtx is the implementation of redis zaddnx command.
func (s *Redis) ZaddnxFloatCtx(ctx context.Context, key string, score float64, value string) (
	bool, error) {
	conn, err := getRedis(s)
	if err != nil {
		return false, err
	}

	v, err := conn.ZAddNX(ctx, key, red.Z{
		Score:  score,
		Member: value,
	}).Result()
	if err != nil {
		return false, err
	}

	return v == 1, nil
}

// Zadds is the implementation of redis zadds command.
func (s *Redis) Zadds(key string, ps ...Pair) (int64, error) {
	return s.ZaddsCtx(context.Background(), key, ps...)
}

// ZaddsCtx is the implementation of redis zadds command.
func (s *Redis) ZaddsCtx(ctx context.Context, key string, ps ...Pair) (int64, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	var zs []red.Z
	for _, p := range ps {
		z := red.Z{Score: float64(p.Score), Member: p.Key}
		zs = append(zs, z)
	}

	return conn.ZAdd(ctx, key, zs...).Result()
}

// Zcard is the implementation of redis zcard command.
func (s *Redis) Zcard(key string) (int, error) {
	return s.ZcardCtx(context.Background(), key)
}

// ZcardCtx is the implementation of redis zcard command.
func (s *Redis) ZcardCtx(ctx context.Context, key string) (int, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	v, err := conn.ZCard(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	return int(v), nil
}

// Zcount is the implementation of redis zcount command.
func (s *Redis) Zcount(key string, start, stop int64) (int, error) {
	return s.ZcountCtx(context.Background(), key, start, stop)
}

// ZcountCtx is the implementation of redis zcount command.
func (s *Redis) ZcountCtx(ctx context.Context, key string, start, stop int64) (int, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	v, err := conn.ZCount(ctx, key, strconv.FormatInt(start, 10),
		strconv.FormatInt(stop, 10)).Result()
	if err != nil {
		return 0, err
	}

	return int(v), nil
}

// Zincrby is the implementation of redis zincrby command.
func (s *Redis) Zincrby(key string, increment int64, field string) (int64, error) {
	return s.ZincrbyCtx(context.Background(), key, increment, field)
}

// ZincrbyCtx is the implementation of redis zincrby command.
func (s *Redis) ZincrbyCtx(ctx context.Context, key string, increment int64, field string) (
	int64, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	v, err := conn.ZIncrBy(ctx, key, float64(increment), field).Result()
	if err != nil {
		return 0, err
	}

	return int64(v), nil
}

// Zscore is the implementation of redis zscore command.
func (s *Redis) Zscore(key, value string) (int64, error) {
	return s.ZscoreCtx(context.Background(), key, value)
}

// ZscoreCtx is the implementation of redis zscore command.
func (s *Redis) ZscoreCtx(ctx context.Context, key, value string) (int64, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	v, err := conn.ZScore(ctx, key, value).Result()
	if err != nil {
		return 0, err
	}

	return int64(v), nil
}

// ZscoreByFloat is the implementation of redis zscore command score by float.
func (s *Redis) ZscoreByFloat(key, value string) (float64, error) {
	return s.ZscoreByFloatCtx(context.Background(), key, value)
}

// ZscoreByFloatCtx is the implementation of redis zscore command score by float.
func (s *Redis) ZscoreByFloatCtx(ctx context.Context, key, value string) (float64, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	return conn.ZScore(ctx, key, value).Result()
}

// Zscan is the implementation of redis zscan command.
func (s *Redis) Zscan(key string, cursor uint64, match string, count int64) (
	[]string, uint64, error) {
	return s.ZscanCtx(context.Background(), key, cursor, match, count)
}

// ZscanCtx is the implementation of redis zscan command.
func (s *Redis) ZscanCtx(ctx context.Context, key string, cursor uint64, match string, count int64) (
	[]string, uint64, error) {
	conn, err := getRedis(s)
	if err != nil {
		return nil, 0, err
	}

	return conn.ZScan(ctx, key, cursor, match, count).Result()
}

// Zrank is the implementation of redis zrank command.
func (s *Redis) Zrank(key, field string) (int64, error) {
	return s.ZrankCtx(context.Background(), key, field)
}

// ZrankCtx is the implementation of redis zrank command.
func (s *Redis) ZrankCtx(ctx context.Context, key, field string) (int64, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	return conn.ZRank(ctx, key, field).Result()
}

// Zrem is the implementation of redis zrem command.
func (s *Redis) Zrem(key string, values ...any) (int, error) {
	return s.ZremCtx(context.Background(), key, values...)
}

// ZremCtx is the implementation of redis zrem command.
func (s *Redis) ZremCtx(ctx context.Context, key string, values ...any) (int, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	v, err := conn.ZRem(ctx, key, values...).Result()
	if err != nil {
		return 0, err
	}

	return int(v), nil
}

// Zremrangebyscore is the implementation of redis zremrangebyscore command.
func (s *Redis) Zremrangebyscore(key string, start, stop int64) (int, error) {
	return s.ZremrangebyscoreCtx(context.Background(), key, start, stop)
}

// ZremrangebyscoreCtx is the implementation of redis zremrangebyscore command.
func (s *Redis) ZremrangebyscoreCtx(ctx context.Context, key string, start, stop int64) (
	int, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	v, err := conn.ZRemRangeByScore(ctx, key, strconv.FormatInt(start, 10),
		strconv.FormatInt(stop, 10)).Result()
	if err != nil {
		return 0, err
	}

	return int(v), nil
}

// Zremrangebyrank is the implementation of redis zremrangebyrank command.
func (s *Redis) Zremrangebyrank(key string, start, stop int64) (int, error) {
	return s.ZremrangebyrankCtx(context.Background(), key, start, stop)
}

// ZremrangebyrankCtx is the implementation of redis zremrangebyrank command.
func (s *Redis) ZremrangebyrankCtx(ctx context.Context, key string, start, stop int64) (
	int, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	v, err := conn.ZRemRangeByRank(ctx, key, start, stop).Result()
	if err != nil {
		return 0, err
	}

	return int(v), nil
}

// Zrange is the implementation of redis zrange command.
func (s *Redis) Zrange(key string, start, stop int64) ([]string, error) {
	return s.ZrangeCtx(context.Background(), key, start, stop)
}

// ZrangeCtx is the implementation of redis zrange command.
func (s *Redis) ZrangeCtx(ctx context.Context, key string, start, stop int64) (
	[]string, error) {
	conn, err := getRedis(s)
	if err != nil {
		return nil, err
	}

	return conn.ZRange(ctx, key, start, stop).Result()
}

// ZrangeWithScores is the implementation of redis zrange command with scores.
func (s *Redis) ZrangeWithScores(key string, start, stop int64) ([]Pair, error) {
	return s.ZrangeWithScoresCtx(context.Background(), key, start, stop)
}

// ZrangeWithScoresCtx is the implementation of redis zrange command with scores.
func (s *Redis) ZrangeWithScoresCtx(ctx context.Context, key string, start, stop int64) (
	[]Pair, error) {
	conn, err := getRedis(s)
	if err != nil {
		return nil, err
	}

	v, err := conn.ZRangeWithScores(ctx, key, start, stop).Result()
	if err != nil {
		return nil, err
	}

	return toPairs(v), nil
}

// ZrangeWithScoresByFloat is the implementation of redis zrange command with scores by float64.
func (s *Redis) ZrangeWithScoresByFloat(key string, start, stop int64) ([]FloatPair, error) {
	return s.ZrangeWithScoresByFloatCtx(context.Background(), key, start, stop)
}

// ZrangeWithScoresByFloatCtx is the implementation of redis zrange command with scores by float64.
func (s *Redis) ZrangeWithScoresByFloatCtx(ctx context.Context, key string, start, stop int64) (
	[]FloatPair, error) {
	conn, err := getRedis(s)
	if err != nil {
		return nil, err
	}

	v, err := conn.ZRangeWithScores(ctx, key, start, stop).Result()
	if err != nil {
		return nil, err
	}

	return toFloatPairs(v), nil
}

// ZRevRangeWithScores is the implementation of redis zrevrange command with scores.
// Deprecated: use ZrevrangeWithScores instead.
func (s *Redis) ZRevRangeWithScores(key string, start, stop int64) ([]Pair, error) {
	return s.ZrevrangeWithScoresCtx(context.Background(), key, start, stop)
}

// ZrevrangeWithScores is the implementation of redis zrevrange command with scores.
func (s *Redis) ZrevrangeWithScores(key string, start, stop int64) ([]Pair, error) {
	return s.ZrevrangeWithScoresCtx(context.Background(), key, start, stop)
}

// ZRevRangeWithScoresCtx is the implementation of redis zrevrange command with scores.
// Deprecated: use ZrevrangeWithScoresCtx instead.
func (s *Redis) ZRevRangeWithScoresCtx(ctx context.Context, key string, start, stop int64) (
	[]Pair, error) {
	return s.ZrevrangeWithScoresCtx(ctx, key, start, stop)
}

// ZrevrangeWithScoresCtx is the implementation of redis zrevrange command with scores.
func (s *Redis) ZrevrangeWithScoresCtx(ctx context.Context, key string, start, stop int64) (
	[]Pair, error) {
	conn, err := getRedis(s)
	if err != nil {
		return nil, err
	}

	v, err := conn.ZRevRangeWithScores(ctx, key, start, stop).Result()
	if err != nil {
		return nil, err
	}

	return toPairs(v), nil
}

// ZRevRangeWithScoresByFloat is the implementation of redis zrevrange command with scores by float.
// Deprecated: use ZrevrangeWithScoresByFloat instead.
func (s *Redis) ZRevRangeWithScoresByFloat(key string, start, stop int64) ([]FloatPair, error) {
	return s.ZrevrangeWithScoresByFloatCtx(context.Background(), key, start, stop)
}

// ZrevrangeWithScoresByFloat is the implementation of redis zrevrange command with scores by float.
func (s *Redis) ZrevrangeWithScoresByFloat(key string, start, stop int64) ([]FloatPair, error) {
	return s.ZrevrangeWithScoresByFloatCtx(context.Background(), key, start, stop)
}

// ZRevRangeWithScoresByFloatCtx is the implementation of redis zrevrange command with scores by float.
// Deprecated: use ZrevrangeWithScoresByFloatCtx instead.
func (s *Redis) ZRevRangeWithScoresByFloatCtx(ctx context.Context, key string, start, stop int64) (
	[]FloatPair, error) {
	return s.ZrevrangeWithScoresByFloatCtx(ctx, key, start, stop)
}

// ZrevrangeWithScoresByFloatCtx is the implementation of redis zrevrange command with scores by float.
func (s *Redis) ZrevrangeWithScoresByFloatCtx(ctx context.Context, key string, start, stop int64) (
	[]FloatPair, error) {
	conn, err := getRedis(s)
	if err != nil {
		return nil, err
	}

	v, err := conn.ZRevRangeWithScores(ctx, key, start, stop).Result()
	if err != nil {
		return nil, err
	}

	return toFloatPairs(v), nil
}

// ZrangebyscoreWithScores is the implementation of redis zrangebyscore command with scores.
func (s *Redis) ZrangebyscoreWithScores(key string, start, stop int64) ([]Pair, error) {
	return s.ZrangebyscoreWithScoresCtx(context.Background(), key, start, stop)
}

// ZrangebyscoreWithScoresCtx is the implementation of redis zrangebyscore command with scores.
func (s *Redis) ZrangebyscoreWithScoresCtx(ctx context.Context, key string, start, stop int64) (
	[]Pair, error) {
	conn, err := getRedis(s)
	if err != nil {
		return nil, err
	}

	v, err := conn.ZRangeByScoreWithScores(ctx, key, &red.ZRangeBy{
		Min: strconv.FormatInt(start, 10),
		Max: strconv.FormatInt(stop, 10),
	}).Result()
	if err != nil {
		return nil, err
	}

	return toPairs(v), nil
}

// ZrangebyscoreWithScoresByFloat is the implementation of redis zrangebyscore command with scores by float.
func (s *Redis) ZrangebyscoreWithScoresByFloat(key string, start, stop float64) (
	[]FloatPair, error) {
	return s.ZrangebyscoreWithScoresByFloatCtx(context.Background(), key, start, stop)
}

// ZrangebyscoreWithScoresByFloatCtx is the implementation of redis zrangebyscore command with scores by float.
func (s *Redis) ZrangebyscoreWithScoresByFloatCtx(ctx context.Context, key string, start, stop float64) (
	[]FloatPair, error) {
	conn, err := getRedis(s)
	if err != nil {
		return nil, err
	}

	v, err := conn.ZRangeByScoreWithScores(ctx, key, &red.ZRangeBy{
		Min: strconv.FormatFloat(start, 'f', -1, 64),
		Max: strconv.FormatFloat(stop, 'f', -1, 64),
	}).Result()
	if err != nil {
		return nil, err
	}

	return toFloatPairs(v), nil
}

// ZrangebyscoreWithScoresAndLimit is the implementation of redis zrangebyscore command
// with scores and limit.
func (s *Redis) ZrangebyscoreWithScoresAndLimit(key string, start, stop int64,
	page, size int) ([]Pair, error) {
	return s.ZrangebyscoreWithScoresAndLimitCtx(context.Background(), key, start, stop, page, size)
}

// ZrangebyscoreWithScoresAndLimitCtx is the implementation of redis zrangebyscore command
// with scores and limit.
func (s *Redis) ZrangebyscoreWithScoresAndLimitCtx(ctx context.Context, key string, start,
	stop int64, page, size int) ([]Pair, error) {
	if size <= 0 {
		return nil, nil
	}

	conn, err := getRedis(s)
	if err != nil {
		return nil, err
	}

	v, err := conn.ZRangeByScoreWithScores(ctx, key, &red.ZRangeBy{
		Min:    strconv.FormatInt(start, 10),
		Max:    strconv.FormatInt(stop, 10),
		Offset: int64(page * size),
		Count:  int64(size),
	}).Result()
	if err != nil {
		return nil, err
	}

	return toPairs(v), nil
}

// ZrangebyscoreWithScoresByFloatAndLimit is the implementation of redis zrangebyscore command
// with scores by float and limit.
func (s *Redis) ZrangebyscoreWithScoresByFloatAndLimit(key string, start, stop float64,
	page, size int) ([]FloatPair, error) {
	return s.ZrangebyscoreWithScoresByFloatAndLimitCtx(context.Background(),
		key, start, stop, page, size)
}

// ZrangebyscoreWithScoresByFloatAndLimitCtx is the implementation of redis zrangebyscore command
// with scores by float and limit.
func (s *Redis) ZrangebyscoreWithScoresByFloatAndLimitCtx(ctx context.Context, key string, start,
	stop float64, page, size int) ([]FloatPair, error) {
	if size <= 0 {
		return nil, nil
	}

	conn, err := getRedis(s)
	if err != nil {
		return nil, err
	}

	v, err := conn.ZRangeByScoreWithScores(ctx, key, &red.ZRangeBy{
		Min:    strconv.FormatFloat(start, 'f', -1, 64),
		Max:    strconv.FormatFloat(stop, 'f', -1, 64),
		Offset: int64(page * size),
		Count:  int64(size),
	}).Result()
	if err != nil {
		return nil, err
	}

	return toFloatPairs(v), nil
}

// Zrevrange is the implementation of redis zrevrange command.
func (s *Redis) Zrevrange(key string, start, stop int64) ([]string, error) {
	return s.ZrevrangeCtx(context.Background(), key, start, stop)
}

// ZrevrangeCtx is the implementation of redis zrevrange command.
func (s *Redis) ZrevrangeCtx(ctx context.Context, key string, start, stop int64) (
	[]string, error) {
	conn, err := getRedis(s)
	if err != nil {
		return nil, err
	}

	return conn.ZRevRange(ctx, key, start, stop).Result()
}

// ZrevrangebyscoreWithScores is the implementation of redis zrevrangebyscore command with scores.
func (s *Redis) ZrevrangebyscoreWithScores(key string, start, stop int64) ([]Pair, error) {
	return s.ZrevrangebyscoreWithScoresCtx(context.Background(), key, start, stop)
}

// ZrevrangebyscoreWithScoresCtx is the implementation of redis zrevrangebyscore command with scores.
func (s *Redis) ZrevrangebyscoreWithScoresCtx(ctx context.Context, key string, start, stop int64) (
	[]Pair, error) {
	conn, err := getRedis(s)
	if err != nil {
		return nil, err
	}

	v, err := conn.ZRevRangeByScoreWithScores(ctx, key, &red.ZRangeBy{
		Min: strconv.FormatInt(start, 10),
		Max: strconv.FormatInt(stop, 10),
	}).Result()
	if err != nil {
		return nil, err
	}

	return toPairs(v), nil
}

// ZrevrangebyscoreWithScoresByFloat is the implementation of redis zrevrangebyscore command with scores by float.
func (s *Redis) ZrevrangebyscoreWithScoresByFloat(key string, start, stop float64) (
	[]FloatPair, error) {
	return s.ZrevrangebyscoreWithScoresByFloatCtx(context.Background(), key, start, stop)
}

// ZrevrangebyscoreWithScoresByFloatCtx is the implementation of redis zrevrangebyscore command with scores by float.
func (s *Redis) ZrevrangebyscoreWithScoresByFloatCtx(ctx context.Context, key string,
	start, stop float64) ([]FloatPair, error) {
	conn, err := getRedis(s)
	if err != nil {
		return nil, err
	}

	v, err := conn.ZRevRangeByScoreWithScores(ctx, key, &red.ZRangeBy{
		Min: strconv.FormatFloat(start, 'f', -1, 64),
		Max: strconv.FormatFloat(stop, 'f', -1, 64),
	}).Result()
	if err != nil {
		return nil, err
	}

	return toFloatPairs(v), nil
}

// ZrevrangebyscoreWithScoresAndLimit is the implementation of redis zrevrangebyscore command
// with scores and limit.
func (s *Redis) ZrevrangebyscoreWithScoresAndLimit(key string, start, stop int64,
	page, size int) ([]Pair, error) {
	return s.ZrevrangebyscoreWithScoresAndLimitCtx(context.Background(),
		key, start, stop, page, size)
}

// ZrevrangebyscoreWithScoresAndLimitCtx is the implementation of redis zrevrangebyscore command
// with scores and limit.
func (s *Redis) ZrevrangebyscoreWithScoresAndLimitCtx(ctx context.Context, key string,
	start, stop int64, page, size int) ([]Pair, error) {
	if size <= 0 {
		return nil, nil
	}

	conn, err := getRedis(s)
	if err != nil {
		return nil, err
	}

	v, err := conn.ZRevRangeByScoreWithScores(ctx, key, &red.ZRangeBy{
		Min:    strconv.FormatInt(start, 10),
		Max:    strconv.FormatInt(stop, 10),
		Offset: int64(page * size),
		Count:  int64(size),
	}).Result()
	if err != nil {
		return nil, err
	}

	return toPairs(v), nil
}

// ZrevrangebyscoreWithScoresByFloatAndLimit is the implementation of redis zrevrangebyscore command
// with scores by float and limit.
func (s *Redis) ZrevrangebyscoreWithScoresByFloatAndLimit(key string, start, stop float64,
	page, size int) ([]FloatPair, error) {
	return s.ZrevrangebyscoreWithScoresByFloatAndLimitCtx(context.Background(),
		key, start, stop, page, size)
}

// ZrevrangebyscoreWithScoresByFloatAndLimitCtx is the implementation of redis zrevrangebyscore command
// with scores by float and limit.
func (s *Redis) ZrevrangebyscoreWithScoresByFloatAndLimitCtx(ctx context.Context, key string,
	start, stop float64, page, size int) ([]FloatPair, error) {
	if size <= 0 {
		return nil, nil
	}

	conn, err := getRedis(s)
	if err != nil {
		return nil, err
	}

	v, err := conn.ZRevRangeByScoreWithScores(ctx, key, &red.ZRangeBy{
		Min:    strconv.FormatFloat(start, 'f', -1, 64),
		Max:    strconv.FormatFloat(stop, 'f', -1, 64),
		Offset: int64(page * size),
		Count:  int64(size),
	}).Result()
	if err != nil {
		return nil, err
	}

	return toFloatPairs(v), nil
}

// Zrevrank is the implementation of redis zrevrank command.
func (s *Redis) Zrevrank(key, field string) (int64, error) {
	return s.ZrevrankCtx(context.Background(), key, field)
}

// ZrevrankCtx is the implementation of redis zrevrank command.
func (s *Redis) ZrevrankCtx(ctx context.Context, key, field string) (int64, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	return conn.ZRevRank(ctx, key, field).Result()
}

// Zunionstore is the implementation of redis zunionstore command.
func (s *Redis) Zunionstore(dest string, store *ZStore) (int64, error) {
	return s.ZunionstoreCtx(context.Background(), dest, store)
}

// ZunionstoreCtx is the implementation of redis zunionstore command.
func (s *Redis) ZunionstoreCtx(ctx context.Context, dest string, store *ZStore) (
	int64, error) {
	conn, err := getRedis(s)
	if err != nil {
		return 0, err
	}

	return conn.ZUnionStore(ctx, dest, store).Result()
}

func (s *Redis) checkConnection(pingTimeout time.Duration) error {
	conn, err := getRedis(s)
	if err != nil {
		return err
	}

	timeout := defaultPingTimeout
	if pingTimeout > 0 {
		timeout = pingTimeout
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return conn.Ping(ctx).Err()
}

// Cluster customizes the given Redis as a cluster.
func Cluster() Option {
	return func(r *Redis) {
		r.Type = ClusterType
	}
}

// SetSlowThreshold sets the slow threshold.
func SetSlowThreshold(threshold time.Duration) {
	slowThreshold.Set(threshold)
}

// WithHook customizes the given Redis with given durationHook.
func WithHook(hook Hook) Option {
	return func(r *Redis) {
		r.hooks = append(r.hooks, hook)
	}
}

// WithPass customizes the given Redis with given password.
func WithPass(pass string) Option {
	return func(r *Redis) {
		r.Pass = pass
	}
}

// WithTLS customizes the given Redis with TLS enabled.
func WithTLS() Option {
	return func(r *Redis) {
		r.tls = true
	}
}

// WithUser customizes the given Redis with given username.
func WithUser(user string) Option {
	return func(r *Redis) {
		r.User = user
	}
}

func acceptable(err error) bool {
	return err == nil || errorx.In(err, red.Nil, context.Canceled)
}

func getRedis(r *Redis) (RedisNode, error) {
	switch r.Type {
	case ClusterType:
		return getCluster(r)
	case NodeType:
		return getClient(r)
	default:
		return nil, fmt.Errorf("redis type '%s' is not supported", r.Type)
	}
}

func toPairs(vals []red.Z) []Pair {
	pairs := make([]Pair, len(vals))
	for i, val := range vals {
		switch member := val.Member.(type) {
		case string:
			pairs[i] = Pair{
				Key:   member,
				Score: int64(val.Score),
			}
		default:
			pairs[i] = Pair{
				Key:   mapping.Repr(val.Member),
				Score: int64(val.Score),
			}
		}
	}
	return pairs
}

func toFloatPairs(vals []red.Z) []FloatPair {
	pairs := make([]FloatPair, len(vals))

	for i, val := range vals {
		switch member := val.Member.(type) {
		case string:
			pairs[i] = FloatPair{
				Key:   member,
				Score: val.Score,
			}
		default:
			pairs[i] = FloatPair{
				Key:   mapping.Repr(val.Member),
				Score: val.Score,
			}
		}
	}

	return pairs
}

func toStrings(vals []any) []string {
	ret := make([]string, len(vals))

	for i, val := range vals {
		if val == nil {
			ret[i] = ""
			continue
		}

		switch val := val.(type) {
		case string:
			ret[i] = val
		default:
			ret[i] = mapping.Repr(val)
		}
	}

	return ret
}
