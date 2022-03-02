package mongoc

import (
	"log"

	"github.com/globalsign/mgo"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/mongo"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// A Model is a mongo model that built with cache capability.
type Model struct {
	*mongo.Model
	cache              cache.Cache
	generateCollection func(*mgo.Session) CachedCollection
}

// MustNewNodeModel returns a Model with a cache node, exists on errors.
func MustNewNodeModel(url, collection string, rds *redis.Redis, opts ...cache.Option) *Model {
	model, err := NewNodeModel(url, collection, rds, opts...)
	if err != nil {
		log.Fatal(err)
	}

	return model
}

// MustNewModel returns a Model with a cache cluster, exists on errors.
func MustNewModel(url, collection string, c cache.CacheConf, opts ...cache.Option) *Model {
	model, err := NewModel(url, collection, c, opts...)
	if err != nil {
		log.Fatal(err)
	}

	return model
}

// NewModel returns a Model with a cache cluster.
func NewModel(url, collection string, conf cache.CacheConf, opts ...cache.Option) (*Model, error) {
	c := cache.New(conf, singleFlight, stats, mgo.ErrNotFound, opts...)
	return NewModelWithCache(url, collection, c)
}

// NewModelWithCache returns a Model with a custom cache.
func NewModelWithCache(url, collection string, c cache.Cache) (*Model, error) {
	return createModel(url, collection, c, func(collection mongo.Collection) CachedCollection {
		return newCollection(collection, c)
	})
}

// NewNodeModel returns a Model with a cache node.
func NewNodeModel(url, collection string, rds *redis.Redis, opts ...cache.Option) (*Model, error) {
	c := cache.NewNode(rds, singleFlight, stats, mgo.ErrNotFound, opts...)
	return NewModelWithCache(url, collection, c)
}

// Count returns the count of given query.
func (mm *Model) Count(query interface{}) (int, error) {
	return mm.executeInt(func(c CachedCollection) (int, error) {
		return c.Count(query)
	})
}

// DelCache deletes the cache with given keys.
func (mm *Model) DelCache(keys ...string) error {
	return mm.cache.Del(keys...)
}

// GetCache unmarshal the cache into v with given key.
func (mm *Model) GetCache(key string, v interface{}) error {
	return mm.cache.Get(key, v)
}

// GetCollection returns a cache collection.
func (mm *Model) GetCollection(session *mgo.Session) CachedCollection {
	return mm.generateCollection(session)
}

// FindAllNoCache finds all records without cache.
func (mm *Model) FindAllNoCache(v, query interface{}, opts ...QueryOption) error {
	return mm.execute(func(c CachedCollection) error {
		return c.FindAllNoCache(v, query, opts...)
	})
}

// FindOne unmarshals a record into v with given key and query.
func (mm *Model) FindOne(v interface{}, key string, query interface{}) error {
	return mm.execute(func(c CachedCollection) error {
		return c.FindOne(v, key, query)
	})
}

// FindOneNoCache unmarshals a record into v with query, without cache.
func (mm *Model) FindOneNoCache(v, query interface{}) error {
	return mm.execute(func(c CachedCollection) error {
		return c.FindOneNoCache(v, query)
	})
}

// FindOneId unmarshals a record into v with query.
func (mm *Model) FindOneId(v interface{}, key string, id interface{}) error {
	return mm.execute(func(c CachedCollection) error {
		return c.FindOneId(v, key, id)
	})
}

// FindOneIdNoCache unmarshals a record into v with query, without cache.
func (mm *Model) FindOneIdNoCache(v, id interface{}) error {
	return mm.execute(func(c CachedCollection) error {
		return c.FindOneIdNoCache(v, id)
	})
}

// Insert inserts docs.
func (mm *Model) Insert(docs ...interface{}) error {
	return mm.execute(func(c CachedCollection) error {
		return c.Insert(docs...)
	})
}

// Pipe returns a mongo pipe with given pipeline.
func (mm *Model) Pipe(pipeline interface{}) (mongo.Pipe, error) {
	return mm.pipe(func(c CachedCollection) mongo.Pipe {
		return c.Pipe(pipeline)
	})
}

// Remove removes a record with given selector, and remove it from cache with given keys.
func (mm *Model) Remove(selector interface{}, keys ...string) error {
	return mm.execute(func(c CachedCollection) error {
		return c.Remove(selector, keys...)
	})
}

// RemoveNoCache removes a record with given selector.
func (mm *Model) RemoveNoCache(selector interface{}) error {
	return mm.execute(func(c CachedCollection) error {
		return c.RemoveNoCache(selector)
	})
}

// RemoveAll removes all records with given selector, and removes cache with given keys.
func (mm *Model) RemoveAll(selector interface{}, keys ...string) (*mgo.ChangeInfo, error) {
	return mm.change(func(c CachedCollection) (*mgo.ChangeInfo, error) {
		return c.RemoveAll(selector, keys...)
	})
}

// RemoveAllNoCache removes all records with given selector, and returns a mgo.ChangeInfo.
func (mm *Model) RemoveAllNoCache(selector interface{}) (*mgo.ChangeInfo, error) {
	return mm.change(func(c CachedCollection) (*mgo.ChangeInfo, error) {
		return c.RemoveAllNoCache(selector)
	})
}

// RemoveId removes a record with given id, and removes cache with given keys.
func (mm *Model) RemoveId(id interface{}, keys ...string) error {
	return mm.execute(func(c CachedCollection) error {
		return c.RemoveId(id, keys...)
	})
}

// RemoveIdNoCache removes a record with given id.
func (mm *Model) RemoveIdNoCache(id interface{}) error {
	return mm.execute(func(c CachedCollection) error {
		return c.RemoveIdNoCache(id)
	})
}

// SetCache sets the cache with given key and value.
func (mm *Model) SetCache(key string, v interface{}) error {
	return mm.cache.Set(key, v)
}

// Update updates the record with given selector, and delete cache with given keys.
func (mm *Model) Update(selector, update interface{}, keys ...string) error {
	return mm.execute(func(c CachedCollection) error {
		return c.Update(selector, update, keys...)
	})
}

// UpdateNoCache updates the record with given selector.
func (mm *Model) UpdateNoCache(selector, update interface{}) error {
	return mm.execute(func(c CachedCollection) error {
		return c.UpdateNoCache(selector, update)
	})
}

// UpdateId updates the record with given id, and delete cache with given keys.
func (mm *Model) UpdateId(id, update interface{}, keys ...string) error {
	return mm.execute(func(c CachedCollection) error {
		return c.UpdateId(id, update, keys...)
	})
}

// UpdateIdNoCache updates the record with given id.
func (mm *Model) UpdateIdNoCache(id, update interface{}) error {
	return mm.execute(func(c CachedCollection) error {
		return c.UpdateIdNoCache(id, update)
	})
}

// Upsert upserts a record with given selector, and delete cache with given keys.
func (mm *Model) Upsert(selector, update interface{}, keys ...string) (*mgo.ChangeInfo, error) {
	return mm.change(func(c CachedCollection) (*mgo.ChangeInfo, error) {
		return c.Upsert(selector, update, keys...)
	})
}

// UpsertNoCache upserts a record with given selector.
func (mm *Model) UpsertNoCache(selector, update interface{}) (*mgo.ChangeInfo, error) {
	return mm.change(func(c CachedCollection) (*mgo.ChangeInfo, error) {
		return c.UpsertNoCache(selector, update)
	})
}

func (mm *Model) change(fn func(c CachedCollection) (*mgo.ChangeInfo, error)) (*mgo.ChangeInfo, error) {
	session, err := mm.TakeSession()
	if err != nil {
		return nil, err
	}
	defer mm.PutSession(session)

	return fn(mm.GetCollection(session))
}

func (mm *Model) execute(fn func(c CachedCollection) error) error {
	session, err := mm.TakeSession()
	if err != nil {
		return err
	}
	defer mm.PutSession(session)

	return fn(mm.GetCollection(session))
}

func (mm *Model) executeInt(fn func(c CachedCollection) (int, error)) (int, error) {
	session, err := mm.TakeSession()
	if err != nil {
		return 0, err
	}
	defer mm.PutSession(session)

	return fn(mm.GetCollection(session))
}

func (mm *Model) pipe(fn func(c CachedCollection) mongo.Pipe) (mongo.Pipe, error) {
	session, err := mm.TakeSession()
	if err != nil {
		return nil, err
	}
	defer mm.PutSession(session)

	return fn(mm.GetCollection(session)), nil
}

func createModel(url, collection string, c cache.Cache,
	create func(mongo.Collection) CachedCollection) (*Model, error) {
	model, err := mongo.NewModel(url, collection)
	if err != nil {
		return nil, err
	}

	return &Model{
		Model: model,
		cache: c,
		generateCollection: func(session *mgo.Session) CachedCollection {
			collection := model.GetCollection(session)
			return create(collection)
		},
	}, nil
}
