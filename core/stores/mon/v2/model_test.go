package mon

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/breaker"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/x/mongo/driver/drivertest"
	"go.uber.org/mock/gomock"
)

func TestModel_StartSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMonCollection := NewMockmonCollection(ctrl)
	mockedMonClient := NewMockmonClient(ctrl)
	mockMonSession := NewMockmonSession(ctrl)
	warpSession := &WrappedSession{
		session: mockMonSession,
		name:    "",
		brk:     breaker.GetBreaker("localhost"),
	}
	mockedMonClient.EXPECT().StartSession(gomock.Any()).Return(warpSession, nil)
	m := newTestModel("foo", mockedMonClient, mockMonCollection, breaker.GetBreaker("test"))
	sess, err := m.StartSession()
	assert.Nil(t, err)
	defer sess.EndSession(context.Background())
	mockMonSession.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).Return(nil, nil)
	mockMonSession.EXPECT().CommitTransaction(gomock.Any()).Return(nil)
	mockMonSession.EXPECT().AbortTransaction(gomock.Any()).Return(nil)
	mockMonSession.EXPECT().EndSession(gomock.Any())
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMonCollection := NewMockmonCollection(ctrl)
	mockedMonClient := NewMockmonClient(ctrl)
	mockMonCollection.EXPECT().Aggregate(gomock.Any(), gomock.Any(), gomock.Any()).Return(&mongo.Cursor{Current: make(bson.Raw, 0)}, errors.New("error"))
	m := newTestModel("foo", mockedMonClient, mockMonCollection, breaker.GetBreaker("test"))
	var result bson.D
	err := m.Aggregate(context.Background(), &result, bson.D{})
	assert.Equal(t, errors.New("error"), err)
	triggerBreaker(m)
	assert.Equal(t, errDummy, m.Aggregate(context.Background(), &result, bson.D{}))
}

func TestModel_DeleteMany(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMonCollection := NewMockmonCollection(ctrl)
	mockedMonClient := NewMockmonClient(ctrl)
	mockMonCollection.EXPECT().DeleteMany(gomock.Any(), gomock.Any(), gomock.Any()).Return(&mongo.DeleteResult{}, nil)
	m := newTestModel("foo", mockedMonClient, mockMonCollection, breaker.GetBreaker("test"))
	_, err := m.DeleteMany(context.Background(), bson.D{})
	assert.Nil(t, err)
	triggerBreaker(m)
	_, err = m.DeleteMany(context.Background(), bson.D{})
	assert.Equal(t, errDummy, err)
}

func TestModel_DeleteOne(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMonCollection := NewMockmonCollection(ctrl)
	mockedMonClient := NewMockmonClient(ctrl)
	mockMonCollection.EXPECT().DeleteOne(gomock.Any(), gomock.Any(), gomock.Any()).Return(&mongo.DeleteResult{}, nil)
	m := newTestModel("foo", mockedMonClient, mockMonCollection, breaker.GetBreaker("test"))
	_, err := m.DeleteOne(context.Background(), bson.D{})
	assert.Nil(t, err)
	triggerBreaker(m)
	_, err = m.DeleteOne(context.Background(), bson.D{})
	assert.Equal(t, errDummy, err)
}

func TestModel_Find(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMonCollection := NewMockmonCollection(ctrl)
	mockedMonClient := NewMockmonClient(ctrl)
	mockMonCollection.EXPECT().Find(gomock.Any(), gomock.Any(), gomock.Any()).Return(&mongo.Cursor{Current: make(bson.Raw, 0)}, errors.New("error"))
	m := newTestModel("foo", mockedMonClient, mockMonCollection, breaker.GetBreaker("test"))
	var result bson.D
	err := m.Find(context.Background(), &result, bson.D{})
	assert.Equal(t, errors.New("error"), err)
	triggerBreaker(m)
	assert.Equal(t, errDummy, m.Find(context.Background(), &result, bson.D{}))
}

func TestModel_FindOne(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMonCollection := NewMockmonCollection(ctrl)
	mockedMonClient := NewMockmonClient(ctrl)
	mockMonCollection.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).Return(&mongo.SingleResult{})
	m := newTestModel("foo", mockedMonClient, mockMonCollection, breaker.GetBreaker("test"))
	var result bson.D
	err := m.FindOne(context.Background(), &result, bson.D{})
	assert.Equal(t, mongo.ErrNoDocuments, err)
	triggerBreaker(m)
	assert.Equal(t, errDummy, m.FindOne(context.Background(), &result, bson.D{}))
}

func TestModel_FindOneAndDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMonCollection := NewMockmonCollection(ctrl)
	mockedMonClient := NewMockmonClient(ctrl)
	mockMonCollection.EXPECT().FindOneAndDelete(gomock.Any(), gomock.Any(), gomock.Any()).Return(&mongo.SingleResult{})
	m := newTestModel("foo", mockedMonClient, mockMonCollection, breaker.GetBreaker("test"))
	var result bson.D
	err := m.FindOneAndDelete(context.Background(), &result, bson.D{})
	assert.Equal(t, mongo.ErrNoDocuments, err)
	triggerBreaker(m)
	assert.Equal(t, errDummy, m.FindOneAndDelete(context.Background(), &result, bson.D{}))
}

func TestModel_FindOneAndReplace(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMonCollection := NewMockmonCollection(ctrl)
	mockedMonClient := NewMockmonClient(ctrl)
	mockMonCollection.EXPECT().FindOneAndReplace(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&mongo.SingleResult{})
	m := newTestModel("foo", mockedMonClient, mockMonCollection, breaker.GetBreaker("test"))
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMonCollection := NewMockmonCollection(ctrl)
	mockedMonClient := NewMockmonClient(ctrl)
	mockMonCollection.EXPECT().FindOneAndUpdate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&mongo.SingleResult{})
	m := newTestModel("foo", mockedMonClient, mockMonCollection, breaker.GetBreaker("test"))
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

func triggerBreaker(m *Model) {
	m.Collection.(*decoratedCollection).brk = new(dropBreaker)
}

func TestMustNewModel(t *testing.T) {
	Inject("mongodb://localhost:27017", &mongo.Client{})
	MustNewModel("mongodb://localhost:27017", "test", "test")
}

func TestNewModel(t *testing.T) {
	Inject("mongodb://localhost:27018", &mongo.Client{})
	NewModel("mongodb://localhost:27018", "test", "test")
}

func Test_newModel(t *testing.T) {
	Inject("mongodb://localhost:27019", &mongo.Client{})
	newModel("mongodb://localhost:27019", nil, nil, nil)
}

func Test_mockMonClient_StartSession(t *testing.T) {
	md := drivertest.NewMockDeployment()
	opts := options.Client()
	opts.Deployment = md
	client, err := mongo.Connect(opts)
	assert.Nil(t, err)
	m := mockMonClient{
		c: client,
	}
	_, err = m.StartSession()
	assert.Nil(t, err)
}
