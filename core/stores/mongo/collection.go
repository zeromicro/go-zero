package mongo

import (
	"encoding/json"
	"time"

	"github.com/globalsign/mgo"
	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/mongo/internal"
	"github.com/zeromicro/go-zero/core/timex"
)

const defaultSlowThreshold = time.Millisecond * 500

// ErrNotFound is an alias of mgo.ErrNotFound.
var ErrNotFound = mgo.ErrNotFound

type (
	// Collection interface represents a mongo connection.
	Collection interface {
		Find(query any) Query
		FindId(id any) Query
		Insert(docs ...any) error
		Pipe(pipeline any) Pipe
		Remove(selector any) error
		RemoveAll(selector any) (*mgo.ChangeInfo, error)
		RemoveId(id any) error
		Update(selector, update any) error
		UpdateId(id, update any) error
		Upsert(selector, update any) (*mgo.ChangeInfo, error)
	}

	decoratedCollection struct {
		name       string
		collection internal.MgoCollection
		brk        breaker.Breaker
	}

	keepablePromise struct {
		promise breaker.Promise
		log     func(error)
	}
)

func newCollection(collection *mgo.Collection, brk breaker.Breaker) Collection {
	return &decoratedCollection{
		name:       collection.FullName,
		collection: collection,
		brk:        brk,
	}
}

func (c *decoratedCollection) Find(query any) Query {
	promise, err := c.brk.Allow()
	if err != nil {
		return rejectedQuery{}
	}

	startTime := timex.Now()
	return promisedQuery{
		Query: c.collection.Find(query),
		promise: keepablePromise{
			promise: promise,
			log: func(err error) {
				duration := timex.Since(startTime)
				c.logDuration("find", duration, err, query)
			},
		},
	}
}

func (c *decoratedCollection) FindId(id any) Query {
	promise, err := c.brk.Allow()
	if err != nil {
		return rejectedQuery{}
	}

	startTime := timex.Now()
	return promisedQuery{
		Query: c.collection.FindId(id),
		promise: keepablePromise{
			promise: promise,
			log: func(err error) {
				duration := timex.Since(startTime)
				c.logDuration("findId", duration, err, id)
			},
		},
	}
}

func (c *decoratedCollection) Insert(docs ...any) (err error) {
	return c.brk.DoWithAcceptable(func() error {
		startTime := timex.Now()
		defer func() {
			duration := timex.Since(startTime)
			c.logDuration("insert", duration, err, docs...)
		}()

		return c.collection.Insert(docs...)
	}, acceptable)
}

func (c *decoratedCollection) Pipe(pipeline any) Pipe {
	promise, err := c.brk.Allow()
	if err != nil {
		return rejectedPipe{}
	}

	startTime := timex.Now()
	return promisedPipe{
		Pipe: c.collection.Pipe(pipeline),
		promise: keepablePromise{
			promise: promise,
			log: func(err error) {
				duration := timex.Since(startTime)
				c.logDuration("pipe", duration, err, pipeline)
			},
		},
	}
}

func (c *decoratedCollection) Remove(selector any) (err error) {
	return c.brk.DoWithAcceptable(func() error {
		startTime := timex.Now()
		defer func() {
			duration := timex.Since(startTime)
			c.logDuration("remove", duration, err, selector)
		}()

		return c.collection.Remove(selector)
	}, acceptable)
}

func (c *decoratedCollection) RemoveAll(selector any) (info *mgo.ChangeInfo, err error) {
	err = c.brk.DoWithAcceptable(func() error {
		startTime := timex.Now()
		defer func() {
			duration := timex.Since(startTime)
			c.logDuration("removeAll", duration, err, selector)
		}()

		info, err = c.collection.RemoveAll(selector)
		return err
	}, acceptable)

	return
}

func (c *decoratedCollection) RemoveId(id any) (err error) {
	return c.brk.DoWithAcceptable(func() error {
		startTime := timex.Now()
		defer func() {
			duration := timex.Since(startTime)
			c.logDuration("removeId", duration, err, id)
		}()

		return c.collection.RemoveId(id)
	}, acceptable)
}

func (c *decoratedCollection) Update(selector, update any) (err error) {
	return c.brk.DoWithAcceptable(func() error {
		startTime := timex.Now()
		defer func() {
			duration := timex.Since(startTime)
			c.logDuration("update", duration, err, selector, update)
		}()

		return c.collection.Update(selector, update)
	}, acceptable)
}

func (c *decoratedCollection) UpdateId(id, update any) (err error) {
	return c.brk.DoWithAcceptable(func() error {
		startTime := timex.Now()
		defer func() {
			duration := timex.Since(startTime)
			c.logDuration("updateId", duration, err, id, update)
		}()

		return c.collection.UpdateId(id, update)
	}, acceptable)
}

func (c *decoratedCollection) Upsert(selector, update any) (info *mgo.ChangeInfo, err error) {
	err = c.brk.DoWithAcceptable(func() error {
		startTime := timex.Now()
		defer func() {
			duration := timex.Since(startTime)
			c.logDuration("upsert", duration, err, selector, update)
		}()

		info, err = c.collection.Upsert(selector, update)
		return err
	}, acceptable)

	return
}

func (c *decoratedCollection) logDuration(method string, duration time.Duration, err error, docs ...any) {
	content, e := json.Marshal(docs)
	if e != nil {
		logx.Error(err)
	} else if err != nil {
		if duration > slowThreshold.Load() {
			logx.WithDuration(duration).Slowf("[MONGO] mongo(%s) - slowcall - %s - fail(%s) - %s",
				c.name, method, err.Error(), string(content))
		} else {
			logx.WithDuration(duration).Infof("mongo(%s) - %s - fail(%s) - %s",
				c.name, method, err.Error(), string(content))
		}
	} else {
		if duration > slowThreshold.Load() {
			logx.WithDuration(duration).Slowf("[MONGO] mongo(%s) - slowcall - %s - ok - %s",
				c.name, method, string(content))
		} else {
			logx.WithDuration(duration).Infof("mongo(%s) - %s - ok - %s", c.name, method, string(content))
		}
	}
}

func (p keepablePromise) accept(err error) error {
	p.promise.Accept()
	p.log(err)
	return err
}

func (p keepablePromise) keep(err error) error {
	if acceptable(err) {
		p.promise.Accept()
	} else {
		p.promise.Reject(err.Error())
	}

	p.log(err)
	return err
}

func acceptable(err error) bool {
	return err == nil || err == mgo.ErrNotFound
}
