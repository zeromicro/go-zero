package mongo

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/tal-tech/go-zero/core/executors"
	"github.com/tal-tech/go-zero/core/logx"
)

const (
	flushInterval = time.Second
	maxBulkRows   = 1000
)

type (
	ResultHandler func(*mgo.BulkResult, error)

	BulkInserter struct {
		executor *executors.PeriodicalExecutor
		inserter *dbInserter
	}
)

func NewBulkInserter(session *mgo.Session, dbName string, collectionNamer func() string) *BulkInserter {
	inserter := &dbInserter{
		session:         session,
		dbName:          dbName,
		collectionNamer: collectionNamer,
	}

	return &BulkInserter{
		executor: executors.NewPeriodicalExecutor(flushInterval, inserter),
		inserter: inserter,
	}
}

func (bi *BulkInserter) Flush() {
	bi.executor.Flush()
}

func (bi *BulkInserter) Insert(doc interface{}) {
	bi.executor.Add(doc)
}

func (bi *BulkInserter) SetResultHandler(handler ResultHandler) {
	bi.executor.Sync(func() {
		bi.inserter.resultHandler = handler
	})
}

type dbInserter struct {
	session         *mgo.Session
	dbName          string
	collectionNamer func() string
	documents       []interface{}
	resultHandler   ResultHandler
}

func (in *dbInserter) AddTask(doc interface{}) bool {
	in.documents = append(in.documents, doc)
	return len(in.documents) >= maxBulkRows
}

func (in *dbInserter) Execute(objs interface{}) {
	docs := objs.([]interface{})
	if len(docs) == 0 {
		return
	}

	bulk := in.session.DB(in.dbName).C(in.collectionNamer()).Bulk()
	bulk.Insert(docs...)
	bulk.Unordered()
	result, err := bulk.Run()
	if in.resultHandler != nil {
		in.resultHandler(result, err)
	} else if err != nil {
		logx.Error(err)
	}
}

func (in *dbInserter) RemoveAll() interface{} {
	documents := in.documents
	in.documents = nil
	return documents
}
