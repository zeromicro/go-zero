package monc

import (
	"context"

	"github.com/globalsign/mgo"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/mon"
	omg "github.com/zeromicro/go-zero/core/stores/mongo"
	"github.com/zeromicro/go-zero/core/syncx"
	"github.com/zeromicro/go-zero/core/timex"
	"go.mongodb.org/mongo-driver/mongo"
	mopt "go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// ErrNotFound is an alias of mgo.ErrNotFound.
	ErrNotFound = mongo.ErrNoDocuments

	// can't use one SingleFlight per conn, because multiple conns may share the same cache key.
	singleFlight = syncx.NewSingleFlight()
	stats        = cache.NewStat("monc")
)

// CachedCollection represents a mongo collection with cache.
type CachedCollection struct {
	mon.Collection
	cache cache.Cache
}

func newCollection(collection mon.Collection, c cache.Cache) *CachedCollection {
	return &CachedCollection{
		Collection: collection,
		cache:      c,
	}
}

func (c *CachedCollection) DelCache(keys ...string) error {
	return c.cache.Del(keys...)
}

func (c *CachedCollection) FindOne(v interface{}, key string, query interface{}) error {
	return c.cache.Take(v, key, func(v interface{}) error {
		q := c.Collection.FindOne(query)
		return q.One(v)
	})
}

func (c *CachedCollection) FindOneNoCache(v, query interface{}) error {
	q := c.collection.Find(query)
	return q.One(v)
}

func (c *CachedCollection) FindOneId(v interface{}, key string, id interface{}) error {
	return c.cache.Take(v, key, func(v interface{}) error {
		q := c.collection.FindId(id)
		return q.One(v)
	})
}

func (c *CachedCollection) FindOneIdNoCache(v, id interface{}) error {
	q := c.collection.FindId(id)
	return q.One(v)
}

func (c *CachedCollection) GetCache(key string, v interface{}) error {
	return c.cache.Get(key, v)
}

func (c *CachedCollection) Insert(docs ...interface{}) error {
	return c.collection.Insert(docs...)
}

func (c *CachedCollection) Pipe(pipeline interface{}) omg.Pipe {
	return c.collection.Pipe(pipeline)
}

func (c *CachedCollection) Remove(selector interface{}, keys ...string) error {
	if err := c.RemoveNoCache(selector); err != nil {
		return err
	}

	return c.DelCache(keys...)
}

func (c *CachedCollection) RemoveNoCache(selector interface{}) error {
	return c.collection.Remove(selector)
}

func (c *CachedCollection) RemoveAll(selector interface{}, keys ...string) (*mgo.ChangeInfo, error) {
	info, err := c.RemoveAllNoCache(selector)
	if err != nil {
		return nil, err
	}

	if err := c.DelCache(keys...); err != nil {
		return nil, err
	}

	return info, nil
}

func (c *CachedCollection) RemoveAllNoCache(selector interface{}) (*mgo.ChangeInfo, error) {
	return c.collection.RemoveAll(selector)
}

func (c *CachedCollection) RemoveId(id interface{}, keys ...string) error {
	if err := c.RemoveIdNoCache(id); err != nil {
		return err
	}

	return c.DelCache(keys...)
}

func (c *CachedCollection) RemoveIdNoCache(id interface{}) error {
	return c.collection.RemoveId(id)
}

func (c *CachedCollection) SetCache(key string, v interface{}) error {
	return c.cache.Set(key, v)
}

func (c *CachedCollection) Update(selector, update interface{}, keys ...string) error {
	if err := c.UpdateNoCache(selector, update); err != nil {
		return err
	}

	return c.DelCache(keys...)
}

func (c *CachedCollection) UpdateNoCache(selector, update interface{}) error {
	return c.collection.Update(selector, update)
}

func (c *CachedCollection) UpdateId(id, update interface{}, keys ...string) error {
	if err := c.UpdateIdNoCache(id, update); err != nil {
		return err
	}

	return c.DelCache(keys...)
}

func (c *CachedCollection) UpdateIdNoCache(id, update interface{}) error {
	return c.collection.UpdateId(id, update)
}

func (c *CachedCollection) Upsert(selector, update interface{}, keys ...string) (*mgo.ChangeInfo, error) {
	info, err := c.UpsertNoCache(selector, update)
	if err != nil {
		return nil, err
	}

	if err := c.DelCache(keys...); err != nil {
		return nil, err
	}

	return info, nil
}

func (c *CachedCollection) UpsertNoCache(selector, update interface{}) (*mgo.ChangeInfo, error) {
	return c.collection.Upsert(selector, update)
}

func (c *CachedCollection) Aggregate(ctx context.Context, v interface{}, pipeline interface{},
	opts ...*mopt.AggregateOptions) error {
	cur, err := c.Collection.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return err
	}

	return cur.All(ctx, v)
}

func (c *CachedCollection) DeleteMany(ctx context.Context, filter interface{},
	opts ...*mopt.DeleteOptions) (int64, error) {
	res, err := c.Collection.DeleteMany(ctx, filter, opts...)
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}

func (c *CachedCollection) DeleteOne(ctx context.Context, filter interface{},
	opts ...*mopt.DeleteOptions) (res *mongo.DeleteResult, err error) {
	err = c.brk.DoWithAcceptable(func() error {
		startTime := timex.Now()
		defer func() {
			c.logDuration("DeleteOne", startTime, err, filter)
		}()

		res, err = c.Collection.DeleteOne(ctx, filter, opts...)
		return err
	}, acceptable)
	return
}

func (c *decoratedCollection) Distinct(ctx context.Context, fieldName string, filter interface{},
	opts ...*mopt.DistinctOptions) (val []interface{}, err error) {
	err = c.brk.DoWithAcceptable(func() error {
		startTime := timex.Now()
		defer func() {
			c.logDurationSimple("Distinct", startTime, err)
		}()

		val, err = c.Collection.Distinct(ctx, fieldName, filter, opts...)
		return err
	}, acceptable)
	return
}

func (c *decoratedCollection) EstimatedDocumentCount(ctx context.Context,
	opts ...*mopt.EstimatedDocumentCountOptions) (val int64, err error) {
	err = c.brk.DoWithAcceptable(func() error {
		startTime := timex.Now()
		defer func() {
			c.logDurationSimple("EstimatedDocumentCount", startTime, err)
		}()

		val, err = c.Collection.EstimatedDocumentCount(ctx, opts...)
		return err
	}, acceptable)
	return
}

func (c *decoratedCollection) Find(ctx context.Context, filter interface{},
	opts ...*mopt.FindOptions) (cur *mongo.Cursor, err error) {
	err = c.brk.DoWithAcceptable(func() error {
		startTime := timex.Now()
		defer func() {
			c.logDuration("Find", startTime, err, filter)
		}()

		cur, err = c.Collection.Find(ctx, filter, opts...)
		return err
	}, acceptable)
	return
}

func (c *decoratedCollection) FindOne(ctx context.Context, filter interface{},
	opts ...*mopt.FindOneOptions) (res *mongo.SingleResult, err error) {
	err = c.brk.DoWithAcceptable(func() error {
		startTime := timex.Now()
		defer func() {
			c.logDuration("FindOne", startTime, err, filter)
		}()

		res = c.Collection.FindOne(ctx, filter, opts...)
		err = res.Err()
		return err
	}, acceptable)
	return
}

func (c *decoratedCollection) FindOneAndDelete(ctx context.Context, filter interface{},
	opts ...*mopt.FindOneAndDeleteOptions) (res *mongo.SingleResult, err error) {
	err = c.brk.DoWithAcceptable(func() error {
		startTime := timex.Now()
		defer func() {
			c.logDuration("FindOneAndDelete", startTime, err, filter)
		}()

		res = c.Collection.FindOneAndDelete(ctx, filter, opts...)
		err = res.Err()
		return err
	}, acceptable)
	return
}

func (c *decoratedCollection) FindOneAndReplace(ctx context.Context, filter interface{},
	replacement interface{}, opts ...*mopt.FindOneAndReplaceOptions) (
	res *mongo.SingleResult, err error) {
	err = c.brk.DoWithAcceptable(func() error {
		startTime := timex.Now()
		defer func() {
			c.logDuration("FindOneAndReplace", startTime, err, filter, replacement)
		}()

		res = c.Collection.FindOneAndReplace(ctx, filter, replacement, opts...)
		err = res.Err()
		return err
	}, acceptable)
	return
}

func (c *decoratedCollection) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{},
	opts ...*mopt.FindOneAndUpdateOptions) (res *mongo.SingleResult, err error) {
	err = c.brk.DoWithAcceptable(func() error {
		startTime := timex.Now()
		defer func() {
			c.logDuration("FindOneAndUpdate", startTime, err, filter, update)
		}()

		res = c.Collection.FindOneAndUpdate(ctx, filter, update, opts...)
		err = res.Err()
		return err
	}, acceptable)
	return
}

func (c *decoratedCollection) InsertMany(ctx context.Context, documents []interface{},
	opts ...*mopt.InsertManyOptions) (res *mongo.InsertManyResult, err error) {
	err = c.brk.DoWithAcceptable(func() error {
		startTime := timex.Now()
		defer func() {
			c.logDurationSimple("InsertMany", startTime, err)
		}()

		res, err = c.Collection.InsertMany(ctx, documents, opts...)
		return err
	}, acceptable)
	return
}

func (c *decoratedCollection) InsertOne(ctx context.Context, document interface{},
	opts ...*mopt.InsertOneOptions) (res *mongo.InsertOneResult, err error) {
	err = c.brk.DoWithAcceptable(func() error {
		startTime := timex.Now()
		defer func() {
			c.logDuration("InsertOne", startTime, err, document)
		}()

		res, err = c.Collection.InsertOne(ctx, document, opts...)
		return err
	}, acceptable)
	return
}

func (c *decoratedCollection) ReplaceOne(ctx context.Context, filter interface{}, replacement interface{},
	opts ...*mopt.ReplaceOptions) (res *mongo.UpdateResult, err error) {
	err = c.brk.DoWithAcceptable(func() error {
		startTime := timex.Now()
		defer func() {
			c.logDuration("ReplaceOne", startTime, err, filter, replacement)
		}()

		res, err = c.Collection.ReplaceOne(ctx, filter, replacement, opts...)
		return err
	}, acceptable)
	return
}

func (c *decoratedCollection) UpdateByID(ctx context.Context, id interface{}, update interface{},
	opts ...*mopt.UpdateOptions) (res *mongo.UpdateResult, err error) {
	err = c.brk.DoWithAcceptable(func() error {
		startTime := timex.Now()
		defer func() {
			c.logDuration("UpdateByID", startTime, err, id, update)
		}()

		res, err = c.Collection.UpdateByID(ctx, id, update, opts...)
		return err
	}, acceptable)
	return
}

func (c *decoratedCollection) UpdateMany(ctx context.Context, filter interface{}, update interface{},
	opts ...*mopt.UpdateOptions) (res *mongo.UpdateResult, err error) {
	err = c.brk.DoWithAcceptable(func() error {
		startTime := timex.Now()
		defer func() {
			c.logDurationSimple("UpdateMany", startTime, err)
		}()

		res, err = c.Collection.UpdateMany(ctx, filter, update, opts...)
		return err
	}, acceptable)
	return
}

func (c *decoratedCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{},
	opts ...*mopt.UpdateOptions) (res *mongo.UpdateResult, err error) {
	err = c.brk.DoWithAcceptable(func() error {
		startTime := timex.Now()
		defer func() {
			c.logDuration("UpdateOne", startTime, err, filter, update)
		}()

		res, err = c.Collection.UpdateOne(ctx, filter, update, opts...)
		return err
	}, acceptable)
	return
}
