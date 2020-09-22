package mongoc

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/stat"
	"github.com/tal-tech/go-zero/core/stores/internal"
	"github.com/tal-tech/go-zero/core/stores/mongo"
	"github.com/tal-tech/go-zero/core/stores/redis"
)

func init() {
	stat.SetReporter(nil)
}

func TestStat(t *testing.T) {
	resetStats()
	s, err := miniredis.Run()
	if err != nil {
		t.Error(err)
	}

	r := redis.NewRedis(s.Addr(), redis.NodeType)
	cach := internal.NewCacheNode(r, sharedCalls, stats, mgo.ErrNotFound)
	c := newCollection(dummyConn{}, cach)

	for i := 0; i < 10; i++ {
		var str string
		if err = c.cache.Take(&str, "name", func(v interface{}) error {
			*v.(*string) = "zero"
			return nil
		}); err != nil {
			t.Error(err)
		}
	}

	assert.Equal(t, uint64(10), atomic.LoadUint64(&stats.Total))
	assert.Equal(t, uint64(9), atomic.LoadUint64(&stats.Hit))
}

func TestStatCacheFails(t *testing.T) {
	resetStats()
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	r := redis.NewRedis("localhost:59999", redis.NodeType)
	cach := internal.NewCacheNode(r, sharedCalls, stats, mgo.ErrNotFound)
	c := newCollection(dummyConn{}, cach)

	for i := 0; i < 20; i++ {
		var str string
		err := c.FindOne(&str, "name", bson.M{})
		assert.NotNil(t, err)
	}

	assert.Equal(t, uint64(20), atomic.LoadUint64(&stats.Total))
	assert.Equal(t, uint64(0), atomic.LoadUint64(&stats.Hit))
	assert.Equal(t, uint64(20), atomic.LoadUint64(&stats.Miss))
	assert.Equal(t, uint64(0), atomic.LoadUint64(&stats.DbFails))
}

func TestStatDbFails(t *testing.T) {
	resetStats()
	s, err := miniredis.Run()
	if err != nil {
		t.Error(err)
	}

	r := redis.NewRedis(s.Addr(), redis.NodeType)
	cach := internal.NewCacheNode(r, sharedCalls, stats, mgo.ErrNotFound)
	c := newCollection(dummyConn{}, cach)

	for i := 0; i < 20; i++ {
		var str string
		err := c.cache.Take(&str, "name", func(v interface{}) error {
			return errors.New("db failed")
		})
		assert.NotNil(t, err)
	}

	assert.Equal(t, uint64(20), atomic.LoadUint64(&stats.Total))
	assert.Equal(t, uint64(0), atomic.LoadUint64(&stats.Hit))
	assert.Equal(t, uint64(20), atomic.LoadUint64(&stats.DbFails))
}

func TestStatFromMemory(t *testing.T) {
	resetStats()
	s, err := miniredis.Run()
	if err != nil {
		t.Error(err)
	}

	r := redis.NewRedis(s.Addr(), redis.NodeType)
	cach := internal.NewCacheNode(r, sharedCalls, stats, mgo.ErrNotFound)
	c := newCollection(dummyConn{}, cach)

	var all sync.WaitGroup
	var wait sync.WaitGroup
	all.Add(10)
	wait.Add(4)
	go func() {
		var str string
		if err := c.cache.Take(&str, "name", func(v interface{}) error {
			*v.(*string) = "zero"
			return nil
		}); err != nil {
			t.Error(err)
		}
		wait.Wait()
		runtime.Gosched()
		all.Done()
	}()

	for i := 0; i < 4; i++ {
		go func() {
			var str string
			wait.Done()
			if err := c.cache.Take(&str, "name", func(v interface{}) error {
				*v.(*string) = "zero"
				return nil
			}); err != nil {
				t.Error(err)
			}
			all.Done()
		}()
	}
	for i := 0; i < 5; i++ {
		go func() {
			var str string
			if err := c.cache.Take(&str, "name", func(v interface{}) error {
				*v.(*string) = "zero"
				return nil
			}); err != nil {
				t.Error(err)
			}
			all.Done()
		}()
	}
	all.Wait()

	assert.Equal(t, uint64(10), atomic.LoadUint64(&stats.Total))
	assert.Equal(t, uint64(9), atomic.LoadUint64(&stats.Hit))
}

func resetStats() {
	atomic.StoreUint64(&stats.Total, 0)
	atomic.StoreUint64(&stats.Hit, 0)
	atomic.StoreUint64(&stats.Miss, 0)
	atomic.StoreUint64(&stats.DbFails, 0)
}

type dummyConn struct {
}

func (c dummyConn) Find(query interface{}) mongo.Query {
	return dummyQuery{}
}

func (c dummyConn) FindId(id interface{}) mongo.Query {
	return dummyQuery{}
}

func (c dummyConn) Insert(docs ...interface{}) error {
	return nil
}

func (c dummyConn) Remove(selector interface{}) error {
	return nil
}

func (dummyConn) Pipe(pipeline interface{}) mongo.Pipe {
	return nil
}

func (c dummyConn) RemoveAll(selector interface{}) (*mgo.ChangeInfo, error) {
	return nil, nil
}

func (c dummyConn) RemoveId(id interface{}) error {
	return nil
}

func (c dummyConn) Update(selector, update interface{}) error {
	return nil
}

func (c dummyConn) UpdateId(id, update interface{}) error {
	return nil
}
func (c dummyConn) Upsert(selector, update interface{}) (*mgo.ChangeInfo, error) {
	return nil, nil
}

type dummyQuery struct {
}

func (d dummyQuery) All(result interface{}) error {
	return nil
}

func (d dummyQuery) Apply(change mgo.Change, result interface{}) (*mgo.ChangeInfo, error) {
	return nil, nil
}

func (d dummyQuery) Count() (int, error) {
	return 0, nil
}

func (d dummyQuery) Distinct(key string, result interface{}) error {
	return nil
}

func (d dummyQuery) Explain(result interface{}) error {
	return nil
}

func (d dummyQuery) For(result interface{}, f func() error) error {
	return nil
}

func (d dummyQuery) MapReduce(job *mgo.MapReduce, result interface{}) (*mgo.MapReduceInfo, error) {
	return nil, nil
}

func (d dummyQuery) One(result interface{}) error {
	return nil
}

func (d dummyQuery) Batch(n int) mongo.Query {
	return d
}

func (d dummyQuery) Collation(collation *mgo.Collation) mongo.Query {
	return d
}

func (d dummyQuery) Comment(comment string) mongo.Query {
	return d
}

func (d dummyQuery) Hint(indexKey ...string) mongo.Query {
	return d
}

func (d dummyQuery) Iter() mongo.Iter {
	return &mgo.Iter{}
}

func (d dummyQuery) Limit(n int) mongo.Query {
	return d
}

func (d dummyQuery) LogReplay() mongo.Query {
	return d
}

func (d dummyQuery) Prefetch(p float64) mongo.Query {
	return d
}

func (d dummyQuery) Select(selector interface{}) mongo.Query {
	return d
}

func (d dummyQuery) SetMaxScan(n int) mongo.Query {
	return d
}

func (d dummyQuery) SetMaxTime(duration time.Duration) mongo.Query {
	return d
}

func (d dummyQuery) Skip(n int) mongo.Query {
	return d
}

func (d dummyQuery) Snapshot() mongo.Query {
	return d
}

func (d dummyQuery) Sort(fields ...string) mongo.Query {
	return d
}

func (d dummyQuery) Tail(timeout time.Duration) mongo.Iter {
	return &mgo.Iter{}
}
