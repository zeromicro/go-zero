package mon

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestBulkInserter(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{{Key: "ok", Value: 1}}...))
		bulk, err := NewBulkInserter(createModel(mt).Collection)
		assert.Equal(t, err, nil)
		bulk.SetResultHandler(func(result *mongo.InsertManyResult, err error) {
			assert.Nil(t, err)
			assert.Equal(t, 2, len(result.InsertedIDs))
		})
		bulk.Insert(bson.D{{Key: "foo", Value: "bar"}})
		bulk.Insert(bson.D{{Key: "foo", Value: "baz"}})
		bulk.Flush()
	})
}
