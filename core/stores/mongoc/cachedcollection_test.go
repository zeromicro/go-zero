package mongoc

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stat"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/mongo"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/redis/redistest"
)

const dummyCount = 10

func init() {
	stat.SetReporter(nil)
}

func TestCollection_Count(t *testing.T) {
	resetStats()
	r, clean, err := redistest.CreateRedis()
	assert.Nil(t, err)
	defer clean()

	cach := cache.NewNode(r, singleFlight, stats, mgo.ErrNotFound)
	c := newCollection(dummyConn{}, cach)
	val, err := c.Count("any")
	assert.Nil(t, err)
	assert.Equal(t, dummyCount, val)

	var value string
	assert.Nil(t, r.Set("any", `"foo"`))
	assert.Nil(t, c.GetCache("any", &value))
	assert.Equal(t, "foo", value)
	assert.Nil(t, c.DelCache("any"))

	assert.Nil(t, c.SetCache("any", "bar"))
	assert.Nil(t, c.FindAllNoCache(&value, "any", func(query mongo.Query) mongo.Query {
		return query
	}))
	assert.Nil(t, c.FindOne(&value, "any", "foo"))
	assert.Equal(t, "bar", value)
	assert.Nil(t, c.DelCache("any"))
	c = newCollection(dummyConn{val: `"bar"`}, cach)
	assert.Nil(t, c.FindOne(&value, "any", "foo"))
	assert.Equal(t, "bar", value)
	assert.Nil(t, c.FindOneNoCache(&value, "foo"))
	assert.Equal(t, "bar", value)
	assert.Nil(t, c.FindOneId(&value, "anyone", "foo"))
	assert.Equal(t, "bar", value)
	assert.Nil(t, c.FindOneIdNoCache(&value, "foo"))
	assert.Equal(t, "bar", value)
	assert.Nil(t, c.Insert("foo"))
	assert.Nil(t, c.Pipe("foo"))
	assert.Nil(t, c.Remove("any"))
	assert.Nil(t, c.RemoveId("any"))
	_, err = c.RemoveAll("any")
	assert.Nil(t, err)
	assert.Nil(t, c.Update("foo", "bar"))
	assert.Nil(t, c.UpdateId("foo", "bar"))
	_, err = c.Upsert("foo", "bar")
	assert.Nil(t, err)

	c = newCollection(dummyConn{
		val:       `"bar"`,
		removeErr: errors.New("any"),
	}, cach)
	assert.NotNil(t, c.Remove("any"))
	_, err = c.RemoveAll("any", "bar")
	assert.NotNil(t, err)
	assert.NotNil(t, c.RemoveId("any"))

	c = newCollection(dummyConn{
		val:       `"bar"`,
		updateErr: errors.New("any"),
	}, cach)
	assert.NotNil(t, c.Update("foo", "bar"))
	assert.NotNil(t, c.UpdateId("foo", "bar"))
	_, err = c.Upsert("foo", "bar")
	assert.NotNil(t, err)
}

func TestStat(t *testing.T) {
	resetStats()
	r, clean, err := redistest.CreateRedis()
	assert.Nil(t, err)
	defer clean()

	cach := cache.NewNode(r, singleFlight, stats, mgo.ErrNotFound)
	c := newCollection(dummyConn{}, cach).(*cachedCollection)

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

	r := redis.New("localhost:59999")
	cach := cache.NewNode(r, singleFlight, stats, mgo.ErrNotFound)
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
	r, clean, err := redistest.CreateRedis()
	assert.Nil(t, err)
	defer clean()

	cach := cache.NewNode(r, singleFlight, stats, mgo.ErrNotFound)
	c := newCollection(dummyConn{}, cach).(*cachedCollection)

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
	r, clean, err := redistest.CreateRedis()
	assert.Nil(t, err)
	defer clean()

	cach := cache.NewNode(r, singleFlight, stats, mgo.ErrNotFound)
	c := newCollection(dummyConn{}, cach).(*cachedCollection)

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
	val       string
	removeErr error
	updateErr error
}

func (c dummyConn) Find(query interface{}) mongo.Query {
	return dummyQuery{val: c.val}
}

func (c dummyConn) FindId(id interface{}) mongo.Query {
	return dummyQuery{val: c.val}
}

func (c dummyConn) Insert(docs ...interface{}) error {
	return nil
}

func (c dummyConn) Remove(selector interface{}) error {
	return c.removeErr
}

func (dummyConn) Pipe(pipeline interface{}) mongo.Pipe {
	return nil
}

func (c dummyConn) RemoveAll(selector interface{}) (*mgo.ChangeInfo, error) {
	return nil, c.removeErr
}

func (c dummyConn) RemoveId(id interface{}) error {
	return c.removeErr
}

func (c dummyConn) Update(selector, update interface{}) error {
	return c.updateErr
}

func (c dummyConn) UpdateId(id, update interface{}) error {
	return c.updateErr
}

func (c dummyConn) Upsert(selector, update interface{}) (*mgo.ChangeInfo, error) {
	return nil, c.updateErr
}

type dummyQuery struct {
	val string
}

func (d dummyQuery) All(result interface{}) error {
	return nil
}

func (d dummyQuery) Apply(change mgo.Change, result interface{}) (*mgo.ChangeInfo, error) {
	return nil, nil
}

func (d dummyQuery) Count() (int, error) {
	return dummyCount, nil
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
	return json.Unmarshal([]byte(d.val), result)
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
