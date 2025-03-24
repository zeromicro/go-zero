package mon

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/mock/gomock"
)

func TestBulkInserter_InsertAndFlush(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCollection := NewMockCollection(ctrl)
	mockCollection.EXPECT().Clone().Return(&mongo.Collection{})
	bulkInserter, err := NewBulkInserter(mockCollection, time.Second)
	assert.NoError(t, err)
	bulkInserter.SetResultHandler(func(result *mongo.InsertManyResult, err error) {
		assert.Nil(t, err)
		assert.Equal(t, 2, len(result.InsertedIDs))
	})
	doc := map[string]interface{}{"name": "test"}
	bulkInserter.Insert(doc)
	bulkInserter.Flush()
}

func TestBulkInserter_SetResultHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCollection := NewMockCollection(ctrl)
	mockCollection.EXPECT().Clone().Return(nil)
	bulkInserter, err := NewBulkInserter(mockCollection)
	assert.NoError(t, err)
	mockHandler := func(result *mongo.InsertManyResult, err error) {}
	bulkInserter.SetResultHandler(mockHandler)
}

func TestDbInserter_RemoveAll(t *testing.T) {
	inserter := &dbInserter{}
	inserter.documents = []interface{}{}
	docs := inserter.RemoveAll()
	assert.NotNil(t, docs)
	assert.Empty(t, inserter.documents)
}

func Test_dbInserter_Execute(t *testing.T) {
	type fields struct {
		collection    collectionInserter
		documents     []any
		resultHandler ResultHandler
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCollection := NewMockcollectionInserter(ctrl)
	type args struct {
		objs any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		mock   func()
	}{
		{
			name: "empty doc",
			fields: fields{
				collection:    nil,
				documents:     nil,
				resultHandler: nil,
			},
			args: args{
				objs: make([]any, 0),
			},
			mock: func() {},
		},
		{
			name: "result handler",
			fields: fields{
				collection: mockCollection,
				resultHandler: func(result *mongo.InsertManyResult, err error) {
					assert.NotNil(t, err)
				},
			},
			args: args{
				objs: make([]any, 1),
			},
			mock: func() {
				mockCollection.EXPECT().InsertMany(gomock.Any(), gomock.Any()).Return(&mongo.InsertManyResult{}, errors.New("error"))
			},
		},
		{
			name: "normal error handler",
			fields: fields{
				collection:    mockCollection,
				resultHandler: nil,
			},
			args: args{
				objs: make([]any, 1),
			},
			mock: func() {
				mockCollection.EXPECT().InsertMany(gomock.Any(), gomock.Any()).Return(&mongo.InsertManyResult{}, errors.New("error"))
			},
		},
		{
			name: "no error",
			fields: fields{
				collection:    mockCollection,
				resultHandler: nil,
			},
			args: args{
				objs: make([]any, 1),
			},
			mock: func() {
				mockCollection.EXPECT().InsertMany(gomock.Any(), gomock.Any()).Return(&mongo.InsertManyResult{}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			in := &dbInserter{
				collection:    tt.fields.collection,
				documents:     tt.fields.documents,
				resultHandler: tt.fields.resultHandler,
			}
			in.Execute(tt.args.objs)
		})
	}
}
