package mon

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func TestModel_StartSession(t *testing.T) {
	m := createTestModel()
	sess, err := m.StartSession()
	assert.Nil(t, err)
	defer sess.EndSession(context.Background())

	_, err = sess.WithTransaction(context.Background(), func(sessCtx context.Context) (any, error) {
		//_ = sessCtx.StartTransaction()
		//sessCtx.Client().Database("1")
		//sessCtx.EndSession(context.Background())
		return nil, nil
	})
	assert.Nil(t, err)
	assert.NoError(t, sess.CommitTransaction(context.Background()))
	assert.NoError(t, sess.AbortTransaction(context.Background()))
}

func TestModel_Aggregate(t *testing.T) {
	//TODO test this func
}

func TestModel_DeleteMany(t *testing.T) {
	m := createTestModel()
	_, err := m.DeleteMany(context.Background(), bson.D{})
	assert.Nil(t, err)
	triggerBreaker(m)
	_, err = m.DeleteMany(context.Background(), bson.D{})
	assert.Equal(t, errDummy, err)
}

func TestModel_DeleteOne(t *testing.T) {
	m := createTestModel()
	_, err := m.DeleteOne(context.Background(), bson.D{})
	assert.Nil(t, err)
	triggerBreaker(m)
	_, err = m.DeleteOne(context.Background(), bson.D{})
	assert.Equal(t, errDummy, err)
}

func TestModel_Find(t *testing.T) {
	//TODO test this func
}

func TestModel_FindOne(t *testing.T) {
	m := createTestModel()
	var result bson.D
	err := m.FindOne(context.Background(), &result, bson.D{})
	assert.Equal(t, mongo.ErrNoDocuments, err)

	triggerBreaker(m)
	assert.Equal(t, errDummy, m.FindOne(context.Background(), &result, bson.D{}))
}

func TestModel_FindOneAndDelete(t *testing.T) {
	m := createTestModel()
	var result bson.D
	err := m.FindOneAndDelete(context.Background(), &result, bson.D{})
	assert.Equal(t, mongo.ErrNoDocuments, err)
	triggerBreaker(m)
	assert.Equal(t, errDummy, m.FindOneAndDelete(context.Background(), &result, bson.D{}))
}

func TestModel_FindOneAndReplace(t *testing.T) {
	m := createTestModel()
	var result bson.D
	err := m.FindOneAndReplace(context.Background(), &result, bson.D{}, bson.D{
		{Key: "name", Value: "Mary"},
	})
	assert.Equal(t, mongo.ErrNoDocuments, err)
	triggerBreaker(m)
	assert.Equal(t, errDummy, m.FindOneAndReplace(context.Background(), &result, bson.D{}, bson.D{
		{Key: "name", Value: "Mary"},
	}))
}

func TestModel_FindOneAndUpdate(t *testing.T) {
	m := createTestModel()
	var result bson.D
	err := m.FindOneAndUpdate(context.Background(), &result, bson.D{}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "name", Value: "Mary"}}},
	})
	assert.Equal(t, mongo.ErrNoDocuments, err)

	triggerBreaker(m)
	assert.Equal(t, errDummy, m.FindOneAndUpdate(context.Background(), &result, bson.D{}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "name", Value: "Mary"}}},
	}))

}

func createModel(mt *mtest.T) *Model {
	//Inject(mt.Name(), mt.Client)
	return MustNewModel(mt.Name(), mt.DB.Name(), mt.Coll.Name())
}

func triggerBreaker(m *Model) {
	m.Collection.(*decoratedCollection).brk = new(dropBreaker)
}

func createTestModel() *Model {
	return mustNewTestModel("a", "b", "c")
}
