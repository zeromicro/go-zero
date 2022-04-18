package mon

import (
	"context"
	"encoding/json"
	"time"

	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/timex"
	"go.mongodb.org/mongo-driver/mongo"
	mopt "go.mongodb.org/mongo-driver/mongo/options"
)

const defaultSlowThreshold = time.Millisecond * 500

// ErrNotFound is an alias of mongo.ErrNoDocuments
var ErrNotFound = mongo.ErrNoDocuments

type (
	// Collection defines a MongoDB collection.
	Collection interface {
		// Aggregate executes an aggregation pipeline.
		Aggregate(ctx context.Context, pipeline interface{}, opts ...*mopt.AggregateOptions) (
			*mongo.Cursor, error)
		// BulkWrite performs a bulk write operation.
		BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...*mopt.BulkWriteOptions) (
			*mongo.BulkWriteResult, error)
		// Clone creates a copy of this collection with the same settings.
		Clone(opts ...*mopt.CollectionOptions) (*mongo.Collection, error)
		// CountDocuments returns the number of documents in the collection that match the filter.
		CountDocuments(ctx context.Context, filter interface{}, opts ...*mopt.CountOptions) (int64, error)
		// Database returns the database that this collection is a part of.
		Database() *mongo.Database
		// DeleteMany deletes documents from the collection that match the filter.
		DeleteMany(ctx context.Context, filter interface{}, opts ...*mopt.DeleteOptions) (
			*mongo.DeleteResult, error)
		// DeleteOne deletes at most one document from the collection that matches the filter.
		DeleteOne(ctx context.Context, filter interface{}, opts ...*mopt.DeleteOptions) (
			*mongo.DeleteResult, error)
		// Distinct returns a list of distinct values for the given key across the collection.
		Distinct(ctx context.Context, fieldName string, filter interface{},
			opts ...*mopt.DistinctOptions) ([]interface{}, error)
		// Drop drops this collection from database.
		Drop(ctx context.Context) error
		// EstimatedDocumentCount returns an estimate of the count of documents in a collection
		// using collection metadata.
		EstimatedDocumentCount(ctx context.Context, opts ...*mopt.EstimatedDocumentCountOptions) (int64, error)
		// Find finds the documents matching the provided filter.
		Find(ctx context.Context, filter interface{}, opts ...*mopt.FindOptions) (*mongo.Cursor, error)
		// FindOne returns up to one document that matches the provided filter.
		FindOne(ctx context.Context, filter interface{}, opts ...*mopt.FindOneOptions) (
			*mongo.SingleResult, error)
		// FindOneAndDelete returns at most one document that matches the filter. If the filter
		// matches multiple documents, only the first document is deleted.
		FindOneAndDelete(ctx context.Context, filter interface{}, opts ...*mopt.FindOneAndDeleteOptions) (
			*mongo.SingleResult, error)
		// FindOneAndReplace returns at most one document that matches the filter. If the filter
		// matches multiple documents, FindOneAndReplace returns the first document in the
		// collection that matches the filter.
		FindOneAndReplace(ctx context.Context, filter interface{}, replacement interface{},
			opts ...*mopt.FindOneAndReplaceOptions) (*mongo.SingleResult, error)
		// FindOneAndUpdate returns at most one document that matches the filter. If the filter
		// matches multiple documents, FindOneAndUpdate returns the first document in the
		// collection that matches the filter.
		FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{},
			opts ...*mopt.FindOneAndUpdateOptions) (*mongo.SingleResult, error)
		// Indexes returns the index view for this collection.
		Indexes() mongo.IndexView
		// InsertMany inserts the provided documents.
		InsertMany(ctx context.Context, documents []interface{}, opts ...*mopt.InsertManyOptions) (
			*mongo.InsertManyResult, error)
		// InsertOne inserts the provided document.
		InsertOne(ctx context.Context, document interface{}, opts ...*mopt.InsertOneOptions) (
			*mongo.InsertOneResult, error)
		// ReplaceOne replaces at most one document that matches the filter.
		ReplaceOne(ctx context.Context, filter interface{}, replacement interface{},
			opts ...*mopt.ReplaceOptions) (*mongo.UpdateResult, error)
		// UpdateByID updates a single document matching the provided filter.
		UpdateByID(ctx context.Context, id interface{}, update interface{},
			opts ...*mopt.UpdateOptions) (*mongo.UpdateResult, error)
		// UpdateMany updates the provided documents.
		UpdateMany(ctx context.Context, filter interface{}, update interface{},
			opts ...*mopt.UpdateOptions) (*mongo.UpdateResult, error)
		// UpdateOne updates a single document matching the provided filter.
		UpdateOne(ctx context.Context, filter interface{}, update interface{},
			opts ...*mopt.UpdateOptions) (*mongo.UpdateResult, error)
		// Watch returns a change stream cursor used to receive notifications of changes to the collection.
		Watch(ctx context.Context, pipeline interface{}, opts ...*mopt.ChangeStreamOptions) (
			*mongo.ChangeStream, error)
	}

	decoratedCollection struct {
		*mongo.Collection
		name string
		brk  breaker.Breaker
	}

	keepablePromise struct {
		promise breaker.Promise
		log     func(error)
	}
)

func newCollection(collection *mongo.Collection, brk breaker.Breaker) Collection {
	return &decoratedCollection{
		Collection: collection,
		name:       collection.Name(),
		brk:        brk,
	}
}

func (c *decoratedCollection) Aggregate(ctx context.Context, pipeline interface{},
	opts ...*mopt.AggregateOptions) (cur *mongo.Cursor, err error) {
	err = c.brk.DoWithAcceptable(func() error {
		starTime := timex.Now()
		defer func() {
			c.logDurationSimple("Aggregate", starTime, err)
		}()

		cur, err = c.Collection.Aggregate(ctx, pipeline, opts...)
		return err
	}, acceptable)
	return
}

func (c *decoratedCollection) BulkWrite(ctx context.Context, models []mongo.WriteModel,
	opts ...*mopt.BulkWriteOptions) (res *mongo.BulkWriteResult, err error) {
	err = c.brk.DoWithAcceptable(func() error {
		startTime := timex.Now()
		defer func() {
			c.logDurationSimple("BulkWrite", startTime, err)
		}()

		res, err = c.Collection.BulkWrite(ctx, models, opts...)
		return err
	}, acceptable)
	return
}

func (c *decoratedCollection) CountDocuments(ctx context.Context, filter interface{},
	opts ...*mopt.CountOptions) (count int64, err error) {
	err = c.brk.DoWithAcceptable(func() error {
		startTime := timex.Now()
		defer func() {
			c.logDurationSimple("CountDocuments", startTime, err)
		}()

		count, err = c.Collection.CountDocuments(ctx, filter, opts...)
		return err
	}, acceptable)
	return
}

func (c *decoratedCollection) DeleteMany(ctx context.Context, filter interface{},
	opts ...*mopt.DeleteOptions) (res *mongo.DeleteResult, err error) {
	err = c.brk.DoWithAcceptable(func() error {
		startTime := timex.Now()
		defer func() {
			c.logDurationSimple("DeleteMany", startTime, err)
		}()

		res, err = c.Collection.DeleteMany(ctx, filter, opts...)
		return err
	}, acceptable)
	return
}

func (c *decoratedCollection) DeleteOne(ctx context.Context, filter interface{},
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

func (c *decoratedCollection) logDuration(method string, startTime time.Duration, err error,
	docs ...interface{}) {
	duration := timex.Since(startTime)
	content, e := json.Marshal(docs)
	if e != nil {
		logx.Error(err)
	} else if err != nil {
		if duration > slowThreshold.Load() {
			logx.WithDuration(duration).Slowf("[MONGO] mongo(%s) - slowcall - %s - fail(%s) - %s",
				c.name, method, err.Error(), string(content))
		} else {
			logx.WithDuration(duration).Infof("mongo(%s) - %s - fail(%s) - %s",
				c.name, method, err.Error(), string(content))
		}
	} else {
		if duration > slowThreshold.Load() {
			logx.WithDuration(duration).Slowf("[MONGO] mongo(%s) - slowcall - %s - ok - %s",
				c.name, method, string(content))
		} else {
			logx.WithDuration(duration).Infof("mongo(%s) - %s - ok - %s", c.name, method, string(content))
		}
	}
}

func (c *decoratedCollection) logDurationSimple(method string, startTime time.Duration, err error) {
	logDuration(c.name, method, startTime, err)
}

func (p keepablePromise) accept(err error) error {
	p.promise.Accept()
	p.log(err)
	return err
}

func (p keepablePromise) keep(err error) error {
	if acceptable(err) {
		p.promise.Accept()
	} else {
		p.promise.Reject(err.Error())
	}

	p.log(err)
	return err
}

func acceptable(err error) bool {
	return err == nil || err == mongo.ErrNoDocuments || err == mongo.ErrNilValue ||
		err == mongo.ErrNilDocument || err == mongo.ErrNilCursor || err == mongo.ErrEmptySlice
}
