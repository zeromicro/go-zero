package mongo

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/tal-tech/go-zero/core/breaker"
)

type (
	Query interface {
		All(result interface{}) error
		Apply(change mgo.Change, result interface{}) (*mgo.ChangeInfo, error)
		Batch(n int) Query
		Collation(collation *mgo.Collation) Query
		Comment(comment string) Query
		Count() (int, error)
		Distinct(key string, result interface{}) error
		Explain(result interface{}) error
		For(result interface{}, f func() error) error
		Hint(indexKey ...string) Query
		Iter() Iter
		Limit(n int) Query
		LogReplay() Query
		MapReduce(job *mgo.MapReduce, result interface{}) (*mgo.MapReduceInfo, error)
		One(result interface{}) error
		Prefetch(p float64) Query
		Select(selector interface{}) Query
		SetMaxScan(n int) Query
		SetMaxTime(d time.Duration) Query
		Skip(n int) Query
		Snapshot() Query
		Sort(fields ...string) Query
		Tail(timeout time.Duration) Iter
	}

	promisedQuery struct {
		*mgo.Query
		promise keepablePromise
	}

	rejectedQuery struct{}
)

func (q promisedQuery) All(result interface{}) error {
	return q.promise.keep(q.Query.All(result))
}

func (q promisedQuery) Apply(change mgo.Change, result interface{}) (*mgo.ChangeInfo, error) {
	info, err := q.Query.Apply(change, result)
	return info, q.promise.keep(err)
}

func (q promisedQuery) Batch(n int) Query {
	return promisedQuery{
		Query:   q.Query.Batch(n),
		promise: q.promise,
	}
}

func (q promisedQuery) Collation(collation *mgo.Collation) Query {
	return promisedQuery{
		Query:   q.Query.Collation(collation),
		promise: q.promise,
	}
}

func (q promisedQuery) Comment(comment string) Query {
	return promisedQuery{
		Query:   q.Query.Comment(comment),
		promise: q.promise,
	}
}

func (q promisedQuery) Count() (int, error) {
	v, err := q.Query.Count()
	return v, q.promise.keep(err)
}

func (q promisedQuery) Distinct(key string, result interface{}) error {
	return q.promise.keep(q.Query.Distinct(key, result))
}

func (q promisedQuery) Explain(result interface{}) error {
	return q.promise.keep(q.Query.Explain(result))
}

func (q promisedQuery) For(result interface{}, f func() error) error {
	var ferr error
	err := q.Query.For(result, func() error {
		ferr = f()
		return ferr
	})
	if ferr == err {
		return q.promise.accept(err)
	}

	return q.promise.keep(err)
}

func (q promisedQuery) Hint(indexKey ...string) Query {
	return promisedQuery{
		Query:   q.Query.Hint(indexKey...),
		promise: q.promise,
	}
}

func (q promisedQuery) Iter() Iter {
	return promisedIter{
		Iter:    q.Query.Iter(),
		promise: q.promise,
	}
}

func (q promisedQuery) Limit(n int) Query {
	return promisedQuery{
		Query:   q.Query.Limit(n),
		promise: q.promise,
	}
}

func (q promisedQuery) LogReplay() Query {
	return promisedQuery{
		Query:   q.Query.LogReplay(),
		promise: q.promise,
	}
}

func (q promisedQuery) MapReduce(job *mgo.MapReduce, result interface{}) (*mgo.MapReduceInfo, error) {
	info, err := q.Query.MapReduce(job, result)
	return info, q.promise.keep(err)
}

func (q promisedQuery) One(result interface{}) error {
	return q.promise.keep(q.Query.One(result))
}

func (q promisedQuery) Prefetch(p float64) Query {
	return promisedQuery{
		Query:   q.Query.Prefetch(p),
		promise: q.promise,
	}
}

func (q promisedQuery) Select(selector interface{}) Query {
	return promisedQuery{
		Query:   q.Query.Select(selector),
		promise: q.promise,
	}
}

func (q promisedQuery) SetMaxScan(n int) Query {
	return promisedQuery{
		Query:   q.Query.SetMaxScan(n),
		promise: q.promise,
	}
}

func (q promisedQuery) SetMaxTime(d time.Duration) Query {
	return promisedQuery{
		Query:   q.Query.SetMaxTime(d),
		promise: q.promise,
	}
}

func (q promisedQuery) Skip(n int) Query {
	return promisedQuery{
		Query:   q.Query.Skip(n),
		promise: q.promise,
	}
}

func (q promisedQuery) Snapshot() Query {
	return promisedQuery{
		Query:   q.Query.Snapshot(),
		promise: q.promise,
	}
}

func (q promisedQuery) Sort(fields ...string) Query {
	return promisedQuery{
		Query:   q.Query.Sort(fields...),
		promise: q.promise,
	}
}

func (q promisedQuery) Tail(timeout time.Duration) Iter {
	return promisedIter{
		Iter:    q.Query.Tail(timeout),
		promise: q.promise,
	}
}

func (q rejectedQuery) All(result interface{}) error {
	return breaker.ErrServiceUnavailable
}

func (q rejectedQuery) Apply(change mgo.Change, result interface{}) (*mgo.ChangeInfo, error) {
	return nil, breaker.ErrServiceUnavailable
}

func (q rejectedQuery) Batch(n int) Query {
	return q
}

func (q rejectedQuery) Collation(collation *mgo.Collation) Query {
	return q
}

func (q rejectedQuery) Comment(comment string) Query {
	return q
}

func (q rejectedQuery) Count() (int, error) {
	return 0, breaker.ErrServiceUnavailable
}

func (q rejectedQuery) Distinct(key string, result interface{}) error {
	return breaker.ErrServiceUnavailable
}

func (q rejectedQuery) Explain(result interface{}) error {
	return breaker.ErrServiceUnavailable
}

func (q rejectedQuery) For(result interface{}, f func() error) error {
	return breaker.ErrServiceUnavailable
}

func (q rejectedQuery) Hint(indexKey ...string) Query {
	return q
}

func (q rejectedQuery) Iter() Iter {
	return rejectedIter{}
}

func (q rejectedQuery) Limit(n int) Query {
	return q
}

func (q rejectedQuery) LogReplay() Query {
	return q
}

func (q rejectedQuery) MapReduce(job *mgo.MapReduce, result interface{}) (*mgo.MapReduceInfo, error) {
	return nil, breaker.ErrServiceUnavailable
}

func (q rejectedQuery) One(result interface{}) error {
	return breaker.ErrServiceUnavailable
}

func (q rejectedQuery) Prefetch(p float64) Query {
	return q
}

func (q rejectedQuery) Select(selector interface{}) Query {
	return q
}

func (q rejectedQuery) SetMaxScan(n int) Query {
	return q
}

func (q rejectedQuery) SetMaxTime(d time.Duration) Query {
	return q
}

func (q rejectedQuery) Skip(n int) Query {
	return q
}

func (q rejectedQuery) Snapshot() Query {
	return q
}

func (q rejectedQuery) Sort(fields ...string) Query {
	return q
}

func (q rejectedQuery) Tail(timeout time.Duration) Iter {
	return rejectedIter{}
}
