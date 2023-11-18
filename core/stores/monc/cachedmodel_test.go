package monc

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/mon"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestNewModel(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		_, err := newModel("foo", mt.DB.Name(), mt.Coll.Name(), nil)
		assert.NotNil(mt, err)
	})
}

func TestModel_DelCache(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		m := createModel(t, mt)
		assert.Nil(t, m.cache.Set("foo", "bar"))
		assert.Nil(t, m.cache.Set("bar", "baz"))
		assert.Nil(t, m.DelCache(context.Background(), "foo", "bar"))
		var v string
		assert.True(t, m.cache.IsNotFound(m.cache.Get("foo", &v)))
		assert.True(t, m.cache.IsNotFound(m.cache.Get("bar", &v)))
	})
}

func TestModel_DeleteOne(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{{Key: "n", Value: 1}}...))
		m := createModel(t, mt)
		assert.Nil(t, m.cache.Set("foo", "bar"))
		val, err := m.DeleteOne(context.Background(), "foo", bson.D{{Key: "foo", Value: "bar"}})
		assert.Nil(t, err)
		assert.Equal(t, int64(1), val)
		var v string
		assert.True(t, m.cache.IsNotFound(m.cache.Get("foo", &v)))
		_, err = m.DeleteOne(context.Background(), "foo", bson.D{{Key: "foo", Value: "bar"}})
		assert.NotNil(t, err)

		m.cache = mockedCache{m.cache}
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{{Key: "n", Value: 1}}...))
		_, err = m.DeleteOne(context.Background(), "foo", bson.D{{Key: "foo", Value: "bar"}})
		assert.Equal(t, errMocked, err)
	})
}

func TestModel_DeleteOneNoCache(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
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
	mt.Run("test", func(mt *mtest.T) {
		resp := mtest.CreateCursorResponse(
			1,
			"DBName.CollectionName",
			mtest.FirstBatch,
			bson.D{
				{Key: "foo", Value: "bar"},
			})
		mt.AddMockResponses(resp)
		m := createModel(t, mt)
		var v struct {
			Foo string `bson:"foo"`
		}
		assert.Nil(t, m.FindOne(context.Background(), "foo", &v, bson.D{}))
		assert.Equal(t, "bar", v.Foo)
		assert.Nil(t, m.cache.Set("foo", "bar"))
	})
}

func TestModel_FindOneNoCache(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		resp := mtest.CreateCursorResponse(
			1,
			"DBName.CollectionName",
			mtest.FirstBatch,
			bson.D{
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
		assert.NotNil(t, m.FindOneAndDelete(context.Background(), "foo", &v, bson.D{}))

		m.cache = mockedCache{m.cache}
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{
			{Key: "value", Value: bson.D{{Key: "foo", Value: "bar"}}},
		}...))
		assert.Equal(t, errMocked, m.FindOneAndDelete(context.Background(), "foo", &v, bson.D{}))
	})
}

func TestModel_FindOneAndDeleteNoCache(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
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
		assert.NotNil(t, m.FindOneAndReplace(context.Background(), "foo", &v, bson.D{}, bson.D{
			{Key: "name", Value: "Mary"},
		}))

		m.cache = mockedCache{m.cache}
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{
			{Key: "value", Value: bson.D{{Key: "foo", Value: "bar"}}},
		}...))
		assert.Equal(t, errMocked, m.FindOneAndReplace(context.Background(), "foo", &v, bson.D{}, bson.D{
			{Key: "name", Value: "Mary"},
		}))
	})
}

func TestModel_FindOneAndReplaceNoCache(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{
			{Key: "value", Value: bson.D{{Key: "foo", Value: "bar"}}},
		}...))
		m := createModel(t, mt)
		v := struct {
			Foo string `bson:"foo"`
		}{}
		assert.Nil(t, m.FindOneAndReplaceNoCache(context.Background(), &v, bson.D{}, bson.D{
			{Key: "name", Value: "Mary"},
		}))
		assert.Equal(t, "bar", v.Foo)
	})
}

func TestModel_FindOneAndUpdate(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{
			{Key: "value", Value: bson.D{{Key: "foo", Value: "bar"}}},
		}...))
		m := createModel(t, mt)
		assert.Nil(t, m.cache.Set("foo", "bar"))
		v := struct {
			Foo string `bson:"foo"`
		}{}
		assert.Nil(t, m.FindOneAndUpdate(context.Background(), "foo", &v, bson.D{}, bson.D{
			{Key: "$set", Value: bson.D{{Key: "name", Value: "Mary"}}},
		}))
		assert.Equal(t, "bar", v.Foo)
		assert.True(t, m.cache.IsNotFound(m.cache.Get("foo", &v)))
		assert.NotNil(t, m.FindOneAndUpdate(context.Background(), "foo", &v, bson.D{}, bson.D{
			{Key: "$set", Value: bson.D{{Key: "name", Value: "Mary"}}},
		}))

		m.cache = mockedCache{m.cache}
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{
			{Key: "value", Value: bson.D{{Key: "foo", Value: "bar"}}},
		}...))
		assert.Equal(t, errMocked, m.FindOneAndUpdate(context.Background(), "foo", &v, bson.D{}, bson.D{
			{Key: "$set", Value: bson.D{{Key: "name", Value: "Mary"}}},
		}))
	})
}

func TestModel_FindOneAndUpdateNoCache(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{
			{Key: "value", Value: bson.D{{Key: "foo", Value: "bar"}}},
		}...))
		m := createModel(t, mt)
		v := struct {
			Foo string `bson:"foo"`
		}{}
		assert.Nil(t, m.FindOneAndUpdateNoCache(context.Background(), &v, bson.D{}, bson.D{
			{Key: "$set", Value: bson.D{{Key: "name", Value: "Mary"}}},
		}))
		assert.Equal(t, "bar", v.Foo)
	})
}

func TestModel_GetCache(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		m := createModel(t, mt)
		assert.NotNil(t, m.cache)
		assert.Nil(t, m.cache.Set("foo", "bar"))
		var s string
		assert.Nil(t, m.cache.Get("foo", &s))
		assert.Equal(t, "bar", s)
	})
}

func TestModel_InsertOne(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{
			{Key: "value", Value: bson.D{{Key: "foo", Value: "bar"}}},
		}...))
		m := createModel(t, mt)
		assert.Nil(t, m.cache.Set("foo", "bar"))
		resp, err := m.InsertOne(context.Background(), "foo", bson.D{
			{Key: "name", Value: "Mary"},
		})
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		var v string
		assert.True(t, m.cache.IsNotFound(m.cache.Get("foo", &v)))
		_, err = m.InsertOne(context.Background(), "foo", bson.D{
			{Key: "name", Value: "Mary"},
		})
		assert.NotNil(t, err)

		m.cache = mockedCache{m.cache}
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{
			{Key: "value", Value: bson.D{{Key: "foo", Value: "bar"}}},
		}...))
		_, err = m.InsertOne(context.Background(), "foo", bson.D{
			{Key: "name", Value: "Mary"},
		})
		assert.Equal(t, errMocked, err)
	})
}

func TestModel_InsertOneNoCache(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{
			{Key: "value", Value: bson.D{{Key: "foo", Value: "bar"}}},
		}...))
		m := createModel(t, mt)
		resp, err := m.InsertOneNoCache(context.Background(), bson.D{
			{Key: "name", Value: "Mary"},
		})
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})
}

func TestModel_ReplaceOne(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{
			{Key: "value", Value: bson.D{{Key: "foo", Value: "bar"}}},
		}...))
		m := createModel(t, mt)
		assert.Nil(t, m.cache.Set("foo", "bar"))
		resp, err := m.ReplaceOne(context.Background(), "foo", bson.D{}, bson.D{
			{Key: "foo", Value: "baz"},
		})
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		var v string
		assert.True(t, m.cache.IsNotFound(m.cache.Get("foo", &v)))
		_, err = m.ReplaceOne(context.Background(), "foo", bson.D{}, bson.D{
			{Key: "foo", Value: "baz"},
		})
		assert.NotNil(t, err)

		m.cache = mockedCache{m.cache}
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{
			{Key: "value", Value: bson.D{{Key: "foo", Value: "bar"}}},
		}...))
		_, err = m.ReplaceOne(context.Background(), "foo", bson.D{}, bson.D{
			{Key: "foo", Value: "baz"},
		})
		assert.Equal(t, errMocked, err)
	})
}

func TestModel_ReplaceOneNoCache(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{
			{Key: "value", Value: bson.D{{Key: "foo", Value: "bar"}}},
		}...))
		m := createModel(t, mt)
		resp, err := m.ReplaceOneNoCache(context.Background(), bson.D{}, bson.D{
			{Key: "foo", Value: "baz"},
		})
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})
}

func TestModel_SetCache(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		m := createModel(t, mt)
		assert.Nil(t, m.SetCache("foo", "bar"))
		var v string
		assert.Nil(t, m.GetCache("foo", &v))
		assert.Equal(t, "bar", v)
	})
}

func TestModel_UpdateByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{
			{Key: "value", Value: bson.D{{Key: "foo", Value: "bar"}}},
		}...))
		m := createModel(t, mt)
		assert.Nil(t, m.cache.Set("foo", "bar"))
		resp, err := m.UpdateByID(context.Background(), "foo", bson.D{}, bson.D{
			{Key: "$set", Value: bson.D{{Key: "foo", Value: "baz"}}},
		})
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		var v string
		assert.True(t, m.cache.IsNotFound(m.cache.Get("foo", &v)))
		_, err = m.UpdateByID(context.Background(), "foo", bson.D{}, bson.D{
			{Key: "$set", Value: bson.D{{Key: "foo", Value: "baz"}}},
		})
		assert.NotNil(t, err)

		m.cache = mockedCache{m.cache}
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{
			{Key: "value", Value: bson.D{{Key: "foo", Value: "bar"}}},
		}...))
		_, err = m.UpdateByID(context.Background(), "foo", bson.D{}, bson.D{
			{Key: "$set", Value: bson.D{{Key: "foo", Value: "baz"}}},
		})
		assert.Equal(t, errMocked, err)
	})
}

func TestModel_UpdateByIDNoCache(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{
			{Key: "value", Value: bson.D{{Key: "foo", Value: "bar"}}},
		}...))
		m := createModel(t, mt)
		resp, err := m.UpdateByIDNoCache(context.Background(), bson.D{}, bson.D{
			{Key: "$set", Value: bson.D{{Key: "foo", Value: "baz"}}},
		})
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})
}

func TestModel_UpdateMany(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{
			{Key: "value", Value: bson.D{{Key: "foo", Value: "bar"}}},
		}...))
		m := createModel(t, mt)
		assert.Nil(t, m.cache.Set("foo", "bar"))
		assert.Nil(t, m.cache.Set("bar", "baz"))
		resp, err := m.UpdateMany(context.Background(), []string{"foo", "bar"}, bson.D{}, bson.D{
			{Key: "$set", Value: bson.D{{Key: "foo", Value: "baz"}}},
		})
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		var v string
		assert.True(t, m.cache.IsNotFound(m.cache.Get("foo", &v)))
		assert.True(t, m.cache.IsNotFound(m.cache.Get("bar", &v)))
		_, err = m.UpdateMany(context.Background(), []string{"foo", "bar"}, bson.D{}, bson.D{
			{Key: "$set", Value: bson.D{{Key: "foo", Value: "baz"}}},
		})
		assert.NotNil(t, err)

		m.cache = mockedCache{m.cache}
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{
			{Key: "value", Value: bson.D{{Key: "foo", Value: "bar"}}},
		}...))
		_, err = m.UpdateMany(context.Background(), []string{"foo", "bar"}, bson.D{}, bson.D{
			{Key: "$set", Value: bson.D{{Key: "foo", Value: "baz"}}},
		})
		assert.Equal(t, errMocked, err)
	})
}

func TestModel_UpdateManyNoCache(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{
			{Key: "value", Value: bson.D{{Key: "foo", Value: "bar"}}},
		}...))
		m := createModel(t, mt)
		resp, err := m.UpdateManyNoCache(context.Background(), bson.D{}, bson.D{
			{Key: "$set", Value: bson.D{{Key: "foo", Value: "baz"}}},
		})
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})
}

func TestModel_UpdateOne(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{
			{Key: "value", Value: bson.D{{Key: "foo", Value: "bar"}}},
		}...))
		m := createModel(t, mt)
		assert.Nil(t, m.cache.Set("foo", "bar"))
		resp, err := m.UpdateOne(context.Background(), "foo", bson.D{}, bson.D{
			{Key: "$set", Value: bson.D{{Key: "foo", Value: "baz"}}},
		})
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		var v string
		assert.True(t, m.cache.IsNotFound(m.cache.Get("foo", &v)))
		_, err = m.UpdateOne(context.Background(), "foo", bson.D{}, bson.D{
			{Key: "$set", Value: bson.D{{Key: "foo", Value: "baz"}}},
		})
		assert.NotNil(t, err)

		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{
			{Key: "value", Value: bson.D{{Key: "foo", Value: "bar"}}},
		}...))
		m.cache = mockedCache{m.cache}
		_, err = m.UpdateOne(context.Background(), "foo", bson.D{}, bson.D{
			{Key: "$set", Value: bson.D{{Key: "foo", Value: "baz"}}},
		})
		assert.Equal(t, errMocked, err)
	})
}

func TestModel_UpdateOneNoCache(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{
			{Key: "value", Value: bson.D{{Key: "foo", Value: "bar"}}},
		}...))
		m := createModel(t, mt)
		resp, err := m.UpdateOneNoCache(context.Background(), bson.D{}, bson.D{
			{Key: "$set", Value: bson.D{{Key: "foo", Value: "baz"}}},
		})
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})
}

func createModel(t *testing.T, mt *mtest.T) *Model {
	s, err := miniredis.Run()
	assert.Nil(t, err)
	mon.Inject(mt.Name(), mt.Client)
	if atomic.AddInt32(&index, 1)%2 == 0 {
		return MustNewNodeModel(mt.Name(), mt.DB.Name(), mt.Coll.Name(), redis.New(s.Addr()))
	} else {
		return MustNewModel(mt.Name(), mt.DB.Name(), mt.Coll.Name(), cache.CacheConf{
			cache.NodeConf{
				RedisConf: redis.RedisConf{
					Host: s.Addr(),
					Type: redis.NodeType,
				},
				Weight: 100,
			},
		})
	}
}

var (
	errMocked = errors.New("mocked error")
	index     int32
)

type mockedCache struct {
	cache.Cache
}

func (m mockedCache) DelCtx(_ context.Context, _ ...string) error {
	return errMocked
}
