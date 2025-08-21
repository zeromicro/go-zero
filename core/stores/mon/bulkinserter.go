//go:generate mockgen -package mon -destination collectioninserter_mock.go -source bulkinserter.go collectionInserter
package mon

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/executors"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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
	cloneColl := coll.Clone()

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

type collectionInserter interface {
	InsertMany(
		ctx context.Context,
		documents interface{},
		opts ...options.Lister[options.InsertManyOptions],
	) (*mongo.InsertManyResult, error)
}

type dbInserter struct {
	collection    collectionInserter
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
