package mon

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/executors"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	flushInterval = time.Second
	maxBulkRows   = 1000
)

type (
	// ResultHandler is a handler that used to handle results.
	ResultHandler func(*mongo.InsertManyResult, error)

	// A BulkInserter is used to insert bulk of mongo records.
	BulkInserter struct {
		executor *executors.PeriodicalExecutor
		inserter *dbInserter
	}
)

// NewBulkInserter returns a BulkInserter.
func NewBulkInserter(coll Collection, interval ...time.Duration) (*BulkInserter, error) {
	cloneColl, err := coll.Clone()
	if err != nil {
		return nil, err
	}

	inserter := &dbInserter{
		collection: cloneColl,
	}

	duration := flushInterval
	if len(interval) > 0 {
		duration = interval[0]
	}

	return &BulkInserter{
		executor: executors.NewPeriodicalExecutor(duration, inserter),
		inserter: inserter,
	}, nil
}

// Flush flushes the inserter, writes all pending records.
func (bi *BulkInserter) Flush() {
	bi.executor.Flush()
}

// Insert inserts doc.
func (bi *BulkInserter) Insert(doc any) {
	bi.executor.Add(doc)
}

// SetResultHandler sets the result handler.
func (bi *BulkInserter) SetResultHandler(handler ResultHandler) {
	bi.executor.Sync(func() {
		bi.inserter.resultHandler = handler
	})
}

type dbInserter struct {
	collection    *mongo.Collection
	documents     []any
	resultHandler ResultHandler
}

func (in *dbInserter) AddTask(doc any) bool {
	in.documents = append(in.documents, doc)
	return len(in.documents) >= maxBulkRows
}

func (in *dbInserter) Execute(objs any) {
	docs := objs.([]any)
	if len(docs) == 0 {
		return
	}

	result, err := in.collection.InsertMany(context.Background(), docs)
	if in.resultHandler != nil {
		in.resultHandler(result, err)
	} else if err != nil {
		logx.Error(err)
	}
}

func (in *dbInserter) RemoveAll() any {
	documents := in.documents
	in.documents = nil
	return documents
}
