package kv

import (
	"errors"
	"log"

	"github.com/tal-tech/go-zero/core/errorx"
	"github.com/tal-tech/go-zero/core/hash"
	"github.com/tal-tech/go-zero/core/stores/internal"
	"github.com/tal-tech/go-zero/core/stores/redis"
)

var ErrNoRedisNode = errors.New("no redis node")

type (
	Store interface {
		Del(keys ...string) (int, error)
		Eval(script string, key string, args ...interface{}) (interface{}, error)
		Exists(key string) (bool, error)
		Expire(key string, seconds int) error
		Expireat(key string, expireTime int64) error
		Get(key string) (string, error)
		Hdel(key, field string) (bool, error)
		Hexists(key, field string) (bool, error)
		Hget(key, field string) (string, error)
		Hgetall(key string) (map[string]string, error)
		Hincrby(key, field string, increment int) (int, error)
		Hkeys(key string) ([]string, error)
		Hlen(key string) (int, error)
		Hmget(key string, fields ...string) ([]string, error)
		Hset(key, field, value string) error
		Hsetnx(key, field, value string) (bool, error)
		Hmset(key string, fieldsAndValues map[string]string) error
		Hvals(key string) ([]string, error)
		Incr(key string) (int64, error)
		Incrby(key string, increment int64) (int64, error)
		Llen(key string) (int, error)
		Lpop(key string) (string, error)
		Lpush(key string, values ...interface{}) (int, error)
		Lrange(key string, start int, stop int) ([]string, error)
		Lrem(key string, count int, value string) (int, error)
		Persist(key string) (bool, error)
		Pfadd(key string, values ...interface{}) (bool, error)
		Pfcount(key string) (int64, error)
		Rpush(key string, values ...interface{}) (int, error)
		Sadd(key string, values ...interface{}) (int, error)
		Scard(key string) (int64, error)
		Set(key string, value string) error
		Setex(key, value string, seconds int) error
		Setnx(key, value string) (bool, error)
		SetnxEx(key, value string, seconds int) (bool, error)
		Sismember(key string, value interface{}) (bool, error)
		Smembers(key string) ([]string, error)
		Spop(key string) (string, error)
		Srandmember(key string, count int) ([]string, error)
		Srem(key string, values ...interface{}) (int, error)
		Sscan(key string, cursor uint64, match string, count int64) (keys []string, cur uint64, err error)
		Ttl(key string) (int, error)
		Zadd(key string, score int64, value string) (bool, error)
		Zadds(key string, ps ...redis.Pair) (int64, error)
		Zcard(key string) (int, error)
		Zcount(key string, start, stop int64) (int, error)
		Zincrby(key string, increment int64, field string) (int64, error)
		Zrange(key string, start, stop int64) ([]string, error)
		ZrangeWithScores(key string, start, stop int64) ([]redis.Pair, error)
		ZrangebyscoreWithScores(key string, start, stop int64) ([]redis.Pair, error)
		ZrangebyscoreWithScoresAndLimit(key string, start, stop int64, page, size int) ([]redis.Pair, error)
		Zrank(key, field string) (int64, error)
		Zrem(key string, values ...interface{}) (int, error)
		Zremrangebyrank(key string, start, stop int64) (int, error)
		Zremrangebyscore(key string, start, stop int64) (int, error)
		Zrevrange(key string, start, stop int64) ([]string, error)
		ZrevrangebyscoreWithScores(key string, start, stop int64) ([]redis.Pair, error)
		ZrevrangebyscoreWithScoresAndLimit(key string, start, stop int64, page, size int) ([]redis.Pair, error)
		Zscore(key string, value string) (int64, error)
	}

	clusterStore struct {
		dispatcher *hash.ConsistentHash
	}
)

func NewStore(c KvConf) Store {
	if len(c) == 0 || internal.TotalWeights(c) <= 0 {
		log.Fatal("no cache nodes")
	}

	// even if only one node, we chose to use consistent hash,
	// because Store and redis.Redis has different methods.
	dispatcher := hash.NewConsistentHash()
	for _, node := range c {
		cn := node.NewRedis()
		dispatcher.AddWithWeight(cn, node.Weight)
	}

	return clusterStore{
		dispatcher: dispatcher,
	}
}

func (cs clusterStore) Del(keys ...string) (int, error) {
	var val int
	var be errorx.BatchError

	for _, key := range keys {
		node, e := cs.getRedis(key)
		if e != nil {
			be.Add(e)
			continue
		}

		if v, e := node.Del(key); e != nil {
			be.Add(e)
		} else {
			val += v
		}
	}

	return val, be.Err()
}

func (cs clusterStore) Eval(script string, key string, args ...interface{}) (interface{}, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return nil, err
	}

	return node.Eval(script, []string{key}, args...)
}

func (cs clusterStore) Exists(key string) (bool, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return false, err
	}

	return node.Exists(key)
}

func (cs clusterStore) Expire(key string, seconds int) error {
	node, err := cs.getRedis(key)
	if err != nil {
		return err
	}

	return node.Expire(key, seconds)
}

func (cs clusterStore) Expireat(key string, expireTime int64) error {
	node, err := cs.getRedis(key)
	if err != nil {
		return err
	}

	return node.Expireat(key, expireTime)
}

func (cs clusterStore) Get(key string) (string, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return "", err
	}

	return node.Get(key)
}

func (cs clusterStore) Hdel(key, field string) (bool, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return false, err
	}

	return node.Hdel(key, field)
}

func (cs clusterStore) Hexists(key, field string) (bool, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return false, err
	}

	return node.Hexists(key, field)
}

func (cs clusterStore) Hget(key, field string) (string, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return "", err
	}

	return node.Hget(key, field)
}

func (cs clusterStore) Hgetall(key string) (map[string]string, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return nil, err
	}

	return node.Hgetall(key)
}

func (cs clusterStore) Hincrby(key, field string, increment int) (int, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.Hincrby(key, field, increment)
}

func (cs clusterStore) Hkeys(key string) ([]string, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return nil, err
	}

	return node.Hkeys(key)
}

func (cs clusterStore) Hlen(key string) (int, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.Hlen(key)
}

func (cs clusterStore) Hmget(key string, fields ...string) ([]string, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return nil, err
	}

	return node.Hmget(key, fields...)
}

func (cs clusterStore) Hset(key, field, value string) error {
	node, err := cs.getRedis(key)
	if err != nil {
		return err
	}

	return node.Hset(key, field, value)
}

func (cs clusterStore) Hsetnx(key, field, value string) (bool, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return false, err
	}

	return node.Hsetnx(key, field, value)
}

func (cs clusterStore) Hmset(key string, fieldsAndValues map[string]string) error {
	node, err := cs.getRedis(key)
	if err != nil {
		return err
	}

	return node.Hmset(key, fieldsAndValues)
}

func (cs clusterStore) Hvals(key string) ([]string, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return nil, err
	}

	return node.Hvals(key)
}

func (cs clusterStore) Incr(key string) (int64, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.Incr(key)
}

func (cs clusterStore) Incrby(key string, increment int64) (int64, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.Incrby(key, increment)
}

func (cs clusterStore) Llen(key string) (int, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.Llen(key)
}

func (cs clusterStore) Lpop(key string) (string, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return "", err
	}

	return node.Lpop(key)
}

func (cs clusterStore) Lpush(key string, values ...interface{}) (int, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.Lpush(key, values...)
}

func (cs clusterStore) Lrange(key string, start int, stop int) ([]string, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return nil, err
	}

	return node.Lrange(key, start, stop)
}

func (cs clusterStore) Lrem(key string, count int, value string) (int, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.Lrem(key, count, value)
}

func (cs clusterStore) Persist(key string) (bool, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return false, err
	}

	return node.Persist(key)
}

func (cs clusterStore) Pfadd(key string, values ...interface{}) (bool, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return false, err
	}

	return node.Pfadd(key, values...)
}

func (cs clusterStore) Pfcount(key string) (int64, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.Pfcount(key)
}

func (cs clusterStore) Rpush(key string, values ...interface{}) (int, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.Rpush(key, values...)
}

func (cs clusterStore) Sadd(key string, values ...interface{}) (int, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.Sadd(key, values...)
}

func (cs clusterStore) Scard(key string) (int64, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.Scard(key)
}

func (cs clusterStore) Set(key string, value string) error {
	node, err := cs.getRedis(key)
	if err != nil {
		return err
	}

	return node.Set(key, value)
}

func (cs clusterStore) Setex(key, value string, seconds int) error {
	node, err := cs.getRedis(key)
	if err != nil {
		return err
	}

	return node.Setex(key, value, seconds)
}

func (cs clusterStore) Setnx(key, value string) (bool, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return false, err
	}

	return node.Setnx(key, value)
}

func (cs clusterStore) SetnxEx(key, value string, seconds int) (bool, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return false, err
	}

	return node.SetnxEx(key, value, seconds)
}

func (cs clusterStore) Sismember(key string, value interface{}) (bool, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return false, err
	}

	return node.Sismember(key, value)
}

func (cs clusterStore) Smembers(key string) ([]string, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return nil, err
	}

	return node.Smembers(key)
}

func (cs clusterStore) Spop(key string) (string, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return "", err
	}

	return node.Spop(key)
}

func (cs clusterStore) Srandmember(key string, count int) ([]string, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return nil, err
	}

	return node.Srandmember(key, count)
}

func (cs clusterStore) Srem(key string, values ...interface{}) (int, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.Srem(key, values...)
}

func (cs clusterStore) Sscan(key string, cursor uint64, match string, count int64) (
	keys []string, cur uint64, err error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return nil, 0, err
	}

	return node.Sscan(key, cursor, match, count)
}

func (cs clusterStore) Ttl(key string) (int, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.Ttl(key)
}

func (cs clusterStore) Zadd(key string, score int64, value string) (bool, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return false, err
	}

	return node.Zadd(key, score, value)
}

func (cs clusterStore) Zadds(key string, ps ...redis.Pair) (int64, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.Zadds(key, ps...)
}

func (cs clusterStore) Zcard(key string) (int, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.Zcard(key)
}

func (cs clusterStore) Zcount(key string, start, stop int64) (int, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.Zcount(key, start, stop)
}

func (cs clusterStore) Zincrby(key string, increment int64, field string) (int64, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.Zincrby(key, increment, field)
}

func (cs clusterStore) Zrank(key, field string) (int64, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.Zrank(key, field)
}

func (cs clusterStore) Zrange(key string, start, stop int64) ([]string, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return nil, err
	}

	return node.Zrange(key, start, stop)
}

func (cs clusterStore) ZrangeWithScores(key string, start, stop int64) ([]redis.Pair, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return nil, err
	}

	return node.ZrangeWithScores(key, start, stop)
}

func (cs clusterStore) ZrangebyscoreWithScores(key string, start, stop int64) ([]redis.Pair, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return nil, err
	}

	return node.ZrangebyscoreWithScores(key, start, stop)
}

func (cs clusterStore) ZrangebyscoreWithScoresAndLimit(key string, start, stop int64, page, size int) (
	[]redis.Pair, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return nil, err
	}

	return node.ZrangebyscoreWithScoresAndLimit(key, start, stop, page, size)
}

func (cs clusterStore) Zrem(key string, values ...interface{}) (int, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.Zrem(key, values...)
}

func (cs clusterStore) Zremrangebyrank(key string, start, stop int64) (int, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.Zremrangebyrank(key, start, stop)
}

func (cs clusterStore) Zremrangebyscore(key string, start, stop int64) (int, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.Zremrangebyscore(key, start, stop)
}

func (cs clusterStore) Zrevrange(key string, start, stop int64) ([]string, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return nil, err
	}

	return node.Zrevrange(key, start, stop)
}

func (cs clusterStore) ZrevrangebyscoreWithScores(key string, start, stop int64) ([]redis.Pair, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return nil, err
	}

	return node.ZrevrangebyscoreWithScores(key, start, stop)
}

func (cs clusterStore) ZrevrangebyscoreWithScoresAndLimit(key string, start, stop int64, page, size int) (
	[]redis.Pair, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return nil, err
	}

	return node.ZrevrangebyscoreWithScoresAndLimit(key, start, stop, page, size)
}

func (cs clusterStore) Zscore(key string, value string) (int64, error) {
	node, err := cs.getRedis(key)
	if err != nil {
		return 0, err
	}

	return node.Zscore(key, value)
}

func (cs clusterStore) getRedis(key string) (*redis.Redis, error) {
	if val, ok := cs.dispatcher.Get(key); !ok {
		return nil, ErrNoRedisNode
	} else {
		return val.(*redis.Redis), nil
	}
}
