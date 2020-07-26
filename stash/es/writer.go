package es

import (
	"context"
	"time"

	"zero/core/executors"
	"zero/core/logx"
	"zero/stash/config"

	"github.com/olivere/elastic"
)

const docType = "doc"

type (
	Writer struct {
		client   *elastic.Client
		indexer  *Index
		inserter *executors.ChunkExecutor
	}

	valueWithTime struct {
		t   time.Time
		val string
	}
)

func NewWriter(c config.ElasticSearchConf, indexer *Index) (*Writer, error) {
	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(c.Hosts...),
		elastic.SetGzip(c.Compress),
	)
	if err != nil {
		return nil, err
	}

	writer := Writer{
		client:  client,
		indexer: indexer,
	}
	writer.inserter = executors.NewChunkExecutor(writer.execute, executors.WithChunkBytes(c.MaxChunkBytes))
	return &writer, nil
}

func (w *Writer) Write(t time.Time, val string) error {
	return w.inserter.Add(valueWithTime{
		t:   t,
		val: val,
	}, len(val))
}

func (w *Writer) execute(vals []interface{}) {
	var bulk = w.client.Bulk()
	for _, val := range vals {
		pair := val.(valueWithTime)
		req := elastic.NewBulkIndexRequest().Index(w.indexer.GetIndex(pair.t)).Type(docType).Doc(pair.val)
		bulk.Add(req)
	}
	_, err := bulk.Do(context.Background())
	if err != nil {
		logx.Error(err)
	}
}
