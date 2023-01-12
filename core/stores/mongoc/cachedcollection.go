package mongoc

import (
	"github.com/globalsign/mgo"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/mongo"
	"github.com/zeromicro/go-zero/core/syncx"
)

var (
	// ErrNotFound is an alias of mgo.ErrNotFound.
	ErrNotFound = mgo.ErrNotFound

	// can't use one SingleFlight per conn, because multiple conns may share the same cache key.
	singleFlight = syncx.NewSingleFlight()
	stats        = cache.NewStat("mongoc")
)

type (
	// QueryOption defines the method to customize a mongo query.
	QueryOption func(query mongo.Query) mongo.Query

	// CachedCollection interface represents a mongo collection with cache.
	CachedCollection interface {
		Count(query any) (int, error)
		DelCache(keys ...string) error
		FindAllNoCache(v, query any, opts ...QueryOption) error
		FindOne(v any, key string, query any) error
		FindOneNoCache(v, query any) error
		FindOneId(v any, key string, id any) error
		FindOneIdNoCache(v, id any) error
		GetCache(key string, v any) error
		Insert(docs ...any) error
		Pipe(pipeline any) mongo.Pipe
		Remove(selector any, keys ...string) error
		RemoveNoCache(selector any) error
		RemoveAll(selector any, keys ...string) (*mgo.ChangeInfo, error)
		RemoveAllNoCache(selector any) (*mgo.ChangeInfo, error)
		RemoveId(id any, keys ...string) error
		RemoveIdNoCache(id any) error
		SetCache(key string, v any) error
		Update(selector, update any, keys ...string) error
		UpdateNoCache(selector, update any) error
		UpdateId(id, update any, keys ...string) error
		UpdateIdNoCache(id, update any) error
		Upsert(selector, update any, keys ...string) (*mgo.ChangeInfo, error)
		UpsertNoCache(selector, update any) (*mgo.ChangeInfo, error)
	}

	cachedCollection struct {
		collection mongo.Collection
		cache      cache.Cache
	}
)

func newCollection(collection mongo.Collection, c cache.Cache) CachedCollection {
	return &cachedCollection{
		collection: collection,
		cache:      c,
	}
}

func (c *cachedCollection) Count(query any) (int, error) {
	return c.collection.Find(query).Count()
}

func (c *cachedCollection) DelCache(keys ...string) error {
	return c.cache.Del(keys...)
}

func (c *cachedCollection) FindAllNoCache(v, query any, opts ...QueryOption) error {
	q := c.collection.Find(query)
	for _, opt := range opts {
		q = opt(q)
	}
	return q.All(v)
}

func (c *cachedCollection) FindOne(v any, key string, query any) error {
	return c.cache.Take(v, key, func(v any) error {
		q := c.collection.Find(query)
		return q.One(v)
	})
}

func (c *cachedCollection) FindOneNoCache(v, query any) error {
	q := c.collection.Find(query)
	return q.One(v)
}

func (c *cachedCollection) FindOneId(v any, key string, id any) error {
	return c.cache.Take(v, key, func(v any) error {
		q := c.collection.FindId(id)
		return q.One(v)
	})
}

func (c *cachedCollection) FindOneIdNoCache(v, id any) error {
	q := c.collection.FindId(id)
	return q.One(v)
}

func (c *cachedCollection) GetCache(key string, v any) error {
	return c.cache.Get(key, v)
}

func (c *cachedCollection) Insert(docs ...any) error {
	return c.collection.Insert(docs...)
}

func (c *cachedCollection) Pipe(pipeline any) mongo.Pipe {
	return c.collection.Pipe(pipeline)
}

func (c *cachedCollection) Remove(selector any, keys ...string) error {
	if err := c.RemoveNoCache(selector); err != nil {
		return err
	}

	return c.DelCache(keys...)
}

func (c *cachedCollection) RemoveNoCache(selector any) error {
	return c.collection.Remove(selector)
}

func (c *cachedCollection) RemoveAll(selector any, keys ...string) (*mgo.ChangeInfo, error) {
	info, err := c.RemoveAllNoCache(selector)
	if err != nil {
		return nil, err
	}

	if err := c.DelCache(keys...); err != nil {
		return nil, err
	}

	return info, nil
}

func (c *cachedCollection) RemoveAllNoCache(selector any) (*mgo.ChangeInfo, error) {
	return c.collection.RemoveAll(selector)
}

func (c *cachedCollection) RemoveId(id any, keys ...string) error {
	if err := c.RemoveIdNoCache(id); err != nil {
		return err
	}

	return c.DelCache(keys...)
}

func (c *cachedCollection) RemoveIdNoCache(id any) error {
	return c.collection.RemoveId(id)
}

func (c *cachedCollection) SetCache(key string, v any) error {
	return c.cache.Set(key, v)
}

func (c *cachedCollection) Update(selector, update any, keys ...string) error {
	if err := c.UpdateNoCache(selector, update); err != nil {
		return err
	}

	return c.DelCache(keys...)
}

func (c *cachedCollection) UpdateNoCache(selector, update any) error {
	return c.collection.Update(selector, update)
}

func (c *cachedCollection) UpdateId(id, update any, keys ...string) error {
	if err := c.UpdateIdNoCache(id, update); err != nil {
		return err
	}

	return c.DelCache(keys...)
}

func (c *cachedCollection) UpdateIdNoCache(id, update any) error {
	return c.collection.UpdateId(id, update)
}

func (c *cachedCollection) Upsert(selector, update any, keys ...string) (*mgo.ChangeInfo, error) {
	info, err := c.UpsertNoCache(selector, update)
	if err != nil {
		return nil, err
	}

	if err := c.DelCache(keys...); err != nil {
		return nil, err
	}

	return info, nil
}

func (c *cachedCollection) UpsertNoCache(selector, update any) (*mgo.ChangeInfo, error) {
	return c.collection.Upsert(selector, update)
}
