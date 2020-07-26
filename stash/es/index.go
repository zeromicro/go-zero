package es

import (
	"context"
	"sync"
	"time"

	"zero/core/fx"
	"zero/core/logx"
	"zero/core/syncx"

	"github.com/olivere/elastic"
)

const sharedCallsKey = "ensureIndex"

type (
	IndexFormat func(time.Time) string
	IndexFunc   func() string

	Index struct {
		client      *elastic.Client
		indexFormat IndexFormat
		index       string
		lock        sync.RWMutex
		sharedCalls syncx.SharedCalls
	}
)

func NewIndex(client *elastic.Client, indexFormat IndexFormat) *Index {
	return &Index{
		client:      client,
		indexFormat: indexFormat,
		sharedCalls: syncx.NewSharedCalls(),
	}
}

func (idx *Index) GetIndex(t time.Time) string {
	index := idx.indexFormat(t)
	if err := idx.ensureIndex(index); err != nil {
		logx.Error(err)
	}
	return index
}

func (idx *Index) ensureIndex(index string) error {
	idx.lock.RLock()
	if index == idx.index {
		idx.lock.RUnlock()
		return nil
	}
	idx.lock.RUnlock()

	_, err := idx.sharedCalls.Do(sharedCallsKey, func() (i interface{}, err error) {
		idx.lock.Lock()
		defer idx.lock.Unlock()

		existsService := elastic.NewIndicesExistsService(idx.client)
		existsService.Index([]string{index})
		exist, err := existsService.Do(context.Background())
		if err != nil {
			return nil, err
		}
		if exist {
			idx.index = index
			return nil, nil
		}

		createService := idx.client.CreateIndex(index)
		if err := fx.DoWithRetries(func() error {
			// is it necessary to check the result?
			_, err := createService.Do(context.Background())
			return err
		}); err != nil {
			return nil, err
		}

		idx.index = index
		return nil, nil
	})
	return err
}
