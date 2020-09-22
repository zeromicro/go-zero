package mongo

import (
	"encoding/json"
	"time"

	"github.com/globalsign/mgo"
	"github.com/tal-tech/go-zero/core/breaker"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/timex"
)

const slowThreshold = time.Millisecond * 500

var ErrNotFound = mgo.ErrNotFound

type (
	Collection interface {
		Find(query interface{}) Query
		FindId(id interface{}) Query
		Insert(docs ...interface{}) error
		Pipe(pipeline interface{}) Pipe
		Remove(selector interface{}) error
		RemoveAll(selector interface{}) (*mgo.ChangeInfo, error)
		RemoveId(id interface{}) error
		Update(selector, update interface{}) error
		UpdateId(id, update interface{}) error
		Upsert(selector, update interface{}) (*mgo.ChangeInfo, error)
	}

	decoratedCollection struct {
		*mgo.Collection
		brk breaker.Breaker
	}

	keepablePromise struct {
		promise breaker.Promise
		log     func(error)
	}
)

func newCollection(collection *mgo.Collection) Collection {
	return &decoratedCollection{
		Collection: collection,
		brk:        breaker.NewBreaker(),
	}
}

func (c *decoratedCollection) Find(query interface{}) Query {
	promise, err := c.brk.Allow()
	if err != nil {
		return rejectedQuery{}
	}

	startTime := timex.Now()
	return promisedQuery{
		Query: c.Collection.Find(query),
		promise: keepablePromise{
			promise: promise,
			log: func(err error) {
				duration := timex.Since(startTime)
				c.logDuration("find", duration, err, query)
			},
		},
	}
}

func (c *decoratedCollection) FindId(id interface{}) Query {
	promise, err := c.brk.Allow()
	if err != nil {
		return rejectedQuery{}
	}

	startTime := timex.Now()
	return promisedQuery{
		Query: c.Collection.FindId(id),
		promise: keepablePromise{
			promise: promise,
			log: func(err error) {
				duration := timex.Since(startTime)
				c.logDuration("findId", duration, err, id)
			},
		},
	}
}

func (c *decoratedCollection) Insert(docs ...interface{}) (err error) {
	return c.brk.DoWithAcceptable(func() error {
		startTime := timex.Now()
		defer func() {
			duration := timex.Since(startTime)
			c.logDuration("insert", duration, err, docs...)
		}()

		return c.Collection.Insert(docs...)
	}, acceptable)
}

func (c *decoratedCollection) Pipe(pipeline interface{}) Pipe {
	promise, err := c.brk.Allow()
	if err != nil {
		return rejectedPipe{}
	}

	startTime := timex.Now()
	return promisedPipe{
		Pipe: c.Collection.Pipe(pipeline),
		promise: keepablePromise{
			promise: promise,
			log: func(err error) {
				duration := timex.Since(startTime)
				c.logDuration("pipe", duration, err, pipeline)
			},
		},
	}
}

func (c *decoratedCollection) Remove(selector interface{}) (err error) {
	return c.brk.DoWithAcceptable(func() error {
		startTime := timex.Now()
		defer func() {
			duration := timex.Since(startTime)
			c.logDuration("remove", duration, err, selector)
		}()

		return c.Collection.Remove(selector)
	}, acceptable)
}

func (c *decoratedCollection) RemoveAll(selector interface{}) (info *mgo.ChangeInfo, err error) {
	err = c.brk.DoWithAcceptable(func() error {
		startTime := timex.Now()
		defer func() {
			duration := timex.Since(startTime)
			c.logDuration("removeAll", duration, err, selector)
		}()

		info, err = c.Collection.RemoveAll(selector)
		return err
	}, acceptable)

	return
}

func (c *decoratedCollection) RemoveId(id interface{}) (err error) {
	return c.brk.DoWithAcceptable(func() error {
		startTime := timex.Now()
		defer func() {
			duration := timex.Since(startTime)
			c.logDuration("removeId", duration, err, id)
		}()

		return c.Collection.RemoveId(id)
	}, acceptable)
}

func (c *decoratedCollection) Update(selector, update interface{}) (err error) {
	return c.brk.DoWithAcceptable(func() error {
		startTime := timex.Now()
		defer func() {
			duration := timex.Since(startTime)
			c.logDuration("update", duration, err, selector, update)
		}()

		return c.Collection.Update(selector, update)
	}, acceptable)
}

func (c *decoratedCollection) UpdateId(id, update interface{}) (err error) {
	return c.brk.DoWithAcceptable(func() error {
		startTime := timex.Now()
		defer func() {
			duration := timex.Since(startTime)
			c.logDuration("updateId", duration, err, id, update)
		}()

		return c.Collection.UpdateId(id, update)
	}, acceptable)
}

func (c *decoratedCollection) Upsert(selector, update interface{}) (info *mgo.ChangeInfo, err error) {
	err = c.brk.DoWithAcceptable(func() error {
		startTime := timex.Now()
		defer func() {
			duration := timex.Since(startTime)
			c.logDuration("upsert", duration, err, selector, update)
		}()

		info, err = c.Collection.Upsert(selector, update)
		return err
	}, acceptable)

	return
}

func (c *decoratedCollection) logDuration(method string, duration time.Duration, err error, docs ...interface{}) {
	content, e := json.Marshal(docs)
	if e != nil {
		logx.Error(err)
	} else if err != nil {
		if duration > slowThreshold {
			logx.WithDuration(duration).Slowf("[MONGO] mongo(%s) - slowcall - %s - fail(%s) - %s",
				c.FullName, method, err.Error(), string(content))
		} else {
			logx.WithDuration(duration).Infof("mongo(%s) - %s - fail(%s) - %s",
				c.FullName, method, err.Error(), string(content))
		}
	} else {
		if duration > slowThreshold {
			logx.WithDuration(duration).Slowf("[MONGO] mongo(%s) - slowcall - %s - ok - %s",
				c.FullName, method, string(content))
		} else {
			logx.WithDuration(duration).Infof("mongo(%s) - %s - ok - %s", c.FullName, method, string(content))
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
