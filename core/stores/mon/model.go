//go:generate mockgen -package mon -destination model_mock.go -source model.go monClient monSession
package mon

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/timex"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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
		cli  monClient
		brk  breaker.Breaker
		opts []Option
	}

	Session struct {
		session monSession
		name    string
		brk     breaker.Breaker
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
		cli:        &wrappedMonClient{c: cli},
		brk:        brk,
		opts:       opts,
	}
}

// StartSession starts a new session.
func (m *Model) StartSession(opts ...options.Lister[options.SessionOptions]) (sess *Session, err error) {
	starTime := timex.Now()
	defer func() {
		logDuration(context.Background(), m.name, startSession, starTime, err)
	}()

	session, sessionErr := m.cli.StartSession(opts...)
	if sessionErr != nil {
		return nil, sessionErr
	}

	return &Session{
		session: session,
		name:    m.name,
		brk:     m.brk,
	}, nil
}

// Aggregate executes an aggregation pipeline.
func (m *Model) Aggregate(ctx context.Context, v, pipeline any,
	opts ...options.Lister[options.AggregateOptions]) error {
	cur, err := m.Collection.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return err
	}
	defer cur.Close(ctx)

	return cur.All(ctx, v)
}

// DeleteMany deletes documents that match the filter.
func (m *Model) DeleteMany(ctx context.Context, filter any,
	opts ...options.Lister[options.DeleteManyOptions]) (int64, error) {
	res, err := m.Collection.DeleteMany(ctx, filter, opts...)
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}

// DeleteOne deletes the first document that matches the filter.
func (m *Model) DeleteOne(ctx context.Context, filter any,
	opts ...options.Lister[options.DeleteOneOptions]) (int64, error) {
	res, err := m.Collection.DeleteOne(ctx, filter, opts...)
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}

// Find finds documents that match the filter.
func (m *Model) Find(ctx context.Context, v, filter any,
	opts ...options.Lister[options.FindOptions]) error {
	cur, err := m.Collection.Find(ctx, filter, opts...)
	if err != nil {
		return err
	}
	defer cur.Close(ctx)

	return cur.All(ctx, v)
}

// FindOne finds the first document that matches the filter.
func (m *Model) FindOne(ctx context.Context, v, filter any,
	opts ...options.Lister[options.FindOneOptions]) error {
	res, err := m.Collection.FindOne(ctx, filter, opts...)
	if err != nil {
		return err
	}

	return res.Decode(v)
}

// FindOneAndDelete finds a single document and deletes it.
func (m *Model) FindOneAndDelete(ctx context.Context, v, filter any,
	opts ...options.Lister[options.FindOneAndDeleteOptions]) error {
	res, err := m.Collection.FindOneAndDelete(ctx, filter, opts...)
	if err != nil {
		return err
	}

	return res.Decode(v)
}

// FindOneAndReplace finds a single document and replaces it.
func (m *Model) FindOneAndReplace(ctx context.Context, v, filter, replacement any,
	opts ...options.Lister[options.FindOneAndReplaceOptions]) error {
	res, err := m.Collection.FindOneAndReplace(ctx, filter, replacement, opts...)
	if err != nil {
		return err
	}

	return res.Decode(v)
}

// FindOneAndUpdate finds a single document and updates it.
func (m *Model) FindOneAndUpdate(ctx context.Context, v, filter, update any,
	opts ...options.Lister[options.FindOneAndUpdateOptions]) error {
	res, err := m.Collection.FindOneAndUpdate(ctx, filter, update, opts...)
	if err != nil {
		return err
	}

	return res.Decode(v)
}

// AbortTransaction implements the mongo.session interface.
func (w *Session) AbortTransaction(ctx context.Context) (err error) {
	ctx, span := startSpan(ctx, abortTransaction)
	defer func() {
		endSpan(span, err)
	}()

	return w.brk.DoWithAcceptableCtx(ctx, func() error {
		starTime := timex.Now()
		defer func() {
			logDuration(ctx, w.name, abortTransaction, starTime, err)
		}()

		return w.session.AbortTransaction(ctx)
	}, acceptable)
}

// CommitTransaction implements the mongo.session interface.
func (w *Session) CommitTransaction(ctx context.Context) (err error) {
	ctx, span := startSpan(ctx, commitTransaction)
	defer func() {
		endSpan(span, err)
	}()

	return w.brk.DoWithAcceptableCtx(ctx, func() error {
		starTime := timex.Now()
		defer func() {
			logDuration(ctx, w.name, commitTransaction, starTime, err)
		}()

		return w.session.CommitTransaction(ctx)
	}, acceptable)
}

// WithTransaction implements the mongo.session interface.
func (w *Session) WithTransaction(
	ctx context.Context,
	fn func(sessCtx context.Context) (any, error),
	opts ...options.Lister[options.TransactionOptions],
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

		res, err = w.session.WithTransaction(ctx, fn, opts...)
		return err
	}, acceptable)

	return
}

// EndSession implements the mongo.session interface.
func (w *Session) EndSession(ctx context.Context) {
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

		w.session.EndSession(ctx)
		return nil
	}, acceptable)
}

type (
	// for unit test
	monClient interface {
		StartSession(opts ...options.Lister[options.SessionOptions]) (monSession, error)
	}

	monSession interface {
		AbortTransaction(ctx context.Context) error
		CommitTransaction(ctx context.Context) error
		EndSession(ctx context.Context)
		WithTransaction(ctx context.Context, fn func(sessCtx context.Context) (any, error),
			opts ...options.Lister[options.TransactionOptions]) (any, error)
	}
)

type wrappedMonClient struct {
	c *mongo.Client
}

// StartSession starts a new session using the underlying *mongo.Client.
// It implements the monClient interface.
// This is used to allow mocking in unit tests.
func (m *wrappedMonClient) StartSession(opts ...options.Lister[options.SessionOptions]) (
	monSession, error) {
	return m.c.StartSession(opts...)
}
