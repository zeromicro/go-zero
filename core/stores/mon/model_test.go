package mon

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestModel_StartSession(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		m := createModel(mt)
		sess, err := m.StartSession()
		assert.Nil(t, err)
		defer sess.EndSession(context.Background())

		_, err = sess.WithTransaction(context.Background(), func(sessCtx mongo.SessionContext) (any, error) {
			_ = sessCtx.StartTransaction()
			sessCtx.Client().Database("1")
			sessCtx.EndSession(context.Background())
			return nil, nil
		})
		assert.Nil(t, err)
		assert.NoError(t, sess.CommitTransaction(context.Background()))
		assert.Error(t, sess.AbortTransaction(context.Background()))
	})
}

func TestModel_Aggregate(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		m := createModel(mt)
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
		var result []any
		err := m.Aggregate(context.Background(), &result, mongo.Pipeline{})
		assert.Nil(t, err)
		assert.Equal(t, 2, len(result))
		assert.Equal(t, "John", result[0].(bson.D).Map()["name"])
		assert.Equal(t, "Mary", result[1].(bson.D).Map()["name"])

		triggerBreaker(m)
		assert.Equal(t, errDummy, m.Aggregate(context.Background(), &result, mongo.Pipeline{}))
	})
}

func TestModel_DeleteMany(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		m := createModel(mt)
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{{Key: "n", Value: 1}}...))
		val, err := m.DeleteMany(context.Background(), bson.D{})
		assert.Nil(t, err)
		assert.Equal(t, int64(1), val)

		triggerBreaker(m)
		_, err = m.DeleteMany(context.Background(), bson.D{})
		assert.Equal(t, errDummy, err)
	})
}

func TestModel_DeleteOne(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		m := createModel(mt)
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{{Key: "n", Value: 1}}...))
		val, err := m.DeleteOne(context.Background(), bson.D{})
		assert.Nil(t, err)
		assert.Equal(t, int64(1), val)

		triggerBreaker(m)
		_, err = m.DeleteOne(context.Background(), bson.D{})
		assert.Equal(t, errDummy, err)
	})
}

func TestModel_Find(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		m := createModel(mt)
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
		var result []any
		err := m.Find(context.Background(), &result, bson.D{})
		assert.Nil(t, err)
		assert.Equal(t, 2, len(result))
		assert.Equal(t, "John", result[0].(bson.D).Map()["name"])
		assert.Equal(t, "Mary", result[1].(bson.D).Map()["name"])

		triggerBreaker(m)
		assert.Equal(t, errDummy, m.Find(context.Background(), &result, bson.D{}))
	})
}

func TestModel_FindOne(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		m := createModel(mt)
		find := mtest.CreateCursorResponse(
			1,
			"DBName.CollectionName",
			mtest.FirstBatch,
			bson.D{
				{Key: "name", Value: "John"},
			})
		killCursors := mtest.CreateCursorResponse(
			0,
			"DBName.CollectionName",
			mtest.NextBatch)
		mt.AddMockResponses(find, killCursors)
		var result bson.D
		err := m.FindOne(context.Background(), &result, bson.D{})
		assert.Nil(t, err)
		assert.Equal(t, "John", result.Map()["name"])

		triggerBreaker(m)
		assert.Equal(t, errDummy, m.FindOne(context.Background(), &result, bson.D{}))
	})
}

func TestModel_FindOneAndDelete(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		m := createModel(mt)
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{
			{Key: "value", Value: bson.D{{Key: "name", Value: "John"}}},
		}...))
		var result bson.D
		err := m.FindOneAndDelete(context.Background(), &result, bson.D{})
		assert.Nil(t, err)
		assert.Equal(t, "John", result.Map()["name"])

		triggerBreaker(m)
		assert.Equal(t, errDummy, m.FindOneAndDelete(context.Background(), &result, bson.D{}))
	})
}

func TestModel_FindOneAndReplace(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		m := createModel(mt)
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{
			{Key: "value", Value: bson.D{{Key: "name", Value: "John"}}},
		}...))
		var result bson.D
		err := m.FindOneAndReplace(context.Background(), &result, bson.D{}, bson.D{
			{Key: "name", Value: "Mary"},
		})
		assert.Nil(t, err)
		assert.Equal(t, "John", result.Map()["name"])

		triggerBreaker(m)
		assert.Equal(t, errDummy, m.FindOneAndReplace(context.Background(), &result, bson.D{}, bson.D{
			{Key: "name", Value: "Mary"},
		}))
	})
}

func TestModel_FindOneAndUpdate(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		m := createModel(mt)
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{
			{Key: "value", Value: bson.D{{Key: "name", Value: "John"}}},
		}...))
		var result bson.D
		err := m.FindOneAndUpdate(context.Background(), &result, bson.D{}, bson.D{
			{Key: "$set", Value: bson.D{{Key: "name", Value: "Mary"}}},
		})
		assert.Nil(t, err)
		assert.Equal(t, "John", result.Map()["name"])

		triggerBreaker(m)
		assert.Equal(t, errDummy, m.FindOneAndUpdate(context.Background(), &result, bson.D{}, bson.D{
			{Key: "$set", Value: bson.D{{Key: "name", Value: "Mary"}}},
		}))
	})
}

func createModel(mt *mtest.T) *Model {
	Inject(mt.Name(), mt.Client)
	return MustNewModel(mt.Name(), mt.DB.Name(), mt.Coll.Name())
}

func triggerBreaker(m *Model) {
	m.Collection.(*decoratedCollection).brk = new(dropBreaker)
}
