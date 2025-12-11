package monc

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/mon"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/mock/gomock"
)

func TestMustNewModel(t *testing.T) {
	s, err := miniredis.Run()
	assert.Nil(t, err)
	original := logx.ExitOnFatal.True()
	logx.ExitOnFatal.Set(false)
	defer logx.ExitOnFatal.Set(original)

	assert.Panics(t, func() {
		MustNewModel("foo", "db", "collectino", cache.CacheConf{
			cache.NodeConf{
				RedisConf: redis.RedisConf{
					Host: s.Addr(),
					Type: redis.NodeType,
				},
				Weight: 100,
			}})
	})
}

func TestMustNewNodeModel(t *testing.T) {
	s, err := miniredis.Run()
	assert.Nil(t, err)
	original := logx.ExitOnFatal.True()
	logx.ExitOnFatal.Set(false)
	defer logx.ExitOnFatal.Set(original)

	assert.Panics(t, func() {
		MustNewNodeModel("foo", "db", "collectino", redis.New(s.Addr()))
	})
}

func TestNewModel(t *testing.T) {
	s, err := miniredis.Run()
	assert.Nil(t, err)
	_, err = NewModel("foo", "db", "coll", cache.CacheConf{
		cache.NodeConf{
			RedisConf: redis.RedisConf{
				Host: s.Addr(),
				Type: redis.NodeType,
			},
			Weight: 100,
		},
	})
	assert.Error(t, err)
}

func TestNewNodeModel(t *testing.T) {
	_, err := NewNodeModel("foo", "db", "coll", nil)
	assert.NotNil(t, err)
}

func TestNewModelWithCache(t *testing.T) {
	_, err := NewModelWithCache("foo", "db", "coll", nil)
	assert.NotNil(t, err)
}

func Test_newModel(t *testing.T) {
	mon.Inject("mongodb://localhost:27018", &mongo.Client{})
	model, err := newModel("mongodb://localhost:27018", "db", "collection", nil)
	assert.Nil(t, err)
	assert.NotNil(t, model)
}

func TestModel_DelCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCollection := mon.NewMockCollection(ctrl)
	m := createModel(t, mockCollection)
	assert.Nil(t, m.cache.Set("foo", "bar"))
	assert.Nil(t, m.cache.Set("bar", "baz"))
	assert.Nil(t, m.DelCache(context.Background(), "foo", "bar"))
	var v string
	assert.True(t, m.cache.IsNotFound(m.cache.Get("foo", &v)))
	assert.True(t, m.cache.IsNotFound(m.cache.Get("bar", &v)))
}

func TestModel_DeleteOne(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCollection := mon.NewMockCollection(ctrl)
	mockCollection.EXPECT().DeleteOne(gomock.Any(), gomock.Any(), gomock.Any()).Return(&mongo.DeleteResult{DeletedCount: 1}, nil)
	m := createModel(t, mockCollection)
	assert.Nil(t, m.cache.Set("foo", "bar"))
	val, err := m.DeleteOne(context.Background(), "foo", bson.D{{Key: "foo", Value: "bar"}})
	assert.Nil(t, err)
	assert.Equal(t, int64(1), val)
	var v string
	assert.True(t, m.cache.IsNotFound(m.cache.Get("foo", &v)))
	mockCollection.EXPECT().DeleteOne(gomock.Any(), gomock.Any(), gomock.Any()).Return(&mongo.DeleteResult{}, errMocked)
	_, err = m.DeleteOne(context.Background(), "foo", bson.D{{Key: "foo", Value: "bar"}})
	assert.NotNil(t, err)

	m.cache = mockedCache{m.cache}
	mockCollection.EXPECT().DeleteOne(gomock.Any(), gomock.Any(), gomock.Any()).Return(&mongo.DeleteResult{}, nil)
	_, err = m.DeleteOne(context.Background(), "foo", bson.D{{Key: "foo", Value: "bar"}})
	assert.Equal(t, errMocked, err)
}

func TestModel_DeleteOneNoCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCollection := mon.NewMockCollection(ctrl)
	mockCollection.EXPECT().DeleteOne(gomock.Any(), gomock.Any(), gomock.Any()).Return(&mongo.DeleteResult{DeletedCount: 1}, nil)
	m := createModel(t, mockCollection)
	assert.Nil(t, m.cache.Set("foo", "bar"))
	val, err := m.DeleteOneNoCache(context.Background(), bson.D{{Key: "foo", Value: "bar"}})
	assert.Nil(t, err)
	assert.Equal(t, int64(1), val)
	var v string
	assert.Nil(t, m.cache.Get("foo", &v))
}

func TestModel_FindOne(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCollection := mon.NewMockCollection(ctrl)
	mockCollection.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mongo.NewSingleResultFromDocument(bson.M{"foo": "bar"}, nil, nil), nil)
	m := createModel(t, mockCollection)
	var v struct {
		Foo string `bson:"foo"`
	}
	assert.Nil(t, m.FindOne(context.Background(), "foo", &v, bson.D{}))
	assert.Equal(t, "bar", v.Foo)
	assert.Nil(t, m.cache.Set("foo", "bar"))
}

func TestModel_FindOneNoCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCollection := mon.NewMockCollection(ctrl)
	mockCollection.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mongo.NewSingleResultFromDocument(bson.M{"foo": "bar"}, nil, nil), nil)
	m := createModel(t, mockCollection)
	v := struct {
		Foo string `bson:"foo"`
	}{}
	assert.Nil(t, m.FindOneNoCache(context.Background(), &v, bson.D{}))
	assert.Equal(t, "bar", v.Foo)
}

func TestModel_FindOneAndDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCollection := mon.NewMockCollection(ctrl)
	mockCollection.EXPECT().FindOneAndDelete(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mongo.NewSingleResultFromDocument(bson.M{"foo": "bar"}, nil, bson.NewRegistry()), nil)
	m := createModel(t, mockCollection)
	assert.Nil(t, m.cache.Set("foo", "bar"))
	v := struct {
		Foo string `bson:"foo"`
	}{}
	assert.Nil(t, m.FindOneAndDelete(context.Background(), "foo", &v, bson.D{}))
	assert.Equal(t, "bar", v.Foo)
	assert.True(t, m.cache.IsNotFound(m.cache.Get("foo", &v)))
	mockCollection.EXPECT().FindOneAndDelete(gomock.Any(), gomock.Any(), gomock.Any()).Return(mongo.NewSingleResultFromDocument(bson.M{"foo": "bar"}, nil, bson.NewRegistry()), errMocked)
	assert.NotNil(t, m.FindOneAndDelete(context.Background(), "foo", &v, bson.D{}))
	m.cache = mockedCache{m.cache}
	mockCollection.EXPECT().FindOneAndDelete(gomock.Any(), gomock.Any(), gomock.Any()).Return(mongo.NewSingleResultFromDocument(bson.M{"foo": "bar"}, nil, bson.NewRegistry()), nil)
	assert.Equal(t, errMocked, m.FindOneAndDelete(context.Background(), "foo", &v, bson.D{}))
}

func TestModel_FindOneAndDeleteNoCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCollection := mon.NewMockCollection(ctrl)
	mockCollection.EXPECT().FindOneAndDelete(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mongo.NewSingleResultFromDocument(bson.M{"foo": "bar"}, nil, nil), nil)
	m := createModel(t, mockCollection)
	v := struct {
		Foo string `bson:"foo"`
	}{}
	assert.Nil(t, m.FindOneAndDeleteNoCache(context.Background(), &v, bson.D{}))
	assert.Equal(t, "bar", v.Foo)
}

func TestModel_FindOneAndReplace(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCollection := mon.NewMockCollection(ctrl)
	m := createModel(t, mockCollection)
	assert.Nil(t, m.cache.Set("foo", "bar"))
	v := struct {
		Foo string `bson:"foo"`
	}{}
	mockCollection.EXPECT().FindOneAndReplace(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mongo.NewSingleResultFromDocument(bson.M{"foo": "bar"}, nil, nil), nil)
	assert.Nil(t, m.FindOneAndReplace(context.Background(), "foo", &v, bson.D{}, bson.D{
		{Key: "name", Value: "Mary"},
	}))
	assert.Equal(t, "bar", v.Foo)
	assert.True(t, m.cache.IsNotFound(m.cache.Get("foo", &v)))
	mockCollection.EXPECT().FindOneAndReplace(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mongo.NewSingleResultFromDocument(bson.M{"name": "Mary"}, nil, nil), errMocked)
	assert.NotNil(t, m.FindOneAndReplace(context.Background(), "foo", &v, bson.D{}, bson.D{
		{Key: "name", Value: "Mary"},
	}))

	m.cache = mockedCache{m.cache}
	mockCollection.EXPECT().FindOneAndReplace(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mongo.NewSingleResultFromDocument(bson.M{"foo": "bar"}, nil, nil), nil)
	assert.Equal(t, errMocked, m.FindOneAndReplace(context.Background(), "foo", &v, bson.D{}, bson.D{
		{Key: "name", Value: "Mary"},
	}))
}

func TestModel_FindOneAndReplaceNoCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCollection := mon.NewMockCollection(ctrl)
	m := createModel(t, mockCollection)
	v := struct {
		Foo string `bson:"foo"`
	}{}
	mockCollection.EXPECT().FindOneAndReplace(gomock.Any(), gomock.Any(), gomock.Any()).Return(mongo.NewSingleResultFromDocument(bson.M{"foo": "bar"}, nil, nil), nil)
	assert.Nil(t, m.FindOneAndReplaceNoCache(context.Background(), &v, bson.D{}, bson.D{
		{Key: "name", Value: "Mary"},
	}))
	assert.Equal(t, "bar", v.Foo)
}

func TestModel_FindOneAndUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCollection := mon.NewMockCollection(ctrl)
	m := createModel(t, mockCollection)
	assert.Nil(t, m.cache.Set("foo", "bar"))
	v := struct {
		Foo string `bson:"foo"`
	}{}
	mockCollection.EXPECT().FindOneAndUpdate(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mongo.NewSingleResultFromDocument(bson.M{"foo": "bar"}, nil, nil), nil)
	assert.Nil(t, m.FindOneAndUpdate(context.Background(), "foo", &v, bson.D{}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "name", Value: "Mary"}}},
	}))
	assert.Equal(t, "bar", v.Foo)
	assert.True(t, m.cache.IsNotFound(m.cache.Get("foo", &v)))
	mockCollection.EXPECT().FindOneAndUpdate(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mongo.NewSingleResultFromDocument(bson.M{"foo": "bar"}, nil, nil), errMocked)
	assert.NotNil(t, m.FindOneAndUpdate(context.Background(), "foo", &v, bson.D{}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "name", Value: "Mary"}}},
	}))

	m.cache = mockedCache{m.cache}
	mockCollection.EXPECT().FindOneAndUpdate(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mongo.NewSingleResultFromDocument(bson.M{"foo": "bar"}, nil, nil), nil)
	assert.Equal(t, errMocked, m.FindOneAndUpdate(context.Background(), "foo", &v, bson.D{}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "name", Value: "Mary"}}},
	}))
}

func TestModel_FindOneAndUpdateNoCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCollection := mon.NewMockCollection(ctrl)
	m := createModel(t, mockCollection)
	v := struct {
		Foo string `bson:"foo"`
	}{}
	mockCollection.EXPECT().FindOneAndUpdate(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mongo.NewSingleResultFromDocument(bson.M{"foo": "bar"}, nil, nil), nil)
	assert.Nil(t, m.FindOneAndUpdateNoCache(context.Background(), &v, bson.D{}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "name", Value: "Mary"}}},
	}))
	assert.Equal(t, "bar", v.Foo)
}

func TestModel_GetCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCollection := mon.NewMockCollection(ctrl)
	m := createModel(t, mockCollection)
	assert.NotNil(t, m.cache)
	assert.Nil(t, m.cache.Set("foo", "bar"))
	var s string
	assert.Nil(t, m.cache.Get("foo", &s))
	assert.Equal(t, "bar", s)
}

func TestModel_InsertOne(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCollection := mon.NewMockCollection(ctrl)
	m := createModel(t, mockCollection)
	assert.Nil(t, m.cache.Set("foo", "bar"))
	mockCollection.EXPECT().InsertOne(gomock.Any(), gomock.Any(), gomock.Any()).Return(&mongo.InsertOneResult{}, nil)
	resp, err := m.InsertOne(context.Background(), "foo", bson.D{
		{Key: "name", Value: "Mary"},
	})
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	var v string
	assert.True(t, m.cache.IsNotFound(m.cache.Get("foo", &v)))
	mockCollection.EXPECT().InsertOne(gomock.Any(), gomock.Any(), gomock.Any()).Return(&mongo.InsertOneResult{}, errMocked)
	_, err = m.InsertOne(context.Background(), "foo", bson.D{
		{Key: "name", Value: "Mary"},
	})
	assert.NotNil(t, err)

	m.cache = mockedCache{m.cache}
	mockCollection.EXPECT().InsertOne(gomock.Any(), gomock.Any(), gomock.Any()).Return(&mongo.InsertOneResult{}, nil)
	_, err = m.InsertOne(context.Background(), "foo", bson.D{
		{Key: "name", Value: "Mary"},
	})
	assert.Equal(t, errMocked, err)
}

func TestModel_InsertOneNoCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCollection := mon.NewMockCollection(ctrl)
	m := createModel(t, mockCollection)
	mockCollection.EXPECT().InsertOne(gomock.Any(), gomock.Any(), gomock.Any()).Return(&mongo.InsertOneResult{}, nil)
	resp, err := m.InsertOneNoCache(context.Background(), bson.D{
		{Key: "name", Value: "Mary"},
	})
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func TestModel_ReplaceOne(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCollection := mon.NewMockCollection(ctrl)
	m := createModel(t, mockCollection)
	assert.Nil(t, m.cache.Set("foo", "bar"))
	mockCollection.EXPECT().ReplaceOne(gomock.Any(), gomock.Any(), gomock.Any()).Return(&mongo.UpdateResult{}, nil)
	resp, err := m.ReplaceOne(context.Background(), "foo", bson.D{}, bson.D{
		{Key: "foo", Value: "baz"},
	})
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	var v string
	assert.True(t, m.cache.IsNotFound(m.cache.Get("foo", &v)))
	mockCollection.EXPECT().ReplaceOne(gomock.Any(), gomock.Any(), gomock.Any()).Return(&mongo.UpdateResult{}, errMocked)
	_, err = m.ReplaceOne(context.Background(), "foo", bson.D{}, bson.D{
		{Key: "foo", Value: "baz"},
	})
	assert.NotNil(t, err)

	m.cache = mockedCache{m.cache}
	mockCollection.EXPECT().ReplaceOne(gomock.Any(), gomock.Any(), gomock.Any()).Return(&mongo.UpdateResult{}, nil)
	_, err = m.ReplaceOne(context.Background(), "foo", bson.D{}, bson.D{
		{Key: "foo", Value: "baz"},
	})
	assert.Equal(t, errMocked, err)
}

func TestModel_ReplaceOneNoCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCollection := mon.NewMockCollection(ctrl)
	m := createModel(t, mockCollection)
	mockCollection.EXPECT().ReplaceOne(gomock.Any(), gomock.Any(), gomock.Any()).Return(&mongo.UpdateResult{}, nil)
	resp, err := m.ReplaceOneNoCache(context.Background(), bson.D{}, bson.D{
		{Key: "foo", Value: "baz"},
	})
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func TestModel_SetCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCollection := mon.NewMockCollection(ctrl)
	m := createModel(t, mockCollection)
	assert.Nil(t, m.SetCache("foo", "bar"))
	var v string
	assert.Nil(t, m.GetCache("foo", &v))
	assert.Equal(t, "bar", v)
}

func TestModel_UpdateByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCollection := mon.NewMockCollection(ctrl)
	m := createModel(t, mockCollection)
	assert.Nil(t, m.cache.Set("foo", "bar"))
	mockCollection.EXPECT().UpdateByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(&mongo.UpdateResult{}, nil)
	resp, err := m.UpdateByID(context.Background(), "foo", bson.D{}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "foo", Value: "baz"}}},
	})
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	var v string
	assert.True(t, m.cache.IsNotFound(m.cache.Get("foo", &v)))
	mockCollection.EXPECT().UpdateByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(&mongo.UpdateResult{}, errMocked)
	_, err = m.UpdateByID(context.Background(), "foo", bson.D{}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "foo", Value: "baz"}}},
	})
	assert.NotNil(t, err)

	m.cache = mockedCache{m.cache}
	mockCollection.EXPECT().UpdateByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(&mongo.UpdateResult{}, nil)
	_, err = m.UpdateByID(context.Background(), "foo", bson.D{}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "foo", Value: "baz"}}},
	})
	assert.Equal(t, errMocked, err)
}

func TestModel_UpdateByIDNoCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCollection := mon.NewMockCollection(ctrl)
	m := createModel(t, mockCollection)
	mockCollection.EXPECT().UpdateByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(&mongo.UpdateResult{}, nil)
	resp, err := m.UpdateByIDNoCache(context.Background(), bson.D{}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "foo", Value: "baz"}}},
	})
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func TestModel_UpdateMany(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCollection := mon.NewMockCollection(ctrl)
	m := createModel(t, mockCollection)
	assert.Nil(t, m.cache.Set("foo", "bar"))
	assert.Nil(t, m.cache.Set("bar", "baz"))
	mockCollection.EXPECT().UpdateMany(gomock.Any(), gomock.Any(), gomock.Any()).Return(&mongo.UpdateResult{}, nil)
	resp, err := m.UpdateMany(context.Background(), []string{"foo", "bar"}, bson.D{}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "foo", Value: "baz"}}},
	})
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	var v string
	assert.True(t, m.cache.IsNotFound(m.cache.Get("foo", &v)))
	assert.True(t, m.cache.IsNotFound(m.cache.Get("bar", &v)))
	mockCollection.EXPECT().UpdateMany(gomock.Any(), gomock.Any(), gomock.Any()).Return(&mongo.UpdateResult{}, errMocked)
	_, err = m.UpdateMany(context.Background(), []string{"foo", "bar"}, bson.D{}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "foo", Value: "baz"}}},
	})
	assert.NotNil(t, err)

	m.cache = mockedCache{m.cache}
	mockCollection.EXPECT().UpdateMany(gomock.Any(), gomock.Any(), gomock.Any()).Return(&mongo.UpdateResult{}, nil)
	_, err = m.UpdateMany(context.Background(), []string{"foo", "bar"}, bson.D{}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "foo", Value: "baz"}}},
	})
	assert.Equal(t, errMocked, err)
}

func TestModel_UpdateManyNoCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCollection := mon.NewMockCollection(ctrl)
	m := createModel(t, mockCollection)
	mockCollection.EXPECT().UpdateMany(gomock.Any(), gomock.Any(), gomock.Any()).Return(&mongo.UpdateResult{}, nil)
	resp, err := m.UpdateManyNoCache(context.Background(), bson.D{}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "foo", Value: "baz"}}},
	})
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func TestModel_UpdateOne(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCollection := mon.NewMockCollection(ctrl)
	m := createModel(t, mockCollection)
	assert.Nil(t, m.cache.Set("foo", "bar"))
	mockCollection.EXPECT().UpdateOne(gomock.Any(), gomock.Any(), gomock.Any()).Return(&mongo.UpdateResult{}, nil)
	resp, err := m.UpdateOne(context.Background(), "foo", bson.D{}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "foo", Value: "baz"}}},
	})
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	var v string
	mockCollection.EXPECT().UpdateOne(gomock.Any(), gomock.Any(), gomock.Any()).Return(&mongo.UpdateResult{}, errMocked)
	assert.True(t, m.cache.IsNotFound(m.cache.Get("foo", &v)))
	_, err = m.UpdateOne(context.Background(), "foo", bson.D{}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "foo", Value: "baz"}}},
	})
	assert.NotNil(t, err)

	mockCollection.EXPECT().UpdateOne(gomock.Any(), gomock.Any(), gomock.Any()).Return(&mongo.UpdateResult{}, nil)
	m.cache = mockedCache{m.cache}
	_, err = m.UpdateOne(context.Background(), "foo", bson.D{}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "foo", Value: "baz"}}},
	})
	assert.Equal(t, errMocked, err)
}

func TestModel_UpdateOneNoCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCollection := mon.NewMockCollection(ctrl)
	m := createModel(t, mockCollection)
	mockCollection.EXPECT().UpdateOne(gomock.Any(), gomock.Any(), gomock.Any()).Return(&mongo.UpdateResult{}, nil)
	resp, err := m.UpdateOneNoCache(context.Background(), bson.D{}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "foo", Value: "baz"}}},
	})
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func createModel(t *testing.T, coll mon.Collection) *Model {
	s, err := miniredis.Run()
	assert.Nil(t, err)
	if atomic.AddInt32(&index, 1)%2 == 0 {
		return mustNewTestNodeModel(coll, redis.New(s.Addr()))
	} else {
		return mustNewTestModel(coll, cache.CacheConf{
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

// mustNewTestModel returns a test Model with the given cache.
func mustNewTestModel(collection mon.Collection, c cache.CacheConf, opts ...cache.Option) *Model {
	return &Model{
		Model: &mon.Model{
			Collection: collection,
		},
		cache: cache.New(c, singleFlight, stats, mongo.ErrNoDocuments, opts...),
	}
}

// NewNodeModel returns a test Model with a cache node.
func mustNewTestNodeModel(collection mon.Collection, rds *redis.Redis, opts ...cache.Option) *Model {
	c := cache.NewNode(rds, singleFlight, stats, mongo.ErrNoDocuments, opts...)
	return &Model{
		Model: &mon.Model{
			Collection: collection,
		},
		cache: c,
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
