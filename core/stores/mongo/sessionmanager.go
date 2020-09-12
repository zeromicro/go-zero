package mongo

import (
	"io"
	"time"

	"github.com/globalsign/mgo"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/syncx"
)

const (
	defaultConcurrency = 50
	defaultTimeout     = time.Second
)

var sessionManager = syncx.NewResourceManager()

type concurrentSession struct {
	*mgo.Session
	limit syncx.TimeoutLimit
}

func (cs *concurrentSession) Close() error {
	cs.Session.Close()
	return nil
}

func getConcurrentSession(url string) (*concurrentSession, error) {
	val, err := sessionManager.GetResource(url, func() (io.Closer, error) {
		mgoSession, err := mgo.Dial(url)
		if err != nil {
			return nil, err
		}

		concurrentSess := &concurrentSession{
			Session: mgoSession,
			limit:   syncx.NewTimeoutLimit(defaultConcurrency),
		}

		return concurrentSess, nil
	})
	if err != nil {
		return nil, err
	}

	return val.(*concurrentSession), nil
}

func (cs *concurrentSession) putSession(session *mgo.Session) {
	if err := cs.limit.Return(); err != nil {
		logx.Error(err)
	}

	// anyway, we need to close the session
	session.Close()
}

func (cs *concurrentSession) takeSession(opts ...Option) (*mgo.Session, error) {
	o := &options{
		timeout: defaultTimeout,
	}
	for _, opt := range opts {
		opt(o)
	}

	if err := cs.limit.Borrow(o.timeout); err != nil {
		return nil, err
	} else {
		return cs.Copy(), nil
	}
}
