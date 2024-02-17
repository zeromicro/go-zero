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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	mopt "go.mongodb.org/mongo-driver/mongo/options"
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
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		coll := mt.Coll
		assert.NotNil(t, coll)
		col := newCollection(coll, breaker.GetBreaker("localhost"))
		assert.Equal(t, t.Name()+"/test", col.(*decoratedCollection).name)
	})
}

func TestCollection_Aggregate(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
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
	mt.Run("test", func(mt *mtest.T) {
		c := decoratedCollection{
			Collection: mt.Coll,
			brk:        breaker.NewBreaker(),
		}
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{{Key: "ok", Value: 1}}...))
		res, err := c.BulkWrite(context.Background(), []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(bson.D{{Key: "foo", Value: 1}}),
		})
		assert.Nil(t, err)
		assert.NotNil(t, res)
		c.brk = new(dropBreaker)
		_, err = c.BulkWrite(context.Background(), []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(bson.D{{Key: "foo", Value: 1}}),
		})
		assert.Equal(t, errDummy, err)
	})
}

func TestCollection_CountDocuments(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
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
				{Key: "n", Value: 1},
			}))
		res, err := c.CountDocuments(context.Background(), bson.D{})
		assert.Nil(t, err)
		assert.Equal(t, int64(1), res)

		c.brk = new(dropBreaker)
		_, err = c.CountDocuments(context.Background(), bson.D{{Key: "foo", Value: 1}})
		assert.Equal(t, errDummy, err)
	})
}

func TestDecoratedCollection_DeleteMany(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		c := decoratedCollection{
			Collection: mt.Coll,
			brk:        breaker.NewBreaker(),
		}
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{{Key: "n", Value: 1}}...))
		res, err := c.DeleteMany(context.Background(), bson.D{})
		assert.Nil(t, err)
		assert.Equal(t, int64(1), res.DeletedCount)

		c.brk = new(dropBreaker)
		_, err = c.DeleteMany(context.Background(), bson.D{{Key: "foo", Value: 1}})
		assert.Equal(t, errDummy, err)
	})
}

func TestCollection_Distinct(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		c := decoratedCollection{
			Collection: mt.Coll,
			brk:        breaker.NewBreaker(),
		}
		mt.AddMockResponses(bson.D{{Key: "ok", Value: 1}, {Key: "values", Value: []int{1}}})
		resp, err := c.Distinct(context.Background(), "foo", bson.D{})
		assert.Nil(t, err)
		assert.Equal(t, 1, len(resp))

		c.brk = new(dropBreaker)
		_, err = c.Distinct(context.Background(), "foo", bson.D{{Key: "foo", Value: 1}})
		assert.Equal(t, errDummy, err)
	})
}

func TestCollection_EstimatedDocumentCount(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		c := decoratedCollection{
			Collection: mt.Coll,
			brk:        breaker.NewBreaker(),
		}
		mt.AddMockResponses(bson.D{{Key: "ok", Value: 1}, {Key: "n", Value: 1}})
		res, err := c.EstimatedDocumentCount(context.Background())
		assert.Nil(t, err)
		assert.Equal(t, int64(1), res)

		c.brk = new(dropBreaker)
		_, err = c.EstimatedDocumentCount(context.Background())
		assert.Equal(t, errDummy, err)
	})
}

func TestCollection_Find(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
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
				{Key: "name", Value: "John"},
			})
		getMore := mtest.CreateCursorResponse(
			1,
			"DBName.CollectionName",
			mtest.NextBatch,
			bson.D{
				{Key: "name", Value: "Mary"},
			})
		killCursors := mtest.CreateCursorResponse(
			0,
			"DBName.CollectionName",
			mtest.NextBatch)
		mt.AddMockResponses(find, getMore, killCursors)
		filter := bson.D{{Key: "x", Value: 1}}
		cursor, err := c.Find(context.Background(), filter, mopt.Find())
		assert.Nil(t, err)
		defer cursor.Close(context.Background())

		var val []struct {
			ID   primitive.ObjectID `bson:"_id"`
			Name string             `bson:"name"`
		}
		assert.Nil(t, cursor.All(context.Background(), &val))
		assert.Equal(t, 2, len(val))
		assert.Equal(t, "John", val[0].Name)
		assert.Equal(t, "Mary", val[1].Name)

		c.brk = new(dropBreaker)
		_, err = c.Find(context.Background(), filter, mopt.Find())
		assert.Equal(t, errDummy, err)
	})
}

func TestCollection_FindOne(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
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
				{Key: "name", Value: "John"},
			})
		getMore := mtest.CreateCursorResponse(
			1,
			"DBName.CollectionName",
			mtest.NextBatch,
			bson.D{
				{Key: "name", Value: "Mary"},
			})
		killCursors := mtest.CreateCursorResponse(
			0,
			"DBName.CollectionName",
			mtest.NextBatch)
		mt.AddMockResponses(find, getMore, killCursors)
		filter := bson.D{{Key: "x", Value: 1}}
		resp, err := c.FindOne(context.Background(), filter)
		assert.Nil(t, err)
		var val struct {
			ID   primitive.ObjectID `bson:"_id"`
			Name string             `bson:"name"`
		}
		assert.Nil(t, resp.Decode(&val))
		assert.Equal(t, "John", val.Name)

		c.brk = new(dropBreaker)
		_, err = c.FindOne(context.Background(), filter)
		assert.Equal(t, errDummy, err)
	})
}

func TestCollection_FindOneAndDelete(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		c := decoratedCollection{
			Collection: mt.Coll,
			brk:        breaker.NewBreaker(),
		}
		filter := bson.D{}
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{}...))
		_, err := c.FindOneAndDelete(context.Background(), filter, mopt.FindOneAndDelete())
		assert.Equal(t, mongo.ErrNoDocuments, err)

		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{
			{Key: "value", Value: bson.D{{Key: "name", Value: "John"}}},
		}...))
		resp, err := c.FindOneAndDelete(context.Background(), filter, mopt.FindOneAndDelete())
		assert.Nil(t, err)
		var val struct {
			Name string `bson:"name"`
		}
		assert.Nil(t, resp.Decode(&val))
		assert.Equal(t, "John", val.Name)

		c.brk = new(dropBreaker)
		_, err = c.FindOneAndDelete(context.Background(), bson.D{{Key: "foo", Value: "bar"}})
		assert.Equal(t, errDummy, err)
	})
}

func TestCollection_FindOneAndReplace(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		c := decoratedCollection{
			Collection: mt.Coll,
			brk:        breaker.NewBreaker(),
		}
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{}...))
		filter := bson.D{{Key: "x", Value: 1}}
		replacement := bson.D{{Key: "x", Value: 2}}
		opts := mopt.FindOneAndReplace().SetUpsert(true)
		_, err := c.FindOneAndReplace(context.Background(), filter, replacement, opts)
		assert.Equal(t, mongo.ErrNoDocuments, err)
		mt.AddMockResponses(bson.D{{Key: "ok", Value: 1}, {Key: "value", Value: bson.D{
			{Key: "name", Value: "John"},
		}}})
		resp, err := c.FindOneAndReplace(context.Background(), filter, replacement, opts)
		assert.Nil(t, err)
		var val struct {
			Name string `bson:"name"`
		}
		assert.Nil(t, resp.Decode(&val))
		assert.Equal(t, "John", val.Name)

		c.brk = new(dropBreaker)
		_, err = c.FindOneAndReplace(context.Background(), filter, replacement, opts)
		assert.Equal(t, errDummy, err)
	})
}

func TestCollection_FindOneAndUpdate(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		c := decoratedCollection{
			Collection: mt.Coll,
			brk:        breaker.NewBreaker(),
		}
		mt.AddMockResponses(bson.D{{Key: "ok", Value: 1}})
		filter := bson.D{{Key: "x", Value: 1}}
		update := bson.D{{Key: "$x", Value: 2}}
		opts := mopt.FindOneAndUpdate().SetUpsert(true)
		_, err := c.FindOneAndUpdate(context.Background(), filter, update, opts)
		assert.Equal(t, mongo.ErrNoDocuments, err)

		mt.AddMockResponses(bson.D{{Key: "ok", Value: 1}, {Key: "value", Value: bson.D{
			{Key: "name", Value: "John"},
		}}})
		resp, err := c.FindOneAndUpdate(context.Background(), filter, update, opts)
		assert.Nil(t, err)
		var val struct {
			Name string `bson:"name"`
		}
		assert.Nil(t, resp.Decode(&val))
		assert.Equal(t, "John", val.Name)

		c.brk = new(dropBreaker)
		_, err = c.FindOneAndUpdate(context.Background(), filter, update, opts)
		assert.Equal(t, errDummy, err)
	})
}

func TestCollection_InsertOne(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		c := decoratedCollection{
			Collection: mt.Coll,
			brk:        breaker.NewBreaker(),
		}
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{{Key: "ok", Value: 1}}...))
		res, err := c.InsertOne(context.Background(), bson.D{{Key: "foo", Value: "bar"}})
		assert.Nil(t, err)
		assert.NotNil(t, res)

		c.brk = new(dropBreaker)
		_, err = c.InsertOne(context.Background(), bson.D{{Key: "foo", Value: "bar"}})
		assert.Equal(t, errDummy, err)
	})
}

func TestCollection_InsertMany(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		c := decoratedCollection{
			Collection: mt.Coll,
			brk:        breaker.NewBreaker(),
		}
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{{Key: "ok", Value: 1}}...))
		res, err := c.InsertMany(context.Background(), []any{
			bson.D{{Key: "foo", Value: "bar"}},
			bson.D{{Key: "foo", Value: "baz"}},
		})
		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 2, len(res.InsertedIDs))

		c.brk = new(dropBreaker)
		_, err = c.InsertMany(context.Background(), []any{bson.D{{Key: "foo", Value: "bar"}}})
		assert.Equal(t, errDummy, err)
	})
}

func TestCollection_DeleteOne(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		c := decoratedCollection{
			Collection: mt.Coll,
			brk:        breaker.NewBreaker(),
		}
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{{Key: "n", Value: 1}}...))
		res, err := c.DeleteOne(context.Background(), bson.D{{Key: "foo", Value: "bar"}})
		assert.Nil(t, err)
		assert.Equal(t, int64(1), res.DeletedCount)

		c.brk = new(dropBreaker)
		_, err = c.DeleteOne(context.Background(), bson.D{{Key: "foo", Value: "bar"}})
		assert.Equal(t, errDummy, err)
	})
}

func TestCollection_DeleteMany(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		c := decoratedCollection{
			Collection: mt.Coll,
			brk:        breaker.NewBreaker(),
		}
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{{Key: "n", Value: 1}}...))
		res, err := c.DeleteMany(context.Background(), bson.D{{Key: "foo", Value: "bar"}})
		assert.Nil(t, err)
		assert.Equal(t, int64(1), res.DeletedCount)

		c.brk = new(dropBreaker)
		_, err = c.DeleteMany(context.Background(), bson.D{{Key: "foo", Value: "bar"}})
		assert.Equal(t, errDummy, err)
	})
}

func TestCollection_ReplaceOne(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		c := decoratedCollection{
			Collection: mt.Coll,
			brk:        breaker.NewBreaker(),
		}
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{{Key: "n", Value: 1}}...))
		res, err := c.ReplaceOne(context.Background(), bson.D{{Key: "foo", Value: "bar"}},
			bson.D{{Key: "foo", Value: "baz"}},
		)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), res.MatchedCount)

		c.brk = new(dropBreaker)
		_, err = c.ReplaceOne(context.Background(), bson.D{{Key: "foo", Value: "bar"}},
			bson.D{{Key: "foo", Value: "baz"}})
		assert.Equal(t, errDummy, err)
	})
}

func TestCollection_UpdateOne(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		c := decoratedCollection{
			Collection: mt.Coll,
			brk:        breaker.NewBreaker(),
		}
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{{Key: "n", Value: 1}}...))
		resp, err := c.UpdateOne(context.Background(), bson.D{{Key: "foo", Value: "bar"}},
			bson.D{{Key: "$set", Value: bson.D{{Key: "baz", Value: "qux"}}}})
		assert.Nil(t, err)
		assert.Equal(t, int64(1), resp.MatchedCount)

		c.brk = new(dropBreaker)
		_, err = c.UpdateOne(context.Background(), bson.D{{Key: "foo", Value: "bar"}},
			bson.D{{Key: "$set", Value: bson.D{{Key: "baz", Value: "qux"}}}})
		assert.Equal(t, errDummy, err)
	})
}

func TestCollection_UpdateByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		c := decoratedCollection{
			Collection: mt.Coll,
			brk:        breaker.NewBreaker(),
		}
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{{Key: "n", Value: 1}}...))
		resp, err := c.UpdateByID(context.Background(), primitive.NewObjectID(),
			bson.D{{Key: "$set", Value: bson.D{{Key: "baz", Value: "qux"}}}})
		assert.Nil(t, err)
		assert.Equal(t, int64(1), resp.MatchedCount)

		c.brk = new(dropBreaker)
		_, err = c.UpdateByID(context.Background(), primitive.NewObjectID(),
			bson.D{{Key: "$set", Value: bson.D{{Key: "baz", Value: "qux"}}}})
		assert.Equal(t, errDummy, err)
	})
}

func TestCollection_UpdateMany(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		c := decoratedCollection{
			Collection: mt.Coll,
			brk:        breaker.NewBreaker(),
		}
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{{Key: "n", Value: 1}}...))
		resp, err := c.UpdateMany(context.Background(), bson.D{{Key: "foo", Value: "bar"}},
			bson.D{{Key: "$set", Value: bson.D{{Key: "baz", Value: "qux"}}}})
		assert.Nil(t, err)
		assert.Equal(t, int64(1), resp.MatchedCount)

		c.brk = new(dropBreaker)
		_, err = c.UpdateMany(context.Background(), bson.D{{Key: "foo", Value: "bar"}},
			bson.D{{Key: "$set", Value: bson.D{{Key: "baz", Value: "qux"}}}})
		assert.Equal(t, errDummy, err)
	})
}

func TestDecoratedCollection_LogDuration(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	c := decoratedCollection{
		Collection: mt.Coll,
		brk:        breaker.NewBreaker(),
	}

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

func (d *dropBreaker) Do(_ func() error) error {
	return nil
}

func (d *dropBreaker) DoWithAcceptable(_ func() error, _ breaker.Acceptable) error {
	return errDummy
}

func (d *dropBreaker) DoWithFallback(_ func() error, _ breaker.Fallback) error {
	return nil
}

func (d *dropBreaker) DoWithFallbackAcceptable(_ func() error, _ breaker.Fallback,
	_ breaker.Acceptable) error {
	return nil
}
