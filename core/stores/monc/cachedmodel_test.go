package monc

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stores/mon"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestModel_DelCache(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("test", func(mt *mtest.T) {
		m := createModel(t, mt)
		assert.Nil(t, m.cache.Set("foo", "bar"))
		assert.Nil(t, m.cache.Set("bar", "baz"))
		assert.Nil(t, m.DelCache("foo", "bar"))
		var v string
		assert.True(t, m.cache.IsNotFound(m.cache.Get("foo", &v)))
		assert.True(t, m.cache.IsNotFound(m.cache.Get("bar", &v)))
	})
}

func TestModel_DeleteOne(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("test", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{{Key: "n", Value: 1}}...))
		m := createModel(t, mt)
		assert.Nil(t, m.cache.Set("foo", "bar"))
		val, err := m.DeleteOne(context.Background(), "foo", bson.D{{Key: "foo", Value: "bar"}})
		assert.Nil(t, err)
		assert.Equal(t, int64(1), val)
		var v string
		assert.True(t, m.cache.IsNotFound(m.cache.Get("foo", &v)))
	})
}

func TestModel_DeleteOneNoCache(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("test", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{{Key: "n", Value: 1}}...))
		m := createModel(t, mt)
		assert.Nil(t, m.cache.Set("foo", "bar"))
		val, err := m.DeleteOneNoCache(context.Background(), bson.D{{Key: "foo", Value: "bar"}})
		assert.Nil(t, err)
		assert.Equal(t, int64(1), val)
		var v string
		assert.Nil(t, m.cache.Get("foo", &v))
	})
}

func TestModel_FindOne(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("test", func(mt *mtest.T) {
		// not need to add mock response, because it will be returned from cache.
		m := createModel(t, mt)
		assert.Nil(t, m.cache.Set("foo", "bar"))
		var v string
		assert.Nil(t, m.FindOne(context.Background(), "foo", &v, bson.D{}))
		assert.Equal(t, "bar", v)
	})
}

func TestModel_FindOneNoCache(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("test", func(mt *mtest.T) {
		resp := mtest.CreateCursorResponse(
			1,
			"DBName.CollectionName",
			mtest.FirstBatch,
			bson.D{
				{Key: "_id", Value: primitive.NewObjectID()},
				{Key: "foo", Value: "bar"},
			})
		mt.AddMockResponses(resp)
		m := createModel(t, mt)
		v := struct {
			Foo string `bson:"foo"`
		}{}
		assert.Nil(t, m.FindOneNoCache(context.Background(), &v, bson.D{}))
		assert.Equal(t, "bar", v.Foo)
	})
}

func TestModel_FindOneAndDelete(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("test", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{
			{Key: "value", Value: bson.D{{Key: "foo", Value: "bar"}}},
		}...))
		m := createModel(t, mt)
		assert.Nil(t, m.cache.Set("foo", "bar"))
		v := struct {
			Foo string `bson:"foo"`
		}{}
		assert.Nil(t, m.FindOneAndDelete(context.Background(), "foo", &v, bson.D{}))
		assert.Equal(t, "bar", v.Foo)
		assert.True(t, m.cache.IsNotFound(m.cache.Get("foo", &v)))
	})
}

func TestModel_FindOneAndDeleteNoCache(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("test", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{
			{Key: "value", Value: bson.D{{Key: "foo", Value: "bar"}}},
		}...))
		m := createModel(t, mt)
		v := struct {
			Foo string `bson:"foo"`
		}{}
		assert.Nil(t, m.FindOneAndDeleteNoCache(context.Background(), &v, bson.D{}))
		assert.Equal(t, "bar", v.Foo)
	})
}

func TestModel_FindOneAndReplace(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("test", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{
			{Key: "value", Value: bson.D{{Key: "foo", Value: "bar"}}},
		}...))
		m := createModel(t, mt)
		assert.Nil(t, m.cache.Set("foo", "bar"))
		v := struct {
			Foo string `bson:"foo"`
		}{}
		assert.Nil(t, m.FindOneAndReplace(context.Background(), "foo", &v, bson.D{}, bson.D{
			{Key: "name", Value: "Mary"},
		}))
		assert.Equal(t, "bar", v.Foo)
		assert.True(t, m.cache.IsNotFound(m.cache.Get("foo", &v)))
	})
}

func createModel(t *testing.T, mt *mtest.T) *Model {
	s, err := miniredis.Run()
	assert.Nil(t, err)
	mon.Inject(mt.Name(), mt.Client)
	return MustNewNodeModel(mt.Name(), mt.DB.Name(), mt.Coll.Name(), redis.New(s.Addr()))
}
