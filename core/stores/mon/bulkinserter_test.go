package mon

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestBulkInserter_Insert(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("test", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{{"ok", 1}}...))
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{{"ok", 1}}...))
		bulkInserter := NewBulkInserter(mt.Coll, time.Millisecond*10)
		bulkInserter.SetResultHandler(func(result *mongo.InsertManyResult, err error) {
			assert.Nil(t, err)
			assert.Equal(t, maxBulkRows, len(result.InsertedIDs))
		})
		for i := 0; i < 2000; i++ {
			bulkInserter.Insert(bson.D{{"x", i}})
		}
		bulkInserter.Flush()
		time.Sleep(time.Millisecond * 100)
	})
}
