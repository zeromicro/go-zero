package mon

import (
	"context"
	"log"
	"strings"

	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/core/timex"
	"go.mongodb.org/mongo-driver/mongo"
	mopt "go.mongodb.org/mongo-driver/mongo/options"
)

// Model is a mongodb store model that represents a collection.
type Model struct {
	Collection
	name string
	cli  *mongo.Client
	brk  breaker.Breaker
	opts []Option
}

// MustNewModel returns a Model, exits on errors.
func MustNewModel(uri, db, collection string, opts ...Option) *Model {
	model, err := NewModel(uri, db, collection, opts...)
	if err != nil {
		log.Fatal(err)
	}

	return model
}

// NewModel returns a Model.
func NewModel(uri, db, collection string, opts ...Option) (*Model, error) {
	cli, err := getClient(uri)
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
	err = m.brk.DoWithAcceptable(func() error {
		starTime := timex.Now()
		defer func() {
			logDuration(m.name, "StartSession", starTime, err)
		}()

		sess, err = m.cli.StartSession(opts...)
		return err
	}, acceptable)
	return
}

// Aggregate executes an aggregation pipeline.
func (m *Model) Aggregate(ctx context.Context, v, pipeline interface{}, opts ...*mopt.AggregateOptions) error {
	cur, err := m.Collection.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return err
	}
	defer cur.Close(ctx)

	return cur.All(ctx, v)
}

// DeleteMany deletes documents that match the filter.
func (m *Model) DeleteMany(ctx context.Context, filter interface{}, opts ...*mopt.DeleteOptions) (int64, error) {
	res, err := m.Collection.DeleteMany(ctx, filter, opts...)
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}

// DeleteOne deletes the first document that matches the filter.
func (m *Model) DeleteOne(ctx context.Context, filter interface{}, opts ...*mopt.DeleteOptions) (int64, error) {
	res, err := m.Collection.DeleteOne(ctx, filter, opts...)
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}

// Find finds documents that match the filter.
func (m *Model) Find(ctx context.Context, v, filter interface{}, opts ...*mopt.FindOptions) error {
	cur, err := m.Collection.Find(ctx, filter, opts...)
	if err != nil {
		return err
	}
	defer cur.Close(ctx)

	return cur.All(ctx, v)
}

// FindOne finds the first document that matches the filter.
func (m *Model) FindOne(ctx context.Context, v, filter interface{}, opts ...*mopt.FindOneOptions) error {
	res, err := m.Collection.FindOne(ctx, filter, opts...)
	if err != nil {
		return err
	}

	return res.Decode(v)
}

// FindOneAndDelete finds a single document and deletes it.
func (m *Model) FindOneAndDelete(ctx context.Context, v, filter interface{},
	opts ...*mopt.FindOneAndDeleteOptions) error {
	res, err := m.Collection.FindOneAndDelete(ctx, filter, opts...)
	if err != nil {
		return err
	}

	return res.Decode(v)
}

// FindOneAndReplace finds a single document and replaces it.
func (m *Model) FindOneAndReplace(ctx context.Context, v, filter interface{}, replacement interface{},
	opts ...*mopt.FindOneAndReplaceOptions) error {
	res, err := m.Collection.FindOneAndReplace(ctx, filter, replacement, opts...)
	if err != nil {
		return err
	}

	return res.Decode(v)
}

// FindOneAndUpdate finds a single document and updates it.
func (m *Model) FindOneAndUpdate(ctx context.Context, v, filter interface{}, update interface{},
	opts ...*mopt.FindOneAndUpdateOptions) error {
	res, err := m.Collection.FindOneAndUpdate(ctx, filter, update, opts...)
	if err != nil {
		return err
	}

	return res.Decode(v)
}
