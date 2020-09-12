package mongoc

import (
	"log"

	"github.com/globalsign/mgo"
	"github.com/tal-tech/go-zero/core/stores/cache"
	"github.com/tal-tech/go-zero/core/stores/internal"
	"github.com/tal-tech/go-zero/core/stores/mongo"
	"github.com/tal-tech/go-zero/core/stores/redis"
)

type Model struct {
	*mongo.Model
	cache              internal.Cache
	generateCollection func(*mgo.Session) *cachedCollection
}

func MustNewNodeModel(url, database, collection string, rds *redis.Redis, opts ...cache.Option) *Model {
	model, err := NewNodeModel(url, database, collection, rds, opts...)
	if err != nil {
		log.Fatal(err)
	}

	return model
}

func MustNewModel(url, database, collection string, c cache.CacheConf, opts ...cache.Option) *Model {
	model, err := NewModel(url, database, collection, c, opts...)
	if err != nil {
		log.Fatal(err)
	}

	return model
}

func NewNodeModel(url, database, collection string, rds *redis.Redis, opts ...cache.Option) (*Model, error) {
	c := internal.NewCacheNode(rds, sharedCalls, stats, mgo.ErrNotFound, opts...)
	return createModel(url, database, collection, c, func(collection mongo.Collection) *cachedCollection {
		return newCollection(collection, c)
	})
}

func NewModel(url, database, collection string, conf cache.CacheConf, opts ...cache.Option) (*Model, error) {
	c := internal.NewCache(conf, sharedCalls, stats, mgo.ErrNotFound, opts...)
	return createModel(url, database, collection, c, func(collection mongo.Collection) *cachedCollection {
		return newCollection(collection, c)
	})
}

func (mm *Model) Count(query interface{}) (int, error) {
	return mm.executeInt(func(c *cachedCollection) (int, error) {
		return c.Count(query)
	})
}

func (mm *Model) DelCache(keys ...string) error {
	return mm.cache.DelCache(keys...)
}

func (mm *Model) GetCache(key string, v interface{}) error {
	return mm.cache.GetCache(key, v)
}

func (mm *Model) GetCollection(session *mgo.Session) *cachedCollection {
	return mm.generateCollection(session)
}

func (mm *Model) FindAllNoCache(v interface{}, query interface{}, opts ...QueryOption) error {
	return mm.execute(func(c *cachedCollection) error {
		return c.FindAllNoCache(v, query, opts...)
	})
}

func (mm *Model) FindOne(v interface{}, key string, query interface{}) error {
	return mm.execute(func(c *cachedCollection) error {
		return c.FindOne(v, key, query)
	})
}

func (mm *Model) FindOneNoCache(v interface{}, query interface{}) error {
	return mm.execute(func(c *cachedCollection) error {
		return c.FindOneNoCache(v, query)
	})
}

func (mm *Model) FindOneId(v interface{}, key string, id interface{}) error {
	return mm.execute(func(c *cachedCollection) error {
		return c.FindOneId(v, key, id)
	})
}

func (mm *Model) FindOneIdNoCache(v interface{}, id interface{}) error {
	return mm.execute(func(c *cachedCollection) error {
		return c.FindOneIdNoCache(v, id)
	})
}

func (mm *Model) Insert(docs ...interface{}) error {
	return mm.execute(func(c *cachedCollection) error {
		return c.Insert(docs...)
	})
}

func (mm *Model) Pipe(pipeline interface{}) (mongo.Pipe, error) {
	return mm.pipe(func(c *cachedCollection) mongo.Pipe {
		return c.Pipe(pipeline)
	})
}

func (mm *Model) Remove(selector interface{}, keys ...string) error {
	return mm.execute(func(c *cachedCollection) error {
		return c.Remove(selector, keys...)
	})
}

func (mm *Model) RemoveNoCache(selector interface{}) error {
	return mm.execute(func(c *cachedCollection) error {
		return c.RemoveNoCache(selector)
	})
}

func (mm *Model) RemoveAll(selector interface{}, keys ...string) (*mgo.ChangeInfo, error) {
	return mm.change(func(c *cachedCollection) (*mgo.ChangeInfo, error) {
		return c.RemoveAll(selector, keys...)
	})
}

func (mm *Model) RemoveAllNoCache(selector interface{}) (*mgo.ChangeInfo, error) {
	return mm.change(func(c *cachedCollection) (*mgo.ChangeInfo, error) {
		return c.RemoveAllNoCache(selector)
	})
}

func (mm *Model) RemoveId(id interface{}, keys ...string) error {
	return mm.execute(func(c *cachedCollection) error {
		return c.RemoveId(id, keys...)
	})
}

func (mm *Model) RemoveIdNoCache(id interface{}) error {
	return mm.execute(func(c *cachedCollection) error {
		return c.RemoveIdNoCache(id)
	})
}

func (mm *Model) SetCache(key string, v interface{}) error {
	return mm.cache.SetCache(key, v)
}

func (mm *Model) Update(selector, update interface{}, keys ...string) error {
	return mm.execute(func(c *cachedCollection) error {
		return c.Update(selector, update, keys...)
	})
}

func (mm *Model) UpdateNoCache(selector, update interface{}) error {
	return mm.execute(func(c *cachedCollection) error {
		return c.UpdateNoCache(selector, update)
	})
}

func (mm *Model) UpdateId(id, update interface{}, keys ...string) error {
	return mm.execute(func(c *cachedCollection) error {
		return c.UpdateId(id, update, keys...)
	})
}

func (mm *Model) UpdateIdNoCache(id, update interface{}) error {
	return mm.execute(func(c *cachedCollection) error {
		return c.UpdateIdNoCache(id, update)
	})
}

func (mm *Model) Upsert(selector, update interface{}, keys ...string) (*mgo.ChangeInfo, error) {
	return mm.change(func(c *cachedCollection) (*mgo.ChangeInfo, error) {
		return c.Upsert(selector, update, keys...)
	})
}

func (mm *Model) UpsertNoCache(selector, update interface{}) (*mgo.ChangeInfo, error) {
	return mm.change(func(c *cachedCollection) (*mgo.ChangeInfo, error) {
		return c.UpsertNoCache(selector, update)
	})
}

func (mm *Model) change(fn func(c *cachedCollection) (*mgo.ChangeInfo, error)) (*mgo.ChangeInfo, error) {
	session, err := mm.TakeSession()
	if err != nil {
		return nil, err
	}
	defer mm.PutSession(session)

	return fn(mm.GetCollection(session))
}

func (mm *Model) execute(fn func(c *cachedCollection) error) error {
	session, err := mm.TakeSession()
	if err != nil {
		return err
	}
	defer mm.PutSession(session)

	return fn(mm.GetCollection(session))
}

func (mm *Model) executeInt(fn func(c *cachedCollection) (int, error)) (int, error) {
	session, err := mm.TakeSession()
	if err != nil {
		return 0, err
	}
	defer mm.PutSession(session)

	return fn(mm.GetCollection(session))
}

func (mm *Model) pipe(fn func(c *cachedCollection) mongo.Pipe) (mongo.Pipe, error) {
	session, err := mm.TakeSession()
	if err != nil {
		return nil, err
	}
	defer mm.PutSession(session)

	return fn(mm.GetCollection(session)), nil
}

func createModel(url, database, collection string, c internal.Cache,
	create func(mongo.Collection) *cachedCollection) (*Model, error) {
	model, err := mongo.NewModel(url, database, collection)
	if err != nil {
		return nil, err
	}

	return &Model{
		Model: model,
		cache: c,
		generateCollection: func(session *mgo.Session) *cachedCollection {
			collection := model.GetCollection(session)
			return create(collection)
		},
	}, nil
}
