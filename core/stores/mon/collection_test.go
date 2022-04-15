package mon

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stringx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	mopt "go.mongodb.org/mongo-driver/mongo/options"
)

var errDummy = errors.New("dummy")

func init() {
	logx.Disable()
}

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
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	mt.Run("test", func(mt *mtest.T) {
		coll := mt.Coll
		assert.NotNil(t, coll)
		col := newCollection(coll, breaker.GetBreaker("localhost"))
		assert.Equal(t, t.Name()+"/test", col.(*decoratedCollection).name)
	})
}

func TestCollection_Aggregate(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	mt.Run("test", func(mt *mtest.T) {
		coll := mt.Coll
		assert.NotNil(t, coll)
		col := newCollection(coll, breaker.GetBreaker("localhost"))
		ns := mt.Coll.Database().Name() + "." + mt.Coll.Name()
		aggRes := mtest.CreateCursorResponse(1, ns, mtest.FirstBatch)
		mt.AddMockResponses(aggRes)
		assert.Equal(t, t.Name()+"/test", col.(*decoratedCollection).name)
		cursor, err := col.Aggregate(context.Background(), mongo.Pipeline{}, mopt.Aggregate())
		assert.Nil(t, err)
		cursor.Close(context.Background())
	})
}

func TestCollection_BulkWrite(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("test", func(mt *mtest.T) {
		c := decoratedCollection{
			Collection: mt.Coll,
			brk:        breaker.NewBreaker(),
		}
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{{"ok", 1}}...))
		res, err := c.BulkWrite(context.Background(), []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(bson.D{{"foo", 1}})},
		)
		assert.Nil(t, err)
		assert.NotNil(t, res)
		c.brk = new(dropBreaker)
		_, err = c.BulkWrite(context.Background(), []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(bson.D{{"foo", 1}})},
		)
		assert.Equal(t, errDummy, err)
	})
}

func TestCollection_CountDocuments(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("test", func(mt *mtest.T) {
		c := decoratedCollection{
			Collection: mt.Coll,
			brk:        breaker.NewBreaker(),
		}
		mt.AddMockResponses(mtest.CreateCursorResponse(
			1,
			"DBName.CollectionName",
			mtest.FirstBatch,
			bson.D{
				{"_id", primitive.NewObjectID()},
				{"n", 1},
			}))
		res, err := c.CountDocuments(context.Background(), bson.D{})
		assert.Nil(t, err)
		assert.Equal(t, int64(1), res)
		c.brk = new(dropBreaker)
		_, err = c.CountDocuments(context.Background(), bson.D{{"foo", 1}})
		assert.Equal(t, errDummy, err)
	})
}

func TestCollection_Distinct(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("test", func(mt *mtest.T) {
		c := decoratedCollection{
			Collection: mt.Coll,
			brk:        breaker.NewBreaker(),
		}
		ns := mt.Coll.Database().Name() + "." + mt.Coll.Name()
		res := mtest.CreateCursorResponse(1, ns, mtest.FirstBatch, bson.D{{"_id", 1}},
			bson.D{{"_id", 2}})
		mt.AddMockResponses(res)
		resp, err := c.Distinct(context.Background(), "foo", bson.D{})
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		c.brk = new(dropBreaker)
		_, err = c.Distinct(context.Background(), "foo", bson.D{{"foo", 1}})
		assert.Equal(t, errDummy, err)
	})
}

func TestCollection_DistinctMore(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("distinct", func(mt *mtest.T) {
		all := []interface{}{int32(1), int32(2), int32(3), int32(4), int32(5)}
		last3 := []interface{}{int32(3), int32(4), int32(5)}
		testCases := []struct {
			name     string
			filter   bson.D
			opts     *mopt.DistinctOptions
			expected []interface{}
		}{
			{"no options", bson.D{}, nil, all},
			{"filter", bson.D{{"x", bson.D{{"$gt", 2}}}}, nil, last3},
			{"options", bson.D{}, mopt.Distinct().SetMaxTime(5000000000), all},
		}
		for _, tc := range testCases {
			mt.Run(tc.name, func(mt *mtest.T) {
				var docs []bson.D
				for i := 1; i <= 5; i++ {
					docs = append(docs, bson.D{{"x", int32(i)}})
				}

				// _, err := mt.Coll.InsertMany(context.Background(), docs)
				// assert.Nil(mt, err, "InsertMany error for initial data: %v", err)

				// mt.AddMockResponses(mtest.CreateCursorResponse(
				// 	1,
				// 	mt.Coll.Database().Name()+"."+mt.Coll.Name(),
				// 	mtest.FirstBatch,
				// 	docs...,
				// ))
				mt.AddMockResponses(docs...)
				res, err := mt.Coll.Distinct(context.Background(), "x", tc.filter, tc.opts)
				assert.Nil(mt, err, "Distinct error: %v", err)
				assert.Equal(mt, tc.expected, res, "expected result %v, got %v", tc.expected, res)
			})
		}
	})
}

func TestCollectionFind(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("test", func(mt *mtest.T) {
		c := decoratedCollection{
			Collection: mt.Coll,
			brk:        breaker.NewBreaker(),
		}
		find := mtest.CreateCursorResponse(
			1,
			"DBName.CollectionName",
			mtest.FirstBatch,
			bson.D{
				{"_id", primitive.NewObjectID()},
				{"name", "John"},
			})
		getMore := mtest.CreateCursorResponse(
			1,
			"DBName.CollectionName",
			mtest.NextBatch,
			bson.D{
				{"_id", primitive.NewObjectID()},
				{"name", "Mary"},
			})
		killCursors := mtest.CreateCursorResponse(
			0,
			"DBName.CollectionName",
			mtest.NextBatch)
		mt.AddMockResponses(find, getMore, killCursors)
		filter := bson.D{{"x", 1}}
		cursor, err := c.Find(context.Background(), filter, mopt.Find())
		assert.Nil(t, err)
		assert.NotNil(t, cursor)
		cursor.Close(context.Background())
		c.brk = new(dropBreaker)
		_, err = c.Find(context.Background(), filter, mopt.Find())
		assert.Equal(t, errDummy, err)
	})
}

func TestCollectionFindOne(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("test", func(mt *mtest.T) {
		c := decoratedCollection{
			Collection: mt.Coll,
			brk:        breaker.NewBreaker(),
		}
		find := mtest.CreateCursorResponse(
			1,
			"DBName.CollectionName",
			mtest.FirstBatch,
			bson.D{
				{"_id", primitive.NewObjectID()},
				{"name", "John"},
			})
		getMore := mtest.CreateCursorResponse(
			1,
			"DBName.CollectionName",
			mtest.NextBatch,
			bson.D{
				{"_id", primitive.NewObjectID()},
				{"name", "Mary"},
			})
		killCursors := mtest.CreateCursorResponse(
			0,
			"DBName.CollectionName",
			mtest.NextBatch)
		mt.AddMockResponses(find, getMore, killCursors)
		filter := bson.D{{"x", 1}}
		resp, err := c.FindOne(context.Background(), filter)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		c.brk = new(dropBreaker)
		_, err = c.FindOne(context.Background(), filter)
		assert.Equal(t, errDummy, err)
	})
}

func TestCollection_FindOneAndReplace(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("test", func(mt *mtest.T) {
		c := decoratedCollection{
			Collection: mt.Coll,
			brk:        breaker.NewBreaker(),
		}
		find := mtest.CreateCursorResponse(
			1,
			"DBName.CollectionName",
			mtest.FirstBatch,
			bson.D{
				{"_id", primitive.NewObjectID()},
				{"name", "John"},
			})
		getMore := mtest.CreateCursorResponse(
			1,
			"DBName.CollectionName",
			mtest.NextBatch,
			bson.D{
				{"_id", primitive.NewObjectID()},
				{"name", "Mary"},
			})
		killCursors := mtest.CreateCursorResponse(
			0,
			"DBName.CollectionName",
			mtest.NextBatch)
		mt.AddMockResponses(find, getMore, killCursors)
		filter := bson.D{{"x", 1}}
		replacement := bson.D{{"name", "foo"}}
		opts := mopt.FindOneAndReplace().SetUpsert(true)
		resp, err := c.FindOneAndReplace(context.Background(), filter, replacement, opts)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		c.brk = new(dropBreaker)
		_, err = c.FindOneAndReplace(context.Background(), filter, replacement, opts)
		assert.Equal(t, errDummy, err)
	})
}

func TestCollectionInsert(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("test", func(mt *mtest.T) {
		c := decoratedCollection{
			Collection: mt.Coll,
			brk:        breaker.NewBreaker(),
		}
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{{"ok", 1}}...))
		res, err := c.InsertOne(context.Background(), bson.D{{"foo", "bar"}})
		assert.Nil(t, err)
		assert.NotNil(t, res)
		c.brk = new(dropBreaker)
		_, err = c.InsertOne(context.Background(), bson.D{{"foo", "bar"}})
		assert.Equal(t, errDummy, err)
	})
}

func TestCollectionRemove(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("test", func(mt *mtest.T) {
		c := decoratedCollection{
			Collection: mt.Coll,
			brk:        breaker.NewBreaker(),
		}
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{{"ok", 1}}...))
		res, err := c.DeleteOne(context.Background(), bson.D{{"foo", "bar"}})
		assert.Nil(t, err)
		assert.NotNil(t, res)
		c.brk = new(dropBreaker)
		_, err = c.DeleteOne(context.Background(), bson.D{{"foo", "bar"}})
		assert.Equal(t, errDummy, err)
	})
}

func TestCollectionRemoveAll(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("test", func(mt *mtest.T) {
		c := decoratedCollection{
			Collection: mt.Coll,
			brk:        breaker.NewBreaker(),
		}
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{{"ok", 1}}...))
		res, err := c.DeleteMany(context.Background(), bson.D{{"foo", "bar"}})
		assert.Nil(t, err)
		assert.NotNil(t, res)
		c.brk = new(dropBreaker)
		_, err = c.DeleteMany(context.Background(), bson.D{{"foo", "bar"}})
		assert.Equal(t, errDummy, err)
	})
}

func TestCollection_FindOneAndDelete(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("test", func(mt *mtest.T) {
		c := decoratedCollection{
			Collection: mt.Coll,
			brk:        breaker.NewBreaker(),
		}
		filter := bson.D{{"x", 1}}
		ns := mt.Coll.Database().Name() + "." + mt.Coll.Name()
		aggRes := mtest.CreateCursorResponse(1, ns, mtest.FirstBatch)
		mt.AddMockResponses(aggRes)
		res, err := c.FindOneAndDelete(context.Background(), filter, mopt.FindOneAndDelete())
		assert.Equal(t, mongo.ErrNoDocuments, err)
		assert.NotNil(t, res)
		c.brk = new(dropBreaker)
		_, err = c.FindOneAndDelete(context.Background(), bson.D{{"foo", "bar"}})
		assert.Equal(t, errDummy, err)
	})
}

func TestCollectionUpdate(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("test", func(mt *mtest.T) {
		c := decoratedCollection{
			Collection: mt.Coll,
			brk:        breaker.NewBreaker(),
		}
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{{"ok", 1}}...))
		resp, err := c.UpdateOne(context.Background(), bson.D{{"foo", "bar"}},
			bson.D{{"$set", bson.D{{"baz", "qux"}}}})
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		c.brk = new(dropBreaker)
		_, err = c.UpdateOne(context.Background(), bson.D{{"foo", "bar"}},
			bson.D{{"$set", bson.D{{"baz", "qux"}}}})
		assert.Equal(t, errDummy, err)
	})
}

func TestCollectionUpdateId(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("test", func(mt *mtest.T) {
		c := decoratedCollection{
			Collection: mt.Coll,
			brk:        breaker.NewBreaker(),
		}
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{{"ok", 1}}...))
		resp, err := c.UpdateByID(context.Background(), primitive.NewObjectID(),
			bson.D{{"$set", bson.D{{"baz", "qux"}}}})
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		c.brk = new(dropBreaker)
		_, err = c.UpdateByID(context.Background(), primitive.NewObjectID(),
			bson.D{{"$set", bson.D{{"baz", "qux"}}}})
		assert.Equal(t, errDummy, err)
	})
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

func (d *dropBreaker) Do(req func() error) error {
	return nil
}

func (d *dropBreaker) DoWithAcceptable(req func() error, acceptable breaker.Acceptable) error {
	return errDummy
}

func (d *dropBreaker) DoWithFallback(req func() error, fallback func(err error) error) error {
	return nil
}

func (d *dropBreaker) DoWithFallbackAcceptable(_ func() error, _ func(err error) error,
	_ breaker.Acceptable) error {
	return nil
}
