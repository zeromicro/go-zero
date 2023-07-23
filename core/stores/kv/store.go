package kv

import (
	"context"
	"errors"
	"log"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/hash"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// ErrNoRedisNode is an error that indicates no redis node.
var ErrNoRedisNode = errors.New("no redis node")

type (
	// Store interface represents a KV store.
	Store interface {
		Decr(key string) (int64, error)
		DecrCtx(ctx context.Context, key string) (int64, error)
		Decrby(key string, decrement int64) (int64, error)
		DecrbyCtx(ctx context.Context, key string, decrement int64) (int64, error)
		Del(keys ...string) (int, error)
		DelCtx(ctx context.Context, keys ...string) (int, error)
		Eval(script, key string, args ...any) (any, error)
		EvalCtx(ctx context.Context, script, key string, args ...any) (any, error)
		Exists(key string) (bool, error)
		ExistsCtx(ctx context.Context, key string) (bool, error)
		Expire(key string, seconds int) error
		ExpireCtx(ctx context.Context, key string, seconds int) error
		Expireat(key string, expireTime int64) error
		ExpireatCtx(ctx context.Context, key string, expireTime int64) error
		Get(key string) (string, error)
		GetCtx(ctx context.Context, key string) (string, error)
		GetSet(key, value string) (string, error)
		GetSetCtx(ctx context.Context, key, value string) (string, error)
		Hdel(key, field string) (bool, error)
		HdelCtx(ctx context.Context, key, field string) (bool, error)
		Hexists(key, field string) (bool, error)
		HexistsCtx(ctx context.Context, key, field string) (bool, error)
		Hget(key, field string) (string, error)
		HgetCtx(ctx context.Context, key, field string) (string, error)
		Hgetall(key string) (map[string]string, error)
		HgetallCtx(ctx context.Context, key string) (map[string]string, error)
		Hincrby(key, field string, increment int) (int, error)
		HincrbyCtx(ctx context.Context, key, field string, increment int) (int, error)
		Hkeys(key string) ([]string, error)
		HkeysCtx(ctx context.Context, key string) ([]string, error)
		Hlen(key string) (int, error)
		HlenCtx(ctx context.Context, key string) (int, error)
		Hmget(key string, fields ...string) ([]string, error)
		HmgetCtx(ctx context.Context, key string, fields ...string) ([]string, error)
		Hset(key, field, value string) error
		HsetCtx(ctx context.Context, key, field, value string) error
		Hsetnx(key, field, value string) (bool, error)
		HsetnxCtx(ctx context.Context, key, field, value string) (bool, error)
		Hmset(key string, fieldsAndValues map[string]string) error
		HmsetCtx(ctx context.Context, key string, fieldsAndValues map[string]string) error
		Hvals(key string) ([]string, error)
		HvalsCtx(ctx context.Context, key string) ([]string, error)
		Incr(key string) (int64, error)
		IncrCtx(ctx context.Context, key string) (int64, error)
		Incrby(key string, increment int64) (int64, error)
		IncrbyCtx(ctx context.Context, key string, increment int64) (int64, error)
		Lindex(key string, index int64) (string, error)
		LindexCtx(ctx context.Context, key string, index int64) (string, error)
		Llen(key string) (int, error)
		LlenCtx(ctx context.Context, key string) (int, error)
		Lpop(key string) (string, error)
		LpopCtx(ctx context.Context, key string) (string, error)
		Lpush(key string, values ...any) (int, error)
		LpushCtx(ctx context.Context, key string, values ...any) (int, error)
		Lrange(key string, start, stop int) ([]string, error)
		LrangeCtx(ctx context.Context, key string, start, stop int) ([]string, error)
		Lrem(key string, count int, value string) (int, error)
		LremCtx(ctx context.Context, key string, count int, value string) (int, error)
		Persist(key string) (bool, error)
		PersistCtx(ctx context.Context, key string) (bool, error)
		Pfadd(key string, values ...any) (bool, error)
		PfaddCtx(ctx context.Context, key string, values ...any) (bool, error)
		Pfcount(key string) (int64, error)
		PfcountCtx(ctx context.Context, key string) (int64, error)
		Rpush(key string, values ...any) (int, error)
		RpushCtx(ctx context.Context, key string, values ...any) (int, error)
		Sadd(key string, values ...any) (int, error)
		SaddCtx(ctx context.Context, key string, values ...any) (int, error)
		Scard(key string) (int64, error)
		ScardCtx(ctx context.Context, key string) (int64, error)
		Set(key, value string) error
		SetCtx(ctx context.Context, key, value string) error
		Setex(key, value string, seconds int) error
		SetexCtx(ctx context.Context, key, value string, seconds int) error
		Setnx(key, value string) (bool, error)
		SetnxCtx(ctx context.Context, key, value string) (bool, error)
		SetnxEx(key, value string, seconds int) (bool, error)
		SetnxExCtx(ctx context.Context, key, value string, seconds int) (bool, error)
		Sismember(key string, value any) (bool, error)
		SismemberCtx(ctx context.Context, key string, value any) (bool, error)
		Smembers(key string) ([]string, error)
		SmembersCtx(ctx context.Context, key string) ([]string, error)
		Spop(key string) (string, error)
		SpopCtx(ctx context.Context, key string) (string, error)
		Srandmember(key string, count int) ([]string, error)
		SrandmemberCtx(ctx context.Context, key string, count int) ([]string, error)
		Srem(key string, values ...any) (int, error)
		SremCtx(ctx context.Context, key string, values ...any) (int, error)
		Sscan(key string, cursor uint64, match string, count int64) (keys []string, cur uint64, err error)
		SscanCtx(ctx context.Context, key string, cursor uint64, match string, count int64) (keys []string, cur uint64, err error)
		Ttl(key string) (int, error)
		TtlCtx(ctx context.Context, key string) (int, error)
		Zadd(key string, score int64, value string) (bool, error)
		ZaddFloat(key string, score float64, value string) (bool, error)
		ZaddCtx(ctx context.Context, key string, score int64, value string) (bool, error)
		ZaddFloatCtx(ctx context.Context, key string, score float64, value string) (bool, error)
		Zadds(key string, ps ...redis.Pair) (int64, error)
		ZaddsCtx(ctx context.Context, key string, ps ...redis.Pair) (int64, error)
		Zcard(key string) (int, error)
		ZcardCtx(ctx context.Context, key string) (int, error)
		Zcount(key string, start, stop int64) (int, error)
		ZcountCtx(ctx context.Context, key string, start, stop int64) (int, error)
		Zincrby(key string, increment int64, field string) (int64, error)
		ZincrbyCtx(ctx context.Context, key string, increment int64, field string) (int64, error)
		Zrange(key string, start, stop int64) ([]string, error)
		ZrangeCtx(ctx context.Context, key string, start, stop int64) ([]string, error)
		ZrangeWithScores(key string, start, stop int64) ([]redis.Pair, error)
		ZrangeWithScoresCtx(ctx context.Context, key string, start, stop int64) ([]redis.Pair, error)
		ZrangebyscoreWithScores(key string, start, stop int64) ([]redis.Pair, error)
		ZrangebyscoreWithScoresCtx(ctx context.Context, key string, start, stop int64) ([]redis.Pair, error)
		ZrangebyscoreWithScoresAndLimit(key string, start, stop int64, page, size int) ([]redis.Pair, error)
		ZrangebyscoreWithScoresAndLimitCtx(ctx context.Context, key string, start, stop int64, page, size int) ([]redis.Pair, error)
		Zrank(key, field string) (int64, error)
		ZrankCtx(ctx context.Context, key, field string) (int64, error)
		Zrem(key string, values ...any) (int, error)
		ZremCtx(ctx context.Context, key string, values ...any) (int, error)
		Zremrangebyrank(key string, start, stop int64) (int, error)
		ZremrangebyrankCtx(ctx context.Context, key string, start, stop int64) (int, error)
		Zremrangebyscore(key string, start, stop int64) (int, error)
		ZremrangebyscoreCtx(ctx context.Context, key string, start, stop int64) (int, error)
		Zrevrange(key string, start, stop int64) ([]string, error)
		ZrevrangeCtx(ctx context.Context, key string, start, stop int64) ([]string, error)
		ZrevrangebyscoreWithScores(key string, start, stop int64) ([]redis.Pair, error)
		ZrevrangebyscoreWithScoresCtx(ctx context.Context, key string, start, stop int64) ([]redis.Pair, error)
		ZrevrangebyscoreWithScoresAndLimit(key string, start, stop int64, page, size int) ([]redis.Pair, error)
		ZrevrangebyscoreWithScoresAndLimitCtx(ctx context.Context, key string, start, stop int64, page, size int) ([]redis.Pair, error)
		Zscore(key, value string) (int64, error)
		ZscoreCtx(ctx context.Context, key, value string) (int64, error)
		Zrevrank(key, field string) (int64, error)
		ZrevrankCtx(ctx context.Context, key, field string) (int64, error)
	}

	clusterStore struct {
		dispatcher *hash.ConsistentHash
	}
)

// NewStore returns a Store.
func NewStore(c KvConf) Store {
	if len(c) == 0 || cache.TotalWeights(c) <= 0 {
		log.Fatal("no cache nodes")
	}

	// even if only one node, we chose to use consistent hash,
	// because Store and redis.Redis has different methods.
	dispatcher := hash.NewConsistentHash()
	for _, node := range c {
		cn := redis.MustNewRedis(node.RedisConf)
		dispatcher.AddWithWeight(cn, node.Weight)
	}

	return clusterStore{
		dispatcher: dispatcher,
	}
}

func (cs clusterStore) Decr(key string) (int64, error) {
	return cs.DecrCtx(context.Background(), key)
}

func (cs clusterStore) DecrCtx(ctx context.Context, key string) (int64, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.DecrCtx(ctx, key)
}

func (cs clusterStore) Decrby(key string, decrement int64) (int64, error) {
	return cs.DecrbyCtx(context.Background(), key, decrement)
}

func (cs clusterStore) DecrbyCtx(ctx context.Context, key string, decrement int64) (int64, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.DecrbyCtx(ctx, key, decrement)
}

func (cs clusterStore) Del(keys ...string) (int, error) {
	return cs.DelCtx(context.Background(), keys...)
}

func (cs clusterStore) DelCtx(ctx context.Context, keys ...string) (int, error) {
	var val int
	var be errorx.BatchError

	for _, key := range keys {
		node, e := cs.getRedis(key)
		if e != nil {
			be.Add(e)
			continue
		}

		if v, e := node.DelCtx(ctx, key); e != nil {
			be.Add(e)
		} else {
			val += v
		}
	}

	return val, be.Err()
}

func (cs clusterStore) Eval(script, key string, args ...any) (any, error) {
	return cs.EvalCtx(context.Background(), script, key, args...)
}

func (cs clusterStore) EvalCtx(ctx context.Context, script, key string, args ...any) (any, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return nil, err
	}

	return node.EvalCtx(ctx, script, []string{key}, args...)
}

func (cs clusterStore) Exists(key string) (bool, error) {
	return cs.ExistsCtx(context.Background(), key)
}

func (cs clusterStore) ExistsCtx(ctx context.Context, key string) (bool, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return false, err
	}

	return node.ExistsCtx(ctx, key)
}

func (cs clusterStore) Expire(key string, seconds int) error {
	return cs.ExpireCtx(context.Background(), key, seconds)
}

func (cs clusterStore) ExpireCtx(ctx context.Context, key string, seconds int) error {
	node, err := cs.getRedis(key)
	if err != nil {
		return err
	}

	return node.ExpireCtx(ctx, key, seconds)
}

func (cs clusterStore) Expireat(key string, expireTime int64) error {
	return cs.ExpireatCtx(context.Background(), key, expireTime)
}

func (cs clusterStore) ExpireatCtx(ctx context.Context, key string, expireTime int64) error {
	node, err := cs.getRedis(key)
	if err != nil {
		return err
	}

	return node.ExpireatCtx(ctx, key, expireTime)
}

func (cs clusterStore) Get(key string) (string, error) {
	return cs.GetCtx(context.Background(), key)
}

func (cs clusterStore) GetCtx(ctx context.Context, key string) (string, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return "", err
	}

	return node.GetCtx(ctx, key)
}

func (cs clusterStore) Hdel(key, field string) (bool, error) {
	return cs.HdelCtx(context.Background(), key, field)
}

func (cs clusterStore) HdelCtx(ctx context.Context, key, field string) (bool, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return false, err
	}

	return node.HdelCtx(ctx, key, field)
}

func (cs clusterStore) Hexists(key, field string) (bool, error) {
	return cs.HexistsCtx(context.Background(), key, field)
}

func (cs clusterStore) HexistsCtx(ctx context.Context, key, field string) (bool, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return false, err
	}

	return node.HexistsCtx(ctx, key, field)
}

func (cs clusterStore) Hget(key, field string) (string, error) {
	return cs.HgetCtx(context.Background(), key, field)
}

func (cs clusterStore) HgetCtx(ctx context.Context, key, field string) (string, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return "", err
	}

	return node.HgetCtx(ctx, key, field)
}

func (cs clusterStore) Hgetall(key string) (map[string]string, error) {
	return cs.HgetallCtx(context.Background(), key)
}

func (cs clusterStore) HgetallCtx(ctx context.Context, key string) (map[string]string, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return nil, err
	}

	return node.HgetallCtx(ctx, key)
}

func (cs clusterStore) Hincrby(key, field string, increment int) (int, error) {
	return cs.HincrbyCtx(context.Background(), key, field, increment)
}

func (cs clusterStore) HincrbyCtx(ctx context.Context, key, field string, increment int) (int, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.HincrbyCtx(ctx, key, field, increment)
}

func (cs clusterStore) Hkeys(key string) ([]string, error) {
	return cs.HkeysCtx(context.Background(), key)
}

func (cs clusterStore) HkeysCtx(ctx context.Context, key string) ([]string, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return nil, err
	}

	return node.HkeysCtx(ctx, key)
}

func (cs clusterStore) Hlen(key string) (int, error) {
	return cs.HlenCtx(context.Background(), key)
}

func (cs clusterStore) HlenCtx(ctx context.Context, key string) (int, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.HlenCtx(ctx, key)
}

func (cs clusterStore) Hmget(key string, fields ...string) ([]string, error) {
	return cs.HmgetCtx(context.Background(), key, fields...)
}

func (cs clusterStore) HmgetCtx(ctx context.Context, key string, fields ...string) ([]string, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return nil, err
	}

	return node.HmgetCtx(ctx, key, fields...)
}

func (cs clusterStore) Hset(key, field, value string) error {
	return cs.HsetCtx(context.Background(), key, field, value)
}

func (cs clusterStore) HsetCtx(ctx context.Context, key, field, value string) error {
	node, err := cs.getRedis(key)
	if err != nil {
		return err
	}

	return node.HsetCtx(ctx, key, field, value)
}

func (cs clusterStore) Hsetnx(key, field, value string) (bool, error) {
	return cs.HsetnxCtx(context.Background(), key, field, value)
}

func (cs clusterStore) HsetnxCtx(ctx context.Context, key, field, value string) (bool, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return false, err
	}

	return node.HsetnxCtx(ctx, key, field, value)
}

func (cs clusterStore) Hmset(key string, fieldsAndValues map[string]string) error {
	return cs.HmsetCtx(context.Background(), key, fieldsAndValues)
}

func (cs clusterStore) HmsetCtx(ctx context.Context, key string, fieldsAndValues map[string]string) error {
	node, err := cs.getRedis(key)
	if err != nil {
		return err
	}

	return node.HmsetCtx(ctx, key, fieldsAndValues)
}

func (cs clusterStore) Hvals(key string) ([]string, error) {
	return cs.HvalsCtx(context.Background(), key)
}

func (cs clusterStore) HvalsCtx(ctx context.Context, key string) ([]string, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return nil, err
	}

	return node.HvalsCtx(ctx, key)
}

func (cs clusterStore) Incr(key string) (int64, error) {
	return cs.IncrCtx(context.Background(), key)
}

func (cs clusterStore) IncrCtx(ctx context.Context, key string) (int64, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.IncrCtx(ctx, key)
}

func (cs clusterStore) Incrby(key string, increment int64) (int64, error) {
	return cs.IncrbyCtx(context.Background(), key, increment)
}

func (cs clusterStore) IncrbyCtx(ctx context.Context, key string, increment int64) (int64, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.IncrbyCtx(ctx, key, increment)
}

func (cs clusterStore) Llen(key string) (int, error) {
	return cs.LlenCtx(context.Background(), key)
}

func (cs clusterStore) LlenCtx(ctx context.Context, key string) (int, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.LlenCtx(ctx, key)
}

func (cs clusterStore) Lindex(key string, index int64) (string, error) {
	return cs.LindexCtx(context.Background(), key, index)
}

func (cs clusterStore) LindexCtx(ctx context.Context, key string, index int64) (string, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return "", err
	}

	return node.LindexCtx(ctx, key, index)
}

func (cs clusterStore) Lpop(key string) (string, error) {
	return cs.LpopCtx(context.Background(), key)
}

func (cs clusterStore) LpopCtx(ctx context.Context, key string) (string, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return "", err
	}

	return node.LpopCtx(ctx, key)
}

func (cs clusterStore) Lpush(key string, values ...any) (int, error) {
	return cs.LpushCtx(context.Background(), key, values...)
}

func (cs clusterStore) LpushCtx(ctx context.Context, key string, values ...any) (int, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.LpushCtx(ctx, key, values...)
}

func (cs clusterStore) Lrange(key string, start, stop int) ([]string, error) {
	return cs.LrangeCtx(context.Background(), key, start, stop)
}

func (cs clusterStore) LrangeCtx(ctx context.Context, key string, start, stop int) ([]string, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return nil, err
	}

	return node.LrangeCtx(ctx, key, start, stop)
}

func (cs clusterStore) Lrem(key string, count int, value string) (int, error) {
	return cs.LremCtx(context.Background(), key, count, value)
}

func (cs clusterStore) LremCtx(ctx context.Context, key string, count int, value string) (int, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.LremCtx(ctx, key, count, value)
}

func (cs clusterStore) Persist(key string) (bool, error) {
	return cs.PersistCtx(context.Background(), key)
}

func (cs clusterStore) PersistCtx(ctx context.Context, key string) (bool, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return false, err
	}

	return node.PersistCtx(ctx, key)
}

func (cs clusterStore) Pfadd(key string, values ...any) (bool, error) {
	return cs.PfaddCtx(context.Background(), key, values...)
}

func (cs clusterStore) PfaddCtx(ctx context.Context, key string, values ...any) (bool, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return false, err
	}

	return node.PfaddCtx(ctx, key, values...)
}

func (cs clusterStore) Pfcount(key string) (int64, error) {
	return cs.PfcountCtx(context.Background(), key)
}

func (cs clusterStore) PfcountCtx(ctx context.Context, key string) (int64, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.PfcountCtx(ctx, key)
}

func (cs clusterStore) Rpush(key string, values ...any) (int, error) {
	return cs.RpushCtx(context.Background(), key, values...)
}

func (cs clusterStore) RpushCtx(ctx context.Context, key string, values ...any) (int, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.RpushCtx(ctx, key, values...)
}

func (cs clusterStore) Sadd(key string, values ...any) (int, error) {
	return cs.SaddCtx(context.Background(), key, values...)
}

func (cs clusterStore) SaddCtx(ctx context.Context, key string, values ...any) (int, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.SaddCtx(ctx, key, values...)
}

func (cs clusterStore) Scard(key string) (int64, error) {
	return cs.ScardCtx(context.Background(), key)
}

func (cs clusterStore) ScardCtx(ctx context.Context, key string) (int64, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.ScardCtx(ctx, key)
}

func (cs clusterStore) Set(key, value string) error {
	return cs.SetCtx(context.Background(), key, value)
}

func (cs clusterStore) SetCtx(ctx context.Context, key, value string) error {
	node, err := cs.getRedis(key)
	if err != nil {
		return err
	}

	return node.SetCtx(ctx, key, value)
}

func (cs clusterStore) Setex(key, value string, seconds int) error {
	return cs.SetexCtx(context.Background(), key, value, seconds)
}

func (cs clusterStore) SetexCtx(ctx context.Context, key, value string, seconds int) error {
	node, err := cs.getRedis(key)
	if err != nil {
		return err
	}

	return node.SetexCtx(ctx, key, value, seconds)
}

func (cs clusterStore) Setnx(key, value string) (bool, error) {
	return cs.SetnxCtx(context.Background(), key, value)
}

func (cs clusterStore) SetnxCtx(ctx context.Context, key, value string) (bool, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return false, err
	}

	return node.SetnxCtx(ctx, key, value)
}

func (cs clusterStore) SetnxEx(key, value string, seconds int) (bool, error) {
	return cs.SetnxExCtx(context.Background(), key, value, seconds)
}

func (cs clusterStore) SetnxExCtx(ctx context.Context, key, value string, seconds int) (bool, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return false, err
	}

	return node.SetnxExCtx(ctx, key, value, seconds)
}

func (cs clusterStore) GetSet(key, value string) (string, error) {
	return cs.GetSetCtx(context.Background(), key, value)
}

func (cs clusterStore) GetSetCtx(ctx context.Context, key, value string) (string, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return "", err
	}

	return node.GetSetCtx(ctx, key, value)
}

func (cs clusterStore) Sismember(key string, value any) (bool, error) {
	return cs.SismemberCtx(context.Background(), key, value)
}

func (cs clusterStore) SismemberCtx(ctx context.Context, key string, value any) (bool, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return false, err
	}

	return node.SismemberCtx(ctx, key, value)
}

func (cs clusterStore) Smembers(key string) ([]string, error) {
	return cs.SmembersCtx(context.Background(), key)
}

func (cs clusterStore) SmembersCtx(ctx context.Context, key string) ([]string, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return nil, err
	}

	return node.SmembersCtx(ctx, key)
}

func (cs clusterStore) Spop(key string) (string, error) {
	return cs.SpopCtx(context.Background(), key)
}

func (cs clusterStore) SpopCtx(ctx context.Context, key string) (string, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return "", err
	}

	return node.SpopCtx(ctx, key)
}

func (cs clusterStore) Srandmember(key string, count int) ([]string, error) {
	return cs.SrandmemberCtx(context.Background(), key, count)
}

func (cs clusterStore) SrandmemberCtx(ctx context.Context, key string, count int) ([]string, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return nil, err
	}

	return node.SrandmemberCtx(ctx, key, count)
}

func (cs clusterStore) Srem(key string, values ...any) (int, error) {
	return cs.SremCtx(context.Background(), key, values...)
}

func (cs clusterStore) SremCtx(ctx context.Context, key string, values ...any) (int, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.SremCtx(ctx, key, values...)
}

func (cs clusterStore) Sscan(key string, cursor uint64, match string, count int64) (
	keys []string, cur uint64, err error) {
	return cs.SscanCtx(context.Background(), key, cursor, match, count)
}

func (cs clusterStore) SscanCtx(ctx context.Context, key string, cursor uint64, match string, count int64) (
	keys []string, cur uint64, err error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return nil, 0, err
	}

	return node.SscanCtx(ctx, key, cursor, match, count)
}

func (cs clusterStore) Ttl(key string) (int, error) {
	return cs.TtlCtx(context.Background(), key)
}

func (cs clusterStore) TtlCtx(ctx context.Context, key string) (int, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.TtlCtx(ctx, key)
}

func (cs clusterStore) Zadd(key string, score int64, value string) (bool, error) {
	return cs.ZaddCtx(context.Background(), key, score, value)
}

func (cs clusterStore) ZaddFloat(key string, score float64, value string) (bool, error) {
	return cs.ZaddFloatCtx(context.Background(), key, score, value)
}

func (cs clusterStore) ZaddCtx(ctx context.Context, key string, score int64, value string) (bool, error) {
	return cs.ZaddFloatCtx(ctx, key, float64(score), value)
}

func (cs clusterStore) ZaddFloatCtx(ctx context.Context, key string, score float64, value string) (bool, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return false, err
	}

	return node.ZaddFloatCtx(ctx, key, score, value)
}

func (cs clusterStore) Zadds(key string, ps ...redis.Pair) (int64, error) {
	return cs.ZaddsCtx(context.Background(), key, ps...)
}

func (cs clusterStore) ZaddsCtx(ctx context.Context, key string, ps ...redis.Pair) (int64, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.ZaddsCtx(ctx, key, ps...)
}

func (cs clusterStore) Zcard(key string) (int, error) {
	return cs.ZcardCtx(context.Background(), key)
}

func (cs clusterStore) ZcardCtx(ctx context.Context, key string) (int, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.ZcardCtx(ctx, key)
}

func (cs clusterStore) Zcount(key string, start, stop int64) (int, error) {
	return cs.ZcountCtx(context.Background(), key, start, stop)
}

func (cs clusterStore) ZcountCtx(ctx context.Context, key string, start, stop int64) (int, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.ZcountCtx(ctx, key, start, stop)
}

func (cs clusterStore) Zincrby(key string, increment int64, field string) (int64, error) {
	return cs.ZincrbyCtx(context.Background(), key, increment, field)
}

func (cs clusterStore) ZincrbyCtx(ctx context.Context, key string, increment int64, field string) (int64, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.ZincrbyCtx(ctx, key, increment, field)
}

func (cs clusterStore) Zrank(key, field string) (int64, error) {
	return cs.ZrankCtx(context.Background(), key, field)
}

func (cs clusterStore) ZrankCtx(ctx context.Context, key, field string) (int64, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.ZrankCtx(ctx, key, field)
}

func (cs clusterStore) Zrange(key string, start, stop int64) ([]string, error) {
	return cs.ZrangeCtx(context.Background(), key, start, stop)
}

func (cs clusterStore) ZrangeCtx(ctx context.Context, key string, start, stop int64) ([]string, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return nil, err
	}

	return node.ZrangeCtx(ctx, key, start, stop)
}

func (cs clusterStore) ZrangeWithScores(key string, start, stop int64) ([]redis.Pair, error) {
	return cs.ZrangeWithScoresCtx(context.Background(), key, start, stop)
}

func (cs clusterStore) ZrangeWithScoresCtx(ctx context.Context, key string, start, stop int64) ([]redis.Pair, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return nil, err
	}

	return node.ZrangeWithScoresCtx(ctx, key, start, stop)
}

func (cs clusterStore) ZrangebyscoreWithScores(key string, start, stop int64) ([]redis.Pair, error) {
	return cs.ZrangebyscoreWithScoresCtx(context.Background(), key, start, stop)
}

func (cs clusterStore) ZrangebyscoreWithScoresCtx(ctx context.Context, key string, start, stop int64) ([]redis.Pair, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return nil, err
	}

	return node.ZrangebyscoreWithScoresCtx(ctx, key, start, stop)
}

func (cs clusterStore) ZrangebyscoreWithScoresAndLimit(key string, start, stop int64, page, size int) (
	[]redis.Pair, error) {
	return cs.ZrangebyscoreWithScoresAndLimitCtx(context.Background(), key, start, stop, page, size)
}

func (cs clusterStore) ZrangebyscoreWithScoresAndLimitCtx(ctx context.Context, key string, start, stop int64, page, size int) (
	[]redis.Pair, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return nil, err
	}

	return node.ZrangebyscoreWithScoresAndLimitCtx(ctx, key, start, stop, page, size)
}

func (cs clusterStore) Zrem(key string, values ...any) (int, error) {
	return cs.ZremCtx(context.Background(), key, values...)
}

func (cs clusterStore) ZremCtx(ctx context.Context, key string, values ...any) (int, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.ZremCtx(ctx, key, values...)
}

func (cs clusterStore) Zremrangebyrank(key string, start, stop int64) (int, error) {
	return cs.ZremrangebyrankCtx(context.Background(), key, start, stop)
}

func (cs clusterStore) ZremrangebyrankCtx(ctx context.Context, key string, start, stop int64) (int, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.ZremrangebyrankCtx(ctx, key, start, stop)
}

func (cs clusterStore) Zremrangebyscore(key string, start, stop int64) (int, error) {
	return cs.ZremrangebyscoreCtx(context.Background(), key, start, stop)
}

func (cs clusterStore) ZremrangebyscoreCtx(ctx context.Context, key string, start, stop int64) (int, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.ZremrangebyscoreCtx(ctx, key, start, stop)
}

func (cs clusterStore) Zrevrange(key string, start, stop int64) ([]string, error) {
	return cs.ZrevrangeCtx(context.Background(), key, start, stop)
}

func (cs clusterStore) ZrevrangeCtx(ctx context.Context, key string, start, stop int64) ([]string, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return nil, err
	}

	return node.ZrevrangeCtx(ctx, key, start, stop)
}

func (cs clusterStore) ZrevrangebyscoreWithScores(key string, start, stop int64) ([]redis.Pair, error) {
	return cs.ZrevrangebyscoreWithScoresCtx(context.Background(), key, start, stop)
}

func (cs clusterStore) ZrevrangebyscoreWithScoresCtx(ctx context.Context, key string, start, stop int64) ([]redis.Pair, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return nil, err
	}

	return node.ZrevrangebyscoreWithScoresCtx(ctx, key, start, stop)
}

func (cs clusterStore) ZrevrangebyscoreWithScoresAndLimit(key string, start, stop int64, page, size int) (
	[]redis.Pair, error) {
	return cs.ZrevrangebyscoreWithScoresAndLimitCtx(context.Background(), key, start, stop, page, size)
}

func (cs clusterStore) ZrevrangebyscoreWithScoresAndLimitCtx(ctx context.Context, key string, start, stop int64, page, size int) (
	[]redis.Pair, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return nil, err
	}

	return node.ZrevrangebyscoreWithScoresAndLimitCtx(ctx, key, start, stop, page, size)
}

func (cs clusterStore) Zrevrank(key, field string) (int64, error) {
	return cs.ZrevrankCtx(context.Background(), key, field)
}

func (cs clusterStore) ZrevrankCtx(ctx context.Context, key, field string) (int64, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.ZrevrankCtx(ctx, key, field)
}

func (cs clusterStore) Zscore(key, value string) (int64, error) {
	return cs.ZscoreCtx(context.Background(), key, value)
}

func (cs clusterStore) ZscoreCtx(ctx context.Context, key, value string) (int64, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.ZscoreCtx(ctx, key, value)
}

func (cs clusterStore) getRedis(key string) (*redis.Redis, error) {
	val, ok := cs.dispatcher.Get(key)
	if !ok {
		return nil, ErrNoRedisNode
	}

	return val.(*redis.Redis), nil
}
