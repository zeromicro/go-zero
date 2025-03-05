package mon

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/core/logx/logtest"
	"github.com/zeromicro/go-zero/core/stringx"
	"github.com/zeromicro/go-zero/core/timex"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	mopt "go.mongodb.org/mongo-driver/v2/mongo/options"
)

var errDummy = errors.New("dummy")

func TestKeepPromise_accept(t *testing.T) {
	p := new(mockPromise)
	kp := keepablePromise{
		promise: p,
		log:     func(error) {},
	}
	assert.Nil(t, kp.accept(nil))
	assert.Equal(t, ErrNotFound, kp.accept(ErrNotFound))
}

func TestKeepPromise_keep(t *testing.T) {
	tests := []struct {
		err      error
		accepted bool
		reason   string
	}{
		{
			err:      nil,
			accepted: true,
			reason:   "",
		},
		{
			err:      ErrNotFound,
			accepted: true,
			reason:   "",
		},
		{
			err:      errors.New("any"),
			accepted: false,
			reason:   "any",
		},
	}

	for _, test := range tests {
		t.Run(stringx.RandId(), func(t *testing.T) {
			p := new(mockPromise)
			kp := keepablePromise{
				promise: p,
				log:     func(error) {},
			}
			assert.Equal(t, test.err, kp.keep(test.err))
			assert.Equal(t, test.accepted, p.accepted)
			assert.Equal(t, test.reason, p.reason)
		})
	}
}

func TestNewCollection(t *testing.T) {
	_ = newCollection(&mongo.Collection{}, breaker.GetBreaker("localhost"))
}

func TestCollection_Aggregate(t *testing.T) {
	c := newTestCollection(breaker.GetBreaker("localhost"))
	_, err := c.Aggregate(context.Background(), []interface{}{}, mopt.Aggregate())
	assert.Nil(t, err)
}

func TestCollection_BulkWrite(t *testing.T) {
	c := newTestCollection(breaker.GetBreaker("localhost"))
	_, err := c.BulkWrite(context.Background(), []mongo.WriteModel{
		mongo.NewInsertOneModel().SetDocument(bson.D{{Key: "foo", Value: 1}}),
	})
	assert.Nil(t, err)
	c.brk = new(dropBreaker)
	_, err = c.BulkWrite(context.Background(), []mongo.WriteModel{
		mongo.NewInsertOneModel().SetDocument(bson.D{{Key: "foo", Value: 1}}),
	})
	assert.Equal(t, errDummy, err)
}

func TestCollection_CountDocuments(t *testing.T) {
	c := newTestCollection(breaker.GetBreaker("localhost"))
	res, err := c.CountDocuments(context.Background(), bson.D{})
	assert.Nil(t, err)
	assert.Equal(t, int64(0), res)

	c.brk = new(dropBreaker)
	_, err = c.CountDocuments(context.Background(), bson.D{{Key: "foo", Value: 1}})
	assert.Equal(t, errDummy, err)
}

func TestDecoratedCollection_DeleteMany(t *testing.T) {
	c := newTestCollection(breaker.GetBreaker("localhost"))
	_, err := c.DeleteMany(context.Background(), bson.D{})
	assert.Nil(t, err)
	//assert.Equal(t, int64(1), res.DeletedCount)

	c.brk = new(dropBreaker)
	_, err = c.DeleteMany(context.Background(), bson.D{{Key: "foo", Value: 1}})
	assert.Equal(t, errDummy, err)
}

func TestCollection_Distinct(t *testing.T) {
	c := newTestCollection(breaker.GetBreaker("localhost"))
	_, err := c.Distinct(context.Background(), "foo", bson.D{})
	assert.Nil(t, err)
	c.brk = new(dropBreaker)
	_, err = c.Distinct(context.Background(), "foo", bson.D{{Key: "foo", Value: 1}})
	assert.Equal(t, errDummy, err)
}

func TestCollection_EstimatedDocumentCount(t *testing.T) {
	c := newTestCollection(breaker.GetBreaker("localhost"))
	_, err := c.EstimatedDocumentCount(context.Background())
	assert.Nil(t, err)
	c.brk = new(dropBreaker)
	_, err = c.EstimatedDocumentCount(context.Background())
	assert.Equal(t, errDummy, err)
}

func TestCollection_Find(t *testing.T) {
	c := newTestCollection(breaker.GetBreaker("localhost"))
	filter := bson.D{{Key: "x", Value: 1}}
	_, err := c.Find(context.Background(), filter, mopt.Find())
	assert.Nil(t, err)
	c.brk = new(dropBreaker)
	_, err = c.Find(context.Background(), filter, mopt.Find())
	assert.Equal(t, errDummy, err)
}

func TestCollection_FindOne(t *testing.T) {
	c := newTestCollection(breaker.GetBreaker("localhost"))
	filter := bson.D{{Key: "x", Value: 1}}
	_, err := c.FindOne(context.Background(), filter)
	c.brk = new(dropBreaker)
	_, err = c.FindOne(context.Background(), filter)
	assert.Equal(t, errDummy, err)
}

func TestCollection_FindOneAndDelete(t *testing.T) {
	c := newTestCollection(breaker.GetBreaker("localhost"))
	filter := bson.D{}
	_, err := c.FindOneAndDelete(context.Background(), filter, mopt.FindOneAndDelete())
	assert.Equal(t, mongo.ErrNoDocuments, err)
	_, err = c.FindOneAndDelete(context.Background(), filter, mopt.FindOneAndDelete())
	c.brk = new(dropBreaker)
	_, err = c.FindOneAndDelete(context.Background(), bson.D{{Key: "foo", Value: "bar"}})
	assert.Equal(t, errDummy, err)
}

func TestCollection_FindOneAndReplace(t *testing.T) {
	c := newTestCollection(breaker.GetBreaker("localhost"))
	filter := bson.D{{Key: "x", Value: 1}}
	replacement := bson.D{{Key: "x", Value: 2}}
	opts := mopt.FindOneAndReplace().SetUpsert(true)
	_, err := c.FindOneAndReplace(context.Background(), filter, replacement, opts)
	assert.Equal(t, mongo.ErrNoDocuments, err)
	_, err = c.FindOneAndReplace(context.Background(), filter, replacement, opts)
	c.brk = new(dropBreaker)
	_, err = c.FindOneAndReplace(context.Background(), filter, replacement, opts)
	assert.Equal(t, errDummy, err)
}

func TestCollection_FindOneAndUpdate(t *testing.T) {
	c := newTestCollection(breaker.GetBreaker("localhost"))
	filter := bson.D{{Key: "x", Value: 1}}
	update := bson.D{{Key: "$x", Value: 2}}
	opts := mopt.FindOneAndUpdate().SetUpsert(true)
	_, err := c.FindOneAndUpdate(context.Background(), filter, update, opts)
	assert.Equal(t, mongo.ErrNoDocuments, err)
	_, err = c.FindOneAndUpdate(context.Background(), filter, update, opts)
	c.brk = new(dropBreaker)
	_, err = c.FindOneAndUpdate(context.Background(), filter, update, opts)
	assert.Equal(t, errDummy, err)
}

func TestCollection_InsertOne(t *testing.T) {
	c := newTestCollection(breaker.GetBreaker("localhost"))
	res, err := c.InsertOne(context.Background(), bson.D{{Key: "foo", Value: "bar"}})
	assert.Nil(t, err)
	assert.NotNil(t, res)
	c.brk = new(dropBreaker)
	_, err = c.InsertOne(context.Background(), bson.D{{Key: "foo", Value: "bar"}})
	assert.Equal(t, errDummy, err)
}

func TestCollection_InsertMany(t *testing.T) {
	c := newTestCollection(breaker.GetBreaker("localhost"))
	_, err := c.InsertMany(context.Background(), []any{
		bson.D{{Key: "foo", Value: "bar"}},
		bson.D{{Key: "foo", Value: "baz"}},
	})
	assert.Nil(t, err)
	c.brk = new(dropBreaker)
	_, err = c.InsertMany(context.Background(), []any{bson.D{{Key: "foo", Value: "bar"}}})
	assert.Equal(t, errDummy, err)
}

func TestCollection_DeleteOne(t *testing.T) {
	c := newTestCollection(breaker.GetBreaker("localhost"))
	_, err := c.DeleteOne(context.Background(), bson.D{{Key: "foo", Value: "bar"}})
	assert.Nil(t, err)
	c.brk = new(dropBreaker)
	_, err = c.DeleteOne(context.Background(), bson.D{{Key: "foo", Value: "bar"}})
	assert.Equal(t, errDummy, err)
}

func TestCollection_DeleteMany(t *testing.T) {
	c := newTestCollection(breaker.GetBreaker("localhost"))
	_, err := c.DeleteMany(context.Background(), bson.D{{Key: "foo", Value: "bar"}})
	assert.Nil(t, err)
	c.brk = new(dropBreaker)
	_, err = c.DeleteMany(context.Background(), bson.D{{Key: "foo", Value: "bar"}})
	assert.Equal(t, errDummy, err)
}

func TestCollection_ReplaceOne(t *testing.T) {
	c := newTestCollection(breaker.GetBreaker("localhost"))
	_, err := c.ReplaceOne(context.Background(), bson.D{{Key: "foo", Value: "bar"}},
		bson.D{{Key: "foo", Value: "baz"}},
	)
	assert.Nil(t, err)
	c.brk = new(dropBreaker)
	_, err = c.ReplaceOne(context.Background(), bson.D{{Key: "foo", Value: "bar"}},
		bson.D{{Key: "foo", Value: "baz"}})
	assert.Equal(t, errDummy, err)
}

func TestCollection_UpdateOne(t *testing.T) {
	c := newTestCollection(breaker.GetBreaker("localhost"))
	_, err := c.UpdateOne(context.Background(), bson.D{{Key: "foo", Value: "bar"}},
		bson.D{{Key: "$set", Value: bson.D{{Key: "baz", Value: "qux"}}}})
	assert.Nil(t, err)
	c.brk = new(dropBreaker)
	_, err = c.UpdateOne(context.Background(), bson.D{{Key: "foo", Value: "bar"}},
		bson.D{{Key: "$set", Value: bson.D{{Key: "baz", Value: "qux"}}}})
	assert.Equal(t, errDummy, err)
}

func TestCollection_UpdateByID(t *testing.T) {
	c := newTestCollection(breaker.GetBreaker("localhost"))
	_, err := c.UpdateByID(context.Background(), primitive.NewObjectID(),
		bson.D{{Key: "$set", Value: bson.D{{Key: "baz", Value: "qux"}}}})
	assert.Nil(t, err)
	c.brk = new(dropBreaker)
	_, err = c.UpdateByID(context.Background(), primitive.NewObjectID(),
		bson.D{{Key: "$set", Value: bson.D{{Key: "baz", Value: "qux"}}}})
	assert.Equal(t, errDummy, err)
}

func TestCollection_UpdateMany(t *testing.T) {
	c := newTestCollection(breaker.GetBreaker("localhost"))
	_, err := c.UpdateMany(context.Background(), bson.D{{Key: "foo", Value: "bar"}},
		bson.D{{Key: "$set", Value: bson.D{{Key: "baz", Value: "qux"}}}})
	assert.Nil(t, err)
	c.brk = new(dropBreaker)
	_, err = c.UpdateMany(context.Background(), bson.D{{Key: "foo", Value: "bar"}},
		bson.D{{Key: "$set", Value: bson.D{{Key: "baz", Value: "qux"}}}})
	assert.Equal(t, errDummy, err)
}

func TestCollection_Watch(t *testing.T) {
	c := newTestCollection(breaker.GetBreaker("localhost"))
	_, err := c.Watch(context.Background(), bson.D{{Key: "foo", Value: "bar"}})
	assert.Nil(t, err)
}

func TestCollection_Clone(t *testing.T) {
	c := newTestCollection(breaker.GetBreaker("localhost"))
	c.Clone()

}

func TestCollection_Database(t *testing.T) {
	c := newTestCollection(breaker.GetBreaker("localhost"))
	c.Database()
}

func TestCollection_Drop(t *testing.T) {
	c := newTestCollection(breaker.GetBreaker("localhost"))
	_ = c.Drop(context.Background())
}

func TestCollection_Indexes(t *testing.T) {
	c := newTestCollection(breaker.GetBreaker("localhost"))
	c.Indexes()
}

func TestDecoratedCollection_LogDuration(t *testing.T) {
	c := newTestCollection(breaker.GetBreaker("localhost"))

	buf := logtest.NewCollector(t)

	buf.Reset()
	c.logDuration(context.Background(), "foo", timex.Now(), nil, "bar")
	assert.Contains(t, buf.String(), "foo")
	assert.Contains(t, buf.String(), "bar")

	buf.Reset()
	c.logDuration(context.Background(), "foo", timex.Now(), errors.New("bar"), make(chan int))
	assert.Contains(t, buf.String(), "foo")
	assert.Contains(t, buf.String(), "bar")

	buf.Reset()
	c.logDuration(context.Background(), "foo", timex.Now(), nil, make(chan int))
	assert.Contains(t, buf.String(), "foo")

	buf.Reset()
	c.logDuration(context.Background(), "foo", timex.Now()-slowThreshold.Load()*2,
		nil, make(chan int))
	assert.Contains(t, buf.String(), "foo")
	assert.Contains(t, buf.String(), "slowcall")

	buf.Reset()
	c.logDuration(context.Background(), "foo", timex.Now()-slowThreshold.Load()*2,
		errors.New("bar"), make(chan int))
	assert.Contains(t, buf.String(), "foo")
	assert.Contains(t, buf.String(), "bar")

	buf.Reset()
	c.logDuration(context.Background(), "foo", timex.Now()-slowThreshold.Load()*2,
		errors.New("bar"))
	assert.Contains(t, buf.String(), "foo")

	buf.Reset()
	c.logDuration(context.Background(), "foo", timex.Now()-slowThreshold.Load()*2, nil)
	assert.Contains(t, buf.String(), "foo")
	assert.Contains(t, buf.String(), "slowcall")
}

func TestAcceptable(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"NilError", nil, true},
		{"NoDocuments", mongo.ErrNoDocuments, true},
		{"NilValue", mongo.ErrNilValue, true},
		{"NilDocument", mongo.ErrNilDocument, true},
		{"NilCursor", mongo.ErrNilCursor, true},
		{"EmptySlice", mongo.ErrEmptySlice, true},
		{"DuplicateKeyError", mongo.WriteException{WriteErrors: []mongo.WriteError{{Code: duplicateKeyCode}}}, true},
		{"OtherError", errors.New("other error"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, acceptable(tt.err))
		})
	}
}

func TestIsDupKeyError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"NilError", nil, false},
		{"NonDupKeyError", errors.New("some other error"), false},
		{"DupKeyError", mongo.WriteException{WriteErrors: []mongo.WriteError{{Code: duplicateKeyCode}}}, true},
		{"OtherMongoError", mongo.WriteException{WriteErrors: []mongo.WriteError{{Code: 12345}}}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, isDupKeyError(tt.err))
		})
	}
}

type mockPromise struct {
	accepted bool
	reason   string
}

func (p *mockPromise) Accept() {
	p.accepted = true
}

func (p *mockPromise) Reject(reason string) {
	p.reason = reason
}

type dropBreaker struct{}

func (d *dropBreaker) Name() string {
	return "dummy"
}

func (d *dropBreaker) Allow() (breaker.Promise, error) {
	return nil, errDummy
}

func (d *dropBreaker) AllowCtx(_ context.Context) (breaker.Promise, error) {
	return nil, errDummy
}

func (d *dropBreaker) Do(_ func() error) error {
	return nil
}

func (d *dropBreaker) DoCtx(_ context.Context, _ func() error) error {
	return nil
}

func (d *dropBreaker) DoWithAcceptable(_ func() error, _ breaker.Acceptable) error {
	return errDummy
}

func (d *dropBreaker) DoWithAcceptableCtx(_ context.Context, _ func() error, _ breaker.Acceptable) error {
	return errDummy
}

func (d *dropBreaker) DoWithFallback(_ func() error, _ breaker.Fallback) error {
	return nil
}

func (d *dropBreaker) DoWithFallbackCtx(_ context.Context, _ func() error, _ breaker.Fallback) error {
	return nil
}

func (d *dropBreaker) DoWithFallbackAcceptable(_ func() error, _ breaker.Fallback,
	_ breaker.Acceptable) error {
	return nil
}

func (d *dropBreaker) DoWithFallbackAcceptableCtx(_ context.Context, _ func() error,
	_ breaker.Fallback, _ breaker.Acceptable) error {
	return nil
}
