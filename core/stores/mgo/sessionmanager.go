package mgo

import (
	"context"
	"io"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/syncx"
	"go.mongodb.org/mongo-driver/mongo"
	mopt "go.mongodb.org/mongo-driver/mongo/options"
)

const (
	defaultConcurrency = 50
	defaultTimeout     = time.Second
)

var sessionManager = syncx.NewResourceManager()

type concurrentSession struct {
	*mongo.Client
	limit syncx.TimeoutLimit
}

func (cs *concurrentSession) Close() error {
	return cs.Client.Disconnect(context.Background())
}

func getConcurrentSession(url string) (*concurrentSession, error) {
	val, err := sessionManager.GetResource(url, func() (io.Closer, error) {
		cli, err := mongo.Connect(context.Background(), mopt.Client().ApplyURI(url))
		if err != nil {
			return nil, err
		}

		concurrentSess := &concurrentSession{
			Client: cli,
			limit:  syncx.NewTimeoutLimit(defaultConcurrency),
		}

		return concurrentSess, nil
	})
	if err != nil {
		return nil, err
	}

	return val.(*concurrentSession), nil
}

func (cs *concurrentSession) putSession(session mongo.Session) {
	if err := cs.limit.Return(); err != nil {
		logx.Error(err)
	}

	// anyway, we need to close the session
	session.EndSession(context.Background())
}

func (cs *concurrentSession) takeSession(opts ...Option) (mongo.Session, error) {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	if err := cs.limit.Borrow(o.timeout); err != nil {
		return nil, err
	}

	return cs.Client.StartSession()
}
