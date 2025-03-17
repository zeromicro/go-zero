package mon

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/timex"
	"go.mongodb.org/mongo-driver/v2/mongo"
	mopt "go.mongodb.org/mongo-driver/v2/mongo/options"
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
		cli  monCli
		brk  breaker.Breaker
		opts []Option
	}

	WrappedSession struct {
		session monSession
		name    string
		brk     breaker.Breaker
	}

	//for unit test, this is a little annoying
	monCli interface {
		StartSession(opts ...mopt.Lister[mopt.SessionOptions]) (monSession, error)
	}
	monSession interface {
		AbortTransaction(ctx context.Context) error
		CommitTransaction(ctx context.Context) error
		EndSession(ctx context.Context)
		WithTransaction(ctx context.Context, fn func(sessCtx context.Context) (any, error),
			opts ...mopt.Lister[mopt.TransactionOptions]) (any, error)
	}
)

// MustNewModel returns a Model, exits on errors.
func MustNewModel(uri, db, collection string, opts ...Option) *Model {
	model, err := NewModel(uri, db, collection, opts...)
	logx.Must(err)
	return model
}

// mustNewTestModel returns a Model for unit test, exits on errors.
func mustNewTestModel(uri, db, collection string, opts ...Option) *Model {
	model, err := newUnitTestModel(uri, db, collection, opts...)
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

// newUnitTestModel returns a Model for unit test.
func newUnitTestModel(uri, db, collection string, opts ...Option) (*Model, error) {
	//cli, err := getClient(uri, opts...)
	//if err != nil {
	//	return nil, err
	//}

	name := strings.Join([]string{uri, collection}, "/")
	brk := breaker.GetBreaker(uri)
	coll := newTestCollection(brk)
	return newModelTest(name, coll, brk, opts...), nil
}

type mockMonClient struct {
	cli *mongo.Client
}

func (m *mockMonClient) StartSession(opts ...mopt.Lister[mopt.SessionOptions]) (monSession, error) {
	if m.cli != nil {
		return m.cli.StartSession(opts...)
	}
	return &mockSession{}, nil
}

type mockSession struct {
}

func (m *mockSession) AbortTransaction(ctx context.Context) error {
	return nil
}

func (m *mockSession) CommitTransaction(ctx context.Context) error {
	return nil
}

func (m *mockSession) EndSession(ctx context.Context) {

}

func (m *mockSession) WithTransaction(ctx context.Context, fn func(sessCtx context.Context) (any, error),
	opts ...mopt.Lister[mopt.TransactionOptions]) (any, error) {
	return nil, nil
}

func newModelTest(name string, coll Collection, brk breaker.Breaker,
	opts ...Option) *Model {
	return &Model{
		name:       name,
		Collection: coll,
		cli:        &mockMonClient{},
		brk:        brk,
		opts:       opts,
	}
}

func newModel(name string, cli *mongo.Client, coll Collection, brk breaker.Breaker,
	opts ...Option) *Model {
	return &Model{
		name:       name,
		Collection: coll,
		cli:        &mockMonClient{cli: cli},
		brk:        brk,
		opts:       opts,
	}
}

// StartSession starts a new session.
func (m *Model) StartSession(opts ...mopt.Lister[mopt.SessionOptions]) (sess *WrappedSession, err error) {
	starTime := timex.Now()
	defer func() {
		logDuration(context.Background(), m.name, startSession, starTime, err)
	}()

	session, sessionErr := m.cli.StartSession(opts...)
	if sessionErr != nil {
		return nil, sessionErr
	}

	return &WrappedSession{
		session: session,
		name:    m.name,
		brk:     m.brk,
	}, nil
}

// Aggregate executes an aggregation pipeline.
func (m *Model) Aggregate(ctx context.Context, v, pipeline any, opts ...mopt.Lister[mopt.AggregateOptions]) error {
	cur, err := m.Collection.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return err
	}
	defer cur.Close(ctx)

	return cur.All(ctx, v)
}

// DeleteMany deletes documents that match the filter.
func (m *Model) DeleteMany(ctx context.Context, filter any, opts ...mopt.Lister[mopt.DeleteManyOptions]) (int64, error) {
	res, err := m.Collection.DeleteMany(ctx, filter, opts...)
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}

// DeleteOne deletes the first document that matches the filter.
func (m *Model) DeleteOne(ctx context.Context, filter any, opts ...mopt.Lister[mopt.DeleteOneOptions]) (int64, error) {
	res, err := m.Collection.DeleteOne(ctx, filter, opts...)
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}

// Find finds documents that match the filter.
func (m *Model) Find(ctx context.Context, v, filter any, opts ...mopt.Lister[mopt.FindOptions]) error {
	cur, err := m.Collection.Find(ctx, filter, opts...)
	if err != nil {
		return err
	}
	defer cur.Close(ctx)

	return cur.All(ctx, v)
}

// FindOne finds the first document that matches the filter.
func (m *Model) FindOne(ctx context.Context, v, filter any, opts ...mopt.Lister[mopt.FindOneOptions]) error {
	res, err := m.Collection.FindOne(ctx, filter, opts...)
	if err != nil {
		return err
	}

	return res.Decode(v)
}

// FindOneAndDelete finds a single document and deletes it.
func (m *Model) FindOneAndDelete(ctx context.Context, v, filter any,
	opts ...mopt.Lister[mopt.FindOneAndDeleteOptions]) error {
	res, err := m.Collection.FindOneAndDelete(ctx, filter, opts...)
	if err != nil {
		return err
	}

	return res.Decode(v)
}

// FindOneAndReplace finds a single document and replaces it.
func (m *Model) FindOneAndReplace(ctx context.Context, v, filter, replacement any,
	opts ...mopt.Lister[mopt.FindOneAndReplaceOptions]) error {
	res, err := m.Collection.FindOneAndReplace(ctx, filter, replacement, opts...)
	if err != nil {
		return err
	}

	return res.Decode(v)
}

// FindOneAndUpdate finds a single document and updates it.
func (m *Model) FindOneAndUpdate(ctx context.Context, v, filter, update any,
	opts ...mopt.Lister[mopt.FindOneAndUpdateOptions]) error {
	res, err := m.Collection.FindOneAndUpdate(ctx, filter, update, opts...)
	if err != nil {
		return err
	}

	return res.Decode(v)
}

// AbortTransaction implements the mongo.session interface.
func (w *WrappedSession) AbortTransaction(ctx context.Context) (err error) {
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
func (w *WrappedSession) CommitTransaction(ctx context.Context) (err error) {
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
func (w *WrappedSession) WithTransaction(
	ctx context.Context,
	fn func(sessCtx context.Context) (any, error),
	opts ...mopt.Lister[mopt.TransactionOptions],
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
func (w *WrappedSession) EndSession(ctx context.Context) {
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
