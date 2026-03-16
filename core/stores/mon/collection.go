//go:generate mockgen -package mon -destination collection_mock.go -source collection.go Collection,monCollection
package mon

import (
	"context"
	"errors"
	"time"

	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/timex"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/x/mongo/driver/session"
)

const (
	defaultSlowThreshold = time.Millisecond * 500
	// spanName is the span name of the mongo calls.
	spanName         = "mongo"
	duplicateKeyCode = 11000

	// mongodb method names
	aggregate              = "Aggregate"
	bulkWrite              = "BulkWrite"
	countDocuments         = "CountDocuments"
	deleteMany             = "DeleteMany"
	deleteOne              = "DeleteOne"
	distinct               = "Distinct"
	estimatedDocumentCount = "EstimatedDocumentCount"
	find                   = "Find"
	findOne                = "FindOne"
	findOneAndDelete       = "FindOneAndDelete"
	findOneAndReplace      = "FindOneAndReplace"
	findOneAndUpdate       = "FindOneAndUpdate"
	insertMany             = "InsertMany"
	insertOne              = "InsertOne"
	replaceOne             = "ReplaceOne"
	updateByID             = "UpdateByID"
	updateMany             = "UpdateMany"
	updateOne              = "UpdateOne"
)

// ErrNotFound is an alias of mongo.ErrNoDocuments
var ErrNotFound = mongo.ErrNoDocuments

type (
	// Collection defines a MongoDB collection.
	Collection interface {
		// Aggregate executes an aggregation pipeline.
		Aggregate(ctx context.Context, pipeline any, opts ...options.Lister[options.AggregateOptions]) (
			*mongo.Cursor, error)
		// BulkWrite performs a bulk write operation.
		BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...options.Lister[options.BulkWriteOptions]) (
			*mongo.BulkWriteResult, error)
		// Clone creates a copy of this collection with the same settings.
		Clone(opts ...options.Lister[options.CollectionOptions]) *mongo.Collection
		// CountDocuments returns the number of documents in the collection that match the filter.
		CountDocuments(ctx context.Context, filter any, opts ...options.Lister[options.CountOptions]) (int64, error)
		// Database returns the database that this collection is a part of.
		Database() *mongo.Database
		// DeleteMany deletes documents from the collection that match the filter.
		DeleteMany(ctx context.Context, filter any, opts ...options.Lister[options.DeleteManyOptions]) (
			*mongo.DeleteResult, error)
		// DeleteOne deletes at most one document from the collection that matches the filter.
		DeleteOne(ctx context.Context, filter any, opts ...options.Lister[options.DeleteOneOptions]) (
			*mongo.DeleteResult, error)
		// Distinct returns a list of distinct values for the given key across the collection.
		Distinct(ctx context.Context, fieldName string, filter any,
			opts ...options.Lister[options.DistinctOptions]) (*mongo.DistinctResult, error)
		// Drop drops this collection from database.
		Drop(ctx context.Context, opts ...options.Lister[options.DropCollectionOptions]) error
		// EstimatedDocumentCount returns an estimate of the count of documents in a collection
		// using collection metadata.
		EstimatedDocumentCount(ctx context.Context, opts ...options.Lister[options.EstimatedDocumentCountOptions]) (int64, error)
		// Find finds the documents matching the provided filter.
		Find(ctx context.Context, filter any, opts ...options.Lister[options.FindOptions]) (*mongo.Cursor, error)
		// FindOne returns up to one document that matches the provided filter.
		FindOne(ctx context.Context, filter any, opts ...options.Lister[options.FindOneOptions]) (
			*mongo.SingleResult, error)
		// FindOneAndDelete returns at most one document that matches the filter. If the filter
		// matches multiple documents, only the first document is deleted.
		FindOneAndDelete(ctx context.Context, filter any, opts ...options.Lister[options.FindOneAndDeleteOptions]) (
			*mongo.SingleResult, error)
		// FindOneAndReplace returns at most one document that matches the filter. If the filter
		// matches multiple documents, FindOneAndReplace returns the first document in the
		// collection that matches the filter.
		FindOneAndReplace(ctx context.Context, filter, replacement any,
			opts ...options.Lister[options.FindOneAndReplaceOptions]) (*mongo.SingleResult, error)
		// FindOneAndUpdate returns at most one document that matches the filter. If the filter
		// matches multiple documents, FindOneAndUpdate returns the first document in the
		// collection that matches the filter.
		FindOneAndUpdate(ctx context.Context, filter, update any,
			opts ...options.Lister[options.FindOneAndUpdateOptions]) (*mongo.SingleResult, error)
		// Indexes returns the index view for this collection.
		Indexes() mongo.IndexView
		// InsertMany inserts the provided documents.
		InsertMany(ctx context.Context, documents []any, opts ...options.Lister[options.InsertManyOptions]) (
			*mongo.InsertManyResult, error)
		// InsertOne inserts the provided document.
		InsertOne(ctx context.Context, document any, opts ...options.Lister[options.InsertOneOptions]) (
			*mongo.InsertOneResult, error)
		// ReplaceOne replaces at most one document that matches the filter.
		ReplaceOne(ctx context.Context, filter, replacement any,
			opts ...options.Lister[options.ReplaceOptions]) (*mongo.UpdateResult, error)
		// UpdateByID updates a single document matching the provided filter.
		UpdateByID(ctx context.Context, id, update any,
			opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error)
		// UpdateMany updates the provided documents.
		UpdateMany(ctx context.Context, filter, update any,
			opts ...options.Lister[options.UpdateManyOptions]) (*mongo.UpdateResult, error)
		// UpdateOne updates a single document matching the provided filter.
		UpdateOne(ctx context.Context, filter, update any,
			opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error)
		// Watch returns a change stream cursor used to receive notifications of changes to the collection.
		Watch(ctx context.Context, pipeline any, opts ...options.Lister[options.ChangeStreamOptions]) (
			*mongo.ChangeStream, error)
	}

	decoratedCollection struct {
		Collection monCollection
		name       string
		brk        breaker.Breaker
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

func (c *decoratedCollection) Aggregate(ctx context.Context, pipeline any,
	opts ...options.Lister[options.AggregateOptions]) (cur *mongo.Cursor, err error) {
	ctx, span := startSpan(ctx, aggregate)
	defer func() {
		endSpan(span, err)
	}()

	err = c.brk.DoWithAcceptableCtx(ctx, func() error {
		starTime := timex.Now()
		defer func() {
			c.logDurationSimple(ctx, aggregate, starTime, err)
		}()

		cur, err = c.Collection.Aggregate(ctx, pipeline, opts...)
		return err
	}, acceptable)

	return
}

func (c *decoratedCollection) BulkWrite(ctx context.Context, models []mongo.WriteModel,
	opts ...options.Lister[options.BulkWriteOptions]) (res *mongo.BulkWriteResult, err error) {
	ctx, span := startSpan(ctx, bulkWrite)
	defer func() {
		endSpan(span, err)
	}()

	err = c.brk.DoWithAcceptableCtx(ctx, func() error {
		startTime := timex.Now()
		defer func() {
			c.logDurationSimple(ctx, bulkWrite, startTime, err)
		}()

		res, err = c.Collection.BulkWrite(ctx, models, opts...)
		return err
	}, acceptable)

	return
}

func (c *decoratedCollection) Clone(opts ...options.Lister[options.CollectionOptions]) *mongo.Collection {
	return c.Collection.Clone(opts...)
}

func (c *decoratedCollection) CountDocuments(ctx context.Context, filter any,
	opts ...options.Lister[options.CountOptions]) (count int64, err error) {
	ctx, span := startSpan(ctx, countDocuments)
	defer func() {
		endSpan(span, err)
	}()

	err = c.brk.DoWithAcceptableCtx(ctx, func() error {
		startTime := timex.Now()
		defer func() {
			c.logDurationSimple(ctx, countDocuments, startTime, err)
		}()

		count, err = c.Collection.CountDocuments(ctx, filter, opts...)
		return err
	}, acceptable)

	return
}

func (c *decoratedCollection) Database() *mongo.Database {
	return c.Collection.Database()
}

func (c *decoratedCollection) DeleteMany(ctx context.Context, filter any,
	opts ...options.Lister[options.DeleteManyOptions]) (res *mongo.DeleteResult, err error) {
	ctx, span := startSpan(ctx, deleteMany)
	defer func() {
		endSpan(span, err)
	}()

	err = c.brk.DoWithAcceptableCtx(ctx, func() error {
		startTime := timex.Now()
		defer func() {
			c.logDurationSimple(ctx, deleteMany, startTime, err)
		}()

		res, err = c.Collection.DeleteMany(ctx, filter, opts...)
		return err
	}, acceptable)

	return
}

func (c *decoratedCollection) DeleteOne(ctx context.Context, filter any,
	opts ...options.Lister[options.DeleteOneOptions]) (res *mongo.DeleteResult, err error) {
	ctx, span := startSpan(ctx, deleteOne)
	defer func() {
		endSpan(span, err)
	}()

	err = c.brk.DoWithAcceptableCtx(ctx, func() error {
		startTime := timex.Now()
		defer func() {
			c.logDuration(ctx, deleteOne, startTime, err, filter)
		}()

		res, err = c.Collection.DeleteOne(ctx, filter, opts...)
		return err
	}, acceptable)

	return
}

func (c *decoratedCollection) Distinct(ctx context.Context, fieldName string, filter any,
	opts ...options.Lister[options.DistinctOptions]) (res *mongo.DistinctResult, err error) {
	ctx, span := startSpan(ctx, distinct)
	defer func() {
		endSpan(span, err)
	}()

	err = c.brk.DoWithAcceptableCtx(ctx, func() error {
		startTime := timex.Now()
		defer func() {
			c.logDurationSimple(ctx, distinct, startTime, err)
		}()

		res = c.Collection.Distinct(ctx, fieldName, filter, opts...)
		err = res.Err()
		return err
	}, acceptable)

	return
}

func (c *decoratedCollection) Drop(ctx context.Context, opts ...options.Lister[options.DropCollectionOptions]) error {
	return c.Collection.Drop(ctx, opts...)
}

func (c *decoratedCollection) EstimatedDocumentCount(ctx context.Context,
	opts ...options.Lister[options.EstimatedDocumentCountOptions]) (val int64, err error) {
	ctx, span := startSpan(ctx, estimatedDocumentCount)
	defer func() {
		endSpan(span, err)
	}()

	err = c.brk.DoWithAcceptableCtx(ctx, func() error {
		startTime := timex.Now()
		defer func() {
			c.logDurationSimple(ctx, estimatedDocumentCount, startTime, err)
		}()

		val, err = c.Collection.EstimatedDocumentCount(ctx, opts...)
		return err
	}, acceptable)

	return
}

func (c *decoratedCollection) Find(ctx context.Context, filter any,
	opts ...options.Lister[options.FindOptions]) (cur *mongo.Cursor, err error) {
	ctx, span := startSpan(ctx, find)
	defer func() {
		endSpan(span, err)
	}()

	err = c.brk.DoWithAcceptableCtx(ctx, func() error {
		startTime := timex.Now()
		defer func() {
			c.logDuration(ctx, find, startTime, err, filter)
		}()

		cur, err = c.Collection.Find(ctx, filter, opts...)
		return err
	}, acceptable)

	return
}

func (c *decoratedCollection) FindOne(ctx context.Context, filter any,
	opts ...options.Lister[options.FindOneOptions]) (res *mongo.SingleResult, err error) {
	ctx, span := startSpan(ctx, findOne)
	defer func() {
		endSpan(span, err)
	}()

	err = c.brk.DoWithAcceptableCtx(ctx, func() error {
		startTime := timex.Now()
		defer func() {
			c.logDuration(ctx, findOne, startTime, err, filter)
		}()

		res = c.Collection.FindOne(ctx, filter, opts...)
		err = res.Err()
		return err
	}, acceptable)

	return
}

func (c *decoratedCollection) FindOneAndDelete(ctx context.Context, filter any,
	opts ...options.Lister[options.FindOneAndDeleteOptions]) (res *mongo.SingleResult, err error) {
	ctx, span := startSpan(ctx, findOneAndDelete)
	defer func() {
		endSpan(span, err)
	}()

	err = c.brk.DoWithAcceptableCtx(ctx, func() error {
		startTime := timex.Now()
		defer func() {
			c.logDuration(ctx, findOneAndDelete, startTime, err, filter)
		}()

		res = c.Collection.FindOneAndDelete(ctx, filter, opts...)
		err = res.Err()
		return err
	}, acceptable)

	return
}

func (c *decoratedCollection) FindOneAndReplace(ctx context.Context, filter any,
	replacement any, opts ...options.Lister[options.FindOneAndReplaceOptions]) (
	res *mongo.SingleResult, err error) {
	ctx, span := startSpan(ctx, findOneAndReplace)
	defer func() {
		endSpan(span, err)
	}()

	err = c.brk.DoWithAcceptableCtx(ctx, func() error {
		startTime := timex.Now()
		defer func() {
			c.logDuration(ctx, findOneAndReplace, startTime, err, filter, replacement)
		}()

		res = c.Collection.FindOneAndReplace(ctx, filter, replacement, opts...)
		err = res.Err()
		return err
	}, acceptable)

	return
}

func (c *decoratedCollection) FindOneAndUpdate(ctx context.Context, filter, update any,
	opts ...options.Lister[options.FindOneAndUpdateOptions]) (res *mongo.SingleResult, err error) {
	ctx, span := startSpan(ctx, findOneAndUpdate)
	defer func() {
		endSpan(span, err)
	}()

	err = c.brk.DoWithAcceptableCtx(ctx, func() error {
		startTime := timex.Now()
		defer func() {
			c.logDuration(ctx, findOneAndUpdate, startTime, err, filter, update)
		}()

		res = c.Collection.FindOneAndUpdate(ctx, filter, update, opts...)
		err = res.Err()
		return err
	}, acceptable)

	return
}

func (c *decoratedCollection) Indexes() mongo.IndexView {
	return c.Collection.Indexes()
}

func (c *decoratedCollection) InsertMany(ctx context.Context, documents []any,
	opts ...options.Lister[options.InsertManyOptions]) (res *mongo.InsertManyResult, err error) {
	ctx, span := startSpan(ctx, insertMany)
	defer func() {
		endSpan(span, err)
	}()

	err = c.brk.DoWithAcceptableCtx(ctx, func() error {
		startTime := timex.Now()
		defer func() {
			c.logDurationSimple(ctx, insertMany, startTime, err)
		}()

		res, err = c.Collection.InsertMany(ctx, documents, opts...)
		return err
	}, acceptable)

	return
}

func (c *decoratedCollection) InsertOne(ctx context.Context, document any,
	opts ...options.Lister[options.InsertOneOptions]) (res *mongo.InsertOneResult, err error) {
	ctx, span := startSpan(ctx, insertOne)
	defer func() {
		endSpan(span, err)
	}()

	err = c.brk.DoWithAcceptableCtx(ctx, func() error {
		startTime := timex.Now()
		defer func() {
			c.logDuration(ctx, insertOne, startTime, err, document)
		}()

		res, err = c.Collection.InsertOne(ctx, document, opts...)
		return err
	}, acceptable)

	return
}

func (c *decoratedCollection) ReplaceOne(ctx context.Context, filter, replacement any,
	opts ...options.Lister[options.ReplaceOptions]) (res *mongo.UpdateResult, err error) {
	ctx, span := startSpan(ctx, replaceOne)
	defer func() {
		endSpan(span, err)
	}()

	err = c.brk.DoWithAcceptableCtx(ctx, func() error {
		startTime := timex.Now()
		defer func() {
			c.logDuration(ctx, replaceOne, startTime, err, filter, replacement)
		}()

		res, err = c.Collection.ReplaceOne(ctx, filter, replacement, opts...)
		return err
	}, acceptable)

	return
}

func (c *decoratedCollection) UpdateByID(ctx context.Context, id, update any,
	opts ...options.Lister[options.UpdateOneOptions]) (res *mongo.UpdateResult, err error) {
	ctx, span := startSpan(ctx, updateByID)
	defer func() {
		endSpan(span, err)
	}()

	err = c.brk.DoWithAcceptableCtx(ctx, func() error {
		startTime := timex.Now()
		defer func() {
			c.logDuration(ctx, updateByID, startTime, err, id, update)
		}()

		res, err = c.Collection.UpdateByID(ctx, id, update, opts...)
		return err
	}, acceptable)

	return
}

func (c *decoratedCollection) UpdateMany(ctx context.Context, filter, update any,
	opts ...options.Lister[options.UpdateManyOptions]) (res *mongo.UpdateResult, err error) {
	ctx, span := startSpan(ctx, updateMany)
	defer func() {
		endSpan(span, err)
	}()

	err = c.brk.DoWithAcceptableCtx(ctx, func() error {
		startTime := timex.Now()
		defer func() {
			c.logDurationSimple(ctx, updateMany, startTime, err)
		}()

		res, err = c.Collection.UpdateMany(ctx, filter, update, opts...)
		return err
	}, acceptable)

	return
}

func (c *decoratedCollection) UpdateOne(ctx context.Context, filter, update any,
	opts ...options.Lister[options.UpdateOneOptions]) (res *mongo.UpdateResult, err error) {
	ctx, span := startSpan(ctx, updateOne)
	defer func() {
		endSpan(span, err)
	}()

	err = c.brk.DoWithAcceptableCtx(ctx, func() error {
		startTime := timex.Now()
		defer func() {
			c.logDuration(ctx, updateOne, startTime, err, filter, update)
		}()

		res, err = c.Collection.UpdateOne(ctx, filter, update, opts...)
		return err
	}, acceptable)

	return
}

func (c *decoratedCollection) Watch(ctx context.Context, pipeline any, opts ...options.Lister[options.ChangeStreamOptions]) (
	*mongo.ChangeStream, error) {
	return c.Collection.Watch(ctx, pipeline, opts...)
}

func (c *decoratedCollection) logDuration(ctx context.Context, method string,
	startTime time.Duration, err error, docs ...any) {
	logDurationWithDocs(ctx, c.name, method, startTime, err, docs...)
}

func (c *decoratedCollection) logDurationSimple(ctx context.Context, method string,
	startTime time.Duration, err error) {
	logDuration(ctx, c.name, method, startTime, err)
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
	return err == nil || isDupKeyError(err) ||
		errorx.In(err, mongo.ErrNoDocuments, mongo.ErrNilValue,
			mongo.ErrNilDocument, mongo.ErrNilCursor, mongo.ErrEmptySlice,
			// session errors
			session.ErrSessionEnded, session.ErrNoTransactStarted, session.ErrTransactInProgress,
			session.ErrAbortAfterCommit, session.ErrAbortTwice, session.ErrCommitAfterAbort,
			session.ErrUnackWCUnsupported, session.ErrSnapshotTransaction)
}

func isDupKeyError(err error) bool {
	var e mongo.WriteException
	if !errors.As(err, &e) {
		return false
	}

	return e.HasErrorCode(duplicateKeyCode)
}

// monCollection defines a MongoDB collection, used for unit test
type monCollection interface {
	// Aggregate executes an aggregation pipeline.
	Aggregate(ctx context.Context, pipeline any, opts ...options.Lister[options.AggregateOptions]) (
		*mongo.Cursor, error)
	// BulkWrite performs a bulk write operation.
	BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...options.Lister[options.BulkWriteOptions]) (
		*mongo.BulkWriteResult, error)
	// Clone creates a copy of this collection with the same settings.
	Clone(opts ...options.Lister[options.CollectionOptions]) *mongo.Collection
	// CountDocuments returns the number of documents in the collection that match the filter.
	CountDocuments(ctx context.Context, filter any, opts ...options.Lister[options.CountOptions]) (int64, error)
	// Database returns the database that this collection is a part of.
	Database() *mongo.Database
	// DeleteMany deletes documents from the collection that match the filter.
	DeleteMany(ctx context.Context, filter any, opts ...options.Lister[options.DeleteManyOptions]) (
		*mongo.DeleteResult, error)
	// DeleteOne deletes at most one document from the collection that matches the filter.
	DeleteOne(ctx context.Context, filter any, opts ...options.Lister[options.DeleteOneOptions]) (
		*mongo.DeleteResult, error)
	// Distinct returns a list of distinct values for the given key across the collection.
	Distinct(ctx context.Context, fieldName string, filter any,
		opts ...options.Lister[options.DistinctOptions]) *mongo.DistinctResult
	// Drop drops this collection from database.
	Drop(ctx context.Context, opts ...options.Lister[options.DropCollectionOptions]) error
	// EstimatedDocumentCount returns an estimate of the count of documents in a collection
	// using collection metadata.
	EstimatedDocumentCount(ctx context.Context, opts ...options.Lister[options.EstimatedDocumentCountOptions]) (int64, error)
	// Find finds the documents matching the provided filter.
	Find(ctx context.Context, filter any, opts ...options.Lister[options.FindOptions]) (*mongo.Cursor, error)
	// FindOne returns up to one document that matches the provided filter.
	FindOne(ctx context.Context, filter any, opts ...options.Lister[options.FindOneOptions]) *mongo.SingleResult
	// FindOneAndDelete returns at most one document that matches the filter. If the filter
	// matches multiple documents, only the first document is deleted.
	FindOneAndDelete(ctx context.Context, filter any, opts ...options.Lister[options.FindOneAndDeleteOptions]) *mongo.SingleResult
	// FindOneAndReplace returns at most one document that matches the filter. If the filter
	// matches multiple documents, FindOneAndReplace returns the first document in the
	// collection that matches the filter.
	FindOneAndReplace(ctx context.Context, filter, replacement any,
		opts ...options.Lister[options.FindOneAndReplaceOptions]) *mongo.SingleResult
	// FindOneAndUpdate returns at most one document that matches the filter. If the filter
	// matches multiple documents, FindOneAndUpdate returns the first document in the
	// collection that matches the filter.
	FindOneAndUpdate(ctx context.Context, filter, update any,
		opts ...options.Lister[options.FindOneAndUpdateOptions]) *mongo.SingleResult
	// Indexes returns the index view for this collection.
	Indexes() mongo.IndexView
	// InsertMany inserts the provided documents.
	InsertMany(ctx context.Context, documents interface{}, opts ...options.Lister[options.InsertManyOptions]) (*mongo.InsertManyResult, error)
	// InsertOne inserts the provided document.
	InsertOne(ctx context.Context, document any, opts ...options.Lister[options.InsertOneOptions]) (*mongo.InsertOneResult, error)
	// ReplaceOne replaces at most one document that matches the filter.
	ReplaceOne(ctx context.Context, filter, replacement any,
		opts ...options.Lister[options.ReplaceOptions]) (*mongo.UpdateResult, error)
	// UpdateByID updates a single document matching the provided filter.
	UpdateByID(ctx context.Context, id, update any,
		opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error)
	// UpdateMany updates the provided documents.
	UpdateMany(ctx context.Context, filter, update any,
		opts ...options.Lister[options.UpdateManyOptions]) (*mongo.UpdateResult, error)
	// UpdateOne updates a single document matching the provided filter.
	UpdateOne(ctx context.Context, filter, update any,
		opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error)
	// Watch returns a change stream cursor used to receive notifications of changes to the collection.
	Watch(ctx context.Context, pipeline any, opts ...options.Lister[options.ChangeStreamOptions]) (
		*mongo.ChangeStream, error)
}
