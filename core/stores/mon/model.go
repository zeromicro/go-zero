package mon

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/timex"
	"go.mongodb.org/mongo-driver/mongo"
	mopt "go.mongodb.org/mongo-driver/mongo/options"
)

const (
	startSession      = "StartSession"
	abortTransaction  = "AbortTransaction"
	commitTransaction = "CommitTransaction"
	withTransaction   = "WithTransaction"
	endSession        = "EndSession"
)

type (
	// Model is a mongodb store model that represents a collection.
	Model struct {
		Collection
		name string
		cli  *mongo.Client
		brk  breaker.Breaker
		opts []Option
	}

	wrappedSession struct {
		mongo.Session
		name string
		brk  breaker.Breaker
	}
)

// MustNewModel returns a Model, exits on errors.
func MustNewModel(uri, db, collection string, opts ...Option) *Model {
	model, err := NewModel(uri, db, collection, opts...)
	logx.Must(err)
	return model
}

// NewModel returns a Model.
func NewModel(uri, db, collection string, opts ...Option) (*Model, error) {
	cli, err := getClient(uri, opts...)
	if err != nil {
		return nil, err
	}

	name := strings.Join([]string{uri, collection}, "/")
	brk := breaker.GetBreaker(uri)
	coll := newCollection(cli.Database(db).Collection(collection), brk)
	return newModel(name, cli, coll, brk, opts...), nil
}

func newModel(name string, cli *mongo.Client, coll Collection, brk breaker.Breaker,
	opts ...Option) *Model {
	return &Model{
		name:       name,
		Collection: coll,
		cli:        cli,
		brk:        brk,
		opts:       opts,
	}
}

// StartSession starts a new session.
func (m *Model) StartSession(opts ...*mopt.SessionOptions) (sess mongo.Session, err error) {
	starTime := timex.Now()
	defer func() {
		logDuration(context.Background(), m.name, startSession, starTime, err)
	}()

	session, sessionErr := m.cli.StartSession(opts...)
	if sessionErr != nil {
		return nil, sessionErr
	}

	return &wrappedSession{
		Session: session,
		name:    m.name,
		brk:     m.brk,
	}, nil
}

// Aggregate executes an aggregation pipeline.
func (m *Model) Aggregate(ctx context.Context, v, pipeline any, opts ...*mopt.AggregateOptions) error {
	cur, err := m.Collection.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return err
	}
	defer cur.Close(ctx)

	return cur.All(ctx, v)
}

// DeleteMany deletes documents that match the filter.
func (m *Model) DeleteMany(ctx context.Context, filter any, opts ...*mopt.DeleteOptions) (int64, error) {
	res, err := m.Collection.DeleteMany(ctx, filter, opts...)
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}

// DeleteOne deletes the first document that matches the filter.
func (m *Model) DeleteOne(ctx context.Context, filter any, opts ...*mopt.DeleteOptions) (int64, error) {
	res, err := m.Collection.DeleteOne(ctx, filter, opts...)
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}

// Find finds documents that match the filter.
func (m *Model) Find(ctx context.Context, v, filter any, opts ...*mopt.FindOptions) error {
	cur, err := m.Collection.Find(ctx, filter, opts...)
	if err != nil {
		return err
	}
	defer cur.Close(ctx)

	return cur.All(ctx, v)
}

// FindOne finds the first document that matches the filter.
func (m *Model) FindOne(ctx context.Context, v, filter any, opts ...*mopt.FindOneOptions) error {
	res, err := m.Collection.FindOne(ctx, filter, opts...)
	if err != nil {
		return err
	}

	return res.Decode(v)
}

// FindOneAndDelete finds a single document and deletes it.
func (m *Model) FindOneAndDelete(ctx context.Context, v, filter any,
	opts ...*mopt.FindOneAndDeleteOptions) error {
	res, err := m.Collection.FindOneAndDelete(ctx, filter, opts...)
	if err != nil {
		return err
	}

	return res.Decode(v)
}

// FindOneAndReplace finds a single document and replaces it.
func (m *Model) FindOneAndReplace(ctx context.Context, v, filter, replacement any,
	opts ...*mopt.FindOneAndReplaceOptions) error {
	res, err := m.Collection.FindOneAndReplace(ctx, filter, replacement, opts...)
	if err != nil {
		return err
	}

	return res.Decode(v)
}

// FindOneAndUpdate finds a single document and updates it.
func (m *Model) FindOneAndUpdate(ctx context.Context, v, filter, update any,
	opts ...*mopt.FindOneAndUpdateOptions) error {
	res, err := m.Collection.FindOneAndUpdate(ctx, filter, update, opts...)
	if err != nil {
		return err
	}

	return res.Decode(v)
}

// AbortTransaction implements the mongo.Session interface.
func (w *wrappedSession) AbortTransaction(ctx context.Context) (err error) {
	ctx, span := startSpan(ctx, abortTransaction)
	defer func() {
		endSpan(span, err)
	}()

	return w.brk.DoWithAcceptableCtx(ctx, func() error {
		starTime := timex.Now()
		defer func() {
			logDuration(ctx, w.name, abortTransaction, starTime, err)
		}()

		return w.Session.AbortTransaction(ctx)
	}, acceptable)
}

// CommitTransaction implements the mongo.Session interface.
func (w *wrappedSession) CommitTransaction(ctx context.Context) (err error) {
	ctx, span := startSpan(ctx, commitTransaction)
	defer func() {
		endSpan(span, err)
	}()

	return w.brk.DoWithAcceptableCtx(ctx, func() error {
		starTime := timex.Now()
		defer func() {
			logDuration(ctx, w.name, commitTransaction, starTime, err)
		}()

		return w.Session.CommitTransaction(ctx)
	}, acceptable)
}

// WithTransaction implements the mongo.Session interface.
func (w *wrappedSession) WithTransaction(
	ctx context.Context,
	fn func(sessCtx mongo.SessionContext) (any, error),
	opts ...*mopt.TransactionOptions,
) (res any, err error) {
	ctx, span := startSpan(ctx, withTransaction)
	defer func() {
		endSpan(span, err)
	}()

	err = w.brk.DoWithAcceptableCtx(ctx, func() error {
		starTime := timex.Now()
		defer func() {
			logDuration(ctx, w.name, withTransaction, starTime, err)
		}()

		res, err = w.Session.WithTransaction(ctx, fn, opts...)
		return err
	}, acceptable)

	return
}

// EndSession implements the mongo.Session interface.
func (w *wrappedSession) EndSession(ctx context.Context) {
	var err error
	ctx, span := startSpan(ctx, endSession)
	defer func() {
		endSpan(span, err)
	}()

	err = w.brk.DoWithAcceptableCtx(ctx, func() error {
		starTime := timex.Now()
		defer func() {
			logDuration(ctx, w.name, endSession, starTime, err)
		}()

		w.Session.EndSession(ctx)
		return nil
	}, acceptable)
}
